package model

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/skip2/go-qrcode"
)

// CompleteDocument is a generic type for upload
type CompleteDocument struct {
	Meta            *MetaData        `json:"meta,omitempty" bson:"meta" validate:"required"`
	Identities      []Identity       `json:"identities,omitempty" bson:"identities" validate:"required"`
	DocumentDisplay *DocumentDisplay `json:"document_display,omitempty" bson:"document_display" validate:"required"`
	DocumentData    map[string]any   `json:"document_data,omitempty" bson:"document_data" validate:"required"`

	// required: true
	// example: "1.0.0"
	DocumentDataVersion string `json:"document_data_version,omitempty" bson:"document_data_version" validate:"required,semver"`
	QR                  *QR    `json:"qr,omitempty" bson:"qr"`
}

// CompleteDocuments is a array of CompleteDocument
type CompleteDocuments []CompleteDocument

// DocumentList is a generic type for document list
type DocumentList struct {
	Meta            *MetaData        `json:"meta,omitempty" bson:"meta" validate:"required"`
	DocumentDisplay *DocumentDisplay `json:"document_display,omitempty" bson:"document_display"`
	QR              *QR              `json:"qr,omitempty" bson:"qr" validate:"required"`
}

// Document is a generic type for get document
type Document struct {
	Meta         *MetaData `json:"meta,omitempty" bson:"meta" validate:"required"`
	DocumentData any       `json:"document_data" bson:"document_data" validate:"required"`
}

// IDMapping is a generic type for ID mapping
type IDMapping struct {
	AuthenticSourcePersonID string `json:"authentic_source_person_id"`
}

// QRGenerator generates a QR code
func (m *MetaData) QRGenerator(ctx context.Context, issuerBaseURL string, recoveryLevel, size int) (*QR, error) {
	credentialOfferURL, err := url.Parse("https://wallet.dc4eu.eu/cb")
	if err != nil {
		return nil, err
	}
	issuerState := fmt.Sprintf("collect_id=%s&document_type=%s&authentic_source=%s", m.Collect.ID, m.DocumentType, m.AuthenticSource)

	query := credentialOfferURL.Query()

	var credentialConfigurationID string
	switch m.DocumentType {
	case "PDA1":
		credentialConfigurationID = "PDA1Credential"
	case "EHIC":
		credentialConfigurationID = "EHICCredential"
	}

	credentialOffer := &CredentialOffer{
		CredentialIssuer:           issuerBaseURL,
		CredentialConfigurationIDs: []string{credentialConfigurationID},
		Grants: map[string]map[string]string{
			"authorization_code": {
				"issuer_state": issuerState,
			},
		},
	}

	credentialOfferByte, err := credentialOffer.Marshal()
	if err != nil {
		return nil, err
	}

	query.Add("credential_offer", string(credentialOfferByte))
	credentialOfferURL.RawQuery = query.Encode()

	qrPNG, err := qrcode.Encode(credentialOfferURL.String(), qrcode.RecoveryLevel(recoveryLevel), size)
	if err != nil {
		return nil, err
	}

	qrBase64 := base64.StdEncoding.EncodeToString(qrPNG)

	qr := &QR{
		CredentialOfferURL: credentialOfferURL.String(),
		Base64Image:        qrBase64,
	}

	return qr, nil
}

// CSV returns the document as a CSV
func (c *CompleteDocument) CSV() (string, []string, error) {
	if len(c.Identities) == 0 {
		return "", nil, errors.New("no identities found")
	}

	if c.Meta == nil {
		return "", nil, errors.New("no metadata found")
	}

	qr, err := c.Meta.QRGenerator(context.Background(), "https://satosa-test-1.sunet.se", 2, 256)
	if err != nil {
		return "", nil, err
	}
	attributes := []string{
		c.Identities[0].AuthenticSourcePersonID,
		c.Identities[0].GivenName,
		c.Identities[0].FamilyName,
		c.Identities[0].BirthDate,
		c.Identities[0].Schema.Name,
		c.Meta.AuthenticSource,
		c.Meta.Collect.ID,
		c.Meta.DocumentType,
		c.Meta.DocumentID,
		qr.CredentialOfferURL,
	}
	csv := strings.Join(attributes, ",")

	return csv, attributes, nil
}

// CSV return CompleteDocuments as a CSV, string array
func (c *CompleteDocuments) CSV() ([]string, [][]string, error) {
	if len(*c) == 0 {
		return nil, nil, errors.New("no documents found")
	}

	var csvResult []string
	var csvRaw [][]string
	for _, doc := range *c {
		csvRow, raw, err := doc.CSV()
		if err != nil {
			return nil, nil, err
		}
		csvRaw = append(csvRaw, raw)
		csvResult = append(csvResult, fmt.Sprintf("%v\n", csvRow))
	}

	return csvResult, csvRaw, nil
}

