package config

import (
	"log"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv  string `env:"APP_ENV" envDefault:"development"`
	DBPath  string `env:"DB_PATH" envDefault:"harbor.db"`
	Port    int    `env:"PORT" envDefault:"8088"`
	
}

// LoadConfig đọc file .env và parse vào struct
func LoadConfig() (*Config, error) {
	// Load .env nếu có (không báo lỗi nếu file không tồn tại - dùng cho Prod)
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("No .env file found, using system environment variables")
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}