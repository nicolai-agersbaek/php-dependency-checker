package dependency_checker

import (
	"github.com/spf13/cobra"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/cmd"
	. "gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/util/slices"
	"path/filepath"
)

const (
	sourceDirName = "src"
	vendorDirName = "vendor"
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
	root := args[0]

	// TODO: Move this logic to Checker type!
	// region <<- [ Perform analysis ] ->>

	// (Ensure root is valid as a Composer project)
	// ... implement this!

	// (Run composer install)
	// ... implement this!

	vendor, src := filepath.Join(root, vendorDirName), filepath.Join(root, sourceDirName)

	var srcImports, srcExports, vendorExports *Names
	var err error

	// Resolve vendorExports from 'vendor'
	_, vendorExports, err = ResolveImports(vendor)
	cmd.CheckError(err)

	// Resolve srcImports from 'src'
	srcImports, srcExports, err = ResolveImports(src)
	cmd.CheckError(err)

	// Calculate unexported uses
	allExports := vendorExports.Merge(srcExports)
	diff := slices.DiffString(srcImports.Classes, allExports.Classes)

	// endregion [ Perform analysis ]

	p := printer{c}

	p.linesWithTitle("Unexported uses:", diff)

	// Print uses
	//p.linesWithTitle("Imports (functions):", srcImports.Functions)
	//p.linesWithTitle("Imports (classes):", srcImports.Classes[:10])
	//p.linesWithTitle("Exports (functions):", vendorExports.Functions)
	//p.linesWithTitle("Exports (classes):", vendorExports.Classes[:10])
}
