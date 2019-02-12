package dependency_checker

import (
	"fmt"
	"github.com/gosuri/uiprogress"
	"github.com/spf13/cobra"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/cmd"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker/analysis"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker/checker"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker/names"
	"time"
)

type exportsCmdOptions struct {
	parallel bool
}

var exportsCmdOpts = &exportsCmdOptions{}

func init() {
	exportsCmd.Flags().BoolVarP(&exportsCmdOpts.parallel, "parallel", "p", true, "Perform parallel name resolution.")

	rootCmd.AddCommand(exportsCmd)
}

var exportsCmd = &cobra.Command{
	Use:   "exports (<dir>|<file>) [(<dir>|<file>), ...]",
	Short: "Display exports for directories and/or files.",
	Args:  cobra.MinimumNArgs(1),
	Run:   exports,
}

func exports(c *cobra.Command, args []string) {
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

	printUses(p, "Exports", R.Exports.Diff(names.GetBuiltInNames()))
}
