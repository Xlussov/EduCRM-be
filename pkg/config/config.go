package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env        string     `yaml:"env" env-default:"prod" env-required:"true"`
	JWT        JWT        `yaml:"jwt" env-required:"true"`
	HTTPServer HTTPServer `yaml:"http_server"`
	Postgres   Postgres   `yaml:"postgres"`
}

type HTTPServer struct {
	Address      string        `yaml:"address"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

type JWT struct {
	Secret     string        `yaml:"secret" env:"JWT_SECRET" env-required:"true"`
	AccessTTL  time.Duration `yaml:"access_ttl" env:"JWT_ACCESS_TTL" env-default:"15m"`
	RefreshTTL time.Duration `yaml:"refresh_ttl" env:"JWT_REFRESH_TTL" env-default:"720h"`
}

type Postgres struct {
	URL string `yaml:"url" env:"DATABASE_URL" env-required:"true"`
}

func LoadENV() error {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}
	return nil
}

func LoadConfig() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		return nil, fmt.Errorf("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("cannot read config: %w", err)
	}

	return &cfg, nil
}

func MustLoad() *Config {
	if err := LoadENV(); err != nil {
		log.Fatalf("[CONFIG] %v", err)
	}

	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("[CONFIG] %v", err)
	}

	return cfg
}
