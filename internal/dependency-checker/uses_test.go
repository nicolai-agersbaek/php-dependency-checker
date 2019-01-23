package dependency_checker

import "testing"

type isClassNameTestArgs struct {
	s string
}

type isClassNameTestInput struct {
	name string
	args isClassNameTestArgs
	want bool
}

func Test_isClassName(t *testing.T) {
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
				t.Errorf("isClassName(%s) = %v, want %v", tt.args.s, got, tt.want)
			}
		})
	}
}
