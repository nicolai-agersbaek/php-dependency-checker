package checker

import (
	"fmt"
	"github.com/gosuri/uiprogress"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/cmd"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker/names"
)

type Result struct {
	Imports, Exports, Diff                   *names.Names
	ImportsByFile, ExportsByFile, DiffByFile *names.NamesByFile
}

type ResultStats struct {
	UniqueClsErrs, FilesAnalyzed, FilesWithErrs int
}

type Checker struct {
}

func NewChecker() *Checker {
	return &Checker{}
}

func (c *Checker) Run(input *Input, parallel bool, p cmd.VerbosePrinter) (*Result, *ResultStats, error) {
	// Get appropriate resolver
	resolver := getResolver(parallel)

	// Resolve files to analyze
	importsFrom := input.ImportsFromFiles()
	exportsFrom := input.ExportsFromFiles()
	I, E, B := dependency_checker.PartitionFileSets(importsFrom, exportsFrom)

	numFiles := len(I) + len(E) + len(B)

	if numFiles <= 0 {
		return nil, nil, nil
	}

	p.VLine(fmt.Sprintf("Analyzing %d files...", numFiles), cmd.VerbosityDetailed)

	// Add progress bar
	bar := uiprogress.AddBar(numFiles)

	completedCount := func(b *uiprogress.Bar) string {
		return fmt.Sprintf("%d/%d", b.Current(), b.Total)
	}
	bar.PrependCompleted()
	bar.PrependFunc(completedCount)
	bar.AppendElapsed()

	inc := func() {
		bar.Incr()
	}

	r, err := runAnalysis(resolver, inc, p, I, E, B)

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

func runAnalysis(r resolver, inc func(), p cmd.VerbosePrinter, importsFrom, exportsFrom, bothFrom []string) (*Result, error) {
	I, E, err := r(inc, p, importsFrom, exportsFrom, bothFrom)

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
