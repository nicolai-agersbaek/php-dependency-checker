package dependency_checker

import (
	"fmt"
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
	imports, err := ResolveImports(args[0])
	cmd.CheckError(err)

	p := printer{c}

	// Print uses
	//p.linesWithTitle("Functions used:", imports.Imports.Functions)
	p.linesWithTitle("Imports:", imports.Imports.Classes)
	//p.linesWithTitle("Functions provided:", imports.Exports.Functions)
	p.linesWithTitle("Exports:", imports.Exports.Classes)
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

func uses(c *cobra.Command, args []string) {
	// Run the analysis
	fmt.Println("Analysing files...")

	funcUses, err := ResolveUses(args, IsFunctionName)
	cmd.CheckError(err)

	clsUses, err := ResolveUses(args, IsClassName)
	cmd.CheckError(err)

	// Print uses
	c.Println("----------  FUNCTIONS  ----------")
	printUses(c, funcUses)
	c.Println("----------  CLASSES  ----------")
	printUses(c, clsUses)
}

func printUses(c *cobra.Command, usesMap ClassUsesMap) {
	for file, uses := range usesMap {
		c.Printf("%s:\n", file)

		for _, use := range uses {
			c.Printf("%s%s\n", indent, use)
		}
	}
}
