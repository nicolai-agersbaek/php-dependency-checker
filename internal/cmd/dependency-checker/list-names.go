package dependency_checker

import (
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

	c.Printf("Defined functions: %d\n", len(defined.Functions))
	c.Printf("Defined classes: %d\n", len(defined.Classes))
	c.Printf("Defined interfaces: %d\n", len(defined.Interfaces))
	c.Printf("Defined traits: %d\n", len(defined.Traits))
	c.Printf("Defined constants: %d\n", len(defined.Constants))
	c.Printf("Defined namespaces: %d\n", len(defined.Namespaces))
}
