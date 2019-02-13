package checker

import (
	"fmt"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker/analysis"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker/names"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/util/slices"
	"os"
	"strings"
)

type Result struct {
	Imports, Exports, Diff                   *names.Names
	ImportsByFile, ExportsByFile, DiffByFile *names.NamesByFile
}

type ResultStats struct {
	UniqueClsErrs, FilesAnalyzed, FilesWithErrs int
}

type Checker struct {
	progress chan<- int
	errs     chan<- *analysis.ParserErrors
}

func NewChecker(progress chan<- int, errs chan<- *analysis.ParserErrors) *Checker {
	return &Checker{progress, errs}
}

func resolveFiles(input *Input) ([]string, []string, []string) {
	// FIXME: Missing tests!
	importFiles := input.ImportsFromFiles()
	exportFiles := input.ExportsFromFiles()

	I := slices.DiffString(importFiles, exportFiles)
	E := slices.DiffString(exportFiles, importFiles)
	B := slices.IntersectionString(importFiles, exportFiles)

	return I, E, B
}

func (c *Checker) Run(input *Input, serial, ignoreGlobals bool, nFiles func(int)) (*Result, *ResultStats, error) {
	fmt.Printf("Checker.Run: ignoreGlobals = %v\n", ignoreGlobals)

	// Get appropriate resolver
	resolver := getResolver(serial)

	// Resolve files to analyze
	I, E, B := resolveFiles(input)

	numFiles := len(I) + len(E) + len(B)
	nFiles(numFiles)

	if numFiles <= 0 {
		return nil, nil, nil
	}

	inc := func() {
		c.progress <- 1
	}

	r, err := c.runAnalysis(resolver, inc, ignoreGlobals, I, E, B)

	close(c.progress)
	close(c.errs)

	if err != nil {
		return nil, nil, err
	}

	stats := &ResultStats{
		UniqueClsErrs: len(r.Diff.Classes),
		FilesAnalyzed: numFiles,
		FilesWithErrs: len(*r.DiffByFile),
	}

	return r, stats, err
}

type resolver func(inc func(), errs chan<- *analysis.ParserErrors, I, E, B []string) (names.NamesByFile, names.NamesByFile, error)

func getResolver(inSerial bool) resolver {
	if inSerial {
		return dependency_checker.ResolveNamesSerialFromFiles
	}

	return dependency_checker.ResolveNamesParallelFromFiles
}

func (c *Checker) runAnalysis(r resolver, inc func(), ignoreGlobals bool, importsFrom, exportsFrom, bothFrom []string) (*Result, error) {
	I, E, err := r(inc, c.errs, importsFrom, exportsFrom, bothFrom)

	if err != nil {
		return nil, err
	}

	// Filter out globals
	if ignoreGlobals {
		I = removeGlobalsFromNamesByFile(I)
		E = removeGlobalsFromNamesByFile(E)
	}

	// Combine all analyses.
	imports, exports := resolveClassImportsExports(I, E)

	// Calculate unexported uses.
	diff := names.Diff(imports, exports)
	U := make(names.NamesByFile)

	for f, N := range I {
		d := N.Diff(exports)

		if !d.Empty() {
			U[f] = d
		}
	}

	return &Result{imports, exports, diff, &I, &E, &U}, nil
}

func resolveClassImportsExports(imports, exports names.NamesByFile) (*names.Names, *names.Names) {
	// FIXME: Missing tests!
	// Combine all analyses.
	I := names.ConsolidateIntoClasses(names.ConvertToNames(imports))

	E := names.ConvertToNames(exports)
	E = E.Merge(names.GetBuiltInNames())
	E = names.ConsolidateIntoClasses(E)

	return I, E
}

func removeGlobalsFromNamesByFile(N names.NamesByFile) names.NamesByFile {
	filtered := make(names.NamesByFile)

	for k, nn := range N {
		nn = removeGlobalsFromNames(nn)

		if !nn.Empty() {
			filtered[k] = nn

			if len(nn.Classes) != len(N[k].Classes) {
				fmt.Println()
				fmt.Println()
				fmt.Println()

				fmt.Printf("nn.Classes = \n%v\n", nn.Classes)
				fmt.Printf("nn.Classes = \n%v\n", N[k].Classes)

				fmt.Println()
				fmt.Println()
				fmt.Println()

				os.Exit(1)
			}
		}
	}

	return filtered
}

func removeGlobalsFromNames(N *names.Names) *names.Names {
	N.Functions = filterOutGlobals(N.Functions)
	N.Classes = filterOutGlobals(N.Classes)
	N.Interfaces = filterOutGlobals(N.Interfaces)
	N.Traits = filterOutGlobals(N.Traits)
	N.Constants = filterOutGlobals(N.Constants)
	N.Namespaces = filterOutGlobals(N.Namespaces)

	return N
}

func filterOutGlobals(S []string) []string {
	filter := func(n string) bool {
		return strings.Contains(n, names.NamespaceSeparator)
	}

	return slices.FilterString(S, filter)
}
