package dependency_checker

import (
	"github.com/spf13/cobra"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/cmd"
	. "gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker"
)

type verbosityOptions struct {
	v, vv, vvv bool
}

func (o verbosityOptions) GetVerbosity() cmd.Verbosity {
	if o.vvv {
		return cmd.VerbosityDebug
	}

	if o.vv {
		return cmd.VerbosityDetailed
	}

	if o.v {
		return cmd.VerbosityNormal
	}

	return cmd.VerbosityNone
}

var rootOptions = &verbosityOptions{}

func init() {
	rootCmd.PersistentFlags().BoolVar(&rootOptions.v, "v", false, "Output additional information.")
	rootCmd.PersistentFlags().BoolVar(&rootOptions.vv, "vv", false, "Output detailed information.")
	rootCmd.PersistentFlags().BoolVar(&rootOptions.vvv, "vvv", false, "Output debug information.")
}

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

func getVerbosePrinter(c *cobra.Command) cmd.VerbosePrinter {
	return cmd.NewVerbosePrinter(cmd.NewPrinter(c), rootOptions.GetVerbosity())
}

// Execute executes the main command for the dependency checker
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		cmd.CheckError(err)
	}
}
