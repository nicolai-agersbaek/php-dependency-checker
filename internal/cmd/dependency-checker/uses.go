package dependency_checker

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/cmd"
)

func init() {
	rootCmd.AddCommand(usesCmd)
}

var usesCmd = &cobra.Command{
	Use:   "uses (<dir>|<file>) [(<dir>|<file>)] [,...]",
	Short: "Resolve class uses for a file or files in a directory.",
	Args:  cobra.MinimumNArgs(1),
	Run:   imports,
}

func imports(c *cobra.Command, args []string) {
	imports, exports, err := ResolveImportsSerial(args...)
	cmd.CheckError(err)

	p := cmd.NewPrinter(c)

	// Print uses
	printUses(p, "Imports", imports)
	printUses(p, "Exports", exports)
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
