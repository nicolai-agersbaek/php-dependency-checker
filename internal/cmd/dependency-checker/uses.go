package dependency_checker

import (
	"github.com/spf13/cobra"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/cmd"
	. "gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker"
	"strings"
)

const indent = "  "

func init() {
	//generateCmd.Flags().BoolVarP(&conf.DryRun, "dry-run", "d", false, "Simulate a run of the generation")

	//generateCmd.Flags().StringVar(&conf.GoOut, "go_out", "", "Output dir for Go files")
	//cmd.CheckError(generateCmd.MarkFlagRequired("go_out"))

	rootCmd.AddCommand(usesCmd)
}

var usesCmd = &cobra.Command{
	Use:   "uses (<dir>|<file>) [(<dir>|<file>)] [,...]",
	Short: "Resolve class uses for a file or files in a directory.",
	Args:  cobra.MinimumNArgs(1),
	//Run: uses,
	Run: imports,
}

func imports(c *cobra.Command, args []string) {
	imports, exports, err := ResolveImports(args[0])
	cmd.CheckError(err)

	p := printer{c}

	// Print uses
	p.linesWithTitle("Imports (functions):", imports.Functions)
	p.linesWithTitle("Imports (classes):", imports.Classes)
	p.linesWithTitle("Exports (functions):", exports.Functions)
	p.linesWithTitle("Exports (classes):", exports.Classes)
}

type printer struct {
	c *cobra.Command
}

func (p *printer) linesWithTitle(title string, lines []string) {
	if len(lines) > 0 {
		p.title(title)
		p.lines(lines)
	}
}

func (p *printer) title(title string) {
	titleBreak := strings.Repeat("-", len(title))

	p.c.Println(titleBreak)
	p.c.Println(title)
	p.c.Println(titleBreak)
}

func (p *printer) lines(lines []string) {
	for _, line := range lines {
		p.c.Println(line)
	}
}
