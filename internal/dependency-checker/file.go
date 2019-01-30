package dependency_checker

import (
	"os"
	"path/filepath"
	"regexp"
)

func getFilesByExtension(ext string, paths ...string) ([]string, error) {
	// FIXME: Missing tests!
	F := make([]string, 0)

	for _, path := range paths {
		files, err := getFilesInDirByExtension(ext, path)

		if err != nil {
			return nil, err
		}

		F = append(F, files...)
	}

	return F, nil
}

func getFilesInDirByExtension(ext, dir string) ([]string, error) {
	// FIXME: Missing tests!
	var files []string

	fileExtPattern := getFileExtPattern(ext)

	walk := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// If the path describes a file with the requested extension, add it to
		// the slice of files to return
		if !info.IsDir() && fileExtPattern.MatchString(info.Name()) {
			files = append(files, path)
		}

		return nil
	}

	err := filepath.Walk(dir, walk)

	if err != nil {
		return nil, err
	}

	return files, nil
}

func getFileExtPattern(ext string) *regexp.Regexp {
	// FIXME: Missing tests!
	return regexp.MustCompile(".*\\." + ext + "$")
}

func isDir(path string) bool {
	info, err := os.Stat(path)

	if err != nil {
		return false
	}

	return info.IsDir()
}
