package dependency_checker

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/cmd"
	. "gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker"
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
	//Run:   listFiles,
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

	var imports, exports *Names
	var err error

	// Resolve exports from 'vendor'
	_, exports, err = ResolveImports(vendor)
	cmd.CheckError(err)

	// Resolve imports from 'src'
	imports, _, err = ResolveImports(src)
	cmd.CheckError(err)

	// Calculate unexported uses
	//cmd.CheckError(err)

	// endregion [ Perform analysis ]

	p := printer{c}

	// Print uses
	//p.linesWithTitle("Imports (functions):", imports.Functions)
	p.linesWithTitle("Imports (classes):", imports.Classes)
	//p.linesWithTitle("Exports (functions):", exports.Functions)
	p.linesWithTitle("Exports (classes):", exports.Classes[:10])
}

func listFiles(c *cobra.Command, args []string) {
	root := args[0]
	//root, err := filepath.Abs(args[0])
	//cmd.CheckError(err)

	src := filepath.Join(root, sourceDirName)
	pattern := filepath.Join(root, src, "**/*.php")

	p := printer{c}

	printFiles(p, pattern)
}

func printFiles(p printer, pattern string) {
	files, err := filepath.Glob(pattern)
	cmd.CheckError(err)

	p.title(fmt.Sprintf("Files matching [%s]:", pattern))
	p.lines(files)
}
