package dependency_checker

import (
	"fmt"
	"github.com/z7zmey/php-parser/php7"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/cmd"
	. "gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker/names"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/util/slices"
	"os"
)

func ResolveNamesSerial(p cmd.VerbosePrinter, importPaths, exportPaths []string) (FileNames, FileNames, error) {
	var err error

	P := make([][]string, 2, 2)

	for k, paths := range [][]string{importPaths, exportPaths} {
		F, err := getPhpFilesSerial(paths)

		if err != nil {
			break
		}

		P[k] = slices.UniqueString(F)
	}

	numFiles := len(slices.UniqueStrings(P...))
	if numFiles > 0 {
		p.VLine(fmt.Sprintf("Analyzing %d files...", numFiles), cmd.VerbosityDetailed)
	}

	N := []FileNames{make(FileNames), make(FileNames)}
	R := make(map[string]*ImportsExports)
	verbosity := p.GetVerbosity()

	for k, F := range P {
		for _, f := range F {
			if _, ok := R[f]; !ok {
				R[f], err = resolveFileImportsSerialByFile(f, verbosity)
				if err != nil {
					break
				}
			}

			N[k][f] = R[f][k]
		}
	}

	return N[0], N[1], err
}

func resolveFileImportsSerialByFile(path string, v cmd.Verbosity) (*ImportsExports, error) {
	src, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer func() {
		if err := src.Close(); err != nil {
			panic(err)
		}
	}()

	parser := php7.NewParser(src, path)
	parser.Parse()

	parserErrors := ParserErrors(parser.GetErrors())

	resolver := NewNameResolver()

	if v >= cmd.VerbosityDebug && len(parserErrors) > 0 {
		logParserErrors(path, parser.GetErrors())
		err = parserErrors
	} else {
		rootNode := parser.GetRootNode()

		// Resolve imports
		rootNode.Walk(resolver)
		resolver.clean()
	}

	return NewImportsExports(resolver.Imports, resolver.Exports), err
}
