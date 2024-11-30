package config

import (
	"flag"
	"fmt"
	"os"
)

// ErrInitConfigFailed - config initialization error.
var ErrInitConfigFailed = fmt.Errorf("failed to init config")

// Config - application configuration structure.
type Config struct {
	RunAddress           string // Address and port of HTTP server
	DatabaseURI          string // Address for database connection
	AccrualSystemAddress string // Address of accrual system
	SecretKey            string // Authentication secret key

}

// configBuilder - application configuration builder.
type configBuilder struct {
	runAddress           string `env:"RUN_ADDRESS"`
	databaseURI          string `env:"DATABASE_URI"`
	accrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	secretKey            string `env:"SECRET_KEY"`
}

// newConfigBuilder creates new application configuration builder.
func newConfigBuilder() *configBuilder {
	return &configBuilder{}
}

// setDefaults defines application configuration parameters defaults.
func (cb *configBuilder) setDefaults() error {
	cb.runAddress = "localhost:8080"
	cb.databaseURI = ""
	cb.accrualSystemAddress = ""
	cb.secretKey = "secret"

	return nil
}

// setFlags sets application configuration parameters from command line parameters.
func (cb *configBuilder) setFlags() error {
	flag.StringVar(&cb.runAddress, "a", cb.runAddress, "HTTP server address and port")
	flag.StringVar(&cb.databaseURI, "d", cb.databaseURI, "database connection string")
	flag.StringVar(&cb.accrualSystemAddress, "r", cb.accrualSystemAddress, "accrual system address and port")
	flag.Parse()

	return nil
}

// setEnvs sets application configuration parameters from environment variables.
func (cb *configBuilder) setEnvs() error {
	ra := os.Getenv("RUN_ADDRESS")
	if ra != "" {
		cb.runAddress = ra
	}

	dbi := os.Getenv("DATABASE_URI")
	if dbi != "" {
		cb.databaseURI = dbi
	}

	asa := os.Getenv("ACCRUAL_SYSTEM_ADDRESS")
	if asa != "" {
		cb.accrualSystemAddress = asa
	}

	sk := os.Getenv("SECRET_KEY")
	if sk != "" {
		cb.secretKey = sk
	}

	return nil
}

// build builds application cofiguration.
func (cb *configBuilder) build() *Config {
	return &Config{
		RunAddress:           cb.runAddress,
		DatabaseURI:          cb.databaseURI,
		AccrualSystemAddress: cb.accrualSystemAddress,
		SecretKey:            cb.secretKey,
	}
}

// Get returns application configuration.
func Get() (*Config, error) {
	cb := newConfigBuilder()

	confSets := []func() error{
		cb.setDefaults,
		cb.setFlags,
		cb.setEnvs,
	}

	for _, confSet := range confSets {
		err := confSet()
		if err != nil {
			return nil, ErrInitConfigFailed
		}
	}

	return cb.build(), nil
}
