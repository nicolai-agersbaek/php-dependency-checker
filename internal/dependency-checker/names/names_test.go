package names

import (
	"testing"
)

type isClassNameTestArgs struct {
	s string
}

type isClassNameTestInput struct {
	name string
	args isClassNameTestArgs
	want bool
}

type isFunctionNameTestArgs struct {
	s string
}

type isFunctionNameTestInput struct {
	name string
	args isFunctionNameTestArgs
	want bool
}

func TestIsClassName(t *testing.T) {
	tests := []isClassNameTestInput{
		{
			"lc start",
			isClassNameTestArgs{"foo"},
			false,
		},
		{
			"lc start",
			isClassNameTestArgs{"foo\\bar"},
			false,
		},
		{
			"uc start",
			isClassNameTestArgs{"A"},
			true,
		},
		{
			"uc start",
			isClassNameTestArgs{"Foo"},
			true,
		},
		{
			"uc start",
			isClassNameTestArgs{"Foo\\Bar"},
			true,
		},
		{
			"mixed case",
			isClassNameTestArgs{"Foo\\Bar\\foo"},
			false,
		},
		{
			"invalid ending",
			isClassNameTestArgs{"A\\B\\"},
			false,
		},
		{
			"leading slash",
			isClassNameTestArgs{"\\A\\B"},
			true,
		},
		{
			"leading slash",
			isClassNameTestArgs{"\\A\\foo"},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsClassName(tt.args.s); got != tt.want {
				t.Errorf("IsClassName(%s) = %v, want %v", tt.args.s, got, tt.want)
			}
		})
	}
}

func TestIsFunctionName(t *testing.T) {
	tests := []isFunctionNameTestInput{
		{
			"lc start",
			isFunctionNameTestArgs{"foo"},
			true,
		},
		{
			"lc start",
			isFunctionNameTestArgs{"foo\\bar"},
			false,
		},
		{
			"with underscore",
			isFunctionNameTestArgs{"foo_bar"},
			true,
		},
		{
			"uc start",
			isFunctionNameTestArgs{"A"},
			false,
		},
		{
			"uc start",
			isFunctionNameTestArgs{"Foo"},
			false,
		},
		{
			"uc start",
			isFunctionNameTestArgs{"Foo\\Bar"},
			false,
		},
		{
			"mixed case",
			isFunctionNameTestArgs{"Foo\\Bar\\foo"},
			false,
		},
		{
			"invalid ending",
			isFunctionNameTestArgs{"A\\B\\"},
			false,
		},
		{
			"leading slash",
			isFunctionNameTestArgs{"\\A\\B"},
			false,
		},
		{
			"leading slash",
			isFunctionNameTestArgs{"\\A\\foo"},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsFunctionName(tt.args.s); got != tt.want {
				t.Errorf("IsFunctionName(%s) = %v, want %v", tt.args.s, got, tt.want)
			}
		})
	}
}
