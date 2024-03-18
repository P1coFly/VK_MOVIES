package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env"`
	Host_db     string `yaml:"host_db"`
	User_db     string `yaml:"user_db"`
	Password_db string `yaml:"password_db"`
	Name_db     string `yaml:"name_db"`
	Admin       `yaml:"admin"`
}

type Admin struct {
	AuthToken string `yaml:"auth_token"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	log.Printf("%s", configPath)
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("can't read config: %s", err)
	}

	return &cfg

}
