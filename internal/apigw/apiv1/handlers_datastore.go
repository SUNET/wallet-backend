package apiv1

import (
	"context"
	"time"
	"vc/internal/apigw/db"
	"vc/pkg/helpers"
	"vc/pkg/model"

	"go.opentelemetry.io/otel/codes"
)

// UploadRequest is the request for Upload
type UploadRequest struct {
	Meta                *model.MetaData        `json:"meta" validate:"required"`
	Identities          []model.Identity       `json:"identities,omitempty" validate:"dive"`
	DocumentDisplay     *model.DocumentDisplay `json:"document_display,omitempty"`
	DocumentData        map[string]any         `json:"document_data" validate:"required"`
	DocumentDataVersion string                 `json:"document_data_version,omitempty" validate:"required,semver"`
}

// Upload uploads a document with a set of attributes
//
//	@Summary		Upload
//	@ID				generic-upload
//	@Description	Upload endpoint
//	@Tags			dc4eu
//	@Accept			json
//	@Produce		json
//	@Success		200	"Success"
//	@Failure		400	{object}	helpers.ErrorResponse	"Bad Request"
//	@Param			req	body		UploadRequest			true	" "
//	@Router			/upload [post]
func (c *Client) Upload(ctx context.Context, req *UploadRequest) error {
	qr, err := req.Meta.QRGenerator(ctx, c.cfg.Common.QR.IssuingBaseURL, c.cfg.Common.QR.RecoveryLevel, c.cfg.Common.QR.Size)
	if err != nil {
		c.log.Debug("QR code generation failed", "error", err)
		return err
	}

	if req.Meta.Collect == nil || req.Meta.Collect.ID == "" {
		collect := &model.Collect{
			ID: req.Meta.DocumentID,
		}

		req.Meta.Collect = collect
	}

	if req.Meta.Revocation == nil {
		req.Meta.Revocation = &model.Revocation{
			ID: req.Meta.DocumentID,
		}
	} else {
		if req.Meta.Revocation.ID == "" {
			req.Meta.Revocation.ID = req.Meta.DocumentID
		}
	}

	upload := &model.CompleteDocument{
		Meta:                req.Meta,
		DocumentDisplay:     req.DocumentDisplay,
		DocumentData:        req.DocumentData,
		DocumentDataVersion: req.DocumentDataVersion,
		Identities:          req.Identities,
		QR:                  qr,
	}

	if err := helpers.ValidateDocumentData(ctx, upload, c.log); err != nil {
		return err
	}

	if upload.Identities == nil {
		upload.Identities = []model.Identity{}
	}

	if err := c.db.VCDatastoreColl.Save(ctx, upload); err != nil {
		c.log.Debug("Failed to save document", "error", err)
		return err
	}

	return nil
}

// NotificationRequest is the request for Notification
type NotificationRequest struct {
	AuthenticSource string `json:"authentic_source" validate:"required"`
	DocumentType    string `json:"document_type" validate:"required"`
	DocumentID      string `json:"document_id" validate:"required"`
}

// NotificationReply is the reply for a Notification
type NotificationReply struct {
	Data *model.QR `json:"data"`
}

// Notification return QR code and DeepLink for a document
//
//	@Summary		Notification
//	@ID				generic-notification
//	@Description	notification endpoint
//	@Tags			dc4eu
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	NotificationReply		"Success"
//	@Failure		400	{object}	helpers.ErrorResponse	"Bad Request"
//	@Param			req	body		NotificationRequest		true	" "
//	@Router			/notification [post]
func (c *Client) Notification(ctx context.Context, req *NotificationRequest) (*NotificationReply, error) {
	qrCode, err := c.db.VCDatastoreColl.GetQR(ctx, &model.MetaData{
		AuthenticSource: req.AuthenticSource,
		DocumentType:    req.DocumentType,
		DocumentID:      req.DocumentID,
	})
	if err != nil {
		return nil, err
	}

	reply := &NotificationReply{
		Data: qrCode,
	}
	return reply, nil
}

// IdentityMappingRequest is the request for IDMapping
type IdentityMappingRequest struct {
	// required: true
	// example: SUNET
	AuthenticSource string          `json:"authentic_source" validate:"required"`
	Identity        *model.Identity `json:"identity" validate:"required"`
}

