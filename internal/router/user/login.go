package userroute

import (
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/gvso/pro/pkg/database"
	"github.com/gvso/pro/pkg/util"
	"github.com/jonboulle/clockwork"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/gvso/pro/internal/app"
	"github.com/gvso/pro/internal/config"
	"github.com/gvso/pro/internal/httpresponse"
	"github.com/gvso/pro/pkg/auth"
	"github.com/gvso/pro/pkg/user"
)

// LoginRoutes returns the routes for /user/login/
func LoginRoutes(r *gin.Engine, app *app.App) *gin.Engine {
	r.GET("/user/login/google", loginWith(app.Logger, app.GoogleAuthProvider))
	r.GET("/user/login/google/callback", callback(app.Logger, app.Config, app.Database, app.GoogleAuthProvider))

	return r
}

// loginWith handles the requests for login with an external provider such as
// Google and Facebook.
func loginWith(logger *logrus.Entry, provider auth.Provider) gin.HandlerFunc {

	return func(c *gin.Context) {
		logger := logger.WithFields(logrus.Fields{
			"requestId":    requestid.Get(c),
			"authProvider": provider.Name(),
		})

		url := provider.RedirectURL()
		logger.Debugf("redirecting for authentication")
		c.Redirect(302, url)
	}
}

// callback handles the requests after user has authenticated with an external
// provider such as Google and Facebook.
func callback(logger *logrus.Entry, cfg config.Config, db database.Client, provider auth.Provider) gin.HandlerFunc {

	return func(c *gin.Context) {
		logger := logger.WithFields(logrus.Fields{
			"requestId":    requestid.Get(c),
			"authProvider": provider.Name(),
		})

		logger.Debug("processing callback")
		logger.Debug("getting access token")
		userInfo, err := userInfoFromProvider(logger, provider, c.Query("code"))
		if err != nil {
			logger.Errorf("failed to get user information from provider: %v", err)
			c.JSON(500, httpresponse.Error("failed to get user info from provider"))
			return
		}

		// We have information on the user from the provider. We can proceed
		// to finding a linked account in our application.
		logger.Debug("getting or creating user")
		ctx := util.ContextWithLogger(c.Request.Context(), logger)
		u, err := user.GetOrCreate(ctx, db, provider.Name(), userInfo)
		if err != nil {
			logger.Errorf("failed to login user: %v", err)
			c.JSON(500, httpresponse.Error("could not log user in"))
			return
		}

		logger.Debug("generating JWT token")

		// 7 days for expiration.
		durationSeconds := 7 * 24 * 60 * 60
		tokenDuration := time.Duration(durationSeconds) * time.Second
		token, err := u.Token(clockwork.NewRealClock(), tokenDuration, cfg.Get("JWT_KEY"))
		if err != nil {
			logger.Errorf("failed to get JWT token: %v", err)
			c.JSON(500, httpresponse.Error("failed to get JWT token"))
			return
		}

		// Redirect to page after setting cookie.
		c.SetCookie("token", token, durationSeconds, "/", "", true, true)
	}
}

// userInfoFromProvider retrieves a long-lived access token and retrieves the
// user information from the provider
func userInfoFromProvider(logger *logrus.Entry, provider auth.Provider, code string) (u auth.ProviderUser, err error) {
	accessToken, err := provider.AccessToken(code)
	if err != nil {
		return u, errors.Wrap(err, "failed to get access token")
	}

	logger.Debug("retrieving user information")
	u, err = provider.UserInfo(accessToken)
	return u, errors.Wrap(err, "failed to retrieve user information")
}
