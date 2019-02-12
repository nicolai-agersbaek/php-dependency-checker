package analysis

import (
	"github.com/z7zmey/php-parser/php7"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker/resolver"
	"io"
	"os"
)

type FileAnalyzer struct {
}

func NewFileAnalyzer() *FileAnalyzer {
	return &FileAnalyzer{}
}

func (a *FileAnalyzer) Analyze(src io.Reader, path string, r resolver.Resolver) (*Analysis, *ParserErrors, error) {
	parser := php7.NewParser(src, path)
	parser.Parse()
	pErrs := parser.GetErrors()

	var parserErrors *ParserErrors

	if len(pErrs) > 0 {
		parserErrors = NewParserErrors(path, pErrs)
	}

	rootNode := parser.GetRootNode()

	// Resolve imports
	rootNode.Walk(r)

	return &Analysis{r.GetImports(), r.GetExports()}, parserErrors, nil
}

func (a *FileAnalyzer) AnalyzeFile(path string, r resolver.Resolver) (*Analysis, *ParserErrors, error) {
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
