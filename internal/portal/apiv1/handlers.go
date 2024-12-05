package apiv1

import (
	"context"
	apiv1_apigw "vc/internal/apigw/apiv1"
	"vc/internal/gen/status/apiv1_status"
	"vc/pkg/model"
)

// Status return status for each ladok instance
func (c *Client) Status(ctx context.Context, req *apiv1_status.StatusRequest) (*apiv1_status.StatusReply, error) {
	probes := model.Probes{}

	status := probes.Check("registry")

	return status, nil
}

func (c *Client) SearchDocuments(ctx context.Context, req *apiv1_apigw.SearchDocumentsRequest) (*apiv1_apigw.SearchDocumentsReply, error) {
	reply, err := c.apigwClient.SearchDocuments(req)
	if err != nil {
		return nil, err
	}
	return reply, nil
}
