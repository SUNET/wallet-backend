package datastoreclient

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"
	"vc/pkg/helpers"
	"vc/pkg/logger"
)

// Client is the client
type Client struct {
	httpClient *http.Client
	url        string
	log        *logger.Log

	Document *documentHandler
	Identity *identityHandler
}

// Config is the configuration for the client
type Config struct {
	URL string `validate:"required"`
}

// New creates a new client
func New(config *Config) (*Client, error) {
	if err := helpers.CheckSimple(config); err != nil {
		return nil, err
	}
	c := &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		url: config.URL,
		log: logger.NewSimple("datastoreclient"),
	}

	c.Document = &documentHandler{client: c, service: "api/v1/document", log: c.log.New("document")}
	c.Identity = &identityHandler{client: c, service: "api/v1/identity", log: c.log.New("identity")}

	return c, nil
}

// NewRequest make a new request
func (c *Client) newRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	rel, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(c.url)
	if err != nil {
		return nil, err
	}
	url := u.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err = json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// Do does the new request
func (c *Client) do(ctx context.Context, req *http.Request, reply any, prefixReplyJsonWithData bool) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if err := checkResponse(resp); err != nil {
		buf := &bytes.Buffer{}
		if _, err := buf.ReadFrom(resp.Body); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(buf.Bytes(), err); err != nil {
			return nil, err
		}
		return nil, err
	}

	var r any
	if prefixReplyJsonWithData {
		r = &struct {
			Data any `json:"data"`
		}{
			Data: reply,
		}
	} else {
		r = &reply
	}

	if err := json.NewDecoder(resp.Body).Decode(r); err != nil {
		c.log.Error(err, "failed to decode response")
		return nil, err
	}

	return resp, nil

}

// read body and make it reusable
func readBody(body io.ReadCloser) ([]byte, error) {
	buf := &bytes.Buffer{}
	if _, err := buf.ReadFrom(body); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func checkResponse(r *http.Response) error {
	switch r.StatusCode {
	case 200, 201, 202, 204, 304:
		return nil
	case 500:
		return ErrInvalidRequest
	case 401:
		return ErrNotAllowedRequest
	}

	return ErrInvalidRequest
}

func (c *Client) call(ctx context.Context, method, url string, body, reply any, prefixReplyJsonWithData bool) (*http.Response, error) {
	request, err := c.newRequest(
		ctx,
		method,
		url,
		body,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(ctx, request, reply, prefixReplyJsonWithData)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
