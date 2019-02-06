package dependency_checker

import (
	"fmt"
	"github.com/z7zmey/php-parser/php7"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/util/slices"
	"os"
	"sync"
)

type Result struct {
	imports, exports *Names
	err              error
}

func NewResult(imports *Names, exports *Names, err error) *Result {
	return &Result{imports: imports, exports: exports, err: err}
}

func (r *Result) Merge(result *Result) {
	if result.err != nil {
		r.err = result.err
	}

	r.imports.Merge(result.imports)
	r.exports.Merge(result.exports)
}

func ResolveImportsParallel(paths ...string) (*Names, *Names, error) {
	var err error

	F := make(chan string)
	phpFiles, err := getPhpFilesParallel(F, paths)

	fmt.Printf("Found %d files:\n", len(phpFiles))

	if err != nil {
		return nil, nil, err
	}

	var imports, exports *Names

	imports, exports, err = resolveImportsParallel(phpFiles...)
	imports.Clean()
	exports.Clean()

	return imports, exports, err
}

func getPhpFilesParallel(F chan<- string, paths []string) ([]string, error) {
	var files, allFiles []string
	var err error

	for _, path := range paths {
		files, err = getFilesInDirByExtension("php", path)

		if err != nil {
			return nil, err
		}

		allFiles = append(allFiles, files...)
	}

	allFiles = slices.UniqueString(allFiles)

	go func() {
		defer close(F)

		for _, f := range allFiles {
			F <- f
		}
	}()

	return nil, nil
}

func resolveImportsParallel(paths ...string) (*Names, *Names, error) {
	I, E := make([]*Names, 0), make([]*Names, 0)

	const maxResolvers = 10
	results := make(chan *Result, maxResolvers)

	var wg sync.WaitGroup

	fmt.Printf("Analyzing %d files...\n", len(paths))
	wg.Add(len(paths))

	for _, path := range paths {
		go doAnalyse(path, results, &wg)
	}

	go func() {
		for r := range results {
			fmt.Print(".")

			I = append(I, r.imports)
			E = append(E, r.exports)
		}
	}()

	wg.Wait()
	close(results)

	fmt.Println()

	return Merge(I), Merge(E), nil
}

func doAnalyse(path string, results chan<- *Result, wg *sync.WaitGroup) {
	err := analyse(path, results)

	if err != nil {
		fmt.Println(err)
	}

	wg.Done()
}

func analyse(path string, results chan<- *Result) error {
	src, err := os.Open(path)

	if err != nil {
		return err
	}

	parser := php7.NewParser(src, path)
	parser.Parse()

	if err := src.Close(); err != nil {
		panic(err)
	}

	// TODO: Return imports, exports and parserErr as a combined Result
	parserErrors := parser.GetErrors()

	resolver := NewImportsResolver()

	if len(parserErrors) > 0 {
		logParserErrors(path, parser.GetErrors())
	} else {
		rootNode := parser.GetRootNode()

		// Resolve imports
		rootNode.Walk(resolver)
		resolver.clean()
	}

	results <- NewResult(resolver.Imports, resolver.Exports, nil)

	return nil
}
