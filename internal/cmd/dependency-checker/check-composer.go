package dependency_checker

import (
	"github.com/spf13/cobra"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/cmd"
	. "gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/util/slices"
)

var conf = &Config{
	SourceDir: SourceDirName,
	VendorDir: VendorDirName,
}

func init() {
	checkComposerCmd.Flags().StringVar(&conf.SourceDir, "src", "", "name of the source dir.")
	checkComposerCmd.Flags().StringVar(&conf.VendorDir, "vendor", "", "name of the vendor dir.")

	rootCmd.AddCommand(checkComposerCmd)
}

var checkComposerCmd = &cobra.Command{
	Use:   "check-composer <dir>",
	Short: "Check dependencies for Composer project in dir.",
	Args:  cobra.MinimumNArgs(1),
	Run:   checkComposer,
}

func checkComposer(c *cobra.Command, args []string) {
	root := args[0]

	// TODO: Move this logic to Checker type!
	// region <<- [ Perform analysis ] ->>

	// (Ensure root is valid as a Composer project)
	// ... implement this!

	// (Run composer install)
	// ... implement this!

	vendor, src := conf.VendorDirPath(root), conf.SourceDirPath(root)

	var srcImports, srcExports, vendorExports *Names
	var err error

	// Resolve vendorExports from 'vendor'
	_, vendorExports, err = ResolveImports(vendor)
	cmd.CheckError(err)

	// Resolve srcImports from 'src'
	srcImports, srcExports, err = ResolveImports(src)
	cmd.CheckError(err)

	// Calculate unexported uses
	// FIXME: Account for built-in names!
	allExports := vendorExports.Merge(srcExports)
	diff := slices.DiffString(srcImports.Classes, allExports.Classes)

	// endregion [ Perform analysis ]

	p := newPrinter(c)

	p.linesWithTitle("Unexported uses:", diff)
}
