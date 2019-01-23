package dependency_checker

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/cmd"
	. "gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker"
)

const indent = "  "

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
		// Run the analysis
		fmt.Println("Analysing files...")

		funcUses, err := ResolveUses(args, IsFunctionName)
		cmd.CheckError(err)

		clsUses, err := ResolveUses(args, IsClassName)
		cmd.CheckError(err)

		// Print uses
		c.Println("----------  FUNCTIONS  ----------")
		printUses(c, funcUses)
		c.Println("----------  CLASSES  ----------")
		printUses(c, clsUses)
	},
}

func printUses(c *cobra.Command, usesMap ClassUsesMap) {
	for file, uses := range usesMap {
		c.Printf("%s:\n", file)

		for _, use := range uses {
			c.Printf("%s%s\n", indent, use)
		}
	}
}
