package apiv1

import (
	"context"
	"time"
	"vc/internal/apigw/db"
	"vc/pkg/helpers"
	"vc/pkg/model"

	"go.opentelemetry.io/otel/codes"
)

// UploadReply is the reply for a generic upload
type UploadReply struct {
	Data struct {
		Status string `json:"status"`
	} `json:"data"`
}

// Upload uploads a document with a set of attributes
//
//	@Summary		Upload
//	@ID				generic-upload
//	@Description	Upload endpoint
//	@Tags			dc4eu
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	UploadReply				"Success"
//	@Failure		400	{object}	helpers.ErrorResponse	"Bad Request"
//	@Param			req	body		model.Upload			true	" "
//	@Router			/upload [post]
func (c *Client) Upload(ctx context.Context, req *model.Upload) error {
	ctx, span := c.tp.Start(ctx, "apiv1:Upload")
	defer span.End()

	if err := helpers.Check(ctx, c.cfg, req, c.log); err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	if err := req.QRGenerator(ctx, c.cfg.Common.QR.BaseURL, c.cfg.Common.QR.RecoveryLevel, c.cfg.Common.QR.Size); err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	req.Meta.CreatedAt = time.Now().UTC()

	_, err := c.simpleQueue.VCPersistentSave.Enqueue(ctx, req)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	return nil
}

