package analysis

import "github.com/z7zmey/php-parser/errors"

type ParserErrors struct {
	Path   string
	Errors []*errors.Error
}

func NewParserErrors(path string, errors []*errors.Error) *ParserErrors {
	return &ParserErrors{Path: path, Errors: errors}
}
