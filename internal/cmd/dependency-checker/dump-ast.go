package dependency_checker

import (
	"github.com/spf13/cobra"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/cmd"
	. "gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker"
)

func init() {
	rootCmd.AddCommand(dumpAstCmd)
}

var dumpAstCmd = &cobra.Command{
	Use:   "dump-ast <file>",
	Short: "Dump the AST of the given PHP file to stdout.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(c *cobra.Command, args []string) {
		cmd.CheckError(DumpAst(args[0], c.OutOrStdout()))
	},
}
