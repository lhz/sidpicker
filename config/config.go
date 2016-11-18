package config

import (
	"log"

	"github.com/caarlos0/env"
)

type ConfigData struct {
	HvscPath string `env:"HVSC_PATH,required"`
}

var Config *ConfigData

func ReadConfig() {
	Config = &ConfigData{}
	err := env.Parse(Config)
	if err != nil {
		log.Fatalf("Config parsing failed: %+v", err)
	}
}
