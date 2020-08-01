package database

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client represents a database connection.
type Client interface {
	Collection(c string) Collection
}

// Collection represents a MongoDB collection.
type Collection interface {
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) SingleResult
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (interface{}, error)
}

// SingleResult represents a single result returned by an operation.
type SingleResult interface {
	Decode(v interface{}) error
}

// Database is the connection to the database.
type Database struct {
	mongoClient *mongo.Database
}

type mongoCollection struct {
	coll *mongo.Collection
}

type singleResult struct {
	sr *mongo.SingleResult
}

// New initializes a database connection using an URI.
func New(ctx context.Context, uri, database string) (*Database, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to database")
	}

	return &Database{
		mongoClient: client.Database(database),
	}, nil
}

// Collection gets a handle for a collection with the given name.
func (d Database) Collection(name string) Collection {
	collection := d.mongoClient.Collection(name)
	return &mongoCollection{collection}
}

// FindOne returns a single document.
func (c mongoCollection) FindOne(ctx context.Context, filter interface{},
	opts ...*options.FindOneOptions) SingleResult {

	sr := c.coll.FindOne(ctx, filter, opts...)
	return &singleResult{sr}
}

func (c mongoCollection) InsertOne(ctx context.Context, document interface{},
	opts ...*options.InsertOneOptions) (interface{}, error) {

	res, err := c.coll.InsertOne(ctx, document, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert document")
	}
	return res.InsertedID, nil
}

// Decode will unmarshal the document represented by this SingleResult into v.
func (sr singleResult) Decode(v interface{}) error {
	return sr.sr.Decode(v)
}
