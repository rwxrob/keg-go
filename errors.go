package keg

import (
	"fmt"
	"net/http"
)

type ErrFetch struct {
	Resp *http.Response
}

func (e ErrFetch) Error() string { return fmt.Sprintf(_Fetch, e.Resp.Status) }
