package client

import "fmt"

type RequestError struct {
	StatusCode int
	Err        error
}

func (r *RequestError) Error() string {
	return fmt.Sprintf("Status %d: %v", r.StatusCode, r.Err)
}
