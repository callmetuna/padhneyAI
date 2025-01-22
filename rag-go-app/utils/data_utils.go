package utils

import (
	"errors"
	"strings"
)

// IsEmpty checks if a string is empty or consists only of whitespace.
func IsEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

// Contains checks if a slice of strings contains a specific string.
func Contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// SanitizeString sanitizes a string by trimming whitespace and converting to lowercase.
func SanitizeString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// ValidateFileName checks if a filename is valid (not empty and not too long).
func ValidateFileName(filename string) error {
	if IsEmpty(filename) {
		return errors.New("filename cannot be empty")
	}
	if len(filename) > 255 {
		return errors.New("filename is too long")
	}
	return nil
}

// ValidateContentType checks if a content type is valid.
func ValidateContentType(contentType string) error {
	if IsEmpty(contentType) {
		return errors.New("content type cannot be empty")
	}
	// Add more content type validation logic if needed
	return nil
}
