package slices

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

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
