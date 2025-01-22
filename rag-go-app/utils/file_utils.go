package utils

import (
	"errors"
	"strings"
)

// Supported file types
const (
	PDF  = ".pdf"
	DOCX = ".docx"
	PPTX = ".pptx"
)

// IsSupportedFileType checks if the file type is supported for extraction.
func IsSupportedFileType(filePath string) bool {
	return strings.HasSuffix(filePath, PDF) || strings.HasSuffix(filePath, DOCX) || strings.HasSuffix(filePath, PPTX)
}

// GetFileType returns the file type based on the file extension.
func GetFileType(filePath string) (string, error) {
	if strings.HasSuffix(filePath, PDF) {
		return PDF, nil
	} else if strings.HasSuffix(filePath, DOCX) {
		return DOCX, nil
	} else if strings.HasSuffix(filePath, PPTX) {
		return PPTX, nil
	}
	return "", errors.New("unsupported file format")
}
