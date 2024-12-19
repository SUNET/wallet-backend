package datastoreclient

import (
	"context"
	"fmt"
	"net/http"
	"vc/pkg/logger"
	"vc/pkg/model"
)

type identityHandler struct {
	client  *Client
	service string
	log     *logger.Log
}

// IdentityMappingQuery is the query for IdentityMapping
type IdentityMappingQuery struct {
	AuthenticSource string          `json:"authentic_source"`
	Identity        *model.Identity `json:"identity"`
}

// Mapping maps an identity, return authentic_source_person_id
func (s *identityHandler) Mapping(ctx context.Context, query *IdentityMappingQuery) (string, *http.Response, error) {
	s.log.Info("Mapping")

	url := fmt.Sprintf("%s/%s", s.service, "mapping")
	reply := ""
	resp, err := s.client.call(ctx, http.MethodPost, url, nil, reply, true)
	if err != nil {
		return "", resp, err
	}
	return reply, resp, nil
}
