package dependency_checker

import (
	"github.com/spf13/cobra"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/cmd"
	. "gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker"
)

var conf = &Config{
	SourceDir: SourceDirName,
	VendorDir: VendorDirName,
}

func init() {
	checkComposerCmd.Flags().BoolVar(&conf.Install, "install", false, "run composer install.")
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

	res, err := analyze(root)
	cmd.CheckError(err)

	p := cmd.NewPrinter(c)

	p.LinesWithTitle("Unexported uses:", res.UnexportedUses.Namespaces)
}

type result struct {
	UnexportedUses *Names
}

func analyze(root string) (r *result, err error) {
	// TODO: Move this logic to Checker type!
	// region <<- [ Perform analysis ] ->>

	// (Ensure root is valid as a Composer project)
	// ... implement this!

	// (Run composer install)
	// ... implement this!

	vendor, src := conf.VendorDirPath(root), conf.SourceDirPath(root)

	var srcImports, srcExports, vendorExports *Names

	// Resolve vendorExports from 'vendor'
	_, vendorExports, err = ResolveImports(vendor)

	if err != nil {
		return r, err
	}

	// Resolve srcImports from 'src'
	srcImports, srcExports, err = ResolveImports(src)

	if err != nil {
		return r, err
	}

	// Calculate unexported uses
	// FIXME: Account for built-in names!
	allExports := vendorExports.Merge(srcExports)
	r.UnexportedUses = Diff(srcImports, allExports)

	return r, err
}
