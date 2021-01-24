package store

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Config ...
type Config struct {
	DBUser     string
	DBHost     string
	DBPassword string
	DBName     string
	DBMode     string
}

// NewConfig ...
func NewConfig() *Config {
	resConfig := &Config{
		DBUser:     "root",
		DBHost:     "localhost",
		DBPassword: "pass",
		DBName:     "restapi_dev",
		DBMode:     "disable",
	}

	logrus.Info("trying to start custom config for database")

	dbu := os.Getenv("DB_USER")
	if dbu != "" {
		resConfig.DBUser = dbu
	} else {
		logrus.Warn("key DB_USER not found. Setting default value " + resConfig.DBUser)
	}

	dbh := os.Getenv("DB_HOST")
	if dbh != "" {
		resConfig.DBHost = dbh
	} else {
		logrus.Warn("key DB_HOST not found. Setting default value " + resConfig.DBHost)
	}

	dbp := os.Getenv("DB_PASSWORD")
	if dbp != "" {
		resConfig.DBPassword = dbp
	} else {
		logrus.Warn("key DB_PASSWORD not found. Setting default value " + resConfig.DBPassword)
	}

	dbn := os.Getenv("DB_NAME")
	if dbn != "" {
		resConfig.DBName = dbn
	} else {
		logrus.Warn("key DB_NAME not found. Setting default value " + resConfig.DBName)
	}

	dbm := os.Getenv("SSL_MODE")
	if dbm != "" {
		resConfig.DBMode = dbm
	} else {
		logrus.Warn("key SSL_MODE not found. Setting default value " + resConfig.DBMode)
	}

	return resConfig
}
