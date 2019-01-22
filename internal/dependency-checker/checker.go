package dependency_checker

import (
	"fmt"
	"github.com/z7zmey/php-parser/php7"
	"github.com/z7zmey/php-parser/visitor"
	"os"
	"path/filepath"
)

const Name = "dependency-checker"

const Version = "0.1.0"

type Checker struct {
	Config *Config
}

func NewChecker(config *Config) *Checker {
	return &Checker{Config: config}
}

func (c *Checker) Run(path string) error {
	fmt.Println("-----")
	fmt.Printf("Uses: (%s)", path)

	resolved, err := fileUses(path)

	if err != nil {
		return err
	}

	for _, resolvedNs := range resolved.Elements() {
		fmt.Println(resolvedNs)
	}

	fmt.Println("-----")

	return nil
}

// ResolveUses determines the set of classes used by the given path. If the
// given path is a file, it will analyze that file. If the path is a directory,
// it will recursively scan each file in the directory and return a combined set
// of (unique) classes used.
func (c *Checker) ResolveUses(paths ...string) (ClassUsesMap, error) {
	// TODO: Missing tests!
	M := make(ClassUsesMap, 0)

	for _, path := range paths {
		usesMap, err := pathUses(path)

		if err != nil {
			return M, err
		}

		M = M.merge(usesMap)
	}

	return M, nil
}

func pathUses(path string) (ClassUsesMap, error) {
	// TODO: Missing tests!
	// FIXME: Remove duplicates!
	M := make(ClassUsesMap, 0)
	allUses := NewStringSet()

	info, err := os.Stat(path)

	if err != nil {
		return M, err
	}

	if info.IsDir() {
		M, err = dirUses(path)
	} else {
		allUses, err = fileUses(path)
		M[path] = allUses
	}

	return M, err
}

func dirUses(dir string) (ClassUsesMap, error) {
	// TODO: Missing tests!
	// FIXME: Remove duplicates!
	M := make(ClassUsesMap, 0)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			uses, err := fileUses(path)

			if err != nil {
				return err
			}

			M[path] = uses
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return M, nil
}

func fileUses(path string) (*StringSet, error) {
	// TODO: Missing tests!
	// FIXME: Remove duplicates!
	src, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	parser := php7.NewParser(src, path)
	parser.Parse()

	for _, e := range parser.GetErrors() {
		fmt.Println(e)
	}

	nsResolver := visitor.NewNamespaceResolver()
	rootNode := parser.GetRootNode()

	rootNode.Walk(nsResolver)

	resolved := NewStringSet()

	for _, resolvedNs := range nsResolver.ResolvedNames {
		resolved.Put(resolvedNs)
	}

	return resolved, nil
}
