package resolver

import (
	"github.com/z7zmey/php-parser/walker"
	. "gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker/names"
)

type Resolver interface {
	walker.Visitor
	GetImports() *Names
	GetExports() *Names
}
