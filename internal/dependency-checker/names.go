package dependency_checker

import (
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/util/slices"
	"regexp"
	"strings"
)

const phpNativeTypeNames = `
bool
boolean
double
float
int
integer
null
NULL
object
string
true
false
void
self
static
parent
`

var phpNativeTypes = slices.FilterString(strings.Split(phpNativeTypeNames, "\n"), IsEmpty)

type Names struct {
	Functions []string
	Classes   []string
	Constants []string
}

func NewNames() *Names {
	return &Names{
		make([]string, 0),
		make([]string, 0),
		make([]string, 0),
	}
}

func (n *Names) Add(fqn string) {
	// FIXME: Missing tests!
	if IsFunctionName(fqn) {
		n.Functions = append(n.Functions, fqn)
		return
	}

	if IsClassName(fqn) {
		n.Classes = append(n.Classes, fqn)
		return
	}

	if IsConstantName(fqn) {
		n.Constants = append(n.Constants, fqn)
		return
	}
}

func (n *Names) Merge(names ...*Names) *Names {
	// FIXME: Missing tests!
	for _, nn := range names {
		n.Functions = append(n.Functions, nn.Functions...)
		n.Classes = append(n.Classes, nn.Classes...)
		n.Constants = append(n.Constants, nn.Constants...)
	}

	return n
}

func (n *Names) clean() {
	// FIXME: Missing tests!
	n.Functions = cleanResolved(n.Functions)
	n.Classes = cleanResolved(n.Classes)
	n.Constants = cleanResolved(n.Constants)
}

func cleanResolved(resolved []string) []string {
	// FIXME: Missing tests!
	resolved = slices.UniqueString(resolved)
	resolved = removeNativeTypes(resolved)
	resolved = slices.FilterString(resolved, IsEmpty)

	return resolved
}

func removeNativeTypes(uses []string) []string {
	// FIXME: Missing tests!
	return slices.DiffString(uses, phpNativeTypes)
}

var funcNamePattern = regexp.MustCompile("^[a-z](_*[A-Za-z0-9])*$")

func IsFunctionName(s string) bool {
	return funcNamePattern.MatchString(s)
}

var clsNamePattern = regexp.MustCompile("^\\\\*[A-Z][A-Za-z0-9]*(\\\\[A-Z][A-Za-z0-9]*)*$")

func IsClassName(s string) bool {
	return clsNamePattern.MatchString(s)
}

var constNamePattern = regexp.MustCompile("^[A-Z](_*[A-Z0-9])*$")

func IsConstantName(s string) bool {
	// FIXME: Missing tests!
	return constNamePattern.MatchString(s)
}

func IsEmpty(s string) bool {
	return s != ""
}