// NotificationRequest is the request for Notification
type NotificationRequest struct {
	AuthenticSource string `json:"authentic_source" required:"true"`
	DocumentType    string `json:"document_type" required:"true"`
	DocumentID      string `json:"document_id" required:"true"`
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
//	@Success		200	{object}	NotificationReply				"Success"
//	@Failure		400	{object}	helpers.ErrorResponse	"Bad Request"
//	@Param			req	body		NotificationRequest			true	" "
//	@Router			/notification [post]
func (c *Client) Notification(ctx context.Context, req *NotificationRequest) (*NotificationReply, error) {
	ctx, span := c.tp.Start(ctx, "apiv1:Notification")
	defer span.End()

	qrCode, err := c.db.VCDatastoreColl.GetQR(ctx, &model.MetaData{
		AuthenticSource: req.AuthenticSource,
		DocumentType:    req.DocumentType,
		DocumentID:      req.DocumentID,
	})
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	reply := &NotificationReply{
		Data: qrCode,
	}
	return reply, nil
}

// IDMappingRequest is the request for IDMapping
type IDMappingRequest struct {
	AuthenticSource string          `json:"authentic_source" validate:"required"`
	Identity        *model.Identity `json:"identity" validate:"required"`
}

// IDMappingReply is the reply for a IDMapping
type IDMappingReply struct {
	Data struct {
		AuthenticSourcePersonID string `json:"authentic_source_person_id"`
	} `json:"data"`
}

// IDMapping return a mapping between PID and AuthenticSource
//
//	@Summary		IDMapping
//	@ID				id-mapping
//	@Description	ID mapping endpoint
//	@Tags			dc4eu
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	IDMappingReply			"Success"
//	@Failure		400	{object}	helpers.ErrorResponse	"Bad Request"
//	@Param			req	body		model.MetaData			true	" "
//	@Router			/id_mapping [post]
func (c *Client) IDMapping(ctx context.Context, reg *IDMappingRequest) (*IDMappingReply, error) {
	ctx, span := c.tp.Start(ctx, "apiv1:IDMapping")
	defer span.End()

	if err := helpers.Check(ctx, c.cfg, reg, c.log); err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	authenticSourcePersonID, err := c.db.VCDatastoreColl.IDMapping(ctx, &db.IDMappingQuery{
		AuthenticSource: reg.AuthenticSource,
		Identity:        reg.Identity,
	})
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	reply := &IDMappingReply{
		Data: struct {
			AuthenticSourcePersonID string `json:"authentic_source_person_id"`
		}{
			AuthenticSourcePersonID: authenticSourcePersonID,
		},
	}
	return reply, nil

}

// GetDocumentRequest is the request for GetDocument
type GetDocumentRequest struct {
	DocumentID      string `json:"document_id"`
	AuthenticSource string `json:"authentic_source"`
	DocumentType    string `json:"document_type"`
}

// GetDocumentReply is the reply for a generic document
type GetDocumentReply struct {
	Data *model.Upload `json:"data"`
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
	ctx, span := c.tp.Start(ctx, "apiv1:GetDocument")
	defer span.End()

	query := &model.MetaData{
		DocumentID:      req.DocumentID,
		DocumentType:    req.DocumentType,
		AuthenticSource: req.AuthenticSource,
	}
	doc, err := c.db.VCDatastoreColl.GetDocument(ctx, query)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	reply := &GetDocumentReply{
		Data: doc,
	}
	return reply, nil
}

// DeleteDocumentRequest is the request for DeleteDocument
type DeleteDocumentRequest struct {
	// required: true
	// example: skatteverket
	AuthenticSource string `json:"authentic_source" required:"true"`

	// required: true
	// example: 5e7a981c-c03f-11ee-b116-9b12c59362b9
	DocumentID string `json:"document_id" required:"true"`
}

// DeleteDocument deletes a specific document
//
//	@Summary		DeleteDocument
//	@ID				delete-document
//	@Description	delete one document endpoint
//	@Tags			dc4eu
//	@Accept			json
//	@Produce		json
//	@Success		200				"Success"
//	@Failure		400	{object}	helpers.ErrorResponse	"Bad Request"
//	@Param			req	body		DeleteDocumentRequest		true	" "
//	@Router			/document [delete]
func (c *Client) DeleteDocument(ctx context.Context, req *DeleteDocumentRequest) error {
	ctx, span := c.tp.Start(ctx, "apiv1:DeleteDocument")
	defer span.End()

	if err := helpers.Check(ctx, c.cfg, req, c.log); err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	_, err := c.simpleQueue.VCPersistentDelete.Enqueue(ctx, req)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}

// GetDocumentAttestationRequest is the request for GetDocumentAttestation
type GetDocumentAttestationRequest struct {
	Identity *model.Identity `json:"identity"`
	Meta     *model.MetaData `json:"meta"`
}

// GetDocumentAttestationReply is the reply for a generic document
type GetDocumentAttestationReply struct {
	Data *model.Upload `json:"data"`
}

// GetDocumentAttestation return a specific document ??
//
//	@Summary		GetDocumentByCollectCode
//	@ID				get-document-collect-code
//	@Description	Get document by collect code endpoint
//	@Tags			dc4eu
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	GetDocumentReply		"Success"
//	@Failure		400	{object}	helpers.ErrorResponse	"Bad Request"
//	@Param			req	body		model.MetaData			true	" "
//	@Router			/document/collection_code [post]
func (c *Client) GetDocumentAttestation(ctx context.Context, req *GetDocumentAttestationRequest) (*GetDocumentAttestationReply, error) {
	ctx, span := c.tp.Start(ctx, "apiv1:GetDocumentAttestation")
	defer span.End()

	if err := helpers.Check(ctx, c.cfg, req, c.log); err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	query := &db.GetDocumentAttestationQuery{
		Identity: req.Identity,
		Meta:     req.Meta,
	}

	doc, err := c.db.VCDatastoreColl.GetDocumentAttestation(ctx, query)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	reply := &GetDocumentAttestationReply{
		Data: doc,
	}
	return reply, nil
}

// PortalRequest is the request for PortalData
type PortalRequest struct {
	AuthenticSource         string `json:"authentic_source" validate:"required"`
	AuthenticSourcePersonID string `json:"authentic_source_person_id" validate:"required"`
	ValidFrom               string `json:"validity_from"`
	ValidTo                 string `json:"validity_to"`
}

// PortalReply is the reply for PortalData
type PortalReply struct {
	Data []model.Upload `json:"data"`
}

// Portal return a list of metadata for a specific person
//
//	@Summary		Portal
//	@ID				portal
//	@Description	Get portal data endpoint
//	@Tags			dc4eu
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	PortalReply				"Success"
//	@Failure		400	{object}	helpers.ErrorResponse	"Bad Request"
//	@Param			req	body		PortalRequest			true	" "
//	@Router			/portal [post]
func (c *Client) Portal(ctx context.Context, req *PortalRequest) (*PortalReply, error) {
	ctx, span := c.tp.Start(ctx, "apiv1:Portal")
	defer span.End()

	if err := helpers.Check(ctx, c.cfg, req, c.log); err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	query := &db.PortalQuery{
		AuthenticSource:         req.AuthenticSource,
		AuthenticSourcePersonID: req.AuthenticSourcePersonID,
		ValidTo:                 req.ValidTo,
		ValidFrom:               req.ValidFrom,
	}
	portalData, err := c.db.VCDatastoreColl.PortalData(ctx, query)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	reply := &PortalReply{
		Data: portalData,
	}
	return reply, nil
}
