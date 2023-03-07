// Package config main config
package config

import (
	"fmt"

	"github.com/caarlos0/env/v7"
)

// MainConfig with init data
type MainConfig struct {
	PostgresPort     string `env:"POSTGRES_PORT,notEmpty" envDefault:"5433"`
	PostgresHost     string `env:"POSTGRES_HOST,notEmpty" envDefault:"localhost"`
	PostgresPassword string `env:"POSTGRES_PASSWORD,notEmpty" envDefault:"postgres"`
	PostgresUser     string `env:"POSTGRES_USER,notEmpty" envDefault:"postgres"`
	PostgresDB       string `env:"POSTGRES_DB,notEmpty" envDefault:"postgres"`
	Port             string `env:"PORT,notEmpty" envDefault:"2000"`
	Host             string `env:"HOST,notEmpty" envDefault:"localhost"`
}

// NewMainConfig parsing config from environment
func NewMainConfig() (*MainConfig, error) {
	mainConfig := &MainConfig{}

	err := env.Parse(mainConfig)
	if err != nil {
		return nil, fmt.Errorf("config - NewMainConfig - Parse:%w", err)
	}

	return mainConfig, nil
}
