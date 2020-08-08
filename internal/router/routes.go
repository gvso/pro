package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gvso/pro/internal/app"
	userroute "github.com/gvso/pro/internal/router/user"
)

// UserRoutes returns the routes for /user/
func UserRoutes(r *gin.Engine, app *app.App) *gin.Engine {
	r = userroute.LoginRoutes(r, app)
	return r
}
