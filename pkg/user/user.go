package user

import (
	"context"

	"github.com/gvso/pro/pkg/auth"
	"github.com/gvso/pro/pkg/database"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// User is the representation of an user in the system.
type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	Lastname string             `json:"lastname" bson:"lastname"`
	Email    string             `json:"email" bson:"email"`

	// IDs in external authentication providers.
	GoogleID      string `json:"-" bson:"googleId"`
	GoogleToken   string `json:"-" bson:"googleToken"`
	FacebookID    string `json:"-" bson:"facebookId"`
	FacebookToken string `json:"-" bson:"facebookToken"`
}

// GetByProviderID returns a user whose provider's ID matches the given one.
func GetByProviderID(ctx context.Context, db database.Client, provider, ID string) (*User, error) {
	var providerIDKey string
	switch provider {
	case auth.GoogleProvider:
		providerIDKey = "google"
	default:
		return nil, errors.Errorf("invalid provider %v", provider)
	}
	providerIDKey += "Id"

	filter := bson.M{providerIDKey: ID}
	var user User
	err := db.Collection("users").FindOne(ctx, filter).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain user record")
	}
	return &user, nil
}

// Create adds a user record to the database.
func Create(ctx context.Context, db database.Client, provider string, userInfo auth.ProviderUser) (*User, error) {
	user := &User{
		Name:     userInfo.Name,
		Lastname: userInfo.Lastname,
		Email:    userInfo.Email,
	}
	switch provider {
	case auth.GoogleProvider:
		user.GoogleID = userInfo.UserID
		user.GoogleToken = userInfo.AccessToken
	default:
		return nil, errors.Errorf("unsupported provider %s", provider)
	}

	res, err := db.Collection("users").InsertOne(ctx, user)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert user in database")
	}

	id, ok := res.(primitive.ObjectID)
	if !ok {
		return nil, errors.Errorf("invalid user id: %v", res)
	}
	user.ID = id
	return user, nil
}