// IdentityMappingReply is the reply for a IDMapping
type IdentityMappingReply struct {
	Data *model.IDMapping `json:"data"`
}

// IdentityMapping return a mapping between PID and AuthenticSource
//
//	@Summary		IdentityMapping
//	@ID				identity-mapping
//	@Description	Identity mapping endpoint
//	@Tags			dc4eu
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	IdentityMappingReply	"Success"
//	@Failure		400	{object}	helpers.ErrorResponse	"Bad Request"
//	@Param			req	body		IdentityMappingRequest	true	" "
//	@Router			/identity/mapping [post]
func (c *Client) IdentityMapping(ctx context.Context, reg *IdentityMappingRequest) (*IdentityMappingReply, error) {
	authenticSourcePersonID, err := c.db.VCDatastoreColl.IDMapping(ctx, &db.IDMappingQuery{
		AuthenticSource: reg.AuthenticSource,
		Identity:        reg.Identity,
	})
	if err != nil {
		return nil, err
	}

	reply := &IdentityMappingReply{
		Data: &model.IDMapping{
			AuthenticSourcePersonID: authenticSourcePersonID,
		},
	}

	return reply, nil
}

// AddDocumentIdentityRequest is the request for DocumentIdentity
type AddDocumentIdentityRequest struct {
	// required: true
	// example: SUNET
	AuthenticSource string `json:"authentic_source" validate:"required"`

	// required: true
	// example: PDA1
	DocumentType string `json:"document_type" validate:"required"`

	// required: true
	// example: 7a00fe1a-3e1a-11ef-9272-fb906803d1b8
	DocumentID string `json:"document_id" validate:"required"`

	Identities []*model.Identity `json:"identities" validate:"required"`
}

// AddDocumentIdentity adds an identity to a document
//
//	@Summary		AddDocumentIdentity
//	@ID				add-document-identity
//	@Description	Adding array of identities to one document
//	@Tags			dc4eu
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Failure		400	{object}	helpers.ErrorResponse		"Bad Request"
//	@Param			req	body		AddDocumentIdentityRequest	true	" "
//	@Router			/document/identity [put]
func (c *Client) AddDocumentIdentity(ctx context.Context, req *AddDocumentIdentityRequest) error {
	err := c.db.VCDatastoreColl.AddDocumentIdentity(ctx, &db.AddDocumentIdentityQuery{
		AuthenticSource: req.AuthenticSource,
		DocumentType:    req.DocumentType,
		DocumentID:      req.DocumentID,
		Identities:      req.Identities,
	})
	if err != nil {
		return err
	}

	return nil
}

// DeleteDocumentIdentityRequest is the request for DeleteDocumentIdentity
type DeleteDocumentIdentityRequest struct {
	// required: true
	// example: SUNET
	AuthenticSource string `json:"authentic_source" validate:"required"`

	// required: true
	// example: PDA1
	DocumentType string `json:"document_type" validate:"required"`

	// required: true
	// example: 7a00fe1a-3e1a-11ef-9272-fb906803d1b8
	DocumentID string `json:"document_id" validate:"required"`

	// required: true
	// example: 83c1a3c8-3e1a-11ef-9c01-6b6642c8d638
	AuthenticSourcePersonID string `json:"authentic_source_person_id" validate:"required"`
}

// DeleteDocumentIdentity deletes an identity from a document
//
//	@Summary		DeleteDocumentIdentity
//	@ID				delete-document-identity
//	@Description	Delete identity to document endpoint
//	@Tags			dc4eu
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Failure		400	{object}	helpers.ErrorResponse			"Bad Request"
//	@Param			req	body		DeleteDocumentIdentityRequest	true	" "
//	@Router			/document/identity [delete]
func (c *Client) DeleteDocumentIdentity(ctx context.Context, req *DeleteDocumentIdentityRequest) error {
	err := c.db.VCDatastoreColl.DeleteDocumentIdentity(ctx, &db.DeleteDocumentIdentityQuery{
		AuthenticSource:         req.AuthenticSource,
		DocumentType:            req.DocumentType,
		DocumentID:              req.DocumentID,
		AuthenticSourcePersonID: req.AuthenticSourcePersonID,
	})
	if err != nil {
		return err
	}

	return nil
}

