package main

import (
	"log"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	setupViper()
	undo := setupLogger()
	defer undo()
}

func setupViper() {
	viper.SetDefault("developmend", false)
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath("/config")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatalf("Error reading config file, %v", err)
		}
	}
}

func setupLogger() func() {
	var err error
	var logger *zap.Logger
	if viper.GetBool("development") {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	defer logger.Sync()
	if err != nil {
		log.Fatalf("Error setup logger, %s", err)
	}
	return zap.ReplaceGlobals(logger)
}