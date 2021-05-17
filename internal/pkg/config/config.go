package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	Protocol string `json:"protocol"`
	Host     string `json:"host"`
	VHost    string `json:"vhost"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load(".env")

	port, err := strconv.Atoi(os.Getenv("MQ_PORT"))
	if err != nil {
		return nil, err
	}

	config := &Config{
		Protocol: os.Getenv("MQ_PROTOCOL"),
		Port:     port,
		Host:     os.Getenv("MQ_HOST"),
		Username: os.Getenv("MQ_USER"),
		Password: os.Getenv("MQ_PWD"),
		VHost:    os.Getenv("MQ_VHOST"),
	}

	return config, nil
}
