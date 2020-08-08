package user_test

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/dgrijalva/jwt-go"
	"github.com/gvso/pro/pkg/user"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
)

func TestTokenAndParseToken(t *testing.T) {
	jwtKey := "key"
	now := time.Now()

	tests := []struct {
		name               string
		userID             string
		duration           time.Duration
		clockTokenGenerate clockwork.Clock
		clockTokenParse    clockwork.Clock
		expected           user.JWTClaims
		shouldErr          bool
	}{
		{
			name:               "valid token",
			userID:             "5f25af44687b6a1aff90b810",
			duration:           7 * 24 * time.Hour,
			clockTokenGenerate: clockwork.NewFakeClockAt(now),
			clockTokenParse:    clockwork.NewFakeClockAt(now.Add(time.Hour)),
			expected: user.JWTClaims{
				ID: "5f25af44687b6a1aff90b810",
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: now.Add(7 * 24 * time.Hour).Unix(),
				},
			},
		},
		{
			name:               "expired token",
			userID:             "5f25af44687b6a1aff90b810",
			duration:           7 * 24 * time.Hour,
			clockTokenGenerate: clockwork.NewFakeClockAt(now.Add(-8 * 24 * time.Hour)),
			clockTokenParse:    clockwork.NewFakeClockAt(now),
			shouldErr:          true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			userID, err := primitive.ObjectIDFromHex(test.userID)
			assert.Nil(tt, err)
			u := user.User{ID: userID}

			tokenStr, err := u.Token(test.clockTokenGenerate, test.duration, jwtKey)
			assert.Nil(tt, err)

			claims, err := user.ParseToken(test.clockTokenParse, tokenStr, jwtKey)
			if test.shouldErr {
				assert.NotNil(tt, err)
			} else {
				assert.Nil(tt, err)
				assert.Equal(tt, test.expected, claims)
			}
		})
	}
}
