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

var phpNativeTypes = slices.FilterAllString(strings.Split(phpNativeTypeNames, "\n"), IsEmpty)

var phpNativeTypesPattern = anyOf(phpNativeTypes)

func notPhpNativeType(s string) bool {
	return !phpNativeTypesPattern.MatchString(s)
}

func anyOf(S []string) *regexp.Regexp {
	return regexp.MustCompile(anyOfPattern(S))
}

func anyOfPattern(S []string) string {
	filtered := slices.FilterString(S, IsEmpty)
	quoted := slices.MapString(filtered, regexp.QuoteMeta)

	return "^(" + strings.Join(quoted, "|") + ")$"
}

type Names struct {
	Functions  []string
	Classes    []string
	Constants  []string
	Namespaces []string
}

func NewNames() *Names {
	return &Names{
		make([]string, 0),
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
		n.Namespaces = append(n.Namespaces, fqn)
		return
	}

	if IsConstantName(fqn) {
		n.Constants = append(n.Constants, fqn)
		return
	}
}

func (n *Names) AddNs(ns string) {
	n.Namespaces = append(n.Namespaces, ns)
}

func (n *Names) Merge(names ...*Names) *Names {
	// FIXME: Missing tests!
	for _, nn := range names {
		n.Functions = append(n.Functions, nn.Functions...)
		n.Classes = append(n.Classes, nn.Classes...)
		n.Constants = append(n.Constants, nn.Constants...)
		n.Namespaces = append(n.Namespaces, nn.Namespaces...)
	}

	return n
}

func (n *Names) clean() {
	// FIXME: Missing tests!
	n.Functions = cleanResolved(n.Functions)
	n.Classes = cleanResolved(n.Classes)
	n.Constants = cleanResolved(n.Constants)
	n.Namespaces = cleanNamespaces(n.Namespaces)
}

func mergeNames(names []*Names) *Names {
	merged := NewNames()

	for _, n := range names {
		merged = merged.Merge(n)
	}

	return merged
}

func cleanResolved(resolved []string) []string {
	// FIXME: Missing tests!
	resolved = slices.FilterAllString(resolved, IsEmpty)
	resolved = removeNativeTypes(resolved)
	resolved = slices.UniqueString(resolved)

	return resolved
}

func removeNativeTypes(uses []string) []string {
	// FIXME: Missing tests!
	return slices.DiffString(uses, phpNativeTypes)
	//return slices.FilterAllString(uses, notPhpNativeType)
}

func cleanNamespaces(namespaces []string) []string {
	// FIXME: Missing tests!
	namespaces = slices.FilterAllString(namespaces, IsEmpty)
	namespaces = slices.UniqueString(namespaces)

	return namespaces
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
