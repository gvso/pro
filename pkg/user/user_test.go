package user_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gvso/pro/pkg/database/mock"
	"github.com/gvso/pro/pkg/user"
)

func TestGetByProviderID(t *testing.T) {
	tests := []struct {
		name          string
		provider      string
		userID        string
		srDecodeError error
		expected      user.User
		shouldErr     bool
	}{
		{
			name:     "successful retrieval from database",
			provider: "Google",
			userID:   "googleid123",
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
		t.Run(test.name, func(tt *testing.T) {
			dbClient := mock.ClientMock(test.expected, test.srDecodeError, test.expected.ID)
			user, err := user.GetByProviderID(context.TODO(), dbClient, test.provider, test.userID)

			if test.shouldErr {
				assert.NotNil(tt, err)
			} else {
				assert.Nil(tt, err)
				assert.Equal(tt, test.expected, *user)
			}
		})
	}
}
