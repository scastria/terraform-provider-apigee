package client

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
)

const (
	FormEncoded = "application/x-www-form-urlencoded"
	IdSeparator = ":"
)

type Client struct {
	username     string
	password     string
	server       string
	port         int
	Organization string
	httpClient   *http.Client
}

func NewClient(username string, password string, server string, port int, organization string) *Client {
	return &Client{
		username:     username,
		password:     password,
		server:       server,
		port:         port,
		Organization: organization,
		httpClient:   &http.Client{},
	}
}

func (c *Client) HttpRequest(method string, path string, query url.Values, headerMap http.Header, body *bytes.Buffer) (closer io.ReadCloser, err error) {
	req, err := http.NewRequest(method, c.requestPath(path), body)
	if err != nil {
		return nil, &RequestError{StatusCode: http.StatusInternalServerError, Err: err}
	}
	req.SetBasicAuth(c.username, c.password)
	if query != nil {
		requestQuery := req.URL.Query()
		for key, values := range query {
			for _, value := range values {
				requestQuery.Add(key, value)
			}
		}
		req.URL.RawQuery = requestQuery.Encode()
	}
	if headerMap != nil {
		for key, values := range headerMap {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
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

// TODO: Allow non-SSL
func (c *Client) requestPath(path string) string {
	return fmt.Sprintf("https://%s:%d/v1/%s", c.server, c.port, path)
}

func GetMultiPartBuffer(filename string, key string) (*multipart.Writer, *bytes.Buffer, error) {
	buf := bytes.Buffer{}
	mp := multipart.NewWriter(&buf)
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()
	fw, err := mp.CreateFormFile(key, filename)
	if err != nil {
		return nil, nil, err
	}
	_, err = io.Copy(fw, file)
	if err != nil {
		return nil, nil, err
	}
	mp.Close()
	return mp, &buf, nil
}
