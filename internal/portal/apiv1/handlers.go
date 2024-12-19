package apiv1

import (
	"context"
	"errors"
	"net/http"
	"vc/internal/gen/status/apiv1_status"
	"vc/pkg/model"
)

// Status return status for each ladok instance
func (c *Client) Status(ctx context.Context, req *apiv1_status.StatusRequest) (*apiv1_status.StatusReply, error) {
	probes := model.Probes{}

	status := probes.Check("registry")

	return status, nil
}

// SearchDocuments search for documents
func (c *Client) SearchDocuments(ctx context.Context, req *model.SearchDocumentsRequest) (*model.SearchDocumentsReply, error) {
	//reply, err := c.apigwClient.SearchDocuments(req)
	reply, httpResponse, err := c.datastoreClient.Document.Search(ctx, req)
	if err != nil {
		return nil, err
	}
	if httpResponse.StatusCode != http.StatusOK {
		return nil, errors.New(httpResponse.Status)
	}
	return reply, nil
}
