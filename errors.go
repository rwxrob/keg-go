package libkeg

import "fmt"

type ErrNotFound struct {
	It string
}

func (e ErrNotFound) Error() string { return fmt.Sprintf(_NotFound, e.It) }
