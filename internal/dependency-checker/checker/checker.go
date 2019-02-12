package checker

import (
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/cmd"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker/names"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/util/slices"
)

type Result struct {
	Imports, Exports, Diff                   *names.Names
	ImportsByFile, ExportsByFile, DiffByFile *names.NamesByFile
}

type ResultStats struct {
	UniqueClsErrs, FilesAnalyzed, FilesWithErrs int
}

type Checker struct {
	printer  cmd.VerbosePrinter
	progress chan<- int
}

func NewChecker(printer cmd.VerbosePrinter, progress chan<- int) *Checker {
	return &Checker{printer, progress}
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

func (c *Checker) Run(input *Input, parallel bool, nFiles func(int)) (*Result, *ResultStats, error) {
	// Get appropriate resolver
	resolver := getResolver(parallel)

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

	r, err := c.runAnalysis(resolver, inc, I, E, B)

	close(c.progress)

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

type resolver func(inc func(), p cmd.VerbosePrinter, I, E, B []string) (names.NamesByFile, names.NamesByFile, error)

func getResolver(inParallel bool) resolver {
	if inParallel {
		return dependency_checker.ResolveNamesParallelFromFiles
	}

	return dependency_checker.ResolveNamesSerialFromFiles
}

func (c *Checker) runAnalysis(r resolver, inc func(), importsFrom, exportsFrom, bothFrom []string) (*Result, error) {
	I, E, err := r(inc, c.printer, importsFrom, exportsFrom, bothFrom)

	if err != nil {
		return nil, err
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
	// Combine all analyses.
	I := names.ConsolidateIntoClasses(names.ConvertToNames(imports))

	E := names.ConvertToNames(exports)
	E = E.Merge(names.GetBuiltInNames())
	E = names.ConsolidateIntoClasses(E)

	return I, E
}
