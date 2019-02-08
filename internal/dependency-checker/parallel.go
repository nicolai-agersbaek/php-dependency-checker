package dependency_checker

import (
	"fmt"
	"github.com/z7zmey/php-parser/errors"
	"github.com/z7zmey/php-parser/php7"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/cmd"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker/files"
	. "gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker/names"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/util/slices"
	"os"
	"sync"
)

func ResolveNamesParallel(p cmd.VerbosePrinter, importPaths, exportPaths []string) (FileNames, FileNames, error) {
	var err error

	importFiles, exportFiles, err := resolveFiles(importPaths, exportPaths)

	if err != nil {
		return nil, nil, err
	}

	I, E, B := partitionFileSets(importFiles, exportFiles)

	numFiles := len(I) + len(E) + len(B)
	if numFiles > 0 {
		p.VLine(fmt.Sprintf("Analyzing %d files...", numFiles), cmd.VerbosityDetailed)
	}

	const maxConcurrentFiles = 20

	importFilesChan := make(chan string, maxConcurrentFiles)
	exportFilesChan := make(chan string, maxConcurrentFiles)
	bothFilesChan := make(chan string, maxConcurrentFiles)
	importResultsChan := make(chan *FileAnalysis)
	exportResultsChan := make(chan *FileAnalysis)
	bothResultsChan := make(chan *FileAnalysis)
	errChan := make(chan error)

	var wg sync.WaitGroup
	wg.Add(numFiles)

	// Start broadcasting paths
	go broadcast(I, importFilesChan)
	go broadcast(E, exportFilesChan)
	go broadcast(B, bothFilesChan)

	// Start analyses
	go startAnalyses(wg, importFilesChan, importResultsChan, errChan, p)
	go startAnalyses(wg, exportFilesChan, exportResultsChan, errChan, p)
	go startAnalyses(wg, bothFilesChan, bothResultsChan, errChan, p)

	done := make(chan bool)

	go func() {
		wg.Wait()
		done <- true
	}()

	// Collect analyses
	imports := make(FileNames)
	exports := make(FileNames)

	select {
	case <-done:
	case err = <-errChan:
	case r := <-importResultsChan:
		r.Imports.Clean()
		imports[r.Path] = r.Imports
	case r := <-exportResultsChan:
		r.Exports.Clean()
		exports[r.Path] = r.Exports
	case r := <-bothResultsChan:
		r.Imports.Clean()
		imports[r.Path] = r.Imports

		r.Exports.Clean()
		exports[r.Path] = r.Exports
	}

	if err != nil {
		return nil, nil, err
	}

	return imports, exports, err
}

func startAnalyses(wg sync.WaitGroup, F <-chan string, results chan<- *FileAnalysis, errs chan<- error, p cmd.VerbosePrinter) {
	for f := range F {
		analyze(f, results, errs, p)
		wg.Done()
	}
}

func broadcast(F []string, c chan<- string) {
	for _, f := range F {
		c <- f
	}
}

func analyze(file string, results chan<- *FileAnalysis, errs chan<- error, p cmd.VerbosePrinter) {
	I, E, err := analyzeFile(file, p)

	if err != nil {
		//errs <- err
		panic(err)
	}

	results <- NewFileAnalysis(file, I, E)
}

func analyzeFile(path string, p cmd.VerbosePrinter) (*Names, *Names, error) {
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

	// TODO: Return imports, exports and parserErr as a combined Result
	parserErrors := parser.GetErrors()

	resolver := NewNameResolver()

	if len(parserErrors) > 0 {
		logParserErrorsV(path, parser.GetErrors(), p)
	} else {
		rootNode := parser.GetRootNode()

		// Resolve imports
		rootNode.Walk(resolver)
		resolver.clean()
	}

	return resolver.Imports, resolver.Exports, nil
}

func logParserErrorsV(path string, errors []*errors.Error, p cmd.VerbosePrinter) {
	indent := "   "
	p.Line("")
	p.Line(path + ":")

	for _, e := range errors {
		p.Line(indent + e.String())
	}
}

func resolveFiles(importPaths, exportPaths []string) ([]string, []string, error) {
	// Resolve the files given by import and export paths
	P := make([][]string, 2, 2)

	for k, paths := range [][]string{importPaths, exportPaths} {
		F, err := getPhpFilesSerial(paths)

		if err != nil {
			return nil, nil, err
		}

		P[k] = slices.UniqueString(F)
	}

	return P[0], P[1], nil
}

// partitionFileSets partitions importFiles and exportFiles into disjunct sets I,
// E and B, representing the files to be imported, exported and both, respectively.
func partitionFileSets(importFiles, exportFiles []string) (I, E, B []string) {
	// FIXME: Missing tests!
	I = slices.DiffString(importFiles, exportFiles)
	E = slices.DiffString(exportFiles, importFiles)
	B = slices.IntersectionString(importFiles, exportFiles)

	return I, E, B
}

func getPhpFilesParallel(F chan<- string, paths []string) ([]string, error) {
	var fs, Fs []string
	var err error

	for _, path := range paths {
		fs, err = files.GetFilesInDirByExtension("php", path)

		if err != nil {
			return nil, err
		}

		Fs = append(Fs, fs...)
	}

	Fs = slices.UniqueString(Fs)

	go func() {
		defer close(F)

		for _, f := range Fs {
			F <- f
		}
	}()

	return nil, nil
}
