package dependency_checker

import (
	"github.com/spf13/cobra"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/cmd"
	"os"
	"os/exec"
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

func listNames(c *cobra.Command, args []string) {
	listFunctions := command(c, "php", "-r", "'foreach(get_defined_functions(true)[\"internal\"] as $f){printf($f);}'")

	cmd.CheckError(listFunctions.Run())
}

func command(c *cobra.Command, name string, args ...string) *exec.Cmd {
	command := exec.Command(name, args...)
	command.Stdout = c.OutOrStdout()
	command.Stderr = c.OutOrStderr()
	command.Stdin = os.Stdin

	return command
}
