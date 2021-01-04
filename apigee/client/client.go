package client

import (
	"bytes"
	"fmt"
	"github.com/go-http-utils/headers"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

const (
	FormEncoded        = "application/x-www-form-urlencoded"
	IdSeparator        = ":"
	Bearer             = "Bearer"
	PublicApigeeServer = "api.enterprise.apigee.com"
)

type Client struct {
	username     string
	password     string
	accessToken  string
	server       string
	port         int
	Organization string
	httpClient   *http.Client
}

func NewClient(username string, password string, accessToken string, server string, port int, organization string) *Client {
	return &Client{
		username:     username,
		password:     password,
		accessToken:  accessToken,
		server:       server,
		port:         port,
		Organization: organization,
		httpClient:   &http.Client{},
	}
}

func (c *Client) IsPublic() bool {
	return c.server == PublicApigeeServer
}

func (c *Client) HttpRequest(method string, path string, query url.Values, headerMap http.Header, body *bytes.Buffer) (closer io.ReadCloser, err error) {
	req, err := http.NewRequest(method, c.requestPath(path), body)
	if err != nil {
		return nil, &RequestError{StatusCode: http.StatusInternalServerError, Err: err}
	}
	//Handle query values
	if query != nil {
		requestQuery := req.URL.Query()
		for key, values := range query {
			for _, value := range values {
				requestQuery.Add(key, value)
			}
		}
		req.URL.RawQuery = requestQuery.Encode()
	}
	//Handle header values
	if headerMap != nil {
		for key, values := range headerMap {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
	}
	//Handle authentication
	if c.accessToken != "" {
		req.Header.Set(headers.Authorization, Bearer+" "+c.accessToken)
	} else {
		req.SetBasicAuth(c.username, c.password)
	}
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Print("Apigee Management API:")
		log.Print(err)
	} else {
		log.Print("Apigee Management API: " + string(requestDump))
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

func GetBuffer(filename string) (*bytes.Buffer, error) {
	buf := bytes.Buffer{}
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	_, err = buf.ReadFrom(file)
	if err != nil {
		return nil, err
	}
	return &buf, nil
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
