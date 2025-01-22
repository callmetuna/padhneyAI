package repositories

import (
	"errors"
	"rag-go-app/models"
)

// FileRepository defines the interface for file-related data operations.
type FileRepository interface {
	SaveFile(file *models.File) error
	FindFileByID(id int64) (*models.File, error)
	DeleteFile(id int64) error
}

// InMemoryFileRepository is an in-memory implementation of FileRepository for demonstration purposes.
type InMemoryFileRepository struct {
	files  map[int64]*models.File
	nextID int64
}

// NewInMemoryFileRepository creates a new instance of InMemoryFileRepository.
func NewInMemoryFileRepository() *InMemoryFileRepository {
	return &InMemoryFileRepository{
		files:  make(map[int64]*models.File),
		nextID: 1,
	}
}

// SaveFile saves a new file to the repository.
func (r *InMemoryFileRepository) SaveFile(file *models.File) error {
	file.ID = r.nextID
	r.files[r.nextID] = file
	r.nextID++
	return nil
}

// FindFileByID retrieves a file by its ID.
func (r *InMemoryFileRepository) FindFileByID(id int64) (*models.File, error) {
	file, exists := r.files[id]
	if !exists {
		return nil, errors.New("file not found")
	}
	return file, nil
}

// DeleteFile deletes a file from the repository.
func (r *InMemoryFileRepository) DeleteFile(id int64) error {
	if _, exists := r.files[id]; !exists {
		return errors.New("file not found")
	}
	delete(r.files, id)
	return nil
}
