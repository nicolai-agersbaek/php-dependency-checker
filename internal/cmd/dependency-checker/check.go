package dependency_checker

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/cmd"
	. "gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker/checker"
	. "gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker/names"
	"time"
)

var checkInput = &checker.Input{}

var parallelMode = false

func init() {
	excludeDesc := `Directory or file to exclude from analysis. May be
specified multiple times.`
	checkCmd.Flags().StringArrayVarP(&checkInput.Excludes, "exclude", "x", nil, excludeDesc)

	additionalExportsDesc := `Directory or file in which to search for additional
exports. May be specified multiple times.`
	checkCmd.Flags().StringArrayVarP(&checkInput.AdditionalExports, "exports", "e", nil, additionalExportsDesc)

	excludedExportsDesc := `Directory or file to exclude exports from. May be
specified multiple times.`
	checkCmd.Flags().StringArrayVarP(&checkInput.ExcludedExports, "exclude-exports", "E", nil, excludedExportsDesc)

	additionalImportsDesc := `Directory or file in which to search for additional
imports. May be specified multiple times.`
	checkCmd.Flags().StringArrayVarP(&checkInput.AdditionalImports, "imports", "i", nil, additionalImportsDesc)

	excludedImportsDesc := `Directory or file to exclude imports from. May be
specified multiple times.`
	checkCmd.Flags().StringArrayVarP(&checkInput.ExcludedImports, "exclude-imports", "I", nil, excludedImportsDesc)

	checkCmd.Flags().BoolVarP(&parallelMode, "parallel", "p", false, "Perform parallel name resolution.")

	rootCmd.AddCommand(checkCmd)
}

var checkCmd = &cobra.Command{
	Use:   "check (<dir>|<file>) [(<dir>|<file>)] [,...]",
	Short: "Check dependencies for directories and/or files.",
	Args:  cobra.MinimumNArgs(1),
	Run:   check,
}

type errUnexportedClasses error

func newUnexportedClsErr() error {
	return errUnexportedClasses(errors.New("unexported classes"))
}

func check(c *cobra.Command, args []string) {
	checkInput.Sources = args

	// Resolve files to analyze
	importPaths := checkInput.ImportPaths()
	exportPaths := checkInput.ExportPaths()

	p := getVerbosePrinter(c)

	start := time.Now()

	// Calculate unexported uses.
	diff := doCheck(getResolver(parallelMode), p, importPaths, exportPaths)

	elapsed := time.Now().Sub(start)

	avgDuration := elapsed / time.Duration(len(importPaths)+len(exportPaths))
	p.VLine(fmt.Sprintf("Elapsed: %s (avg. %s)", elapsed.String(), formatDuration(avgDuration)), cmd.VerbosityNormal)

	const maxLines = 15

	if len(diff.Classes) > 0 {
		p.VLinesWithTitleMax("Unexported uses (classes):", diff.Classes, maxLines, cmd.VerbosityNone)
		panic(newUnexportedClsErr())
	} else {
		p.VLine("No unexported uses found!", cmd.VerbosityNormal)
	}
}

func formatDuration(d time.Duration) string {
	suffix := "ns"

	if d > time.Second {
		d = d / time.Second
		suffix = "s"
	} else if d > 10*time.Millisecond {
		d = d / time.Millisecond
		suffix = "ms"
	} else if d > 10*time.Microsecond {
		d = d / time.Microsecond
		suffix = "μs"
	}

	return fmt.Sprintf("%d"+suffix, d)
}

type resolver func(p cmd.VerbosePrinter, importPaths, exportPaths []string) (NamesByFile, NamesByFile, error)

func getResolver(inParallel bool) resolver {
	if inParallel {
		return ResolveNamesParallel
	}

	return ResolveNamesSerial
}

func doCheck(r resolver, p cmd.VerbosePrinter, importPaths, exportPaths []string) *Names {
	I, E, err := r(p, importPaths, exportPaths)
	cmd.CheckError(err)

	// Combine all analyses.
	imports := consolidateIntoClasses(convertToNames(I))

	exports := convertToNames(E)
	exports = exports.Merge(GetBuiltInNames())
	exports = consolidateIntoClasses(exports)

	// Calculate unexported uses.
	return Diff(imports, exports)
}

func convertToNames(F NamesByFile) *Names {
	N := NewNames()

	for _, nn := range F {
		N = N.Merge(nn)
	}

	N.Clean()

	return N
}

func consolidateIntoClasses(n *Names) *Names {
	n.Classes = append(n.Classes, n.Interfaces...)
	n.Classes = append(n.Classes, n.Traits...)

	return n
}
