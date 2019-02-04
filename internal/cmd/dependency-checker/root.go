package dependency_checker

import (
	"github.com/spf13/cobra"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/cmd"
	. "gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker"
)

// FIXME: Fix incomplete descriptions!
var rootCmd = &cobra.Command{
	Use:   Name,
	Short: "dp",
	Long: `Some
                long
                description`,
	Version: Version,
	Run: func(c *cobra.Command, args []string) {
		// Print the help information if command is invoked without any arguments
		cmd.CheckError(c.Help())
	},
}

// Execute executes the main command for the dependency checker
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		cmd.CheckError(err)
	}
}
