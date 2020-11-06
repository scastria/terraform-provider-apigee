package client

import (
	"fmt"
	"net/http"
)

type RequestError struct {
	StatusCode int
	Err        error
}

func (r *RequestError) Error() string {
	return fmt.Sprintf("Status %d: Message: %s: %v", r.StatusCode, http.StatusText(r.StatusCode), r.Err)
}
