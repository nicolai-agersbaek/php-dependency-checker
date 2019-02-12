package dependency_checker

import (
	"fmt"
	"github.com/gosuri/uiprogress"
	"github.com/spf13/cobra"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/cmd"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker/analysis"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker/checker"
	. "gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker/names"
	"time"
)

type importsCmdOptions struct {
	parallel       bool
	excludeBuiltIn bool
}

var importsCmdOpts = &importsCmdOptions{}

func init() {
	importsCmd.Flags().BoolVarP(&importsCmdOpts.excludeBuiltIn, "exclude-native", "x", true, "Exclude built-in PHP names from results.")
	importsCmd.Flags().BoolVarP(&importsCmdOpts.parallel, "parallel", "p", true, "Perform parallel name resolution.")

	rootCmd.AddCommand(importsCmd)
}

var importsCmd = &cobra.Command{
	Use:   "imports (<dir>|<file>) [(<dir>|<file>), ...]",
	Short: "Display imports for directories and/or files.",
	Args:  cobra.MinimumNArgs(1),
	Run:   imports,
}

func imports(c *cobra.Command, args []string) {
	checkInput.Sources = args

	p := getVerbosePrinter(c)
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

			bar = progressBar(numFiles)

			go func(bar *uiprogress.Bar, progress <-chan int) {
				for range progress {
					bar.Incr()
				}
			}(bar, progressChan)
		}
	}

	go printParserErrors(p, parserErrs)

	// Calculate unexported uses.
	start := time.Now()
	R, S, err := ch.Run(checkInput, importsCmdOpts.parallel, nFiles)
	elapsed := time.Now().Sub(start)

	uiprogress.Stop()

	cmd.CheckError(err)

	avgDuration := elapsed / time.Duration(S.FilesAnalyzed)
	p.VLine(fmt.Sprintf("Elapsed: %.2fs (avg. %s)", elapsed.Seconds(), cmd.FormatDuration(avgDuration)), cmd.VerbosityNormal)

	imports := R.Imports

	if importsCmdOpts.excludeBuiltIn {
		imports = R.Imports.Diff(GetBuiltInNames())
	}

	printUses(p, "Imports", imports)
}

func printUses(p cmd.Printer, title string, names *Names) {
	printUsesOf(p, title, "functions", names.Functions)
	printUsesOf(p, title, "classes", names.Classes)
	printUsesOf(p, title, "constants", names.Constants)
	printUsesOf(p, title, "namespaces", names.Namespaces)
}

func printUsesOf(p cmd.Printer, title, nameType string, names []string) {
	t := fmt.Sprintf("%s [%s]:", title, nameType)
	p.LinesWithTitle(t, names)
}
