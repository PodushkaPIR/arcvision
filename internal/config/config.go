package config

import (
	"log"
	"time"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	HttpServer        `yaml:"http_server"`
	AI                `yaml:"ai"`
	Deck              `yaml:"deck"`
}

type HttpServer struct {
    Address 	 string        `yaml:"address" env-default:"localhost:8082"`
    Timeout      time.Duration `yaml:"timeout" env-default:"5s"`
    IdleTimeout  time.Duration `yaml:"iddle_timeout" env-default:"60s"`
}

type AI struct {
	AIProvider string
	AIAPIKey   string
	AIBaseURL  string
	AIModel    string
}

type Deck struct {
	DeckPath   string `yaml:"deck_path" env:"DECK_PATH" env-default:"./web/data/deck.json"`
	SpreadPath string `yaml:"spread_path" env:"SPREAD_PATH" env-default:"./web/data/spreads.json"`
}

func MustLoad() *Config {

	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
