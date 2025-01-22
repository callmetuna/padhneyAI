package service

import (
	"errors"
	"os"

	"rag-go-app/models"
	"rag-go-app/parser"
	"rag-go-app/repositories"
)

type DataService struct {
	repo      repositories.DataRepository
	pdfParser *parser.PDFParser // Add PDFParser instance
}

// NewDataService creates a new instance of DataService.
func NewDataService(repo repositories.DataRepository, parser *parser.PDFParser) *DataService { // Add parser dependency
	return &DataService{repo: repo, pdfParser: parser} // Initialize parser
}

// CreateDocument processes a file and creates a new document in the repository.
func (s *DataService) CreateDocument(filePath string) (*models.Document, error) {
	// Read the file content
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Parse the PDF content
	doc, err := s.pdfParser.Parse(data) // Use the Parse method
	if err != nil {
		return nil, err
	}

	// Save the document to the repository
	if err := s.repo.SaveDocument(doc); err != nil {
		return nil, err
	}

	return doc, nil
}

// GetDocument retrieves a document by its ID.
func (s *DataService) GetDocument(id int64) (*models.Document, error) {
	doc, err := s.repo.FindDocumentByID(id)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// UpdateDocument updates an existing document.
func (s *DataService) UpdateDocument(id int64, text string, metadata models.Metadata) error {
	doc, err := s.repo.FindDocumentByID(id)
	if err != nil {
		return err // Return the error directly
	}
	if doc == nil {
		return errors.New("document not found") // Handle the case where the document is not found
	}

	// Update the document
	if err := doc.Update(text, metadata); err != nil {
		return err
	}

	// Save the updated document back to the repository
	return s.repo.UpdateDocument(doc)
}
