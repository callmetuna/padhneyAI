package repositories

import (
	"errors"
	"rag-go-app/models"
)

// DataRepository defines the interface for document-related data operations.
type DataRepository interface {
	SaveDocument(doc *models.Document) error
	FindDocumentByID(id int64) (*models.Document, error)
	UpdateDocument(doc *models.Document) error
}

// InMemoryDataRepository is an in-memory implementation of DataRepository for demonstration purposes.
type InMemoryDataRepository struct {
	documents map[int64]*models.Document
	nextID    int64
}

// NewInMemoryDataRepository creates a new instance of InMemoryDataRepository.
func NewInMemoryDataRepository() *InMemoryDataRepository {
	return &InMemoryDataRepository{
		documents: make(map[int64]*models.Document),
		nextID:    1,
	}
}

// SaveDocument saves a new document to the repository.
func (r *InMemoryDataRepository) SaveDocument(doc *models.Document) error {
	doc.ID = r.nextID
	r.documents[r.nextID] = doc
	r.nextID++
	return nil
}

// FindDocumentByID retrieves a document by its ID.
func (r *InMemoryDataRepository) FindDocumentByID(id int64) (*models.Document, error) {
	doc, exists := r.documents[id]
	if !exists {
		return nil, errors.New("document not found")
	}
	return doc, nil
}

// UpdateDocument updates an existing document in the repository.
func (r *InMemoryDataRepository) UpdateDocument(doc *models.Document) error {
	if _, exists := r.documents[doc.ID]; !exists {
		return errors.New("document not found")
	}
	r.documents[doc.ID] = doc
	return nil
}
