package app

import (
	"github.com/gvso/pro/internal/config"
	"github.com/gvso/pro/pkg/auth"
	"github.com/gvso/pro/pkg/database"
	"github.com/sirupsen/logrus"
)

// App contains dependencies used across the application.
type App struct {
	Database           *database.Database
	Config             config.Config
	Logger             *logrus.Entry
	GoogleAuthProvider auth.Google
}
