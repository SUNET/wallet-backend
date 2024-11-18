package apiv1

import (
	"context"
	"time"
	"vc/pkg/helpers"
	"vc/pkg/model"

	"github.com/brianvoe/gofakeit/v6"
)

// MockInputData is the input data for the mock function
type MockInputData struct {
	DocumentType            string `json:"document_type"`
	DocumentID              string `json:"document_id"`
	AuthenticSource         string `json:"authentic_source"`
	AuthenticSourcePersonID string `json:"authentic_source_person_id"`
	GivenName               string `json:"given_name"`
	FamilyName              string `json:"family_name"`
	BirthDate               string `json:"birth_date"`
	CollectID               string `json:"collect_id"`
	IdentitySchemaName      string `json:"identity_schema_name"`
}

type uploadMock struct {
	Meta                *model.MetaData        `json:"meta" validate:"required"`
	Identities          []model.Identity       `json:"identities,omitempty" validate:"required,dive"`
	DocumentDisplay     *model.DocumentDisplay `json:"document_display,omitempty" validate:"required"`
	DocumentData        map[string]any         `json:"document_data" validate:"required"`
	DocumentDataVersion string                 `json:"document_data_version,omitempty" validate:"required,semver"`
}

func (c *Client) mockOne(ctx context.Context, data MockInputData) (*uploadMock, error) {
	c.log.Debug("mockOne")
	person := gofakeit.Person()

	if data.AuthenticSourcePersonID == "" {
		data.AuthenticSourcePersonID = gofakeit.UUID()
	}

	if data.GivenName == "" {
		data.GivenName = person.FirstName
	}

	if data.FamilyName == "" {
		data.FamilyName = person.LastName
	}

	if data.BirthDate == "" {
		data.BirthDate = gofakeit.Date().Format("2006-01-02")
	}

	if data.CollectID == "" {
		data.CollectID = gofakeit.UUID()
	}

	if data.DocumentID == "" {
		data.DocumentID = gofakeit.UUID()
	}

	if data.IdentitySchemaName == "" {
		data.IdentitySchemaName = "SE"
	}

	if data.AuthenticSource == "" {
		data.AuthenticSource = gofakeit.Company()
	}

	meta := &model.MetaData{
		AuthenticSource: data.AuthenticSource,
		DocumentType:    data.DocumentType,
		DocumentID:      data.DocumentID,
		DocumentVersion: "1.0.0",
		RealData:        false,
		Collect: &model.Collect{
			ID:         data.CollectID,
			ValidUntil: time.Now().Add(10 * 24 * time.Hour).Unix(),
		},
		CredentialValidFrom: gofakeit.Date().Unix(),
		CredentialValidTo:   gofakeit.Date().Unix(),
		Revocation: &model.Revocation{
			ID:      gofakeit.UUID(),
			Revoked: false,
			Reference: model.RevocationReference{
				AuthenticSource: data.AuthenticSource,
				DocumentType:    data.DocumentType,
				DocumentID:      data.DocumentID,
			},
		},
	}

	identities := []model.Identity{
		{
			AuthenticSourcePersonID: data.AuthenticSourcePersonID,
			Schema: &model.IdentitySchema{
				Name:    data.IdentitySchemaName,
				Version: "1.0.0",
			},
			FamilyName: data.FamilyName,
			GivenName:  data.GivenName,
			BirthDate:  data.BirthDate,
		},
	}

	documentDisplay := &model.DocumentDisplay{
		Version: "1.0.0",
		Type:    data.DocumentType,
		DescriptionStructured: map[string]any{
			"en": "issuer",
			"sv": "utfärdare",
		},
	}

	mockUpload := &uploadMock{
		Meta:            meta,
		Identities:      identities,
		DocumentDisplay: documentDisplay,
	}

	switch data.DocumentType {
	case "PDA1":
		mockUpload.DocumentData = c.PDA1.random(ctx, person)
		mockUpload.Meta.DocumentDataValidationRef = "file://../../standards/schema_pda1.json"
	case "EHIC":
		mockUpload.DocumentData = c.EHIC.random(ctx, person)
		mockUpload.Meta.DocumentDataValidationRef = "file://../../standards/schema_ehic.json"
	default:
		return nil, helpers.ErrNoKnownDocumentType
	}

	mockUpload.DocumentDataVersion = "1.0.0"

	c.log.Debug("2")
	if err := helpers.CheckSimple(mockUpload); err != nil {
		c.log.Debug("mockOne", "error", err)
		return nil, err
	}

	c.log.Debug("3")

	return mockUpload, nil
}
