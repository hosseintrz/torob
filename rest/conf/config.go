package conf

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
)

func GetEnv(key string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", err
	}
	secret, ok := os.LookupEnv(key)
	if !ok {
		return "", errors.New("env var doesn't exists")
	}
	return secret, nil
}

type Config struct {
	Host string
	Port string
}
