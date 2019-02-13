package dependency_checker

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/cmd"
	. "gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker/names"
)

type importsCmdOptions struct {
	excludeBuiltIn bool
}

var importsCmdOpts = &importsCmdOptions{}

func init() {
	importsCmd.Flags().BoolVarP(&importsCmdOpts.excludeBuiltIn, "exclude-native", "x", true, "Exclude built-in PHP names from results.")
	importsCmd.Flags().BoolVar(&printOpts.disableProgressBar, "no-progress", printOpts.disableProgressBar, "Disable progress-bar in output.")

	addAnalysisOptions(importsCmd)

	rootCmd.AddCommand(importsCmd)
}

var importsCmd = &cobra.Command{
	Use:   "imports (<dir>|<file>) [(<dir>|<file>), ...]",
	Short: "Display imports for directories and/or files.",
	Args:  cobra.MinimumNArgs(1),
	Run:   imports,
}

func imports(c *cobra.Command, args []string) {
	p := getVerbosePrinter(c)
	R, _ := runCheck(args, p, analysisCmdOpts, printOpts.disableProgressBar)

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
