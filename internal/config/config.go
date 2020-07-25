package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

// Config is the configuration handler.
type Config struct{}

// New loads a .env file.
func New(filepath ...string) (*Config, error) {
	err := godotenv.Overload(filepath...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load .env file")
	}

	return &Config{}, nil
}

// Get retrieves an environment variable.
func (Config) Get(key string) string {
	return os.Getenv(key)
}

// Set sets an environment variable.
func (Config) Set(key, value string) {
	os.Setenv(key, value)
}
