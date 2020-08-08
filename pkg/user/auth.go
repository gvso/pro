package user

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jonboulle/clockwork"
	"github.com/pkg/errors"
)

// JWTClaims is the user information embedded in the JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type JWTClaims struct {
	ID string `bson:"_id" json:"id"`
	jwt.StandardClaims
}

// Token returns the token for a user.
func (u User) Token(clock clockwork.Clock, duration time.Duration, jwtKey string) (string, error) {
	expiration := clock.Now().Add(duration)
	claims := &JWTClaims{
		ID: u.ID.Hex(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiration.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	key := []byte(jwtKey)
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", errors.Wrap(err, "could not get signed token")
	}

	return tokenString, nil
}

// ParseToken parses a JWT string to the claims struct.
func ParseToken(clock clockwork.Clock, tokenStr string, jwtKey string) (claims JWTClaims, err error) {
	jwt.TimeFunc = clock.Now
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return *claims, nil
	}
	return claims, errors.New("could not parse token claims")
}
