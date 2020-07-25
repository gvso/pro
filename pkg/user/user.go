package user

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User is the representation of an user in the system.
type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	Lastname string             `json:"lastname" bson:"lastname"`

	// IDs in external authentication providers.
	GoogleID   string `json:"googleId" bson:"googleId"`
	FacebookID string `json:"facebookId" bson:"facebookId"`
}
