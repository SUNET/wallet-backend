package apiv1

import (
	"context"
	"errors"
	"vc/pkg/model"
)

// MockNextRequest holds the request
type MockNextRequest struct {
	DocumentType    string `json:"document_type"`
	AuthenticSource string `json:"authentic_source"`
}

// MockNextReply is the reply
type MockNextReply struct {
	Upload *model.Upload `json:"upload"`
}

// MockNext sends one mock upload to the datastore
func (c *Client) MockNext(ctx context.Context, inData *MockNextRequest) (*MockNextReply, error) {
	// send to datastore
	mockUpload, err := c.mockOne(ctx, inData.AuthenticSource, inData.DocumentType)
	if err != nil {
		c.log.Debug("mocknext", "mockone", err)
		return nil, err
	}

	_, _, err = c.uploader(ctx, mockUpload)
	if err != nil {
		c.log.Debug("uploader", "error", err)
		return nil, err
	}

	reply := &MockNextReply{
		Upload: mockUpload,
	}

	return reply, nil
}

// MockBulkRequest holds the request
type MockBulkRequest struct {
	DocumentType    string `json:"document_type"`
	AuthenticSource string `json:"authentic_source"`
	N               int    `form:"n"`
}

// MockBulkReply is the reply
type MockBulkReply struct {
	DocumentIDS []string `json:"document_ids"`
}

// MockBulk sends N mock uploads to the datastore
func (c *Client) MockBulk(ctx context.Context, inData *MockBulkRequest) (*MockBulkReply, error) {
	documentIDS := []string{}

	if inData.N < 1 {
		return nil, errors.New("n must be greater than 0")
	}

	for i := 0; i < inData.N; i++ {
		mockUpload, err := c.mockOne(ctx, inData.AuthenticSource, inData.DocumentType)
		if err != nil {
			return nil, err
		}
		documentIDS = append(documentIDS, mockUpload.Meta.DocumentID)

		_, _, err = c.uploader(ctx, mockUpload)
		if err != nil {
			return nil, err
		}
	}

	return &MockBulkReply{
		DocumentIDS: documentIDS,
	}, nil
}
