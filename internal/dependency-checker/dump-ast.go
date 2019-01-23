package dependency_checker

import (
	"fmt"
	"github.com/z7zmey/php-parser/php7"
	"github.com/z7zmey/php-parser/visitor"
	"io"
	"os"
)

func DumpAst(path string, writer io.Writer) error {
	src, err := os.Open(path)

	if err != nil {
		return err
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

	dumper := visitor.Dumper{
		Writer: writer,
		Indent: "",
		//Comments:  parser.GetComments(),
		//Positions: parser.GetPositions(),
	}

	rootNode := parser.GetRootNode()
	rootNode.Walk(dumper)

	return nil
}
