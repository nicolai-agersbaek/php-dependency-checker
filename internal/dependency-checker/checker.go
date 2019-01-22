package dependency_checker

import (
	"fmt"
	"github.com/z7zmey/php-parser/php7"
	"github.com/z7zmey/php-parser/visitor"
	"os"
)

const Name = "dependency-checker"

const Version = "0.1.0"

type Checker struct {
	Config *Config
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

	F, err := getFilesByExtension("php", paths...)

	if err != nil {
		return nil, err
	}

	for _, f := range uniqueStr(F) {
		uses, err := fileUses(f)

		if err != nil {
			return nil, err
		}

		M[f] = uses
	}

	return M, nil
}

func uniqueStr(strings []string) []string {
	U := make([]string, 0, len(strings))
	M := make(map[string]bool, len(strings))

	for _, str := range strings {
		if _, ok := M[str]; !ok {
			U = append(U, str)
			M[str] = true
		}
	}

	return U
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
