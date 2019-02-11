package analysis

import (
	"github.com/z7zmey/php-parser/errors"
	"github.com/z7zmey/php-parser/php7"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker/resolver"
	"io"
	"os"
)

type FileAnalyzer struct {
}

func (a *FileAnalyzer) Analyze(src io.Reader, path string, r resolver.Resolver) (*Analysis, chan<- *ParserError, error) {
	parser := php7.NewParser(src, path)
	parser.Parse()

	parserErrors := parser.GetErrors()
	parserErrorsChan := make(chan *ParserError)

	if len(parserErrors) > 0 {
		go broadcastParserErrors(parserErrorsChan, path, parserErrors)
	} else {
		rootNode := parser.GetRootNode()

		// Resolve imports
		rootNode.Walk(r)
		r.Clean()
	}

	return &Analysis{r.GetImports(), r.GetExports()}, parserErrorsChan, nil
}

func (a *FileAnalyzer) AnalyzeFile(path string, r resolver.Resolver) (*Analysis, chan<- *ParserError, error) {
	src, err := os.Open(path)

	if err != nil {
		return nil, nil, err
	}

	defer func() {
		if err := src.Close(); err != nil {
			panic(err)
		}
	}()

	return a.Analyze(src, path, r)
}

func broadcastParserErrors(errsChan chan<- *ParserError, path string, errs []*errors.Error) {
	for _, e := range errs {
		errsChan <- NewParserError(path, e)
	}
	close(errsChan)
}
