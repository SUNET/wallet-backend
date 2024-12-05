package httpserver

import (
	"context"
	apiv1_apigw "vc/internal/apigw/apiv1"
	"vc/internal/gen/status/apiv1_status"
)

// Apiv1 interface
type Apiv1 interface {
	Status(ctx context.Context, req *apiv1_status.StatusRequest) (*apiv1_status.StatusReply, error)

	SearchDocuments(ctx context.Context, request *apiv1_apigw.SearchDocumentsRequest) (*apiv1_apigw.SearchDocumentsReply, error)
}
