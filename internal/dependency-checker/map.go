package dependency_checker

type StringSet struct {
	elements []string
	contains map[string]bool
}

func NewStringSet() *StringSet {
	return &StringSet{
		elements: make([]string, 0),
		contains: make(map[string]bool, 0),
	}
}

func (s *StringSet) Has(e string) bool {
	_, ok := s.contains[e]

	return ok
}

func (s *StringSet) Put(E ...string) *StringSet {
	for _, e := range E {
		if !s.Has(e) {
			s.elements = append(s.elements, e)
		}
	}

	return s
}

func (s *StringSet) Elements() []string {
	return s.elements
}

func (s *StringSet) Merge(S *StringSet) *StringSet {
	return s.Put(S.Elements()...)
}

type ClassUsesMap map[string]*StringSet

func (m ClassUsesMap) has(key string) bool {
	// TODO: Missing tests!
	_, ok := m[key]

	return ok
}

func (m ClassUsesMap) merge(M ClassUsesMap) ClassUsesMap {
	// TODO: Missing tests!
	for p, uses := range M {
		if m.has(p) {
			m[p] = m[p].Merge(uses)
		}
	}

	return m
}
