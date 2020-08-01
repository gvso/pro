package mock

import (
	"context"
	"reflect"

	"github.com/gvso/pro/pkg/database"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CollectionFunc is the function signature for Collection.
type CollectionFunc func(name string) database.Collection

// FindOneFunc is the function signature for FindOne.
type FindOneFunc func(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) database.SingleResult

// Client is a mock implementation of database.Client.
type Client struct {
	CollectionFn      CollectionFunc
	CollectionInvoked bool
}

// Collection is a mock implementation of database.Collection.
type Collection struct {
	FindOneFn      FindOneFunc
	FindOneInvoked bool

	InsertOneFn      func(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (interface{}, error)
	InsertOneInvoked bool
}

// SingleResult is a mock implementation of database.SingleResult.
type SingleResult struct {
	DecodeFn      func(v interface{}) error
	DecodeInvoked bool
}

// Collection returns a Collection mock and marks the function as invoked.
func (c Client) Collection(name string) database.Collection {
	c.CollectionInvoked = true
	return c.CollectionFn(name)
}

// FindOne invokes the mock implementation and marks the function as invoked.
func (c Collection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) database.SingleResult {
	c.FindOneInvoked = true
	return c.FindOneFn(ctx, filter, opts...)
}

// InsertOne invokes the mock implementation and marks the function as invoked.
func (c Collection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (interface{}, error) {
	c.InsertOneInvoked = true
	return c.InsertOneFn(ctx, document, opts...)
}

// Decode will unmarshal the document represented by this SingleResult into v.
func (sr SingleResult) Decode(v interface{}) error {
	sr.DecodeInvoked = true
	return sr.DecodeFn(v)
}

// ClientMock returns a basic mock for the entire database.
func ClientMock(srValue interface{}, srDecodeErr error, insertOneID interface{}) Client {
	dbClient := Client{}
	dbClient.CollectionFn = func(name string) database.Collection {
		collectionMock := Collection{}

		collectionMock.FindOneFn = func(ctx context.Context, filter interface{},
			opts ...*options.FindOneOptions) database.SingleResult {

			singleResultMock := SingleResult{}
			singleResultMock.DecodeFn = func(v interface{}) error {
				d := reflect.ValueOf(srValue)
				if d.IsValid() {
					s := reflect.ValueOf(v).Elem()
					s.Set(d)
				}

				return srDecodeErr
			}

			return singleResultMock
		}

		collectionMock.InsertOneFn = func(ctx context.Context, document interface{},
			opts ...*options.InsertOneOptions) (interface{}, error) {

			return insertOneID, nil
		}

		return collectionMock
	}

	return dbClient
}