// DeleteDocumentRequest is the request for DeleteDocument
type DeleteDocumentRequest struct {
	// required: true
	// example: skatteverket
	AuthenticSource string `json:"authentic_source" validate:"required"`

	// required: true
	// example: 5e7a981c-c03f-11ee-b116-9b12c59362b9
	DocumentID string `json:"document_id" validate:"required"`

	// required: true
	// example: PDA1
	DocumentType string `json:"document_type" validate:"required"`
}

// DeleteDocument deletes a specific document
//
//	@Summary		DeleteDocument
//	@ID				delete-document
//	@Description	delete one document endpoint
//	@Tags			dc4eu
//	@Accept			json
//	@Produce		json
//	@Success		200	"Success"
//	@Failure		400	{object}	helpers.ErrorResponse	"Bad Request"
//	@Param			req	body		DeleteDocumentRequest	true	" "
//	@Router			/document [delete]
func (c *Client) DeleteDocument(ctx context.Context, req *DeleteDocumentRequest) error {
	err := c.db.VCDatastoreColl.Delete(ctx, &model.MetaData{
		AuthenticSource: req.AuthenticSource,
		DocumentType:    req.DocumentType,
		DocumentID:      req.DocumentID,
	})
	if err != nil {
		return err
	}

	return nil
}

// GetDocumentRequest is the request for GetDocument
type GetDocumentRequest struct {
	AuthenticSource string `json:"authentic_source" validate:"required"`
	DocumentType    string `json:"document_type" validate:"required"`
	DocumentID      string `json:"document_id" validate:"required"`
}

// GetDocumentReply is the reply for a generic document
type GetDocumentReply struct {
	Data *model.Document `json:"data"`
}

// GetDocument return a specific document
//
//	@Summary		GetDocument
//	@ID				get-document
//	@Description	Get document endpoint
//	@Tags			dc4eu
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	GetDocumentReply		"Success"
//	@Failure		400	{object}	helpers.ErrorResponse	"Bad Request"
//	@Param			req	body		GetDocumentRequest		true	" "
//	@Router			/document [post]
func (c *Client) GetDocument(ctx context.Context, req *GetDocumentRequest) (*GetDocumentReply, error) {
	query := &db.GetDocumentQuery{
		Meta: &model.MetaData{
			AuthenticSource: req.AuthenticSource,
			DocumentType:    req.DocumentType,
			DocumentID:      req.DocumentID,
		},
	}
	doc, err := c.db.VCDatastoreColl.GetDocument(ctx, query)
	if err != nil {
		return nil, err
	}
	reply := &GetDocumentReply{
		Data: doc,
	}

	return reply, nil
}

// DocumentListRequest is the request for DocumentList
type DocumentListRequest struct {
	AuthenticSource string          `json:"authentic_source"`
	Identity        *model.Identity `json:"identity" validate:"required"`
	DocumentType    string          `json:"document_type"`
	ValidFrom       int64           `json:"valid_from"`
	ValidTo         int64           `json:"valid_to"`
}

// DocumentListReply is the reply for a list of documents
type DocumentListReply struct {
	Data []*model.DocumentList `json:"data"`
}

// DocumentList return a list of metadata for a specific identity
//
//	@Summary		DocumentList
//	@ID				document-list
//	@Description	List documents for an identity
//	@Tags			dc4eu
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	DocumentListReply		"Success"
//	@Failure		400	{object}	helpers.ErrorResponse	"Bad Request"
//	@Param			req	body		DocumentListRequest		true	" "
//	@Router			/document/list [post]
func (c *Client) DocumentList(ctx context.Context, req *DocumentListRequest) (*DocumentListReply, error) {
	docs, err := c.db.VCDatastoreColl.DocumentList(ctx, &db.DocumentListQuery{
		AuthenticSource: req.AuthenticSource,
		Identity:        req.Identity,
		DocumentType:    req.DocumentType,
		ValidFrom:       req.ValidFrom,
		ValidTo:         req.ValidTo,
	})
	if err != nil {
		return nil, err
	}
	resp := &DocumentListReply{
		Data: docs,
	}
	return resp, nil
}

// GetDocumentCollectIDRequest is the request for GetDocumentAttestation
type GetDocumentCollectIDRequest struct {
	AuthenticSource string          `json:"authentic_source" validate:"required"`
	DocumentType    string          `json:"document_type" validate:"required"`
	CollectID       string          `json:"collect_id" validate:"required"`
	Identity        *model.Identity `json:"identity" validate:"required"`
}

