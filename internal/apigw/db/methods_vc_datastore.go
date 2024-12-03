package db

import (
	"context"
	"vc/pkg/helpers"
	"vc/pkg/logger"
	"vc/pkg/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/otel/codes"
)

// VCDatastoreColl is the generic collection
type VCDatastoreColl struct {
	Service *Service
	Coll    *mongo.Collection
	log     *logger.Log
}

func (c *VCDatastoreColl) createIndex(ctx context.Context) error {
	ctx, span := c.Service.tracer.Start(ctx, "db:vc:datastore:createIndex")
	defer span.End()

	indexDocumentIDInAuthenticSourceUniq := mongo.IndexModel{
		Keys: bson.D{
			primitive.E{Key: "meta.document_id", Value: 1},
			primitive.E{Key: "meta.authentic_source", Value: 1},
			primitive.E{Key: "meta.document_type", Value: 1},
		},
		Options: options.Index().SetName("document_unique_within_namespace").SetUnique(true),
	}
	_, err := c.Coll.Indexes().CreateMany(ctx, []mongo.IndexModel{indexDocumentIDInAuthenticSourceUniq})
	if err != nil {
		return err
	}
	return nil
}

// Save saves one document to the generic collection
func (c *VCDatastoreColl) Save(ctx context.Context, doc *model.CompleteDocument) error {
	ctx, span := c.Service.tracer.Start(ctx, "db:vc:datastore:save")
	defer span.End()

	_, err := c.Coll.InsertOne(ctx, doc)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	return nil
}

// IDMappingQuery is the query to get authentic source person id
type IDMappingQuery struct {
	AuthenticSource string
	Identity        *model.Identity
}

// IDMapping return authentic source person id if any
func (c *VCDatastoreColl) IDMapping(ctx context.Context, query *IDMappingQuery) (string, error) {
	filter := bson.M{
		"meta.authentic_source":     bson.M{"$eq": query.AuthenticSource},
		"identities.schema.version": bson.M{"$eq": query.Identity.Schema.Version},
		"identities.family_name":    bson.M{"$eq": query.Identity.FamilyName},
		"identities.given_name":     bson.M{"$eq": query.Identity.GivenName},
		"identities.birth_date":     bson.M{"$eq": query.Identity.BirthDate},
	}

	opts := options.FindOne().SetProjection(bson.M{
		"identities.authentic_source_person_id": 1,
	})
	res := &model.CompleteDocument{}
	if err := c.Coll.FindOne(ctx, filter, opts).Decode(&res); err != nil {
		return "", err
	}
	if res.Identities == nil || len(res.Identities) == 0 {
		return "", helpers.ErrNoIdentityFound
	}

	return res.Identities[0].AuthenticSourcePersonID, nil
}

// AddDocumentIdentityQuery is the query to add document identity
type AddDocumentIdentityQuery struct {
	AuthenticSource string            `json:"authentic_source" bson:"authentic_source"`
	DocumentType    string            `json:"document_type" bson:"document_type"`
	DocumentID      string            `json:"document_id" bson:"document_id"`
	Identities      []*model.Identity `json:"identities" bson:"identities"`
}

// AddDocumentIdentity adds document identity
func (c *VCDatastoreColl) AddDocumentIdentity(ctx context.Context, query *AddDocumentIdentityQuery) error {
	filter := bson.M{
		"meta.authentic_source": bson.M{"$eq": query.AuthenticSource},
		"meta.document_id":      bson.M{"$eq": query.DocumentID},
		"meta.document_type":    bson.M{"$eq": query.DocumentType},
	}

	// This needs to make sure no duplicate authentic_source_person_id is added in the future
	update := bson.M{"$addToSet": bson.M{"identities": bson.M{"$each": query.Identities}}}

	result, err := c.Coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.ModifiedCount == 0 {
		return helpers.ErrNoDocumentFound
	}

	return nil
}

// DeleteDocumentIdentityQuery is the query to delete identity in document
type DeleteDocumentIdentityQuery struct {
	AuthenticSource         string `json:"authentic_source" bson:"authentic_source"`
	DocumentType            string `json:"document_type" bson:"document_type"`
	DocumentID              string `json:"document_id" bson:"document_id"`
	AuthenticSourcePersonID string `json:"authentic_source_person_id" bson:"authentic_source_person_id"`
}

