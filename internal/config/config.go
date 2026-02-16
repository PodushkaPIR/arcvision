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
	Storage           `yaml:"storage"`
	AI                `yaml:"ai"`
}

type HttpServer struct {
    Address 	 string        `yaml:"address" env-default:"localhost:8082"`
    Timeout      time.Duration `yaml:"timeout" env-default:"5s"`
    IdleTimeout  time.Duration `yaml:"iddle_timeout" env-default:"60s"`
}

type Storage struct {
	DeckPath   string `yaml:"deck_path" env:"DECK_PATH" env-default:"./web/data/deck.json"`
	SpreadPath string `yaml:"spread_path" env:"SPREAD_PATH" env-default:"./web/data/spreads.json"`
}

type AIConfig struct {
	Mode     string `yaml:"mode" env:"AI_MODE" env-default:"local"`

	Cloud struct {
		URL      string `yaml:"url" env:"AI_CLOUD_URL"`
		Key      string `env:"AI_CLOUD_KEY"`
		FolderID string `env:"AI_CLOUD_FOLDER_ID"`
		Model    string `yalm:"model" env-default:"yandexgpt-lite"`
	} `yaml:"cloud"`

	Local struct {
		URL   string `yaml:"url" env:"AI_LOCAL_URL" env-default:"https://localhost:11434/api/generate`
		Model string `yaml:"model" env-default:"llama3"`

	} `yaml:"local"`

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
