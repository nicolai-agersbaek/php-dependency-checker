package dependency_checker

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/cmd"
	. "gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker"
)

func init() {
	//generateCmd.Flags().BoolVarP(&conf.DryRun, "dry-run", "d", false, "Simulate a run of the generation")

	//generateCmd.Flags().StringVar(&conf.GoOut, "go_out", "", "Output dir for Go files")
	//cmd.CheckError(generateCmd.MarkFlagRequired("go_out"))

	rootCmd.AddCommand(usesCmd)
}

var usesCmd = &cobra.Command{
	Use:   "uses (<dir>|<file>) [(<dir>|<file>)] [,...]",
	Short: "Resolve class uses for a file or files in a directory.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(c *cobra.Command, args []string) {
		// Create a Checker that will perform the analysis
		checker := NewChecker(conf)

		// Run the analysis
		fmt.Println("Analysing files...")

		uses, err := checker.ResolveUses(args...)

		cmd.CheckError(err)

		// Print uses
		printUses(c, uses)
	},
}

func printUses(c *cobra.Command, uses []string) {
	c.Println("----------  USES  ----------")
	//c.Printf(strings.Join(uses, "\n"))

	for _, use := range uses {
		c.Println(use)
	}
}
