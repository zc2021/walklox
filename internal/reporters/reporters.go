package reporters

import (
	"errors"
	"fmt"
	"io"
	"os"
)

type Accumulator struct {
	errors   []*loxError
	comments []*loxComment
}

type loxError struct {
	line    int
	message string
}

type loxComment struct {
	line    int
	message string
}

func (le loxError) String() string {
	return fmt.Sprintf("[line %d] Error: %s", le.line, le.message)
}

func (lc loxComment) String() string {
	return fmt.Sprintf("[line %d] Comment: %s", lc.line, lc.message)
}

func (a *Accumulator) AddError(ln int, msg string) error {
	err := loxError{line: ln, message: msg}
	return a.add(err)
}

func (a *Accumulator) AddComment(ln int, msg string) error {
	com := loxComment{line: ln, message: msg}
	return a.add(com)
}

func (a *Accumulator) add(item interface{}) error {
	switch it := item.(type) {
	case loxError:
		le := item.(loxError)
		a.errors = append(a.errors, &le)
	case loxComment:
		lc := item.(loxComment)
		a.comments = append(a.comments, &lc)
	default:
		errString := fmt.Sprintf("Accumulator cannot collect %T", it)
		return errors.New(errString)
	}
	return nil
}

func (a *Accumulator) HasErrors() bool {
	return len(a.errors) > 0
}

func (a *Accumulator) PrintErrors() {
	for _, err := range a.errors {
		report(os.Stderr, err)
	}
}

func (a *Accumulator) PrintComments() {
	for _, com := range a.comments {
		report(os.Stdout, com)
	}
}

func report(w io.Writer, i interface{}) {
	fmt.Fprintln(w, i)
}