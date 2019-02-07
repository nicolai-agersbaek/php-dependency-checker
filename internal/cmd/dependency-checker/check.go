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
	Sources,
	Excludes,
	AdditionalExports,
	ExcludedExports,
	AdditionalImports,
	ExcludedImports []string
}

func NewCheckerInput() *CheckerInput {
	return &CheckerInput{}
}

var checkInput = &CheckerInput{}

func init() {
	excludeDesc := `Directory or file to exclude from analysis. May be
specified multiple times.`
	checkCmd.Flags().StringArrayVarP(&checkInput.Excludes, "exclude", "x", nil, excludeDesc)

	additionalExportsDesc := `Directory or file in which to search for additional
exports. May be specified multiple times.`
	checkCmd.Flags().StringArrayVarP(&checkInput.AdditionalExports, "exports", "e", nil, additionalExportsDesc)

	excludedExportsDesc := `Directory or file to exclude exports from. May be
specified multiple times.`
	checkCmd.Flags().StringArrayVarP(&checkInput.ExcludedExports, "exclude-exports", "E", nil, excludedExportsDesc)

	additionalImportsDesc := `Directory or file in which to search for additional
imports. May be specified multiple times.`
	checkCmd.Flags().StringArrayVarP(&checkInput.AdditionalImports, "imports", "i", nil, additionalImportsDesc)

	excludedImportsDesc := `Directory or file to exclude imports from. May be
specified multiple times.`
	checkCmd.Flags().StringArrayVarP(&checkInput.ExcludedImports, "exclude-imports", "I", nil, excludedImportsDesc)

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

	p := getVerbosePrinter(c)

	start := time.Now()

	// Calculate unexported uses.
	diff := doCheck(p, checkInput)

	stop := time.Now()
	elapsed := stop.Sub(start)

	p.VLine("Elapsed: "+elapsed.String(), cmd.VerbosityNormal)

	const maxLines = 15

	if len(diff.Classes) > 0 {
		p.VLinesWithTitleMax("Unexported uses (classes):", diff.Classes, maxLines, cmd.VerbosityNone)
		os.Exit(1)
	} else {
		p.VLine("No unexported uses found!", cmd.VerbosityNormal)
	}
}

func doCheck(p cmd.VerbosePrinter, input *CheckerInput) *Names {
	// Resolve imports and exports from sources.
	srcImports, srcExports, err := ResolveImportsSerial(p, input.Sources...)
	cmd.CheckError(err)

	// Resolve exports from specifically provided exporters.
	_, additionalExports, err := ResolveImportsSerial(p, input.AdditionalExports...)
	cmd.CheckError(err)

	// Resolve excluded exports from specifically provided exporters.
	_, excludedExports, err := ResolveImportsSerial(p, input.ExcludedExports...)
	cmd.CheckError(err)

	// Resolve imports from specifically provided importers.
	additionalImports, _, err := ResolveImportsSerial(p, input.AdditionalImports...)
	cmd.CheckError(err)

	// Resolve excluded imports from specifically provided importers.
	excludedImports, _, err := ResolveImportsSerial(p, input.ExcludedImports...)
	cmd.CheckError(err)

	// Resolve imports and exports to exclude from analysis.
	alsoExcludedImports, alsoExcludedExports, err := ResolveImportsSerial(p, input.Excludes...)
	cmd.CheckError(err)

	// Combine all analyses.
	imports := srcImports
	imports = imports.Diff(excludedImports, alsoExcludedImports).Merge(additionalImports)
	imports = consolidateIntoClasses(imports)

	exports := srcExports
	exports = exports.Diff(excludedExports, alsoExcludedExports).Merge(names.GetBuiltInNames(), additionalExports)
	exports = consolidateIntoClasses(exports)

	// Calculate unexported uses.
	return Diff(imports, exports)
}

func consolidateIntoClasses(n *Names) *Names {
	n.Classes = append(n.Classes, n.Interfaces...)
	n.Classes = append(n.Classes, n.Traits...)

	return n
}
