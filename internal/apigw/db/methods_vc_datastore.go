package db

import (
	"context"
	"vc/pkg/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// VCDatastoreColl is the generic collection
type VCDatastoreColl struct {
	Service *Service
	Coll    *mongo.Collection
}

func newVCDatastoreColl(ctx context.Context, service *Service, coll *mongo.Collection) (*VCDatastoreColl, error) {
	c := &VCDatastoreColl{
		Service: service,
		Coll:    coll,
	}

	if err := c.createIndex(ctx); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *VCDatastoreColl) createIndex(ctx context.Context) error {
	ctx, span := c.Service.tp.Start(ctx, "db:vc:datastore:createIndex")
	defer span.End()

	indexModel := mongo.IndexModel{
		Keys: bson.M{
			"document_id": 1,
		},
	}
	_, err := c.Coll.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return err
	}
	return nil
}

// Save saves one document to the generic collection
func (c *VCDatastoreColl) Save(ctx context.Context, doc *model.Upload) error {
	_, err := c.Coll.InsertOne(ctx, doc)
	return err
}

// IDMappingQuery is the query to get authentic source person id
type IDMappingQuery struct {
	AuthenticSource string
	Identity        *model.Identity
}

// IDMapping return authentic source person id if any
func (c *VCDatastoreColl) IDMapping(ctx context.Context, query *IDMappingQuery) (string, error) {
	filter := bson.M{
		"meta.authentic_source": bson.M{"$eq": query.AuthenticSource},
		"identity.version":      bson.M{"$eq": query.Identity.Version},
		"identity.family_name":  bson.M{"$eq": query.Identity.FamilyName},
		"identity.given_name":   bson.M{"$eq": query.Identity.GivenName},
		"identity.uid":          bson.M{"$eq": query.Identity.UID},
	}
	opts := options.FindOne().SetProjection(bson.M{
		"meta.authentic_source_person_id": 1,
	})
	res := &model.Upload{}
	if err := c.Coll.FindOne(ctx, filter, opts).Decode(&res); err != nil {
		return "", err
	}

	return res.Meta.AuthenticSourcePersonID, nil
}

// GetDocument return matching document if any, or error
func (c *VCDatastoreColl) GetDocument(ctx context.Context, attr *model.MetaData) (*model.Upload, error) {
	filter := bson.M{
		"meta.authentic_source": bson.M{"$eq": attr.AuthenticSource},
		"meta.document_type":    bson.M{"$eq": attr.DocumentType},
		"meta.document_id":      bson.M{"$eq": attr.DocumentID},
	}
	opt := options.FindOne().SetProjection(bson.M{
		"meta.revocation_id": 1,
		"meta.collect_id":    1,
		"meta.document_data": 1,
	})

	res := &model.Upload{}
	if err := c.Coll.FindOne(ctx, filter, opt).Decode(res); err != nil {
		return nil, err
	}
	return res, nil
}

// GetQR return matching document and return its QR code, else error
func (c *VCDatastoreColl) GetQR(ctx context.Context, attr *model.MetaData) (*model.QR, error) {
	filter := bson.M{
		"meta.authentic_source": bson.M{"$eq": attr.AuthenticSource},
		"meta.document_type":    bson.M{"$eq": attr.DocumentType},
		"meta.document_id":      bson.M{"$eq": attr.DocumentID},
	}
	opt := options.FindOne().SetProjection(bson.M{
		"qr": 1,
	})

	res := &model.Upload{}
	if err := c.Coll.FindOne(ctx, filter, opt).Decode(res); err != nil {
		return nil, err
	}
	return res.QR, nil
}

// GetDocumentAttestationQuery is the query to get document attestation
type GetDocumentAttestationQuery struct {
	Identity *model.Identity
	Meta     *model.MetaData
}

// GetDocumentAttestation return matching document if any, or error
func (c *VCDatastoreColl) GetDocumentAttestation(ctx context.Context, query *GetDocumentAttestationQuery) (*model.Upload, error) {
	filter := bson.M{
		"meta.authentic_source": bson.M{"$eq": query.Meta.AuthenticSource},
		"meta.collect_id":       bson.M{"$eq": query.Meta.CollectID},
		"meta.document_type":    bson.M{"$eq": query.Meta.DocumentType},
		"meta.document_id":      bson.M{"$eq": query.Meta.DocumentID},
		"meta.uid":              bson.M{"$eq": query.Meta.UID},
		"meta.date_of_birth":    bson.M{"$eq": query.Meta.DateOfBirth},
		"meta.last_name":        bson.M{"$eq": query.Meta.LastName},
		"meta.first_name":       bson.M{"$eq": query.Meta.FirstName},
	}

	opt := options.FindOne().SetProjection(bson.M{
		"meta.document_type": 1,
		"meta.document_id":   1,
		"meta.revocation_id": 1,
		"meta.collect_id":    1,
		"meta.document_data": 1,
	})

	res := &model.Upload{}
	if err := c.Coll.FindOne(ctx, filter, opt).Decode(res); err != nil {
		return nil, err
	}
	return res, nil
}

// PortalQuery is the query to get portal data
type PortalQuery struct {
	AuthenticSource         string
	AuthenticSourcePersonID string
	ValidTo                 string
	ValidFrom               string
}

// PortalData returns a list of portal data
func (c *VCDatastoreColl) PortalData(ctx context.Context, query *PortalQuery) ([]model.Upload, error) {
	filter := bson.M{}

	if query.ValidFrom != "" {
		filter["meta.valid_from"] = bson.M{"$eq": query.ValidFrom}
	} else if query.ValidTo != "" {
		filter["meta.valid_to"] = bson.M{"$eq": query.ValidTo}
	}
	c.Service.log.Info("filter", "valid_to", query.ValidTo)

	filter["meta.authentic_source"] = bson.M{"$eq": query.AuthenticSource}
	filter["meta.authentic_source_person_id"] = bson.M{"$eq": query.AuthenticSourcePersonID}

	cursor, err := c.Coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	res := []model.Upload{}
	if err := cursor.All(ctx, &res); err != nil {
		return nil, err
	}
	return res, nil
}
