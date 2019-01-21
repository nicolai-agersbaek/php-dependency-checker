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
	src, err := os.Open(path)

	if err != nil {
		return err
	}

	parser := php7.NewParser(src, path)
	parser.Parse()

	for _, e := range parser.GetErrors() {
		fmt.Println(e)
	}

	v := visitor.Dumper{
		Writer: os.Stdout,
		Indent: "",
	}

	rootNode := parser.GetRootNode()
	rootNode.Walk(v)

	return nil
}
