package set

type StringSet struct {
	m map[string]bool
}

func (S *StringSet) Add(s string) {
	// FIXME: Missing tests!
	S.m[s] = true
}

func (S *StringSet) Elements() []string {
	// FIXME: Missing tests!
	E := make([]string, len(S.m))

	var i uint
	for s := range S.m {
		E[i] = s
		i++
	}

	return E
}
