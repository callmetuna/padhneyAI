package models

import (
	"errors"
	"regexp"
	"time"
	"unicode"

	"github.com/microcosm-cc/bluemonday"
)

type File struct {
	ID          int64     `json:"id"`
	Filename    string    `json:"filename"`
	ContentType string    `json:"content_type"`
	Processed   bool      `json:"processed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Document struct {
	ID         int64     `json:"id"`
	FileID     int64     `json:"file_id"`
	Text       string    `json:"text"`
	Metadata   Metadata  `json:"metadata"`
	Embeddings []float32 `json:"embeddings"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Metadata struct {
	Title     string     `json:"title"`
	Authors   []string   `json:"authors"`
	Keywords  []string   `json:"keywords"`
	Abstract  string     `json:"abstract"`
	Citations []Citation `json:"citations"`
}

type Citation struct {
	Text   string `json:"text"`
	Author string `json:"author"`
	Year   int    `json:"year"`
	Title  string `json:"title"`
}

func NewDocument(fileID int64, text string, metadata Metadata, embeddings []float32) (*Document, error) {
	if err := validateDocumentInput(text, metadata, embeddings); err != nil {
		return nil, err
	}

	now := time.Now()
	return &Document{
		FileID:     fileID,
		Text:       text,
		Metadata:   metadata,
		Embeddings: embeddings,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

func (d *Document) Update(text string, metadata Metadata) error {
	if err := validateDocumentInput(text, metadata, d.Embeddings); err != nil {
		return err
	}

	d.Text = text
	d.Metadata = metadata
	d.UpdatedAt = time.Now()
	return nil
}

func (m *Metadata) AddCitation(citation Citation) error {
	if err := validateCitation(citation); err != nil {
		return err
	}
	m.Citations = append(m.Citations, citation)
	return nil
}

func (m *Metadata) RemoveCitation(index int) error {
	if index < 0 || index >= len(m.Citations) {
		return errors.New("index out of range")
	}
	m.Citations = append(m.Citations[:index], m.Citations[index+1:]...)
	return nil
}

func validateDocumentInput(text string, metadata Metadata, embeddings []float32) error {
	if len(text) == 0 {
		return errors.New("text cannot be empty")
	}
	if err := validateMetadata(metadata); err != nil {
		return err
	}
	if len(embeddings) == 0 {
		return errors.New("embeddings cannot be empty")
	}
	return nil
}

func validateMetadata(metadata Metadata) error {
	if len(metadata.Title) == 0 {
		return errors.New("title cannot be empty")
	}
	if len(metadata.Authors) == 0 {
		return errors.New("authors cannot be empty")
	}
	for i := range metadata.Authors {
		if err := sanitizeInput(&metadata.Authors[i]); err != nil {
			return err
		}
	}
	if err := sanitizeInput(&metadata.Abstract); err != nil {
		return err
	}
	for i, keyword := range metadata.Keywords {
		if err := sanitizeInput(&metadata.Keywords[i]); err != nil {
			return err
		}
		if len(keyword) > 50 {
			return errors.New("keyword is too long")
		}
		if !isValidKeyword(keyword) {
			return errors.New("keyword contains invalid characters")
		}
	}
	return nil
}

func validateCitation(citation Citation) error {
	if len(citation.Text) == 0 {
		return errors.New("citation text cannot be empty")
	}
	if len(citation.Author) == 0 {
		return errors.New("citation author cannot be empty")
	}
	if citation.Year < 0 {
		return errors.New("citation year must not be negative")
	}
	if len(citation.Title) == 0 {
		return errors.New("citation title cannot be empty")
	}
	if err := sanitizeInput(&citation.Text); err != nil {
		return err
	}
	if err := sanitizeInput(&citation.Author); err != nil {
		return err
	}
	if err := sanitizeInput(&citation.Title); err != nil {
		return err
	}
	return nil
}

func sanitizeInput(input *string) error {
	p := bluemonday.UGCPolicy()
	*input = p.Sanitize(*input)

	if len(*input) == 0 {
		return errors.New("input cannot be empty after sanitization")
	}

	return nil
}

func isValidKeyword(keyword string) bool {
	for _, r := range keyword {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

// Example using regex
func isValidKeywordRegex(keyword string) bool {
	validInput := regexp.MustCompile(`^[\w\s]+$`)
	return validInput.MatchString(keyword)
}
