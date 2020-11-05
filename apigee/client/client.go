package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	username     string
	password     string
	server       string
	port         int
	organization string
	httpClient   *http.Client
}

func NewClient(username string, password string, server string, port int, organization string) *Client {
	return &Client{
		username:     username,
		password:     password,
		server:       server,
		port:         port,
		organization: organization,
		httpClient:   &http.Client{},
	}
}

func (c *Client) HttpRequest(path string, method string, body bytes.Buffer) (closer io.ReadCloser, err error) {
	req, err := http.NewRequest(method, c.requestPath(path), &body)
	if err != nil {
		return nil, &RequestError{StatusCode: http.StatusInternalServerError, Err: err}
	}
	req.SetBasicAuth(c.username, c.password)
	switch method {
	case "GET":
	case "DELETE":
	default:
		req.Header.Add("Content-Type", "application/json")
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &RequestError{StatusCode: http.StatusInternalServerError, Err: err}
	}
	if (resp.StatusCode < http.StatusOK) || (resp.StatusCode >= http.StatusMultipleChoices) {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(resp.Body)
		if err != nil {
			return nil, &RequestError{StatusCode: resp.StatusCode, Err: err}
		}
		return nil, &RequestError{StatusCode: resp.StatusCode, Err: fmt.Errorf("%s", respBody.String())}
	}
	return resp.Body, nil
}

func (c *Client) requestPath(path string) string {
	return fmt.Sprintf("https://%s:%d/v1/%s", c.server, c.port, path)
}
