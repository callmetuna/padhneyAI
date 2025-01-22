package service

import (
	"errors"
	"os"
	"time"

	"rag-go-app/models"
	"rag-go-app/repositories"
)

type FileService struct {
	repo repositories.FileRepository
}

// NewFileService creates a new instance of FileService.
func NewFileService(repo repositories.FileRepository) *FileService {
	return &FileService{repo: repo}
}

// UploadFile handles the file upload and creates a new File record.
func (s *FileService) UploadFile(filename string, contentType string) (*models.File, error) {
	if filename == "" || contentType == "" {
		return nil, errors.New("filename and content type cannot be empty")
	}

	// Create a new File model
	file := &models.File{
		Filename:    filename,
		ContentType: contentType,
		Processed:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Save the file record to the repository
	if err := s.repo.SaveFile(file); err != nil {
		return nil, err
	}

	return file, nil
}

// GetFile retrieves a file by its ID.
func (s *FileService) GetFile(id int64) (*models.File, error) {
	file, err := s.repo.FindFileByID(id)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// DeleteFile deletes a file record by its ID.
func (s *FileService) DeleteFile(id int64) error {
	file, err := s.repo.FindFileByID(id)
	if err != nil {
		return err
	}

	// Optionally, you can delete the actual file from the filesystem
	if err := os.Remove(file.Filename); err != nil {
		return err
	}

	return s.repo.DeleteFile(id)
}
