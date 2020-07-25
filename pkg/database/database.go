package database

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database is the connection to the database.
type Database struct {
	mongoClient *mongo.Database
}

// NewF initializes a database connection using an URI.
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
