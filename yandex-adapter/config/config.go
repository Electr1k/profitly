package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env        string `env:"ENV" env-default:"development"`
	HttpServer `env-prefix:"HTTP_"`
	LogConfig  `env-prefix:"LOG_"`
	Yandex     `env-prefix:"YANDEX_"`
}

type HttpServer struct {
	Address         string        `env:"ADDRESS" env-default:""`
	Port            string        `env:"PORT" env-default:"3000"`
	Timeout         time.Duration `env:"TIMEOUT" env-default:"5s"`
	IdleTimeout     time.Duration `env:"IDLE_TIMEOUT" env-default:"60s"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" env-default:"10s"`
}

type LogConfig struct {
	Level  string `env:"LEVEL" env-default:"debug"`
	Format string `env:"FORMAT" env-default:"json"`
}

type Yandex struct {
	APIKey  string        `env:"APIKEY" env-required:"true"`
	BaseURL string        `env:"BASE_URL" env-default:"https://geocode-maps.yandex.ru/v1"`
	Lang    string        `env:"LANG" env-default:"ru_RU"`
	Timeout time.Duration `env:"TIMEOUT" env-default:"5s"`
}

func MustLoad() *Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found, using environment variables")
	}

	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