// CredentialOffer https://openid.net/specs/openid-4-verifiable-credential-issuance-1_0.html 4.1.1 Credential Offer Parameters
type CredentialOffer struct {
	CredentialIssuer           string                       `json:"credential_issuer"`
	CredentialConfigurationIDs []string                     `json:"credential_configuration_ids"`
	Grants                     map[string]map[string]string `json:"grants"`
}

// Marshal marshals the CredentialOffer
func (c *CredentialOffer) Marshal() ([]byte, error) {
	return json.Marshal(c)
}

// Consent is a generic type for consent
type Consent struct {
	// required: true
	// example: "Using my data for research"
	ConsentTo string `json:"consent_to,omitempty" bson:"consent_to" validate:"required"`

	// required: true
	// example: "sess-123"
	SessionID string `json:"session_id,omitempty" bson:"session_id" validate:"required"`

	// required: true
	// example: 509567558
	// format: int64
	CreatedAt int64 `json:"created_at,omitempty" bson:"created_at" validate:"required"`
}

// Collect is a generic type for collect
type Collect struct {
	// required: false
	// example: 98fe67fc-c03f-11ee-bbee-4345224d414f
	ID string `json:"id,omitempty" bson:"id"`

	// required: false
	// example: 509567558
	// format: int64
	ValidUntil int64 `json:"valid_until,omitempty" bson:"valid_until"`
}

// MetaData is a generic type for metadata
type MetaData struct {
	// required: true
	// example: SUNET
	AuthenticSource string `json:"authentic_source,omitempty" bson:"authentic_source" validate:"required"`

	// required: true
	// example: "1.0.0"
	DocumentVersion string `json:"document_version,omitempty" bson:"document_version" validate:"required,semver"`

	// required: true
	// example: PDA1
	DocumentType string `json:"document_type,omitempty" bson:"document_type" validate:"required,oneof=PDA1 EHIC"`

	// required: true
	// example: 5e7a981c-c03f-11ee-b116-9b12c59362b9
	DocumentID string `json:"document_id,omitempty" bson:"document_id" validate:"required"`

	// RealData is a flag to indicate if the document contains real data
	// required: true
	// example: true
	RealData bool `json:"real_data" bson:"real_data"`

	Collect *Collect `json:"collect,omitempty" bson:"collect"`

	// Revocation is a collection of fields representing a revocation
	Revocation *Revocation `json:"revocation,omitempty" bson:"revocation"`

	// required: false
	// example: 509567558
	// format: int64
	CredentialValidFrom int64 `json:"credential_valid_from,omitempty" bson:"valid_from"`

	// required: false
	// example: 509567558
	// format: int64
	CredentialValidTo int64 `json:"credential_valid_to,omitempty" bson:"valid_to"`

	// required: false
	// example: file://path/to/schema.json or http://example.com/schema.json
	// format: string
	DocumentDataValidationRef string `json:"document_data_validation,omitempty" bson:"document_data_validation"`
}

// RevocationReference refer to a document
type RevocationReference struct {
	AuthenticSource string `json:"authentic_source,omitempty" bson:"authentic_source"`
	DocumentType    string `json:"document_type,omitempty" bson:"document_type"`
	DocumentID      string `json:"document_id,omitempty" bson:"document_id"`
}

// Revocation is a collection of fields representing a revocation
type Revocation struct {

	// ID is the ID of the revocation
	// required: false
	// example: 8dbd2680-c03f-11ee-a21b-034aafe41222
	ID string `json:"id,omitempty" bson:"id"`

	// Revoked is a flag to indicate if the document has been revoked
	// required: false
	// example: false
	Revoked bool `json:"revoked,omitempty" bson:"revoked"`

	Reference RevocationReference `json:"reference" bson:"reference"`

	// RevokedAt is the time the document was revoked or going to be revoked
	// required: false
	// example: 509567558
	// format: int64
	RevokedAt int64 `json:"revoked_at,omitempty" bson:"revoked_at"`

	// Reason is the reason for revocation
	// required: false
	// example: lost or stolen
	Reason string `json:"reason,omitempty" bson:"reason"`
}

// IdentitySchema is a collection of fields representing an identity schema
type IdentitySchema struct {
	// required: true
	// example: "SE"
	Name string `json:"name,omitempty" bson:"name" validate:"required"`

	// required: false
	// example: "1.0.0"
	Version string `json:"version,omitempty" bson:"version" validate:"omitempty,semver"`
}

