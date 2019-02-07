package checker

import (
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/cmd"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker/files"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/util/slices"
)

type Input struct {
	Sources,
	Excludes,
	AdditionalExports,
	ExcludedExports,
	AdditionalImports,
	ExcludedImports []string
}

func NewInput() *Input {
	return &Input{}
}

func (i *Input) ImportPaths() []string {
	// FIXME: Missing tests!
	include := slices.UniqueStrings(resolvePhpFiles(i.Sources), resolvePhpFiles(i.AdditionalImports))
	exclude := slices.UniqueStrings(resolvePhpFiles(i.Excludes), resolvePhpFiles(i.ExcludedImports))

	return slices.DiffString(include, exclude)
}

func (i *Input) ExportPaths() []string {
	// FIXME: Missing tests!
	include := slices.UniqueStrings(resolvePhpFiles(i.Sources), resolvePhpFiles(i.AdditionalExports))
	exclude := slices.UniqueStrings(resolvePhpFiles(i.Excludes), resolvePhpFiles(i.ExcludedExports))

	return slices.DiffString(include, exclude)
}

func resolvePhpFiles(paths []string) []string {
	// FIXME: Missing tests!
	var F, Fs []string
	var err error

	for _, path := range paths {
		F, err = files.GetFilesInDirByExtension("php", path)

		cmd.CheckError(err)

		Fs = append(Fs, F...)
	}

	return Fs
}
