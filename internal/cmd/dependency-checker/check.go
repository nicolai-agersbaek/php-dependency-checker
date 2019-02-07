package dependency_checker

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/cmd"
	. "gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker/checker"
	. "gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker/names"
	"os"
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

func check(c *cobra.Command, args []string) {
	checkInput.Sources = args

	// Resolve files to analyze
	importPaths := checkInput.ImportPaths()
	exportPaths := checkInput.ExportPaths()

	p := getVerbosePrinter(c)

	start := time.Now()

	// Calculate unexported uses.
	diff := getCheckFunc(parallelMode)(p, importPaths, exportPaths)

	elapsed := time.Now().Sub(start)

	avgDuration := elapsed / time.Duration(len(importPaths)+len(exportPaths))
	p.VLine(fmt.Sprintf("Elapsed: %s (avg. %s)", elapsed.String(), formatDuration(avgDuration)), cmd.VerbosityNormal)

	const maxLines = 15

	if len(diff.Classes) > 0 {
		p.VLinesWithTitleMax("Unexported uses (classes):", diff.Classes, maxLines, cmd.VerbosityNone)
		os.Exit(1)
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

type checkFunc func(p cmd.VerbosePrinter, importPaths, exportPaths []string) *Names

func getCheckFunc(inParallel bool) checkFunc {
	if inParallel {
		return checkParallel
	}

	return checkSerial
}

func checkSerial(p cmd.VerbosePrinter, importPaths, exportPaths []string) *Names {
	I, E, err := ResolveNamesSerial(p, importPaths, exportPaths)
	cmd.CheckError(err)

	// Combine all analyses.
	imports := consolidateIntoClasses(convertToNames(I))

	exports := convertToNames(E)
	exports = exports.Merge(GetBuiltInNames())
	exports = consolidateIntoClasses(exports)

	// Calculate unexported uses.
	return Diff(imports, exports)
}

func checkParallel(p cmd.VerbosePrinter, importPaths, exportPaths []string) *Names {
	I, E, err := ResolveNamesParallel(p, importPaths, exportPaths)
	cmd.CheckError(err)

	// Combine all analyses.
	imports := consolidateIntoClasses(convertToNames(I))

	exports := convertToNames(E)
	exports = exports.Merge(GetBuiltInNames())
	exports = consolidateIntoClasses(exports)

	// Calculate unexported uses.
	return Diff(imports, exports)
}

func convertToNames(F FileNames) *Names {
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
