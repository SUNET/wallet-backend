package apiv1

import (
	"context"
	"vc/pkg/model"

	"github.com/brianvoe/gofakeit/v6"
)

func (c *Client) mockOne(ctx context.Context, authenticSource, documentType string) (*model.Upload, error) {
	person := gofakeit.Person()

	meta := &model.MetaData{
		AuthenticSource:         authenticSource,
		AuthenticSourcePersonID: gofakeit.UUID(),
		DocumentVersion:         1,
		DocumentType:            documentType,
		DocumentID:              gofakeit.UUID(),
		FirstName:               person.FirstName,
		LastName:                person.LastName,
		DateOfBirth:             gofakeit.Date().String(),
		UID:                     gofakeit.UUID(),
		RevocationID:            gofakeit.UUID(),
		CollectID:               gofakeit.UUID(),
		MemberState:             "SE",
		ValidFrom:               gofakeit.Date().String(),
		ValidTo:                 gofakeit.Date().String(),
	}

	attestation := &model.Attestation{
		Version:          1,
		Type:             documentType,
		DescriptionShort: "a short description",
		DescriptionLong:  "a longer description",
	}

	identity := &model.Identity{
		Version:             "1",
		FamilyName:          gofakeit.LastName(),
		GivenName:           gofakeit.FirstName(),
		BirthDate:           gofakeit.Date().String(),
		UID:                 gofakeit.UUID(),
		FamilyNameAtBirth:   gofakeit.LastName(),
		GivenNameAtBirth:    gofakeit.FirstName(),
		BirthPlace:          gofakeit.City(),
		Gender:              gofakeit.RandomString([]string{"M", "F", "X"}),
		AgeOver18:           "",
		AgeOverNN:           "",
		AgeInYears:          "",
		AgeBirthYear:        "",
		BirthCountry:        "",
		BirthState:          "",
		BirthCity:           "",
		ResidentAddress:     "",
		ResidentCountry:     "",
		ResidentState:       "",
		ResidentCity:        "",
		ResidentPostalCode:  "",
		ResidentStreet:      "",
		ResidentHouseNumber: "",
		Nationality:         "",
	}

	mockUpload := &model.Upload{
		Meta:        meta,
		Attestation: attestation,
		Identity:    identity,
	}

	switch documentType {
	case "PDA1":
		mockUpload.DocumentData = c.PDA1.random(ctx, meta)
	case "EHIC":
		mockUpload.DocumentData = c.EHIC.random(ctx, meta)
	default:
		return nil, model.ErrNoKnownDocumentType
	}

	return mockUpload, nil
}
