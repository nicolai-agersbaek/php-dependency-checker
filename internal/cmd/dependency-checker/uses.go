package dependency_checker

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/cmd"
	. "gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker"
	"strings"
)

func init() {
	rootCmd.AddCommand(usesCmd)
}

var usesCmd = &cobra.Command{
	Use:   "uses (<dir>|<file>) [(<dir>|<file>)] [,...]",
	Short: "Resolve class uses for a file or files in a directory.",
	Args:  cobra.MinimumNArgs(1),
	Run:   imports,
}

func imports(c *cobra.Command, args []string) {
	imports, exports, err := ResolveAllImports(args...)
	cmd.CheckError(err)

	p := newPrinter(c)

	// Print uses
	printUses(p, "Imports", imports)
	printUses(p, "Exports", exports)
}

func printUses(p *printer, title string, names *Names) {
	printUsesOf(p, title, "functions", names.Functions)
	printUsesOf(p, title, "classes", names.Classes)
	printUsesOf(p, title, "constants", names.Constants)
	printUsesOf(p, title, "namespaces", names.Namespaces)
}

func printUsesOf(p *printer, title, nameType string, names []string) {
	t := fmt.Sprintf("%s [%s]:", title, nameType)
	p.linesWithTitle(t, names)
}

type printer struct {
	c *cobra.Command
}

func newPrinter(c *cobra.Command) *printer {
	return &printer{c: c}
}

func (p *printer) linesWithTitle(title string, lines []string) {
	if len(lines) > 0 {
		title = fmt.Sprintf("%s (%d)", title, len(lines))
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
