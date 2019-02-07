package files

import (
	"os"
	"path/filepath"
	"regexp"
)

func GetFilesByExtension(ext string, paths ...string) ([]string, error) {
	// FIXME: Missing tests!
	F := make([]string, 0)

	for _, path := range paths {
		files, err := GetFilesInDirByExtension(ext, path)

		if err != nil {
			return nil, err
		}

		F = append(F, files...)
	}

	return F, nil
}

func GetFilesInDirByExtension(ext, dir string) ([]string, error) {
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
