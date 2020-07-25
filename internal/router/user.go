package router

import (
	"github.com/gvso/pro/pkg/auth"

	"github.com/gin-contrib/requestid"
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/gvso/pro/internal/app"
)

// UserRoutes returns the routes for /user
func UserRoutes(r *gin.Engine, app *app.App) *gin.Engine {
	r.GET("/user/login/google", loginWith(app.Logger, app.GoogleAuthProvider))
	r.GET("/user/login/google/callback", callback(app.Logger, app.GoogleAuthProvider))

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
func callback(logger *logrus.Entry, provider auth.Provider) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := logger.WithFields(logrus.Fields{
			"requestId":    requestid.Get(c),
			"authProvider": provider.Name(),
		})
		logger.Debug("processing callback")

		code := c.Query("code")
		token, err := provider.AccessToken(code)
		if err != nil {
			logger.Errorf("failed to get access token: %v", err)
			c.JSON(500, ServerErrorResponse{500, "could not get token"})
		}

		user, err := provider.UserInfo(token)
		if err != nil {
			logger.Errorf("failed to get user information: %v", err)
			c.JSON(500, ServerErrorResponse{500, "could not get user information"})
		}

		c.JSON(200, user)
	}
}
