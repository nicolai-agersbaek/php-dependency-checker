package dependency_checker

import (
	"fmt"
	"github.com/z7zmey/php-parser/php7"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/cmd"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/util/slices"
	"os"
)

func ResolveImportsSerial(p cmd.VerbosePrinter, paths ...string) (*Names, *Names, error) {
	var err error

	phpFiles, err := getPhpFilesSerial(paths)

	if err != nil {
		return nil, nil, err
	}

	var imports, exports *Names

	imports, exports, err = resolveImportsSerial(p, phpFiles...)
	imports.Clean()
	exports.Clean()

	return imports, exports, err
}

func getPhpFilesSerial(paths []string) ([]string, error) {
	var files, allFiles []string
	var err error

	for _, path := range paths {
		files, err = getFilesInDirByExtension("php", path)

		if err != nil {
			return nil, err
		}

		allFiles = append(allFiles, files...)
	}

	return slices.UniqueString(allFiles), nil
}

func resolveImportsSerial(p cmd.VerbosePrinter, paths ...string) (*Names, *Names, error) {
	I, E := make([]*Names, 0), make([]*Names, 0)

	var imports, exports *Names
	var err error

	p.VLine(fmt.Sprintf("Analyzing %d files...", len(paths)), cmd.VerbosityDebug)

	for _, path := range paths {
		imports, exports, err = resolveFileImportsSerial(path)

		if err != nil {
			return nil, nil, err
		}

		I = append(I, imports)
		E = append(E, exports)
	}

	return Merge(I), Merge(E), nil
}

func resolveFileImportsSerial(path string) (*Names, *Names, error) {
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

	resolver := NewImportsResolver()

	if len(parserErrors) > 0 {
		logParserErrors(path, parser.GetErrors())
	} else {
		rootNode := parser.GetRootNode()

		// Resolve imports
		rootNode.Walk(resolver)
		resolver.clean()
	}

	return resolver.Imports, resolver.Exports, nil
}
