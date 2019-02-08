package dependency_checker

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/cmd"
	. "gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker/checker"
	. "gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker/names"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/util/slices"
	"time"
)

var checkInput = &checker.Input{}

type printOptions struct {
	maxFiles, maxLines int
}

var printOpts = &printOptions{5, 10}

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
	checkCmd.Flags().IntVar(&printOpts.maxFiles, "max-files", printOpts.maxFiles, "Max files to display in error summary.")
	checkCmd.Flags().IntVar(&printOpts.maxLines, "max-lines", printOpts.maxLines, "Max lines per file to display in error summary.")

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
	U, diff := doCheck(getResolver(parallelMode), p, importPaths, exportPaths)

	elapsed := time.Now().Sub(start)

	avgDuration := elapsed / time.Duration(len(importPaths)+len(exportPaths))
	p.VLine(fmt.Sprintf("Elapsed: %s (avg. %s)", elapsed.String(), formatDuration(avgDuration)), cmd.VerbosityNormal)

	numClsErrors := len(diff.Classes)
	if numClsErrors > 0 {
		printLnf("Found %d unique errors in %d files.", numClsErrors, len(U))
		printByFile(p, U, printOpts.maxFiles, printOpts.maxLines)

		//p.VLinesWithTitleMax("Unexported uses (classes):", diff.Classes, printOpts.maxLines, cmd.VerbosityNone)
		panic(newUnexportedClsErr())
	} else {
		p.VLine("No unexported uses found!", cmd.VerbosityNormal)
	}
}

func printLnf(format string, a ...interface{}) {
	fmt.Println(fmt.Sprintf(format, a...))
}

const indent = "    "

func printByFile(p cmd.Printer, N NamesByFile, maxFiles, maxLines int) {
	fmt.Println()
	p.Title(indent + "Errors" + indent)

	var i int
	for f, nn := range N {
		if i >= maxFiles {
			break
		}

		classes := consolidateIntoClasses(nn).Classes
		numCls := len(classes)

		if numCls > 0 {
			fmt.Printf("%s (%d):\n", f, numCls)
			for _, cls := range slices.SliceString(classes, 0, maxLines) {
				fmt.Println(indent + cls)
			}
		}

		i++
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

func doCheck(r resolver, p cmd.VerbosePrinter, importPaths, exportPaths []string) (NamesByFile, *Names) {
	I, E, err := r(p, importPaths, exportPaths)
	cmd.CheckError(err)

	// Combine all analyses.
	imports := consolidateIntoClasses(convertToNames(I))

	exports := convertToNames(E)
	exports = exports.Merge(GetBuiltInNames())
	exports = consolidateIntoClasses(exports)

	// Calculate unexported uses.
	D := Diff(imports, exports)
	U := make(NamesByFile)

	for f, N := range I {
		d := N.Diff(exports)

		if !d.Empty() {
			U[f] = d
		}
	}

	return U, D
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
