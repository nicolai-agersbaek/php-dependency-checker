package dependency_checker

import (
	"github.com/spf13/cobra"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/cmd"
	. "gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker/names"
	"os"
	"time"
)

type CheckerInput struct {
	Sources   []string
	Exporters []string
}

func NewCheckerInput() *CheckerInput {
	return &CheckerInput{}
}

var checkInput = &CheckerInput{}

type checkCommandOptions struct {
	commandOptions
}

var checkOptions = &checkCommandOptions{}

func init() {
	exportsDesc := `Specify the directory in which to search for exports.
May be specified multiple times.`
	checkCmd.Flags().StringArrayVar(&checkInput.Exporters, "exports", nil, exportsDesc)

	checkCmd.Flags().BoolVar(&checkOptions.v, "v", false, "Output additional information.")
	checkCmd.Flags().BoolVar(&checkOptions.vv, "vv", false, "Output detailed information.")
	checkCmd.Flags().BoolVar(&checkOptions.vvv, "vvv", false, "Output debug information.")

	rootCmd.AddCommand(checkCmd)
}

var checkCmd = &cobra.Command{
	Use:   "check (<dir>|<file>) [(<dir>|<file>)] [,...]",
	Short: "Check dependencies for directories and/or files.",
	Args:  cobra.MinimumNArgs(1),
	Run:   check,
}

func check(c *cobra.Command, args []string) {
	checkInput.Sources = args

	start := time.Now()

	// Calculate unexported uses.
	diff := doCheck(checkInput)

	stop := time.Now()
	elapsed := stop.Sub(start)

	p := cmd.NewVerbosePrinter(cmd.NewPrinter(c), checkOptions.GetVerbosity())

	p.VLine("Elapsed: "+elapsed.String(), cmd.VerbosityNormal)

	const maxLines = 15

	if len(diff.Classes) > 0 {
		p.VLinesWithTitleMax("Unexported uses (classes):", diff.Classes, maxLines, cmd.VerbosityNone)
		os.Exit(1)
	} else {
		p.VLine("No unexported uses found!", cmd.VerbosityNormal)
	}
}

func doCheck(input *CheckerInput) *Names {
	// Resolve srcImports and srcExports from sources.
	srcImports, srcExports, err := ResolveImportsSerial(input.Sources...)
	cmd.CheckError(err)

	// Resolve srcExports from specifically provided exporters.
	_, exports, err := ResolveImportsSerial(input.Exporters...)
	cmd.CheckError(err)

	// Add built-in names and additional exports to list of available names.
	srcExports.Merge(names.GetBuiltInNames())
	srcExports.Merge(exports)

	srcExports = consolidateIntoClasses(srcExports)

	// Calculate unexported uses.
	return Diff(srcImports, srcExports)
}

func consolidateIntoClasses(n *Names) *Names {
	n.Classes = append(n.Classes, n.Interfaces...)
	n.Classes = append(n.Classes, n.Traits...)

	return n
}
