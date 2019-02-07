package dependency_checker

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker/names"
)

func init() {
	rootCmd.AddCommand(listNamesCmd)
}

var listNamesCmd = &cobra.Command{
	Use:   "list-names",
	Short: "List built-in PHP functions, classes and constants.",
	Args:  cobra.NoArgs,
	Run:   listNames,
}

//noinspection GoUnusedParameter
func listNames(c *cobra.Command, args []string) {
	defined := names.GetBuiltInNames()

	p := getVerbosePrinter(c)

	p.Line(fmt.Sprintf("Defined functions: %d\n", len(defined.Functions)))
	p.Line(fmt.Sprintf("Defined classes: %d\n", len(defined.Classes)))
	p.Line(fmt.Sprintf("Defined interfaces: %d\n", len(defined.Interfaces)))
	p.Line(fmt.Sprintf("Defined traits: %d\n", len(defined.Traits)))
	p.Line(fmt.Sprintf("Defined constants: %d\n", len(defined.Constants)))
	p.Line(fmt.Sprintf("Defined namespaces: %d\n", len(defined.Namespaces)))
}
