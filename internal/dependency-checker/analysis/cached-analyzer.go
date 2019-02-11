package analysis

import (
	"crypto/md5"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker/resolver"
	"io"
	"io/ioutil"
	"os"
)

type CachedAnalyzer struct {
	cacheDir string
	analyzer Analyzer
}

func NewCachedAnalyzer(cacheDir string, analyzer Analyzer) *CachedAnalyzer {
	return &CachedAnalyzer{cacheDir: cacheDir, analyzer: analyzer}
}

func (a *CachedAnalyzer) Analyze(src io.Reader, path string, r resolver.Resolver) (*Analysis, chan<- *ParserError, error) {
	h, err := hashContent(src)

	if err != nil {
		return nil, nil, err
	}

	if a.cacheHas(h) {
		return a.cacheGet(h), nil, nil
	}

	analysis, ch, err := a.analyzer.Analyze(src, path, r)

	if err != nil {
		return nil, ch, err
	}

	err = a.cacheSave(h, analysis)

	return analysis, ch, err
}

func (a *CachedAnalyzer) AnalyzeFile(path string, r resolver.Resolver) (*Analysis, chan<- *ParserError, error) {
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

func hashContent(reader io.Reader) (string, error) {
	bytes, err := ioutil.ReadAll(reader)

	if err != nil {
		return "", err
	}

	return string(md5.New().Sum(bytes)), nil
}

func (a *CachedAnalyzer) cacheHas(fileHash string) bool {
	return false
}

func (a *CachedAnalyzer) cacheGet(fileHash string) *Analysis {
	return nil
}

func (a *CachedAnalyzer) cacheSave(fileHash string, analysis *Analysis) error {
	// Save parse to file named fileHash

	return nil
}
