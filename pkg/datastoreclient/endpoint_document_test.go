package datastoreclient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"vc/pkg/ehic"
	"vc/pkg/model"

	"github.com/stretchr/testify/assert"
	"gotest.tools/v3/golden"
)

func mockHappyHttServer(t *testing.T, serverReply []byte) *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/document", func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.URL.Path, "/api/v1/document")
		assert.Equal(t, req.Method, http.MethodPost)
		rw.Write(golden.Get(t, "documentGetReplyOK.golden"))
	})

	mux.HandleFunc("/api/v1/document/list", func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.URL.Path, "/api/v1/document/list")
		assert.Equal(t, req.Method, http.MethodPost)
		rw.Write(golden.Get(t, "documentListReplyOK.golden"))
	})

	server := httptest.NewServer(mux)
	return server
}

func mockClient(ctx context.Context, t *testing.T, url string) *Client {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	client, err := New(&Config{
		URL: url,
	})
	if err != nil {
		t.FailNow()
		return nil
	}

	return client
}

func TestGet(t *testing.T) {
	tts := []struct {
		name                 string
		query                *DocumentGetQuery
		expected             *model.Document
		expectedDocumentData *ehic.Document
	}{
		{
			name: "success",
			query: &DocumentGetQuery{
				AuthenticSource: "test_authentic_source",
				DocumentType:    "test_document_type",
				DocumentID:      "test_document_id",
			},
			expected: &model.Document{
				Meta: &model.MetaData{
					AuthenticSource: "SUNET",
					DocumentVersion: "1.0.0",
					DocumentType:    "EHIC",
					DocumentID:      "test_document_id",
					RealData:        false,
					Collect: &model.Collect{
						ID:         "test_collect_id",
						ValidUntil: 1731767173,
					},
					Revocation: &model.Revocation{
						ID:      "test_revocation_id",
						Revoked: false,
						Reference: model.RevocationReference{
							AuthenticSource: "SUNET",
							DocumentType:    "EHIC",
							DocumentID:      "test_document_id",
						},
						RevokedAt: 0,
						Reason:    "",
					},
					CredentialValidFrom:       695706629,
					CredentialValidTo:         -1730367911,
					DocumentDataValidationRef: "",
				},
			},
			expectedDocumentData: &ehic.Document{
				Subject: ehic.Subject{
					Forename:    "test_forename",
					FamilyName:  "test_family_name",
					DateOfBirth: "1986-02-23",
				},
				SocialSecurityPin: "1234",
				PeriodEntitlement: ehic.PeriodEntitlement{
					StartingDate: "1970-01-01",
					EndingDate:   "2038-01-19",
				},
				DocumentID: "test_document_id",
				CompetentInstitution: ehic.CompetentInstitution{
					InstitutionID:      "SE:1234",
					InstitutionName:    "Myndigheten",
					InstitutionCountry: "SE",
				},
			},
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			serverReply := golden.Get(t, "documentGetReplyOK.golden")

			httpServer := mockHappyHttServer(t, serverReply)
			defer httpServer.Close()

			client := mockClient(ctx, t, httpServer.URL)
			got, _, err := client.Document.Get(ctx, tt.query)
			assert.NoError(t, err)

			documentData, err := tt.expectedDocumentData.Marshal()
			assert.NoError(t, err)
			tt.expected.DocumentData = documentData

			assert.Equal(t, tt.expected, got)
		})
	}
}
