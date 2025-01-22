package parser

import (
	"errors"
)

// Parser interface defines the methods that a parser must implement.
type Parser interface {
	SupportedContentTypes() []string
	Parse(data []byte) (string, error)
}

// Global map to hold registered parsers
var parsers = make(map[string]Parser)

// RegisterParser registers a parser for the given content types.
func RegisterParser(p Parser) {
	for _, contentType := range p.SupportedContentTypes() {
		parsers[contentType] = p
	}
}

// GetParser retrieves a parser for the specified content type.
func GetParser(contentType string) (Parser, error) {
	parser, ok := parsers[contentType]
	if !ok {
		return nil, errors.New("no parser found for content type: " + contentType)
	}
	return parser, nil
}

// ExtractText extracts text from the given file data and content type.
func ExtractText(contentType string, data []byte) (string, error) {
	parser, err := GetParser(contentType)
	if err != nil {
		return "", err
	}
	return parser.Parse(data)
}

// Example of a PDF parser implementation
type PDFParser struct{}

func (p *PDFParser) SupportedContentTypes() []string {
	return []string{"application/pdf"}
}

func (p *PDFParser) Parse(data []byte) (string, error) {
	// Implement PDF parsing logic here
	return "Parsed PDF content", nil
}

// Example of a DOCX parser implementation
type DOCXParser struct{}

func (d *DOCXParser) SupportedContentTypes() []string {
	return []string{"application/vnd.openxmlformats-officedocument.wordprocessingml.document"}
}

func (d *DOCXParser) Parse(data []byte) (string, error) {
	// Implement DOCX parsing logic here
	return "Parsed DOCX content", nil
}

// Example of a PPTX parser implementation
type PPTXParser struct{}

func (p *PPTXParser) SupportedContentTypes() []string {
	return []string{"application/vnd.openxmlformats-officedocument.presentationml.presentation"}
}

func (p *PPTXParser) Parse(data []byte) (string, error) {
	// Implement PPTX parsing logic here
	return "Parsed PPTX content", nil
}

// init function to register parsers
func init() {
	RegisterParser(&PDFParser{})
	RegisterParser(&DOCXParser{})
	RegisterParser(&PPTXParser{})
}
