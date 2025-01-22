package parser

import (
	"errors"
)

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

// Example of initializing parsers
func init() {
	// Register your parsers here
	RegisterParser(&PDFParser{})
	RegisterParser(&DOCXParser{})
	RegisterParser(&PPTXParser{})
}
