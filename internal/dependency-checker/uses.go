package dependency_checker

import (
	"regexp"
)

var funcNamePattern = regexp.MustCompile("^[a-z](_*[A-Za-z0-9])*$")

func IsFunctionName(s string) bool {
	return funcNamePattern.MatchString(s)
}

var clsNamePattern = regexp.MustCompile("^\\\\*[A-Z][A-Za-z0-9]*(\\\\[A-Z][A-Za-z0-9]*)*$")

func IsClassName(s string) bool {
	return clsNamePattern.MatchString(s)
}
