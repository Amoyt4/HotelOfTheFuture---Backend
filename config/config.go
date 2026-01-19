package config

import (
	"log/slog"

	"github.com/caarlos0/env/v11"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Env struct {
	DB_HOST     string `env:"DB_HOST"`
	DB_PORT     int    `env:"DB_PORT"`
	DB_NAME     string `env:"DB_NAME"`
	DB_USERNAME string `env:"DB_USERNAME"`
	DB_PASSWORD string `env:"DB_PASSWORD"`
	IP_ADDRESS  string `env:"IP_ADDRESS"`
	API_PORT    string `env:"API_PORT"`
	JWT_SECRET  string `env:"JWT_SECRET"`
}

type Config struct {
	Env    Env
	Client *pgxpool.Pool
}

var config Config

func GetConfig() *Config {
	config.Env = *getEnv()
	slog.Info("Getting config successfully")
	return &config
}

// maybe if we have huge env file we should
// return pointer of Env. but in this situation
// i prefer to test it like this
func getEnv() *Env {
	err := godotenv.Load()
	if err != nil {
		slog.Warn("Error loading .env file")
	}

	var envCfg Env
	err = env.Parse(&envCfg)
	if err != nil {
		slog.Warn("Error parsing .env file")
	}
	return &envCfg
}