// Identity identifies a person
type Identity struct {
	// required: true
	// example: 65636cbc-c03f-11ee-8dc4-67135cc9bd8a
	AuthenticSourcePersonID string `json:"authentic_source_person_id,omitempty" bson:"authentic_source_person_id"`

	Schema *IdentitySchema `json:"schema,omitempty" bson:"schema" validate:"required"`

	// required: true
	// example: Svensson
	FamilyName string `json:"family_name,omitempty" bson:"family_name"`

	// required: true
	// example: Magnus
	GivenName string `json:"given_name,omitempty" bson:"given_name"`

	// required: true
	// example: 1970-01-01
	BirthDate string `json:"birth_date,omitempty" bson:"birth_date" validate:"omitempty,datetime=2006-01-02"`

	// required: false
	// example: Karlsson
	FamilyNameAtBirth string `json:"family_name_at_birth,omitempty" bson:"family_name_at_birth"`

	// required: false
	// example: Magnus
	GivenNameAtBirth string `json:"given_name_at_birth,omitempty" bson:"given_name_at_birth"`

	// required: false
	// example: Stockholm
	BirthPlace string `json:"birth_place,omitempty" bson:"birth_place"`

	// required: false
	// example: male
	Gender string `json:"gender,omitempty" bson:"gender"`

	// TODO(masv): full name or just country code?
	// required: false
	// example: sweden
	BirthCountry string `json:"birth_country,omitempty" bson:"birth_country"`

	// required: false
	// example: Stockholm
	BirthState string `json:"birth_state,omitempty" bson:"birth_state"`

	// required: false
	// example: Stockholm
	BirthCity string `json:"birth_city,omitempty" bson:"birth_city"`

	// required: false
	// example: 221b baker street
	ResidentAddress string `json:"resident_address,omitempty" bson:"resident_address"`

	// required: false
	// example: england
	ResidentCountry string `json:"resident_country,omitempty" bson:"resident_country"`

	// required: false
	// example: england
	ResidentState string `json:"resident_state,omitempty" bson:"resident_state"`

	// required: false
	// example: london
	ResidentCity string `json:"resident_city,omitempty" bson:"resident_city"`

	// required: false
	// example: W1U 6SG
	ResidentPostalCode string `json:"resident_postal_code,omitempty" bson:"resident_postal_code"`

	// required: false
	// example: baker street
	ResidentStreet string `json:"resident_street,omitempty" bson:"resident_street"`

	// required: false
	// example: 221b
	ResidentHouseNumber string `json:"resident_house_number,omitempty" bson:"resident_house_number"`

	// required: false
	// example: swedish
	Nationality string `json:"nationality,omitempty" bson:"nationality"`
}

// DocumentDisplay is a collection of fields representing display of document
type DocumentDisplay struct {
	// required: true
	// example: "1.0.0"
	Version string `json:"version,omitempty" bson:"version" validate:"required,semver"`

	// required: true
	// example: secure
	Type string `json:"type,omitempty" bson:"type" validate:"required"`

	// DescriptionStructured is a map of structured descriptions
	// required: true
	// example: {"en": "European Health Insurance Card", "sv": "Europeiskt sjukförsäkringskortet"}
	DescriptionStructured map[string]any `json:"description_structured,omitempty" bson:"description_structured" validate:"required"`
}

// QR is a collection of fields representing a QR code
type QR struct {
	// required: true
	// example: "ZWFzdGVyIGVnZyE="
	Base64Image string `json:"base64_image,omitempty" bson:"base64_image" validate:"required"`

	// required: true
	CredentialOfferURL string `json:"credential_offer,omitempty" bson:"credential_offer"`
}

// SearchDocumentsReply the reply from search documents
type SearchDocumentsReply struct {
	Documents      []*CompleteDocument `json:"documents"`
	HasMoreResults bool                `json:"has_more_results"`
}

// SearchDocumentsRequest the request to search for documents
type SearchDocumentsRequest struct {
	AuthenticSource string `json:"authentic_source,omitempty"`
	DocumentType    string `json:"document_type,omitempty"`
	DocumentID      string `json:"document_id,omitempty"`
	CollectID       string `json:"collect_id,omitempty"`

	AuthenticSourcePersonID string `json:"authentic_source_person_id,omitempty"`
	FamilyName              string `json:"family_name,omitempty"`
	GivenName               string `json:"given_name,omitempty"`
	BirthDate               string `json:"birth_date,omitempty"`

	Limit      int64          `json:"limit,omitempty"`
	Fields     []string       `json:"fields,omitempty"`
	SortFields map[string]int `json:"sort_fields,omitempty"`
}
