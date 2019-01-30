package dependency_checker

import "regexp"

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

func (n *Names) Merge(names *Names) *Names {
	n.Functions = append(n.Functions, names.Constants...)
	n.Classes = append(n.Classes, names.Constants...)
	n.Constants = append(n.Constants, names.Constants...)

	return n
}

func (n *Names) clean() {
	n.Functions = cleanResolved(n.Functions)
	n.Classes = cleanResolved(n.Classes)
	n.Constants = cleanResolved(n.Constants)
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
