package cmd

import (
	"context"
	"fmt"
	"os"
)

type Recoverable interface {
	error
	GetCode() int
}

type recoverable struct {
	error
	code int
}

func NewRecoverable(error error, code int) Recoverable {
	return &recoverable{error: error, code: code}
}

func (e *recoverable) GetCode() int {
	return e.code
}

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
