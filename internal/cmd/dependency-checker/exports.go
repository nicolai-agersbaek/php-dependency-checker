package dependency_checker

import (
	"github.com/spf13/cobra"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker/names"
)

type exportsCmdOptions struct {
}

var exportsCmdOpts = &exportsCmdOptions{}

func init() {
	addAnalysisOptions(exportsCmd)
	exportsCmd.Flags().BoolVar(&printOpts.disableProgressBar, "no-progress", printOpts.disableProgressBar, "Disable progress-bar in output.")

	rootCmd.AddCommand(exportsCmd)
}

var exportsCmd = &cobra.Command{
	Use:   "exports (<dir>|<file>) [(<dir>|<file>), ...]",
	Short: "Display exports for directories and/or files.",
	Args:  cobra.MinimumNArgs(1),
	Run:   exports,
}

func exports(c *cobra.Command, args []string) {
	p := getVerbosePrinter(c)
	R, _ := runCheck(args, p, analysisCmdOpts, printOpts.disableProgressBar)

	printUses(p, "Exports", R.Exports.Diff(names.GetBuiltInNames()))
}
