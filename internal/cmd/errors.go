package cmd

import (
	"context"
	"fmt"
	"os"
)

// CheckError prints err to stderr and exits with code 1 if err is not nil. Otherwise, it is a
// no-op.
func CheckError(err error) {
	if err != nil {
		if err != context.Canceled {
			_, e := fmt.Fprintf(os.Stderr, fmt.Sprintf("An error occurred: %v\n", err))

			if e != nil {
				panic(e)
			}
		}
		os.Exit(1)
	}
}
