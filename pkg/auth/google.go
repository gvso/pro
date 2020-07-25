package auth

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Google is the provider for Google Authentication.
type Google struct {
	oauth2Config *oauth2.Config
	clientID     string
	clientSecret string
	baseURL      string
	callbackURL  string
	scopes       []string
}

// NewGoogleProvider initializes the provider for Gooogle.
func NewGoogleProvider(clientID, clientSecret, baseURL string) *Google {
	callbackURL := baseURL + "/user/login/google/callback"
	scopes := []string{"openid", "email", "profile"}
	return &Google{
		oauth2Config: newOAuth2Config(clientID, clientSecret, callbackURL, scopes, google.Endpoint),
	}
}

// Name returns the provider's name.
func (Google) Name() string {
	return "Google"
}

// RedirectURL returns the redirect URL for authentication.
func (g *Google) RedirectURL() string {
	return g.oauth2Config.AuthCodeURL("state")
}

// AccessToken retrieves an access token user has returned from Google.
func (g *Google) AccessToken(code string) (*oauth2.Token, error) {
	token, err := g.oauth2Config.Exchange(oauth2.NoContext, code)
	return token, errors.Wrap(err, "failed to obtain access token")
}

// UserInfo retrieves information about the user from the provider.
func (g *Google) UserInfo(token *oauth2.Token) (*ProviderUser, error) {
	client := g.oauth2Config.Client(oauth2.NoContext, token)
	resp, err := client.Get("https://openidconnect.googleapis.com/v1/userinfo")
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain user information")
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read user information")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unsuccessful request to obtain user information")
	}

	obj := map[string]interface{}{}
	if err := json.Unmarshal(data, &obj); err != nil {
		log.Fatal(err)
	}

	return &ProviderUser{
		Email:       obj["email"].(string),
		Name:        obj["given_name"].(string),
		Lastname:    obj["family_name"].(string),
		Fullname:    obj["name"].(string),
		AccessToken: token.AccessToken,
	}, nil
}
