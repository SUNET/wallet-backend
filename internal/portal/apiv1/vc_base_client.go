package apiv1

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
	"vc/pkg/logger"
	"vc/pkg/trace"
)

type VCBaseClient struct {
	serviceName string
	baseUrl     string
	logger      *logger.Log
	httpClient  *http.Client
	tp          *trace.Tracer
}

func NewClient(serviceName string, baseUrl string, tracer *trace.Tracer, logger *logger.Log) *VCBaseClient {
	client := &VCBaseClient{
		serviceName: serviceName,
		baseUrl:     baseUrl,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger: logger,
		tp:     tracer,
	}
	return client
}

func DoPostJSONGeneric[T any](c *VCBaseClient, endpoint string, reqBody any) (*T, error) {
	url := c.url(endpoint)

	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBodyJSON))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer c.closeBody(resp)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var jsonResp T
	if err := json.Unmarshal(body, &jsonResp); err != nil {
		return nil, err
	}

	return &jsonResp, nil
}

func (c *VCBaseClient) url(path string) string {
	return c.baseUrl + path
}

func (c *VCBaseClient) closeBody(resp *http.Response) {
	func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			c.logger.Error(err, "could not close response body")
		}
	}(resp.Body)
}
