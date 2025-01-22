package parser

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
	"os/exec"
	"strings"

	"rag-go-app/models"

	"github.com/otiai10/gosseract/v2"
	"github.com/unidoc/unidoc/v3/pdf/model"
)

type PDFParser struct {
	OCRLanguage string
}

func (p *PDFParser) SupportedContentTypes() []string {
	return []string{"application/pdf"}
}

func (p *PDFParser) Parse(data []byte) (*models.Document, error) {
	doc := &models.Document{}
	reader, err := model.NewPdfReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("creating PDF reader failed: %w", err)
	}

	numPages, err := reader.GetNumPages()
	if err != nil {
		return nil, fmt.Errorf("getting number of pages failed: %w", err)
	}

	// Extract metadata
	metadata, err := extractMetadata(reader)
	if err != nil {
		log.Printf("Metadata extraction failed: %v", err)
	}
	doc.Metadata = metadata

	// Process each page
	for i := 1; i <= numPages; i++ {
		page, err := reader.GetPage(i)
		if err != nil {
			return nil, fmt.Errorf("getting page %d failed: %w", i, err)
		}

		text, err := extractTextFromPage(page, p.OCRLanguage, data)
		if err != nil {
			log.Printf("Text extraction failed for page %d: %v", i, err)
		}
		doc.Text += text

		tables, err := extractTables(data)
		if err != nil {
			log.Printf("Table extraction failed: %v", err)
		}
		doc.Metadata.Tables = append(doc.Metadata.Tables, tables...)

		figures, err := extractFigures(page)
		if err != nil {
			log.Printf("Figure extraction failed: %v", err)
		}
		doc.Metadata.Figures = append(doc.Metadata.Figures, figures...)
	}

	newDoc, err := models.NewDocument(0, doc.Text, doc.Metadata, nil)
	if err != nil {
		return nil, err
	}
	return newDoc, nil
}

func extractTextFromPage(page *model.PdfPage, language string, pdfData []byte) (string, error) {
	// Extract text using UniDoc
	text, err := extractTextUnidoc(page)
	if err == nil && len(text) > 10 {
		return text, nil
	}

	// Fallback to pdfminer.six
	text, err = extractTextPdfminer(pdfData)
	if err == nil && len(text) > 10 {
		return text, nil
	}

	// Fallback to OCR
	images, err := page.GetImages()
	if err != nil {
		return "", fmt.Errorf("getting images for OCR failed: %w", err)
	}
	text, err = extractTextOCR(images, language)
	if err != nil {
		return "", fmt.Errorf("OCR failed: %w", err)
	}
	return text, nil
}

func extractTextUnidoc(page *model.PdfPage) (string, error) {
	contentStreams, err := page.GetContentStreams()
	if err != nil {
		return "", fmt.Errorf("failed to get content streams: %w", err)
	}
	var textContent strings.Builder
	for _, stream := range contentStreams {
		textContent.WriteString(stream) // Custom parsing logic may be required
	}
	return textContent.String(), nil
}

