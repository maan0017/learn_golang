package config

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HttpServer struct {
	URL string `yaml:"url" env-required:"true"`
}

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-required:"true" env-default:"production"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HttpServer  `yaml:"http_server"`
}

func MustLoad() *Config {
	var configPath string

	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		flags := flag.String("config", "config/local.yaml", "path to configuration file.")
		flag.Parse()

		configPath = *flags

		if configPath == "" {
			log.Fatal("Config path is not set")
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exists: %s", configPath)
	}

	var cfg Config

	// main - reads config file
	err := cleanenv.ReadConfig(configPath, &cfg)

	if err != nil {
		log.Fatalf("Cannot read config file: %s", err)
	}

	fmt.Printf("url >%s\n", cfg.URL)

	return &cfg
}
