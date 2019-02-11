package dependency_checker

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/cmd"
	"os"
	"runtime/pprof"
)

const Name = "php-dependency-checker"

const Version = "0.1.0"

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
	cpuProfile, memProfile string
	memProfileWritten      bool
}

var rootOptions = &rootOpts{}

func (o *rootOpts) preRunE() error {
	if rootOptions.cpuProfile != "" {
		f, err := os.Create(rootOptions.cpuProfile)

		if err != nil {
			return err
		}

		err = pprof.StartCPUProfile(f)

		if err != nil {
			return err
		}

		fmt.Println("Saving CPU profile to: " + rootOptions.cpuProfile)
	}

	return nil
}

func (o *rootOpts) postRunE() error {
	pprof.StopCPUProfile()

	return nil
}

func (o *rootOpts) writeMemProfile() error {
	if o.memProfileWritten {
		return nil
	}

	if o.memProfile != "" {
		f, err := os.Create(o.memProfile)
		if err != nil {
			return err
		}

		err = pprof.WriteHeapProfile(f)
		if err != nil {
			return err
		}

		err = f.Close()
		if err != nil {
			return err
		}

		fmt.Println("Saving memory profile to: " + rootOptions.memProfile)
		o.memProfileWritten = true
	}

	return nil
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&rootOptions.v, "v", false, "Output additional information.")
	rootCmd.PersistentFlags().BoolVar(&rootOptions.vv, "vv", false, "Output detailed information.")
	rootCmd.PersistentFlags().BoolVar(&rootOptions.vvv, "vvv", false, "Output debug information.")

	rootCmd.PersistentFlags().StringVar(&rootOptions.cpuProfile, "cpu-profile", "", "Write CPU profile to file.")
	rootCmd.PersistentFlags().StringVar(&rootOptions.memProfile, "mem-profile", "", "Write memory profile to file.")

	rootCmd.PersistentPreRunE = cmd.WrapE(rootOptions.preRunE)
	rootCmd.PersistentPostRunE = cmd.WrapE(rootOptions.postRunE)
}

//noinspection SpellCheckingInspection
var rootCmd = &cobra.Command{
	Use:   Name,
	Short: "phpdp",
	Long: `
Provides tooling for analyzing PHP imports and exports of functions, classes,
constants and namespaces.
Can be used to list these as well as analyze discrepancies in dependencies of
any PHP project, including projects using Composer for dependency management.`,
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

			err := rootOptions.writeMemProfile()
			if err != nil {
				panic(err)
			}

			if _, ok := r.(errUnexportedClasses); ok {
				os.Exit(1)
			}

			panic(r)
		}
	}()

	cmd.CheckError(rootCmd.Execute())
}
