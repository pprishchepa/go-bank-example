package config

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/go-envparse"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	HTTP struct {
		Port int `env:"HTTP_PORT, default=8081"`
	}

	Log struct {
		Level  string `env:"LOG_LEVEL, default=info"`
		Pretty bool   `env:"LOG_PRETTY, default=false"`
	}

	Postgres Postgres `env:", prefix=POSTGRES_"`
	Redis    Redis    `env:", prefix=REDIS_"`
}

type Postgres struct {
	Host        string `env:"HOST, default=localhost"`
	Port        int    `env:"PORT, default=5432"`
	User        string `env:"USER, required"`
	Password    string `env:"PASSWORD, required"`
	Database    string `env:"DB, required"`
	SSLMode     string `env:"SSLMODE, default=verify-full"`
	ConnTimeout int    `env:"CONNTIMEOUT, default=5"`
	MaxConn     int    `env:"MAXCONN, default=8"`
}

type Redis struct {
	Host     string `env:"HOST, default=localhost"`
	Port     int    `env:"PORT, default=6379"`
	Username string `env:"USER"`
	Password string `env:"PASSWORD"`
	DB       int    `env:"DB"`
}

func NewConfig() (Config, error) {
	f, err := os.Open(".env")
	if err != nil && !os.IsNotExist(err) {
		return Config{}, err
	}

	if f != nil {
		envs, err := envparse.Parse(f)
		if err != nil {
			return Config{}, err
		}
		for k, v := range envs {
			if err = os.Setenv(k, v); err != nil {
				return Config{}, err
			}
		}
	}

	var conf Config
	if err := envconfig.Process(context.Background(), &conf); err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}

	return conf, nil
}
