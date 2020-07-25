package auth

import (
	"golang.org/x/oauth2"
)

// ProviderUser is information about an user retrieved from provider.
type ProviderUser struct {
	Email       string
	Name        string
	Lastname    string
	Fullname    string
	AccessToken string
}

// Provider is an authentication provider such as Facebook and Google.
type Provider interface {
	// Returns the provider's name.
	Name() string

	// Returns the redirect URL for authentication with the provider.
	RedirectURL() string

	// Retrieves an access token user has returned from provider.
	AccessToken(code string) (*oauth2.Token, error)

	// Retrieves information about the user from the provider.
	UserInfo(token *oauth2.Token) (*ProviderUser, error)
}

// newOAuth2Config initializes configuration for OAuth2 authentication.
func newOAuth2Config(clientID, clientSecret, callbackURL string, scopes []string, endpoint oauth2.Endpoint) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  callbackURL,
		Scopes:       scopes,
		Endpoint:     endpoint,
	}
}
