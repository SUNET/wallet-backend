package db

import (
	"context"
	"vc/pkg/model"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/otel/codes"
)

// VCDatastoreColl is the collection of datastore
type VCDatastoreColl struct {
	service    *Service
	coll       *mongo.Collection
	metricSave prometheus.Counter
}

func (c *VCDatastoreColl) createIndex(ctx context.Context) error {
	ctx, span := c.service.tp.Start(ctx, "db:vc:datastore:createIndex")
	defer span.End()

	indexDocumentIDInAuthenticSourceUniq := mongo.IndexModel{
		Keys: bson.D{
			primitive.E{Key: "meta.document_id", Value: 1},
			primitive.E{Key: "meta.authentic_source", Value: 1},
		},
		Options: options.Index().SetName("document_id_uniq").SetUnique(true),
	}
	_, err := c.coll.Indexes().CreateMany(ctx, []mongo.IndexModel{indexDocumentIDInAuthenticSourceUniq})
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}

// Save saves one document
func (c *VCDatastoreColl) Save(ctx context.Context, doc *model.Upload) error {
	ctx, span := c.service.tp.Start(ctx, "db:vc:datastore:save")
	defer span.End()

	c.metricSave = promauto.NewCounter(prometheus.CounterOpts{
		Name: "persistent_vc_db_save_total",
	})

	
	res, err := c.coll.InsertOne(ctx, doc)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	c.metricSave.Inc()
	c.service.log.Info("saved document", "document_id", doc.Meta.DocumentID, "inserted_id", res.InsertedID)
	return nil
}

// Delete deletes a document
func (c *VCDatastoreColl) Delete(ctx context.Context, doc *model.MetaData) error {
	ctx, span := c.service.tp.Start(ctx, "db:vc:datastore:delete")
	defer span.End()

	filter := bson.M{
		"meta.document_id":      bson.M{"$eq": doc.DocumentID},
		"meta.authentic_source": bson.M{"$eq": doc.AuthenticSource},
	}
	_, err := c.coll.DeleteOne(ctx, filter)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	c.service.log.Info("deleted document", "document_id", doc.DocumentID, "from authentic_source", doc.AuthenticSource)
	return nil

}
