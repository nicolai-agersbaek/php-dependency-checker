package dependency_checker

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/cmd"
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

	maxFilesDesc := `Maximum number of files to display in error summary.
If negative, all files will be shown.`
	checkCmd.Flags().IntVar(&printOpts.maxFiles, "max-files", printOpts.maxFiles, maxFilesDesc)

	maxLinesDesc := `Maximum number of lines per file to display in error
summary. If negative, all lines will be shown.`
	checkCmd.Flags().IntVar(&printOpts.maxLines, "max-lines", printOpts.maxLines, maxLinesDesc)

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

	p := getVerbosePrinter(c)
	ch := checker.NewChecker()

	start := time.Now()

	// Calculate unexported uses.
	R, S, err := ch.Run(checkInput, parallelMode, p)
	cmd.CheckError(err)
	//U, diff := doCheck(getResolver(parallelMode), p, importPaths, exportPaths)

	elapsed := time.Now().Sub(start)

	avgDuration := elapsed / time.Duration(S.FilesAnalyzed)
	p.VLine(fmt.Sprintf("Elapsed: %s (avg. %s)", elapsed.String(), formatDuration(avgDuration)), cmd.VerbosityNormal)

	if S.UniqueClsErrs > 0 {
		printLnf("Found %d unique errors in %d files.", S.UniqueClsErrs, S.FilesWithErrs)

		printByFile(p, *R.DiffByFile, printOpts.maxFiles, printOpts.maxLines)

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
	p.Title(indent + "Error" + indent)

	if maxFiles < 0 {
		maxFiles = len(N)
	}

	var i, m int
	for f, nn := range N {
		if i >= maxFiles {
			break
		}

		classes := ConsolidateIntoClasses(nn).Classes
		numCls := len(classes)

		if numCls > 0 {
			m = maxLines
			if maxLines < 0 {
				m = numCls
			}

			fmt.Printf("%s (%d):\n", f, numCls)

			for _, cls := range slices.SliceString(classes, 0, m) {
				fmt.Println(indent + cls)
			}

			i++
		}
	}
}

func formatDuration(d time.Duration) string {
	// TODO: Move to cmd.printer
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
