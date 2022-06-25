package config

import (
	"errors"
	"github.com/hosseintrz/torob/product/internal/db"
	"github.com/joho/godotenv"
	"os"
)

const (
	DBTypeDefault = db.MONGODB
	DBConnDefault = "mongodb://127.0.0.1:27017"
)

type ServiceConfig struct {
	DbType db.DBType `json:"db_type"`
	DbConn string    `json:"db_conn"`
}

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

func GetConfig() *ServiceConfig {
	conf := &ServiceConfig{
		DbType: DBTypeDefault,
		DbConn: DBConnDefault,
	}
	dbType, err := GetEnv("DB_TYPE")
	if err == nil {
		conf.DbType = db.DBType(dbType)
	}
	dbConn, err := GetEnv("DB_CONN")
	if err == nil {
		conf.DbConn = dbConn
	}
	return conf
}