func extractTextPdfminer(pdfData []byte) (string, error) {
	tmpFile, err := os.CreateTemp("", "temp.pdf")
	if err != nil {
		return "", fmt.Errorf("creating temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(pdfData); err != nil {
		return "", fmt.Errorf("writing to temp file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return "", fmt.Errorf("closing temp file failed: %w", err)
	}

	cmd := exec.Command("pdf2txt.py", "-t", "xml", tmpFile.Name())
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("pdf2txt.py execution error: %w, stderr: %s", err, stderr.String())
	}

	return out.String(), nil
}

func extractTextOCR(images []*model.PdfImage, language string) (string, error) {
	client := gosseract.NewClient()
	defer client.Close()

	availableLangs, _ := client.GetAvailableLanguages()
	if !contains(availableLangs, language) {
		log.Printf("Language '%s' not available. Using default language.", language)
		language = "eng"
	}
	client.SetLanguage(language)

	var allText strings.Builder
	for _, img := range images {
		imgData, err := img.GetData()
		if err != nil {
			return "", fmt.Errorf("error getting image data: %w", err)
		}
		tmpFile, err := createTempImageFile(imgData, img.Fmt)
		if err != nil {
			return "", fmt.Errorf("error creating temp image file: %w", err)
		}
		defer os.Remove(tmpFile.Name())

		client.SetImage(tmpFile.Name())
		text, err := client.Text()
		if err != nil {
			return "", fmt.Errorf("OCR processing failed: %w", err)
		}
		allText.WriteString(text + "\n")
	}
	return allText.String(), nil
}

func createTempImageFile(imgData []byte, imgFmt model.ImageFormat) (*os.File, error) {
	tmpFile, err := os.CreateTemp("", "tempImage.*."+getImageExtension(imgFmt))
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpFile.Name())

	switch imgFmt {
	case model.ImageFormatJPEG:
		if _, err := jpeg.Decode(bytes.NewReader(imgData)); err != nil {
			return nil, fmt.Errorf("decoding JPEG image failed: %v", err)
		}
	case model.ImageFormatPNG:
		if _, err := png.Decode(bytes.NewReader(imgData)); err != nil {
			return nil, fmt.Errorf("decoding PNG image failed: %v", err)
		}
	default:
		return nil, fmt.Errorf("unsupported image format: %v", imgFmt)
	}

	if _, err := tmpFile.Write(imgData); err != nil {
		return nil, fmt.Errorf("writing image data to temporary file failed: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		return nil, fmt.Errorf("closing temporary file failed: %v", err)
	}
	return tmpFile, nil
}

func getImageExtension(imgFmt model.ImageFormat) string {
	switch imgFmt {
	case model.ImageFormatJPEG:
		return "jpg"
	case model.ImageFormatPNG:
		return "png"
	default:
		return "bin"
	}
}

func extractTables(pdfData []byte) ([]models.Table, error) {
	tmpFile, err := os.CreateTemp("", "temp.pdf")
	if err != nil {
		return nil, fmt.Errorf("creating temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(pdfData); err != nil {
		return nil, fmt.Errorf("writing to temp file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return nil, fmt.Errorf("closing temp file failed: %w", err)
	}

	cmd := exec.Command("pdf2txt.py", "-t", "xml", tmpFile.Name())
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("pdf2txt.py execution error: %w, stderr: %s", err, stderr.String())
	}

	tables, err := parseTablesFromXML(out.Bytes())
	if err != nil {
		return nil, fmt.Errorf("parsing tables from XML: %w", err)
	}
	return tables, nil
}

func parseTablesFromXML(xmlData []byte) ([]models.Table, error) {
	type LTTextBoxHorizontal struct {
		Text string  `xml:",chardata"`
		X0   float64 `xml:"x0,attr"`
		Y0   float64 `xml:"y0,attr"`
		X1   float64 `xml:"x1,attr"`
		Y1   float64 `xml:"y1,attr"`
	}

	type LTPage struct {
		TextBoxes []LTTextBoxHorizontal `xml:"textbox"`
	}

	type XMLOutput struct {
		Pages []LTPage `xml:"page"`
	}

	var output XMLOutput
	err := xml.Unmarshal(xmlData, &output)
	if err != nil {
		return nil, fmt.Errorf("XML parsing failed: %w", err)
	}

	var tables []models.Table
	for _, page := range output.Pages {
		var currentTable models.Table
		var currentRow []string
		prevY := -1000.0 // Initialize with a very small value

		for _, box := range page.TextBoxes {
			if math.Abs(box.Y0-prevY) > 15 {
				if len(currentRow) > 0 {
					currentTable.Data = append(currentTable.Data, currentRow)
					currentRow = []string{}
				}
			}
			currentRow = append(currentRow, box.Text)
			prevY = box.Y0
		}

		if len(currentRow) > 0 {
			currentTable.Data = append(currentTable.Data, currentRow)
		}
		if len(currentTable.Data) > 0 {
			tables = append(tables, currentTable)
		}
	}
	return tables, nil
}

func extractFigures(page *model.PdfPage) ([]models.Figure, error) {
	images, err := page.GetImages()
	if err != nil {
		return nil, err
	}

	var figures []models.Figure
	for _, img := range images {
		imgData, err := img.GetData()
		if err != nil {
			log.Printf("Error getting image data: %v", err)
			continue // Skip this image
		}
		figures = append(figures, models.Figure{ImageData: imgData})
	}
	return figures, nil
}

func extractMetadata(reader *model.PdfReader) (models.Metadata, error) {
	pdfInfo, err := reader.GetPdfInfo()
	if err != nil {
		return models.Metadata{}, fmt.Errorf("failed to extract metadata: %w", err)
	}

	return models.Metadata{
		Title:    pdfInfo.Title,
		Authors:  pdfInfo.Authors,
		Abstract: pdfInfo.Subject,
	}, nil
}
