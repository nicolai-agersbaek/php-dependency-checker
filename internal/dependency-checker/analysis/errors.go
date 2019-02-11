package analysis

import "github.com/z7zmey/php-parser/errors"

type ParserError struct {
	Path  string
	Error *errors.Error
}

func NewParserError(path string, error *errors.Error) *ParserError {
	return &ParserError{Path: path, Error: error}
}
