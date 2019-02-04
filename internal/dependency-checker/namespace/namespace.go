package namespace

import "strings"

const Separator = "\\"

type Namespace struct {
	parts []string
}

func New(parts []string) *Namespace {
	return &Namespace{parts: parts}
}

func (n *Namespace) String() string {
	return strings.Join(n.parts, Separator)
}

func (n *Namespace) Parts() []string {
	return n.parts
}
