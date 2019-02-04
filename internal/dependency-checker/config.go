package dependency_checker

import "path/filepath"

const (
	SourceDirName = "src"
	VendorDirName = "vendor"
)

type Config struct {
	SourceDir, VendorDir string
}

func (c *Config) SourceDirPath(rootPath string) string {
	return filepath.Join(rootPath, c.SourceDir)
}

func (c *Config) VendorDirPath(rootPath string) string {
	return filepath.Join(rootPath, c.VendorDir)
}
