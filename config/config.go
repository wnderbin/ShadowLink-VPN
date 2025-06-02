package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	BotApiKey       string `yaml:"BOT-API-KEY"`
	DebugMode       bool   `yaml:"DEBUG-MODE"`
	WGInterface     string `yaml:"WG-Interface"`
	WGConfigPath    string `yaml:"WG-ConfigPath"`
	ServerPublicIP  string `yaml:"Server-PublicIP"`
	ServerPort      int    `yaml:"ServerPort"`
	ServerPublicKey string `yaml:"ServerPublicKey"`
	AllowedIPs      string `yaml:"AllowedIPs"`
	DNS             string `yaml:"DNS"`
	DB              string `yaml:"DB"`
	RedisHost       string `yaml:"Redis-host"`
	RedisPass       string `yaml:"Redis-pass"`
	RedisDB         int    `yaml:"Redis-DB"`
}

func Load() *Config {
	config_path := os.Getenv("CONFIG_PATH")
	if config_path == "" {
		log.Fatal("[ config.go ] CONFIG_PATH is not set")
	}
	if _, err := os.Stat(config_path); os.IsExist(err) {
		log.Fatalf("[ config.go ] Config is not exist: %s\n", config_path)
	}
	var conf Config

	if err := cleanenv.ReadConfig(config_path, &conf); err != nil {
		log.Fatalf("[ config.go ] Cannot read config: %s\n", config_path)
	}

	return &conf
}
