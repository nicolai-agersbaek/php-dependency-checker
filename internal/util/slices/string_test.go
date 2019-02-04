package slices

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

// region <<- [ DiffString ] ->>

type diffStringTestArgs struct {
	A []string
	B []string
}

type diffStringTestInput struct {
	name string
	args diffStringTestArgs
	want []string
}

func newDiffTestInput(A, B, want string) diffStringTestInput {
	sliceA := strings.Split(A, "")
	sliceB := strings.Split(B, "")
	sliceWant := strings.Split(want, "")

	stringA := strings.Join(sliceA, ",")
	stringB := strings.Join(sliceB, ",")

	name := fmt.Sprintf("A=[%s], B=[%s]", stringA, stringB)

	return diffStringTestInput{
		name,
		diffStringTestArgs{sliceA, sliceB},
		sliceWant,
	}
}

func TestDiffString(t *testing.T) {
	tests := []diffStringTestInput{
		newDiffTestInput("abcd", "cde", "ab"),
		newDiffTestInput("abc", "abcd", ""),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DiffString(tt.args.A, tt.args.B); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DiffString() = %v, want %v", got, tt.want)
			}
		})
	}
}

// endregion [ DiffString ]

// region <<- [ UniqueString ] ->>

type uniqueStringTestArgs struct {
	S []string
}

type uniqueStringTestInput struct {
	name string
	args uniqueStringTestArgs
	want []string
}

func newUniqueStringTestInput(S, want string) uniqueStringTestInput {
	sliceS := strings.Split(S, "")
	sliceWant := strings.Split(want, "")

	stringS := strings.Join(sliceS, ",")

	name := fmt.Sprintf("S=[%s]", stringS)

	return uniqueStringTestInput{
		name,
		uniqueStringTestArgs{sliceS},
		sliceWant,
	}
}

func TestUniqueString(t *testing.T) {
	tests := []uniqueStringTestInput{
		newUniqueStringTestInput("aaabbbccc", "abc"),
		newUniqueStringTestInput("abc c", "abc "),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UniqueString(tt.args.S); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UniqueString() = %v, want %v", got, tt.want)
			}
		})
	}
}

// endregion [ UniqueString ]

// region <<- [ FilterString ] ->>

type filterStringTestArgs struct {
	S []string
	f StringFilter
}

type filterStringTestInput struct {
	name string
	args filterStringTestArgs
	want []string
}

func TestFilterString(t *testing.T) {
	isEmpty := func(s string) bool {
		return s == ""
	}

	tests := []filterStringTestInput{
		{
			"is empty",
			filterStringTestArgs{
				[]string{"", "abc", " ", ""},
				isEmpty,
			},
			[]string{"", ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FilterString(tt.args.S, tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterString() = %v, want %v", got, tt.want)
			}
		})
	}
}

// endregion [ FilterString ]

// region <<- [ StringFilterNot ] ->>

func TestStringFilterNot(t *testing.T) {
	isEmpty := func(s string) bool {
		return s == ""
	}

	tests := []filterStringTestInput{
		{
			"not(is empty)",
			filterStringTestArgs{
				[]string{"", "abc", " ", ""},
				StringFilterNot(isEmpty),
			},
			[]string{"abc", " "},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FilterString(tt.args.S, tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterString() = %v, want %v", got, tt.want)
			}
		})
	}
}

// endregion [ StringFilterNot ]

// region <<- [ FilterStringAll ] ->>

type filterAllStringTestArgs struct {
	S []string
	F []StringFilter
}

type filterAllStringTestInput struct {
	name string
	args filterAllStringTestArgs
	want []string
}

func TestFilterAllString(t *testing.T) {
	longerThan3 := func(s string) bool {
		return len(s) > 3
	}
	shorterThan6 := func(s string) bool {
		return len(s) < 6
	}

	tests := []filterAllStringTestInput{
		{
			"3 < len(s) < 6",
			filterAllStringTestArgs{
				[]string{"", "1", "12", "123", "1234", "12345", "123456", "1234567"},
				[]StringFilter{longerThan3, shorterThan6},
			},
			[]string{"1234", "12345"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FilterAllString(tt.args.S, tt.args.F...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterAllString() = %v, want %v", got, tt.want)
			}
		})
	}
}

// endregion [ FilterStringAll ]

// region <<- [ MapString ] ->>

type mapStringTestArgs struct {
	S []string
	m StringMapping
}

type mapStringTestInput struct {
	name string
	args mapStringTestArgs
	want []string
}

func TestMapString(t *testing.T) {
	tests := []mapStringTestInput{
		{
			"upper-case",
			mapStringTestArgs{
				[]string{"", "a", "aB", "aBc", "a_b"},
				strings.ToUpper,
			},
			[]string{"", "A", "AB", "ABC", "A_B"},
		},
		{
			"lower-case",
			mapStringTestArgs{
				[]string{"", "a", "aB", "aBc", "a_b"},
				strings.ToLower,
			},
			[]string{"", "a", "ab", "abc", "a_b"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MapString(tt.args.S, tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapString() = %v, want %v", got, tt.want)
			}
		})
	}
}

// endregion [ MapString ]
