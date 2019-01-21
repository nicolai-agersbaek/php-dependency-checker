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

	for _, resolvedNs := range resolved {
		fmt.Println(resolvedNs)
	}

	fmt.Println("-----")

	return nil
}

func fileUses(path string) ([]string, error) {
	src, err := os.Open(path)

	if err != nil {
		return make([]string, 0), err
	}

	parser := php7.NewParser(src, path)
	parser.Parse()

	for _, e := range parser.GetErrors() {
		fmt.Println(e)
	}

	nsResolver := visitor.NewNamespaceResolver()
	rootNode := parser.GetRootNode()

	rootNode.Walk(nsResolver)

	resolved := make([]string, len(nsResolver.ResolvedNames))

	for _, resolvedNs := range nsResolver.ResolvedNames {
		resolved = append(resolved, resolvedNs)
	}

	return resolved, nil
}
