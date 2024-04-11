package apiv1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"vc/pkg/model"

	"github.com/masv3971/gosunetca/types"
)

type uploaderReply struct {
	Data struct {
		Status string `json:"status"`
	} `json:"data"`
}

func (c *Client) uploader(ctx context.Context, upload *model.Upload) (*uploaderReply, *http.Response, error) {
	reply := &uploaderReply{}
	resp, err := c.call(
		ctx,
		http.MethodPost,
		"/api/v1/upload",
		upload,
		reply,
	)
	if err != nil {
		c.log.Debug("uploader", "error", err)
		return nil, resp, err
	}

	return reply, resp, nil
}

// NewRequest make a new request
func (c *Client) newRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	rel, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(c.cfg.MockAS.DatastoreURL)
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
	c.log.Debug("newrequest", "method", method, "url", url.String())

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

// Do does the new request
func (c *Client) do(ctx context.Context, req *http.Request, value any) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.log.Debug("httpClient do", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	if err := checkResponse(resp); err != nil {
		buf := &bytes.Buffer{}
		if _, err := buf.ReadFrom(resp.Body); err != nil {
			c.log.Debug("ReadForm", "error", err)
			return nil, err
		}
		caError := &types.ErrorReply{}
		if err := json.Unmarshal(buf.Bytes(), caError); err != nil {
			c.log.Debug("json unmarshal", "error", err)
			return nil, err
		}
		c.log.Debug("checkResponse", "error", err)
		return nil, caError
	}

	if err := json.NewDecoder(resp.Body).Decode(value); err != nil {
		c.log.Debug("json decode", "error", err)
		return nil, err
	}

	return resp, nil
}

func checkResponse(r *http.Response) error {
	switch r.StatusCode {
	case 200, 201, 202, 204, 304:
		return nil
	case 401:
		return errors.New("Unauthorized")
	case 500:
		return errors.New("Invalid")
	}
	return errors.New("Invalid request")
}

func (c *Client) call(ctx context.Context, method, path string, body, reply any) (*http.Response, error) {
	request, err := c.newRequest(
		ctx,
		method,
		path,
		body,
	)
	if err != nil {
		c.log.Debug("call", "error", err)
		return nil, err
	}

	resp, err := c.do(ctx, request, reply)
	if err != nil {
		c.log.Debug("do", "error", err)
		return resp, err
	}

	return resp, nil
}
