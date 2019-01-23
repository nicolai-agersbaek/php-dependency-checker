package dependency_checker

import (
	"fmt"
	"github.com/z7zmey/php-parser/php7"
	"github.com/z7zmey/php-parser/visitor"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/util/slices"
	"os"
)

const Name = "dependency-checker"

const Version = "0.1.0"

var phpNativeTypes = []string{
	"bool",
	"boolean",
	"double",
	"float",
	"int",
	"integer",
	"null",
	"NULL",
	"object",
	"string",

	"true",
	"false",
	"void",

	"self",
	"static",
	"parent",
}

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

	for _, resolvedNs := range resolved {
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

	for _, f := range slices.UniqueString(F) {
		uses, err := fileUses(f)

		if err != nil {
			return nil, err
		}

		M[f] = uses
	}

	return M, nil
}

func removeNativeTypes(uses []string) []string {
	return slices.DiffString(uses, phpNativeTypes)
}

func fileUses(path string) ([]string, error) {
	// TODO: Missing tests!
	src, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer func() {
		if err := src.Close(); err != nil {
			panic(err)
		}
	}()

	parser := php7.NewParser(src, path)
	parser.Parse()

	for _, e := range parser.GetErrors() {
		fmt.Println(e)
	}

	nsResolver := visitor.NewNamespaceResolver()
	rootNode := parser.GetRootNode()

	rootNode.Walk(nsResolver)

	resolved := make([]string, len(nsResolver.ResolvedNames))

	i := 0
	for _, resolvedNs := range nsResolver.ResolvedNames {
		resolved[i] = resolvedNs
		i++
	}

	resolved = slices.UniqueString(resolved)
	resolved = removeNativeTypes(resolved)

	return resolved, nil
}
