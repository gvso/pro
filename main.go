package main

import (
	"context"
	"fmt"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/gvso/pro/internal/app"
	"github.com/gvso/pro/internal/config"
	"github.com/gvso/pro/internal/router"
	"github.com/gvso/pro/internal/uuid"
	"github.com/gvso/pro/pkg/auth"
	"github.com/gvso/pro/pkg/database"
)

func main() {
	ctx := context.Background()
	logger := logrus.New()

	config, err := config.New()
	if err != nil {
		logger.Fatalf("failed to load initialize configuration: %v", err)
	}

	loggingLevel, err := logrus.ParseLevel(config.Get("LOGGING_LEVEL"))
	if err != nil {
		logger.Fatalf("invalid logging level: %v", err)
	}
	logger.SetLevel(loggingLevel)

	// Database.
	uri := config.Get("DB_URI")
	if uri == "" {
		uri = "mongodb://" + config.Get("DB_HOST") + ":" + config.Get("DB_PORT")
	}
	database, err := database.New(ctx, uri, config.Get("DB_NAME"))

	// OAuth2 providers
	googleProvider := auth.NewGoogleProvider(
		config.Get("GOOGLE_CLIENT_ID"),
		config.Get("GOOGLE_CLIENT_SECRET"),
		config.Get("BASE_URL"),
	)

	app := &app.App{
		Database:           database,
		Config:             config,
		Logger:             logrus.NewEntry(logger),
		GoogleAuthProvider: googleProvider,
	}

	port := fmt.Sprintf(":%s", config.Get("HOST_PORT"))
	r := gin.Default()
	r.Use(requestid.New(requestid.Config{
		Generator: func() string {
			return uuid.Get()
		},
	}))

	r = router.UserRoutes(r, app)
	r.Run(port)
}
