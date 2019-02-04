package slices

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
