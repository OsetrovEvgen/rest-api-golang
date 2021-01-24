package server

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Config for APIServer struct
type Config struct {
	BindAddr string
	Logger   *logrus.Logger
}

// NewConfig return default config if .env file was not found
func NewConfig() *Config {
	// default config settings
	resConfig := &Config{
		BindAddr: ":3000",
		Logger:   logrus.New(),
	}
	level, _ := logrus.ParseLevel("debug")
	resConfig.Logger.SetLevel(level)

	logrus.Info("trying to start custom config for apiserver")

	// trying to get key LOG_LEVEL
	ll := os.Getenv("LOG_LEVEL")
	if ll != "" {
		level, err := logrus.ParseLevel(ll)
		if err != nil {
			logrus.Fatal(err)
		}
		resConfig.Logger.SetLevel(level)
	} else {
		logrus.Warn("key LOG_LEVEL not found. Setting default value " + resConfig.Logger.GetLevel().String())
	}

	// trying to get key BIND_ADDRESS
	adr := os.Getenv("BIND_ADDRESS")
	if adr != "" {
		resConfig.BindAddr = adr
	} else {
		logrus.Warn("key BIND_ADDRESS not found. Setting default value " + resConfig.BindAddr)
	}

	return resConfig
}
