package dependency_checker

import (
	"github.com/z7zmey/php-parser/php7"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/cmd"
	. "gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker/names"
	"os"
)

func ResolveNamesSerialFromFiles(p cmd.VerbosePrinter, I, E, B []string) (NamesByFile, NamesByFile, error) {
	var err error

	N := []NamesByFile{make(NamesByFile), make(NamesByFile)}
	verbosity := p.GetVerbosity()

	for k, F := range [][]string{I, E, B} {
		for _, f := range F {
			names, err := resolveFileImportsSerialByFile(f, verbosity)

			if err != nil {
				break
			}

			switch k {
			case 0: // I
				N[0][f] = names[0]
			case 1: // E
				N[1][f] = names[1]
			case 2: // B
				N[0][f] = names[0]
				N[1][f] = names[1]
			}
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
