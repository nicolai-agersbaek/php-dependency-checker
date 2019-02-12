package dependency_checker

import (
	"errors"
	pErrors "github.com/z7zmey/php-parser/errors"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/cmd"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker/analysis"
	. "gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker/names"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker/resolver"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/util/slices"
	"sync"
)

func ResolveNamesParallelFromFiles(inc func(), p cmd.VerbosePrinter, I, E, B []string) (NamesByFile, NamesByFile, error) {
	const numAnalyzers = 5

	analyzer := analysis.NewFileAnalyzer()

	analysisInput := []struct {
		files []string
		mode  analysisMode
	}{
		{I, analyzeImports},
		{E, analyzeExports},
		{B, analyzeBoth},
	}

	C := newCollector(analyzeBoth)

	for _, input := range analysisInput {
		c := newCollector(input.mode)

		err := resolveBothNames(analyzer, inc, p, input.files, numAnalyzers, c)

		if err != nil {
			return nil, nil, err
		}

		C.Merge(c)
	}

	return C.imports.Data(), C.exports.Data(), nil
}

func resolveBothNames(analyzer analysis.Analyzer, inc func(), p cmd.VerbosePrinter, files []string, numAnalyzers int, c *collector) error {
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
		inc()
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

type analysisMode uint

const (
	analyzeImports analysisMode = 1 << iota
	analyzeExports analysisMode = 1 << iota
	analyzeBoth                 = analyzeImports | analyzeExports
)

type collector struct {
	imports, exports *NamesByFileData
}

func newCollector(mode analysisMode) *collector {
	var imports, exports *NamesByFileData

	if mode&analyzeImports != 0 {
		imports = NewNamesByFileData()
	}

	if mode&analyzeExports != 0 {
		exports = NewNamesByFileData()
	}

	return &collector{imports, exports}
}

func (p *collector) Add(r *FileAnalysis) {
	if p.imports != nil {
		p.imports.Put(r.Path, r.Imports)
	}

	if p.exports != nil {
		p.exports.Put(r.Path, r.Exports)
	}
}

func (c *collector) Merge(C *collector) {
	if c.imports != nil {
		c.imports.Merge(C.imports)
	}

	if c.exports != nil {
		c.exports.Merge(C.exports)
	}
}

func walkFiles(done <-chan bool, files []string) (<-chan string, <-chan error) {
	paths := make(chan string)
	errChan := make(chan error, 1)

	pushFiles := func() error {
		for _, f := range files {
			select {
			case paths <- f:
			case <-done:
				return errors.New("walkFiles canceled")
			}
		}

		return nil
	}

	go func() {
		defer close(paths)

		errChan <- pushFiles()
	}()

	return paths, errChan
}

func digester(analyzer analysis.Analyzer, printer cmd.VerbosePrinter, done <-chan bool, paths <-chan string, results chan<- *FileAnalysis) {
	for p := range paths {
		imports, exports, err := analyzeFile(analyzer, p, printer)

		select {
		case results <- NewFileAnalysisExp(p, imports, exports, err):
		case <-done:
			return
		}
	}
}

func analyzeFile(analyzer analysis.Analyzer, path string, p cmd.VerbosePrinter) (*Names, *Names, error) {
	result, parserErrs, err := analyzer.AnalyzeFile(path, resolver.NewImportExportResolver())

	if err != nil {
		return nil, nil, err
	}

	if len(parserErrs.Errors) > 0 {
		logParserErrorsV(parserErrs.Path, parserErrs.Errors, p)
	}

	return result.Imports, result.Exports, nil
}

func logParserErrorsV(path string, errors []*pErrors.Error, p cmd.VerbosePrinter) {
	v := cmd.VerbosityDebug
	indent := "   "
	p.VLine("", v)
	p.VLine(path+":", v)

	for _, e := range errors {
		p.VLine(indent+e.String(), v)
	}
}

// PartitionFileSets partitions importFiles and exportFiles into disjunct sets I,
// E and B, representing the files to be imported, exported and both, respectively.
func PartitionFileSets(importFiles, exportFiles []string) (I, E, B []string) {
	// FIXME: Missing tests!
	I = slices.DiffString(importFiles, exportFiles)
	E = slices.DiffString(exportFiles, importFiles)
	B = slices.IntersectionString(importFiles, exportFiles)

	return I, E, B
}
