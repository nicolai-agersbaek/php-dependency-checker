package analysis

import (
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker/resolver"
	"io"
)

type Analyzer interface {
	Analyze(src io.Reader, path string, r resolver.Resolver) (*Analysis, chan<- *ParserError, error)
	AnalyzeFile(path string, r resolver.Resolver) (*Analysis, chan<- *ParserError, error)
}
