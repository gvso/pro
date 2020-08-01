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

// GetOrCreate registers or gets an user from database.
//
// If the user was not registered yet, it creates an record in the database.
func GetOrCreate(ctx context.Context, db database.Client, provider string,
	userInfo auth.ProviderUser) (user User, err error) {

	var providerIDKey string
	switch provider {
	case auth.GoogleProvider:
		providerIDKey = "google"
	default:
		return user, errors.Errorf("invalid provider %v", provider)
	}
	providerIDKey += "Id"

	filter := bson.M{providerIDKey: userInfo.UserID}
	err = db.Collection("users").FindOne(ctx, filter).Decode(&user)
	if err == mongo.ErrNoDocuments {
		user, err := createUser(ctx, db, provider, userInfo)
		return user, errors.Wrap(err, "failed to create user")

	}
	if err != nil {
		return user, errors.Wrap(err, "error when querying database")
	}

	return user, nil
}

func createUser(ctx context.Context, db database.Client, provider string,
	userInfo auth.ProviderUser) (user User, err error) {

	user = User{
		Name:     userInfo.Name,
		Lastname: userInfo.Lastname,
		Email:    userInfo.Email,
	}
	switch provider {
	case auth.GoogleProvider:
		user.GoogleID = userInfo.UserID
		user.GoogleToken = userInfo.AccessToken
	default:
		return user, errors.Errorf("unsupported provider %s", provider)
	}

	res, err := db.Collection("users").InsertOne(ctx, user)
	if err != nil {
		errors.Wrap(err, "failed to insert user in database")
	}
	user.ID = res.(primitive.ObjectID)
	return user, nil
}
