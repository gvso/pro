package user_test

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gvso/pro/pkg/auth"
	"github.com/gvso/pro/pkg/database/mock"
	"github.com/gvso/pro/pkg/user"
)

func TestSignInOrLogin(t *testing.T) {
	tests := []struct {
		name          string
		provider      string
		userInfo      auth.ProviderUser
		srDecodeError error
		expected      user.User
		shouldErr     bool
	}{
		{
			name:     "user is new, create it",
			provider: "Google",
			userInfo: auth.ProviderUser{
				UserID:      "googleid123",
				Name:        "name",
				Lastname:    "lastname",
				Email:       "email",
				AccessToken: "googletoken",
			},
			srDecodeError: mongo.ErrNoDocuments,
			expected: user.User{
				ID:          primitive.NewObjectID(),
				Name:        "name",
				Lastname:    "lastname",
				Email:       "email",
				GoogleID:    "googleid123",
				GoogleToken: "googletoken",
			},
		},
		{
			name:     "user is not new",
			provider: "Google",
			userInfo: auth.ProviderUser{
				UserID:      "googleid123",
				Name:        "name",
				Lastname:    "lastname",
				Email:       "email",
				AccessToken: "googletoken",
			},
			srDecodeError: nil,
			expected: user.User{
				ID:          primitive.NewObjectID(),
				Name:        "name",
				Lastname:    "lastname",
				Email:       "email",
				GoogleID:    "googleid123",
				GoogleToken: "googletoken",
			},
		},
		{
			name:      "invalid provider",
			provider:  "AProvider",
			shouldErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dbClient := mock.ClientMock(test.expected, test.srDecodeError, test.expected.ID)
			user, err := user.GetOrCreate(context.TODO(), dbClient, test.provider, test.userInfo)

			if test.shouldErr {
				assert.NotNil(t, err)
			} else {
				assert.Equal(t, test.expected, user)
			}
		})
	}
}
