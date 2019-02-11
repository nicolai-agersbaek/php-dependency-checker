package dependency_checker

import (
	"github.com/z7zmey/php-parser/php7"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/cmd"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker/analysis"
	. "gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker/names"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker/resolver"
	"os"
	"sync"
)

func ResolveImportsParallel(analyzer analysis.Analyzer, p cmd.VerbosePrinter, I []string) (NamesByFile, NamesByFile, error) {
	const numAnalyzers = 5
	const mode = analyzeImports

	C := newCollector(mode)

	err := resolveNamesParallel(analyzer, p, I, numAnalyzers, C, mode)

	if err != nil {
		return nil, nil, err
	}

	return C.imports.Data(), C.exports.Data(), nil
}

func resolveNamesParallel(analyzer analysis.Analyzer, p cmd.VerbosePrinter, files []string, numAnalyzers int, c *collector, mode analysisMode) error {
	done := make(chan bool)
	defer close(done)

	fileChan, errChan := walkFiles(done, files)
	resultChan := make(chan *FileAnalysis)

	var wg sync.WaitGroup
	wg.Add(numAnalyzers)

	for i := 0; i < numAnalyzers; i++ {
		go func() {
			digester(analyzer, p, done, fileChan, resultChan)
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect analyses
	for r := range resultChan {
		if r.Error != nil {
			return r.Error
		}

		c.Add(r)
	}

	if err := <-errChan; err != nil {
		return err
	}

	return nil
}

func digesterFlexible(analyze analyzeFileFunc, done <-chan bool, paths <-chan string, results chan<- *FileAnalysis) {
	for p := range paths {
		imports, exports, err := analyze(p)

		select {
		case results <- NewFileAnalysisExp(p, imports, exports, err):
		case <-done:
			return
		}
	}
}

type analyzeFileFunc func(path string) (*Names, *Names, error)

func analyzeFileFlexible(path string, p cmd.VerbosePrinter, r resolver.Resolver) (*Names, *Names, error) {
	src, err := os.Open(path)

	if err != nil {
		return nil, nil, err
	}

	defer func() {
		if err := src.Close(); err != nil {
			panic(err)
		}
	}()

	parser := php7.NewParser(src, path)
	parser.Parse()

	parserErrors := parser.GetErrors()

	if len(parserErrors) > 0 {
		logParserErrorsV(path, parser.GetErrors(), p)
	} else {
		rootNode := parser.GetRootNode()

		// Resolve imports
		rootNode.Walk(r)
	}

	return r.GetImports(), r.GetExports(), nil
}
