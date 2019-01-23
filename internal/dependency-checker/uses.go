package dependency_checker

import (
	"github.com/z7zmey/php-parser/visitor"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/util/slices"
	"regexp"
)

func getUses(nsResolver *visitor.NamespaceResolver) []string {
	resolved := make([]string, len(nsResolver.ResolvedNames))

	i := 0
	for _, resolvedNs := range nsResolver.ResolvedNames {
		resolved[i] = resolvedNs
		i++
	}

	resolved = slices.UniqueString(resolved)
	resolved = removeNativeTypes(resolved)

	return resolved
}

var funcNamePattern = regexp.MustCompile("^(\\\\[A-Z][A-Za-z0-9]*)*\\\\*[a-z][A-Za-z0-9]*$")

func IsFunctionName(s string) bool {
	return funcNamePattern.MatchString(s)
}

func functionUses(nsResolver *visitor.NamespaceResolver) []string {
	return slices.FilterString(getUses(nsResolver), IsFunctionName)
}

var clsNamePattern = regexp.MustCompile("^\\\\*[A-Z][A-Za-z0-9]*(\\\\[A-Z][A-Za-z0-9]*)*$")

func IsClassName(s string) bool {
	return clsNamePattern.MatchString(s)
}

func classUses(nsResolver *visitor.NamespaceResolver) []string {
	return slices.FilterString(getUses(nsResolver), IsClassName)
}
