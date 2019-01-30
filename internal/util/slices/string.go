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

func FilterString(S []string, filters ...StringFilter) []string {
	// FIXME: Missing tests!
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
