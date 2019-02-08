package dependency_checker

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/cmd"
	. "gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker"
	"os"
	"runtime/pprof"
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

type rootOpts struct {
	verbosityOptions
	cpuProfile string
}

var rootOptions = &rootOpts{}

func (o *rootOpts) preRunE() error {
	fmt.Println("CPU Profile: " + rootOptions.cpuProfile)

	if rootOptions.cpuProfile != "" {
		f, err := os.Create(rootOptions.cpuProfile)

		if err != nil {
			return err
		}

		err = pprof.StartCPUProfile(f)

		if err != nil {
			return err
		}
	}

	return nil
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&rootOptions.v, "v", false, "Output additional information.")
	rootCmd.PersistentFlags().BoolVar(&rootOptions.vv, "vv", false, "Output detailed information.")
	rootCmd.PersistentFlags().BoolVar(&rootOptions.vvv, "vvv", false, "Output debug information.")

	rootCmd.PersistentFlags().StringVar(&rootOptions.cpuProfile, "cpu-profile", "", "Write CPU profile to file.")

	rootCmd.PersistentPreRunE = cmd.WrapE(rootOptions.preRunE)
	rootCmd.PersistentPostRun = cmd.Wrap(pprof.StopCPUProfile)
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
	defer func() {
		if r := recover(); r != nil {
			pprof.StopCPUProfile()

			if _, ok := r.(errUnexportedClasses); ok {
				os.Exit(1)
			}

			panic(r)
		}
	}()

	cmd.CheckError(rootCmd.Execute())
}
