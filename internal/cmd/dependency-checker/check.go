package dependency_checker

import (
	"github.com/spf13/cobra"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/cmd"
	. "gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker"
	"os"
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
	Run:   check,
}

func check(c *cobra.Command, args []string) {
	path := args[0]

	imports, exports, err := ResolveImports(path)

	cmd.CheckError(err)

	p := printer{c}

	// Print uses
	p.linesWithTitle("Imports (functions):", imports.Functions)
	p.linesWithTitle("Imports (classes):", imports.Classes)
	p.linesWithTitle("Exports (functions):", exports.Functions)
	p.linesWithTitle("Exports (classes):", exports.Classes)
}

func isDir(path string) bool {
	info, err := os.Stat(path)

	if err != nil {
		return false
	}

	return info.IsDir()
}
