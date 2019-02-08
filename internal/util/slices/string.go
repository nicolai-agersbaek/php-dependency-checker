package slices

func MergeString(Ss ...[]string) []string {
	// FIXME: Missing tests!
	switch len(Ss) {
	case 0:
		return make([]string, 0)
	case 1:
		return Ss[0]
	default:
		for _, S := range Ss[0:] {
			Ss[0] = append(Ss[0], S...)
		}

		return Ss[0]
	}
}

func MergeUniqueString(Ss ...[]string) []string {
	// FIXME: Missing tests!
	// If no arguments were passes, we simply return an empty slice.
	if len(Ss) == 0 {
		return []string{}
	}

	M := make(map[string]bool, 0)

	// TODO: Benchmark performance when using two-pass w/ map, vs. one-pass with looped U.append
	for _, S := range Ss {
		for _, s := range S {
			if _, ok := M[s]; ok {
				M[s] = true
			}
		}
	}

	// Extract keys of M to form the desired slice.
	U := make([]string, len(M))

	var i int
	for k := range M {
		U[i] = k
		i++
	}

	return U
}

func DiffAllString(A, B []string, Ss ...[]string) []string {
	// FIXME: Missing tests!
	M := make(map[string]bool, len(A))

	for _, a := range A {
		M[a] = true
	}

	// Create an exclude set from B and all S ∈ Ss
	exclude := MergeUniqueString(append(Ss, B)...)

	for _, e := range exclude {
		if _, ok := M[e]; ok {
			M[e] = false
		}
	}

	C := make([]string, 0)

	for m, ok := range M {
		if ok {
			C = append(C, m)
		}
	}

	return C
}

func DiffString(A, B []string) []string {
	M := make(map[string]bool, len(A))

	for _, a := range A {
		M[a] = true
	}

	for _, b := range B {
		if _, ok := M[b]; ok {
			M[b] = false
		}
	}

	C := make([]string, 0)

	for m, ok := range M {
		if ok {
			C = append(C, m)
		}
	}

	return C
}

// IntersectionString returns the intersection of elements in the two slices A
// and B. Thus, I = { i | i ∈ A ^ i ∈ B}.
func IntersectionString(A, B []string) []string {
	// FIXME: Missing tests!
	M := make(map[string]bool, len(A))

	for _, a := range A {
		M[a] = true
	}

	for _, b := range B {
		if _, ok := M[b]; !ok {
			M[b] = false
		}
	}

	I := make([]string, 0)

	for m, ok := range M {
		if ok {
			I = append(I, m)
		}
	}

	return I
}

func SliceString(S []string, from, to int) []string {
	// FIXME: Missing tests!
	from = intMax(from, 0)
	to = intMin(to, len(S))

	return S[from:to]
}

func intMin(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func intMax(a, b int) int {
	if a > b {
		return a
	}

	return b
}

// EmptyStrings determines if all slices S ∈ Ss are empty.
// The definition of empty is given by EmptyString.
func EmptyStrings(Ss ...[]string) bool {
	// FIXME: Missing tests!
	for _, S := range Ss {
		if !EmptyString(S) {
			return false
		}
	}

	return true
}

// EmptyString determines if the given slice of strings, S, is empty.
// This means that either len(S) == 0, or s == "" for all s ∈ S.
func EmptyString(S []string) bool {
	// FIXME: Missing tests!
	for _, s := range S {
		if s != "" {
			return false
		}
	}

	return true
}

func UniqueString(S []string) []string {
	U := make([]string, 0, len(S))
	M := make(map[string]bool, len(S))

	for _, str := range S {
		if _, ok := M[str]; !ok {
			U = append(U, str)
			M[str] = true
		}
	}

	return U
}

func UniqueStrings(Ss ...[]string) []string {
	// FIXME: Missing tests!
	U := make([]string, 0)
	M := make(map[string]bool, 0)

	for _, S := range Ss {
		for _, s := range S {
			if _, ok := M[s]; !ok {
				U = append(U, s)
				M[s] = true
			}
		}
	}

	return U
}

type StringFilter func(s string) bool

func StringFilterNot(filter StringFilter) StringFilter {
	return func(s string) bool {
		return !filter(s)
	}
}

func FilterString(S []string, filter StringFilter) []string {
	return FilterAllString(S, filter)
}

func FilterAllString(S []string, filters ...StringFilter) []string {
	F := make([]string, 0, len(S))

InputLoop:
	for _, str := range S {
		// apply filters
		for _, filter := range filters {
			if !filter(str) {
				continue InputLoop
			}
		}

		F = append(F, str)
	}

	return F
}

type StringMapping func(s string) string

func MapString(S []string, m StringMapping) []string {
	M := make([]string, len(S))

	var i int
	for _, str := range S {
		M[i] = m(str)
		i++
	}

	return M
}

func MatchAllString(S []string, filters ...StringFilter) bool {
	// FIXME: Missing tests!
	for _, str := range S {
		// apply filters
		for _, filter := range filters {
			// if ANY filter fails, we exit and return false
			if !filter(str) {
				return false
			}
		}
	}

	// no filters failed, and all strings in S match all filters
	return true
}
