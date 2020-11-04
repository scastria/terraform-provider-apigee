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
		return nil, err
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
		return nil, err
	}
	if (resp.StatusCode < http.StatusOK) || (resp.StatusCode >= http.StatusMultipleChoices) {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("got a non 2XX status code: %v", resp.StatusCode)
		}
		return nil, fmt.Errorf("got a non 2XX status code: %v - %s", resp.StatusCode, respBody.String())
	}
	return resp.Body, nil
}

func (c *Client) requestPath(path string) string {
	return fmt.Sprintf("https://%s:%d/v1/%s", c.server, c.port, path)
}