// GetDocumentCollectIDReply is the reply for a generic document
type GetDocumentCollectIDReply struct {
	Data *model.Document `json:"data"`
}

// GetDocumentCollectID return a specific document ??
//
//	@Summary		GetDocumentByCollectID
//	@ID				get-document-collect-id
//	@Description	Get one document with collect id
//	@Tags			dc4eu
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	GetDocumentCollectIDReply	"Success"
//	@Failure		400	{object}	helpers.ErrorResponse		"Bad Request"
//	@Param			req	body		GetDocumentCollectIDRequest	true	" "
//	@Router			/document/collect_id [post]
func (c *Client) GetDocumentCollectID(ctx context.Context, req *GetDocumentCollectIDRequest) (*GetDocumentCollectIDReply, error) {
	query := &db.GetDocumentCollectIDQuery{
		Identity: req.Identity,
		Meta: &model.MetaData{
			AuthenticSource: req.AuthenticSource,
			DocumentType:    req.DocumentType,
			Collect: &model.Collect{
				ID: req.CollectID,
			},
		},
	}

	doc, err := c.db.VCDatastoreColl.GetDocumentCollectID(ctx, query)
	if err != nil {
		return nil, err
	}

	reply := &GetDocumentCollectIDReply{
		Data: doc,
	}
	return reply, nil
}

// RevokeDocumentRequest is the request for RevokeDocument
type RevokeDocumentRequest struct {
	AuthenticSource string            `json:"authentic_source" validate:"required"`
	DocumentType    string            `json:"document_type" validate:"required"`
	Revocation      *model.Revocation `json:"revocation" validate:"required"`
}

// RevokeDocument revokes a specific document
//
//	@Summary		RevokeDocument
//	@ID				revoke-document
//	@Description	Revoke one document
//	@Tags			dc4eu
//	@Accept			json
//	@Produce		json
//	@Success		200	"Success"
//	@Failure		400	{object}	helpers.ErrorResponse	"Bad Request"
//	@Param			req	body		RevokeDocumentRequest	true	" "
//	@Router			/document/revoke [post]
func (c *Client) RevokeDocument(ctx context.Context, req *RevokeDocumentRequest) error {
	ctx, span := c.tracer.Start(ctx, "db:apigw:datastore:revoke")
	defer span.End()

	if req.Revocation.ID == "" {
		return helpers.ErrNoRevocationID
	}

	doc, err := c.db.VCDatastoreColl.GetByRevocationID(ctx, &model.MetaData{
		AuthenticSource: req.AuthenticSource,
		DocumentType:    req.DocumentType,
		Revocation:      &model.Revocation{ID: req.Revocation.ID},
	})
	if err != nil {
		return err
	}
	c.log.Debug("Document found", "document_id", doc.Meta.DocumentID)

	doc.Meta.Revocation = req.Revocation

	if req.Revocation.RevokedAt == 0 {
		doc.Meta.Revocation.RevokedAt = time.Now().Unix()
		doc.Meta.Revocation.Revoked = true
	}

	if err := c.db.VCDatastoreColl.Replace(ctx, doc); err != nil {
		span.SetStatus(codes.Error, err.Error())
		c.log.Error(err, "replace failed")
		return err
	}
	c.log.Debug("Document enqueued for update", "document_id", doc.Meta.DocumentID)

	return nil
}

// SearchDocuments search for documents
func (c *Client) SearchDocuments(ctx context.Context, req *model.SearchDocumentsRequest) (*model.SearchDocumentsReply, error) {
	docs, hasMore, err := c.db.VCDatastoreColl.SearchDocuments(ctx, &db.SearchDocumentsQuery{
		AuthenticSource: req.AuthenticSource,
		DocumentType:    req.DocumentType,
		DocumentID:      req.DocumentID,
		CollectID:       req.CollectID,

		AuthenticSourcePersonID: req.AuthenticSourcePersonID,
		FamilyName:              req.FamilyName,
		GivenName:               req.GivenName,
		BirthDate:               req.BirthDate,
	}, req.Limit, req.Fields, req.SortFields)

	if err != nil {
		return nil, err
	}
	resp := &model.SearchDocumentsReply{
		Documents:      docs,
		HasMoreResults: hasMore,
	}
	return resp, nil
}
