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

	rootCmd.AddCommand(generateCmd)
}

var generateCmd = &cobra.Command{
	Use:   "check <dir>",
	Short: "Check dependencies for Composer project in dir.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(c *cobra.Command, args []string) {
		// Set the project dir from args

		// Validate the configuration
		//cmd.CheckError(generate.ValidateConfig(conf))

		// Create a generator that will perform code generation
		checker := NewChecker(conf)

		// Run the generation using the constructed configuration
		fmt.Println("Checking files...")
		cmd.CheckError(checker.Run(args[0]))
	},
}
