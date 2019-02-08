package cmd

import "github.com/spf13/cobra"

type Action func(c *cobra.Command, args []string)

func Wrap(f func()) Action {
	return func(c *cobra.Command, args []string) {
		f()
	}
}

type ActionE func(c *cobra.Command, args []string) error

func WrapE(f func() error) ActionE {
	return func(c *cobra.Command, args []string) error {
		return f()
	}
}
