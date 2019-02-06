package dependency_checker

import (
	"github.com/spf13/cobra"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/cmd"
	. "gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker"
	"time"
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
	start := time.Now()

	imports, exports, err := ResolveImportsSerial(args...)
	cmd.CheckError(err)

	// Calculate unexported uses
	diff := Diff(imports, exports)

	stop := time.Now()
	elapsed := stop.Sub(start)

	p := cmd.NewPrinter(c)

	c.Printf("Elapsed: %s\n", elapsed)

	const maxLines = 15

	p.LinesWithTitleMax("Imports (namespaces):", imports.Namespaces, maxLines)
	p.LinesWithTitleMax("Exports (namespaces):", exports.Namespaces, maxLines)

	//p.linesWithTitleMax("Unexported uses (classes):", diff.Classes, maxLines)
	p.LinesWithTitleMax("Unexported uses (namespaces):", diff.Namespaces, maxLines)
}
