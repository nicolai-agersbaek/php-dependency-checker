package dependency_checker

import (
	"fmt"
	"github.com/z7zmey/php-parser/errors"
)

func logParserErrors(path string, errors []*errors.Error) {
	indent := "   "
	fmt.Println()
	fmt.Println(path, ":")

	for _, e := range errors {
		fmt.Println(indent, e)
	}
}
