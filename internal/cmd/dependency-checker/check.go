package dependency_checker

import (
	"errors"
	"fmt"
	"github.com/gosuri/uiprogress"
	"github.com/spf13/cobra"
	pErrors "github.com/z7zmey/php-parser/errors"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/cmd"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker/analysis"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker/checker"
	. "gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker/names"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/util/slices"
	"time"
)

var checkInput = &checker.Input{}

type printOptions struct {
	maxFiles, maxLines int
	disableProgressBar bool
}

var printOpts = &printOptions{5, 10, false}

type analysisCmdOptions struct {
	serial, ignoreGlobals bool
}

var analysisCmdOpts = &analysisCmdOptions{false, false}

type checkCmdOptions struct {
}

var checkCmdOpts = &checkCmdOptions{}

func addAnalysisOptions(c *cobra.Command) {
	c.Flags().BoolVarP(&analysisCmdOpts.serial, "serial", "P", analysisCmdOpts.serial, "Perform serial name resolution.")
	c.Flags().BoolVarP(&analysisCmdOpts.ignoreGlobals, "ignore-globals", "G", analysisCmdOpts.ignoreGlobals, "Ignore 'global' names during analysis.")
}

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

	checkCmd.Flags().BoolVar(&printOpts.disableProgressBar, "no-progress", printOpts.disableProgressBar, "Disable progress-bar in output.")

	maxFilesDesc := `Maximum number of files to display in error summary.
If negative, all files will be shown.`
	checkCmd.Flags().IntVar(&printOpts.maxFiles, "max-files", printOpts.maxFiles, maxFilesDesc)

	maxLinesDesc := `Maximum number of lines per file to display in error
summary. If negative, all lines will be shown.`
	checkCmd.Flags().IntVar(&printOpts.maxLines, "max-lines", printOpts.maxLines, maxLinesDesc)

	addAnalysisOptions(checkCmd)

	rootCmd.AddCommand(checkCmd)
}

var checkCmd = &cobra.Command{
	Use:   "check (<dir>|<file>) [(<dir>|<file>), ...]",
	Short: "Check dependencies for directories and/or files.",
	Args:  cobra.MinimumNArgs(1),
	Run:   check,
}

type errUnexportedClasses cmd.Recoverable

func newUnexportedClsErr() errUnexportedClasses {
	return cmd.NewRecoverable(errors.New("unexported classes"), 1)
}

func check(c *cobra.Command, args []string) {
	p := getVerbosePrinter(c)
	R, S := runCheck(args, p, analysisCmdOpts, printOpts.disableProgressBar)

	if S.UniqueClsErrs > 0 {
		printLnf("Found %d unique errors in %d files.", S.UniqueClsErrs, S.FilesWithErrs)

		printByFile(p, *R.DiffByFile, printOpts.maxFiles, printOpts.maxLines)

		//p.VLinesWithTitleMax("Unexported uses (classes):", diff.Classes, printOpts.maxLines, cmd.VerbosityNone)
		panic(newUnexportedClsErr())
	} else {
		p.VLine("No unexported uses found!", cmd.VerbosityNormal)
	}
}

func runCheck(args []string, p cmd.VerbosePrinter, opts *analysisCmdOptions, noProgress bool) (*checker.Result, *checker.ResultStats) {
	checkInput.Sources = args

	progressChan := make(chan int)
	parserErrs := make(chan *analysis.ParserErrors)
	started := false
	ch := checker.NewChecker(progressChan, parserErrs)

	uiprogress.Start()

	var bar *uiprogress.Bar

	nFiles := func(numFiles int) {
		if !started && numFiles > 0 {
			started = true

			p.VLine(fmt.Sprintf("Analyzing %d files...", numFiles), cmd.VerbosityDetailed)

			if !noProgress {
				bar = progressBar(numFiles)
			}

			go readProgress(bar, progressChan)
		}
	}

	go printParserErrors(p, parserErrs)

	// Calculate unexported uses.
	start := time.Now()
	R, S, err := ch.Run(checkInput, opts.serial, opts.ignoreGlobals, nFiles)
	elapsed := time.Now().Sub(start)

	uiprogress.Stop()

	cmd.CheckError(err)

	avgDuration := elapsed / time.Duration(S.FilesAnalyzed)
	p.VLine(fmt.Sprintf("Elapsed: %.2fs (avg. %s)", elapsed.Seconds(), cmd.FormatDuration(avgDuration)), cmd.VerbosityNormal)

	return R, S
}

func readProgress(bar *uiprogress.Bar, progress <-chan int) {
	if bar == nil {
		for range progress {
		}
		return
	}

	for range progress {
		bar.Incr()
	}
}

func printParserErrors(printer cmd.VerbosePrinter, errs <-chan *analysis.ParserErrors) {
	for e := range errs {
		if e != nil {
			logParserErrorsV(e.Path, e.Errors, printer)
		}
	}
}

func logParserErrorsV(path string, errors []*pErrors.Error, p cmd.VerbosePrinter) {
	v := cmd.VerbosityDebug
	indent := "   "
	p.VLine("", v)
	p.VLine(path+":", v)

	for _, e := range errors {
		p.VLine(indent+e.String(), v)
	}
}

func progressBar(total int) *uiprogress.Bar {
	// Add progress bar
	bar := uiprogress.AddBar(total)

	completedCount := func(b *uiprogress.Bar) string {
		return fmt.Sprintf("%d/%d", b.Current(), b.Total)
	}
	bar.PrependCompleted()
	bar.PrependFunc(completedCount)
	bar.AppendElapsed()

	return bar
}

func printLnf(format string, a ...interface{}) {
	fmt.Println(fmt.Sprintf(format, a...))
}

const indent = "    "

func printByFile(p cmd.Printer, N NamesByFile, maxFiles, maxLines int) {
	fmt.Println()
	p.Title(indent + "Errors" + indent)

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
