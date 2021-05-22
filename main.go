package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
	"github.com/tjololo/linkerd-cert-notifier/pkg/linkerd"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	setupViper()
	undo := setupLogger()
	defer undo()
	config, err := rest.InClusterConfig()
	if err != nil {
		zap.L().Fatal(fmt.Sprintf("Failed to get kubernetes config. %s", err))
		os.Exit(1)
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		zap.L().Fatal(fmt.Sprintf("Failed to get kubernetes client. %s", err))
		os.Exit(1)
	}
	lc := linkerd.LinkerdReader{Client: client}
	ctx := context.Background()
	pem, err := lc.FetchTrustAnchor("linkerd", ctx)
	if err != nil {
		zap.L().Fatal(fmt.Sprintf("Failed to retrive trustAnchorPEM. %s", err))
		os.Exit(1)
	}
	s := string(pem)
	zap.L().Info(fmt.Sprintf("fetched configmap: %s", s))
}

func setupViper() {
	viper.SetDefault("developmend", false)
	viper.SetDefault("namespace", "linkerd")
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