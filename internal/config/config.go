package config

import (
	"log"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	_ "time"
)

type Config struct{
	HTTPServer
	API
	Redis
}

type HTTPServer struct{
	HTTPPort string `env:"HTTP_PORT" env-default:"8080"`
}

type API struct{
	ApiKey string `env:"API_KEY" env-required:"true"`
	CurrencyApiEndpoint string `env:"CURRENCY_API_ENDPOINT" env-required:"true"`
}

type Redis struct {
	Addr     string `env:"ADDR" env-default:"localhost:6379"`
	Password string `env:"PASSWORD"`
	DB       int `env:"DB"`
	TTL      int `env:"TTL" env-default:"2"`
}
func MustLoad() *Config{
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf(".env file don’t want loading in environment variables: %v", err)
	}
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("cannot read env: $v", err)
	}
	return &cfg
}