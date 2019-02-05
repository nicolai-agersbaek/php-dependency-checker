package dependency_checker

import (
	"github.com/spf13/cobra"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/cmd"
	. "gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker"
)

func init() {
	rootCmd.AddCommand(checkCmd)
}

var checkCmd = &cobra.Command{
	Use:   "check (<dir>|<file>) [(<dir>|<file>)] [,...]",
	Short: "Check dependencies for directories and/or files.",
	Args:  cobra.MinimumNArgs(1),
	Run:   check,
}

func check(c *cobra.Command, args []string) {
	imports, exports, err := ResolveAllImports(args...)
	cmd.CheckError(err)

	// Calculate unexported uses
	diff := Diff(imports, exports)

	p := cmd.NewPrinter(c)

	const maxLines = 15

	p.LinesWithTitleMax("Imports (namespaces):", imports.Namespaces, maxLines)
	p.LinesWithTitleMax("Exports (namespaces):", exports.Namespaces, maxLines)

	//p.linesWithTitleMax("Unexported uses (classes):", diff.Classes, maxLines)
	p.LinesWithTitleMax("Unexported uses (namespaces):", diff.Namespaces, maxLines)
}