// DeleteDocumentIdentity deletes identity in document
func (c *VCDatastoreColl) DeleteDocumentIdentity(ctx context.Context, query *DeleteDocumentIdentityQuery) error {
	filter := bson.M{
		"meta.authentic_source": bson.M{"$eq": query.AuthenticSource},
		"meta.document_id":      bson.M{"$eq": query.DocumentID},
		"meta.document_type":    bson.M{"$eq": query.DocumentType},
	}

	update := bson.M{"$pull": bson.M{"identities": bson.M{"authentic_source_person_id": query.AuthenticSourcePersonID}}}
	_, err := c.Coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

// Delete deletes a document
func (c *VCDatastoreColl) Delete(ctx context.Context, doc *model.MetaData) error {
	ctx, span := c.Service.tracer.Start(ctx, "db:vc:datastore:delete")
	defer span.End()

	filter := bson.M{
		"meta.document_id":      bson.M{"$eq": doc.DocumentID},
		"meta.authentic_source": bson.M{"$eq": doc.AuthenticSource},
		"meta.document_type":    bson.M{"$eq": doc.DocumentType},
	}
	_, err := c.Coll.DeleteOne(ctx, filter)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	//c.service.log.Info("deleted document", "document_id", doc.DocumentID, "from authentic_source", doc.AuthenticSource)
	return nil

}

// GetDocumentForCredential is the query to get document attestation
type GetDocumentForCredential struct {
	Meta     *model.MetaData
	Identity *model.Identity
}

// GetDocumentForCredential return matching document if any, or error
func (c *VCDatastoreColl) GetDocumentForCredential(ctx context.Context, query *GetDocumentForCredential) (*model.Document, error) {
	filter := bson.M{
		"meta.authentic_source":                 bson.M{"$eq": query.Meta.AuthenticSource},
		"meta.document_type":                    bson.M{"$eq": query.Meta.DocumentType},
		"identities.authentic_source_person_id": bson.M{"$eq": query.Identity.AuthenticSourcePersonID},
	}
	opt := options.FindOne().SetProjection(bson.M{
		"meta":          1,
		"document_data": 1,
	})

	res := &model.CompleteDocument{}
	if err := c.Coll.FindOne(ctx, filter, opt).Decode(res); err != nil {
		return nil, err
	}

	reply := &model.Document{
		Meta:         res.Meta,
		DocumentData: res.DocumentData,
	}
	return reply, nil
}

// GetDocumentQuery is the query to get document attestation
type GetDocumentQuery struct {
	Meta     *model.MetaData
	Identity *model.Identity
}

// GetDocument return matching document if any, or error
func (c *VCDatastoreColl) GetDocument(ctx context.Context, query *GetDocumentQuery) (*model.Document, error) {
	filter := bson.M{
		"meta.authentic_source": bson.M{"$eq": query.Meta.AuthenticSource},
		"meta.document_type":    bson.M{"$eq": query.Meta.DocumentType},
		"meta.document_id":      bson.M{"$eq": query.Meta.DocumentID},
		//"identities.authentic_source_person_id": bson.M{"$eq": query.Identity.AuthenticSourcePersonID},
	}
	opt := options.FindOne().SetProjection(bson.M{
		"meta":          1,
		"document_data": 1,
	})

	res := &model.CompleteDocument{}
	if err := c.Coll.FindOne(ctx, filter, opt).Decode(res); err != nil {
		return nil, err
	}

	reply := &model.Document{
		Meta:         res.Meta,
		DocumentData: res.DocumentData,
	}
	return reply, nil
}

// DocumentListQuery is the query to get document list
type DocumentListQuery struct {
	AuthenticSource string          `json:"authentic_source" bson:"authentic_source"`
	Identity        *model.Identity `json:"identity" bson:"identity" validate:"required"`
	DocumentType    string          `json:"document_type" bson:"document_type"`
	ValidFrom       int64           `json:"valid_from" bson:"valid_from"`
	ValidTo         int64           `json:"valid_to" bson:"valid_to"`
}

// DocumentList return matching documents if any, or error
func (c *VCDatastoreColl) DocumentList(ctx context.Context, query *DocumentListQuery) ([]*model.DocumentList, error) {
	if err := helpers.Check(ctx, c.Service.cfg, query, c.Service.log); err != nil {
		return nil, err
	}

	filter := bson.M{
		"identities.schema.name": bson.M{"$eq": query.Identity.Schema.Name},
	}

	if query.AuthenticSource != "" {
		filter["meta.authentic_source"] = bson.M{"$eq": query.AuthenticSource}
	}

	if query.DocumentType != "" {
		filter["meta.document_type"] = bson.M{"$eq": query.DocumentType}
	}

	if query.Identity.AuthenticSourcePersonID != "" {
		filter["identities.authentic_source_person_id"] = bson.M{"$eq": query.Identity.AuthenticSourcePersonID}
	} else {
		filter["identities.family_name"] = bson.M{"$eq": query.Identity.FamilyName}
		filter["identities.given_name"] = bson.M{"$eq": query.Identity.GivenName}
		filter["identities.birth_date"] = bson.M{"$eq": query.Identity.BirthDate}
	}

	cursor, err := c.Coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	res := []*model.DocumentList{}
	if err := cursor.All(ctx, &res); err != nil {
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

	res := &model.CompleteDocument{}
	if err := c.Coll.FindOne(ctx, filter, opt).Decode(res); err != nil {
		return nil, err
	}
	return res.QR, nil
}

// GetDocumentCollectIDQuery is the query to get document attestation
type GetDocumentCollectIDQuery struct {
	Identity *model.Identity
	Meta     *model.MetaData
}

// GetDocumentCollectID return matching document if any, or error
func (c *VCDatastoreColl) GetDocumentCollectID(ctx context.Context, query *GetDocumentCollectIDQuery) (*model.Document, error) {
	filter := bson.M{
		"meta.authentic_source":  bson.M{"$eq": query.Meta.AuthenticSource},
		"meta.collect.id":        bson.M{"$eq": query.Meta.Collect.ID},
		"meta.document_type":     bson.M{"$eq": query.Meta.DocumentType},
		"identities.schema.name": bson.M{"$eq": query.Identity.Schema.Name},
	}

	if query.Identity.AuthenticSourcePersonID != "" {
		filter["identities.authentic_source_person_id"] = bson.M{"$eq": query.Identity.AuthenticSourcePersonID}
	} else {
		filter["identities.family_name"] = bson.M{"$eq": query.Identity.FamilyName}
		filter["identities.given_name"] = bson.M{"$eq": query.Identity.GivenName}
		filter["identities.birth_date"] = bson.M{"$eq": query.Identity.BirthDate}
	}

	opts := options.FindOne().SetProjection(bson.M{
		"meta":          1,
		"document_data": 1,
	})

	res := &model.CompleteDocument{}
	if err := c.Coll.FindOne(ctx, filter, opts).Decode(res); err != nil {
		return nil, err
	}

	reply := &model.Document{
		Meta:         res.Meta,
		DocumentData: res.DocumentData,
	}
	return reply, nil
}

// GetByRevocationID gets one document by meta.revocation.id and meta.authentic_source
func (c *VCDatastoreColl) GetByRevocationID(ctx context.Context, q *model.MetaData) (*model.CompleteDocument, error) {
	filter := bson.M{
		"meta.authentic_source": bson.M{"$eq": q.AuthenticSource},
		"meta.document_type":    bson.M{"$eq": q.DocumentType},
		"meta.revocation.id":    bson.M{"$eq": q.Revocation.ID},
	}
	res := &model.CompleteDocument{}
	if err := c.Coll.FindOne(ctx, filter).Decode(res); err != nil {
		return nil, err
	}
	return res, nil
}

// Replace replaces one document
func (c *VCDatastoreColl) Replace(ctx context.Context, doc *model.CompleteDocument) error {
	ctx, span := c.Service.tracer.Start(ctx, "db:vc:datastore:replace")
	defer span.End()

	filter := bson.M{
		"meta.document_id":      bson.M{"$eq": doc.Meta.DocumentID},
		"meta.authentic_source": bson.M{"$eq": doc.Meta.AuthenticSource},
	}

	_, err := c.Coll.ReplaceOne(ctx, filter, doc)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	c.log.Info("updated document", "document_id", doc.Meta.DocumentID)
	return nil
}

type SearchDocumentsQuery struct {
	AuthenticSource string `json:"authentic_source" bson:"authentic_source"`
}

func (c *VCDatastoreColl) SearchDocuments(ctx context.Context, query *SearchDocumentsQuery) ([]*model.CompleteDocument, bool, error) {
	if err := helpers.Check(ctx, c.Service.cfg, query, c.Service.log); err != nil {
		return nil, false, err
	}

	filter := bson.M{}

	if query.AuthenticSource != "" {
		filter["meta.authentic_source"] = bson.M{"$eq": query.AuthenticSource}
	}
	//TODO: add more filters

	findOptions := options.Find()
	var limit int64 = 3
	// One more than wanted to see if there are more results i db
	findOptions.SetLimit(limit + 1)

	cursor, err := c.Coll.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, false, err
	}

	res := []*model.CompleteDocument{}
	if err := cursor.All(ctx, &res); err != nil {
		return nil, false, err
	}

	hasMore := len(res) > int(limit)
	if hasMore {
		// Remove the last entry from the result
		res = res[:limit]
	}

	return res, hasMore, nil
}
