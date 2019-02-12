package dependency_checker

import (
	"github.com/z7zmey/php-parser/php7"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker/analysis"
	. "gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker/names"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker/resolver"
	"os"
)

func ResolveNamesSerialFromFiles(inc func(), errs chan<- *analysis.ParserErrors, I, E, B []string) (NamesByFile, NamesByFile, error) {
	var err error

	N := []NamesByFile{make(NamesByFile), make(NamesByFile)}

	for k, F := range [][]string{I, E, B} {
		for _, f := range F {
			names, err := resolveFileImportsSerialByFile(f, errs)
			inc()

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

func resolveFileImportsSerialByFile(path string, errs chan<- *analysis.ParserErrors) (*ImportsExports, error) {
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
	pErrs := parser.GetErrors()

	if len(pErrs) > 0 {
		errs <- analysis.NewParserErrors(path, pErrs)
		return nil, nil
	}

	r := resolver.NewImportExportResolver()
	rootNode := parser.GetRootNode()

	// Resolve imports
	rootNode.Walk(r)
	r.Clean()

	return NewImportsExports(r.Imports, r.Exports), err
}
