package namespace

import (
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/util/slices"
	"regexp"
)

type Validator func(*Namespace) bool

var nsPartPattern = regexp.MustCompile("^[A-Z][A-Za-z0-9]*$")

func IsValidNamespacePart(s string) bool {
	// FIXME: Missing tests!
	return nsPartPattern.MatchString(s)
}

var ValidCase Validator = func(ns *Namespace) bool {
	// FIXME: Missing tests!
	return slices.MatchAllString(ns.Parts(), IsValidNamespacePart)
}
