package cfg

import (
	"log"
	"github.com/caarlos0/env"
)

type Config struct {
	HvscPath  string  `env:"HVSC_PATH,required"`
}

func ReadConfig() *Config {
	config := Config{}
	err := env.Parse(&config)
	if err != nil {
		log.Fatalf("Config parsing failed: %+v", err)
	}

	return &config
}
