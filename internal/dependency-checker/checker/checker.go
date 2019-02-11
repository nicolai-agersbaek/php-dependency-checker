package checker

import (
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker/names"
)

type Result struct {
	Imports, Exports, Diff                   *names.Names
	ImportsByFile, ExportsByFile, DiffByFile *names.NamesByFile
}

type Checker struct {
}

func (c *Checker) Run(input *dependency_checker.Input, parallel bool) (*Result, error) {
	return nil, nil
}
