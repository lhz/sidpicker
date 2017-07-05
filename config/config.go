package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/caarlos0/env"
)

type ConfigData struct {
	HvscBase string `env:"HVSC_BASE,required"`
	AppBase  string `env:"SIDPICKER_BASE"`
}

var Config *ConfigData

func init() {
	log.SetFlags(log.Ldate | log.Lmicroseconds)
}

func ReadConfig() {
	Config = &ConfigData{}
	err := env.Parse(Config)
	if err != nil {
		log.Fatalf("Config parsing failed: %+v", err)
	}
	ensureAppBase()
}

func ensureAppBase() {
	if len(Config.AppBase) < 1 {
		Config.AppBase = defaultAppBase()
	}
	if _, err := os.Stat(Config.AppBase); os.IsNotExist(err) {
		if err = os.MkdirAll(Config.AppBase, os.ModePerm); err != nil {
			log.Fatalf("Unable to create app directory %q: %v", Config.AppBase, err)
		}
	}
}

func defaultAppBase() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(os.Getenv("AppData"), "sidpicker")
	}
	return filepath.Join(os.Getenv("HOME"), ".sidpicker")
}
