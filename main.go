package main

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/viper"
	"github.com/tjololo/linkerd-cert-notifier/pkg/certificate"
	"github.com/tjololo/linkerd-cert-notifier/pkg/linkerd"
	"github.com/tjololo/linkerd-cert-notifier/pkg/notification"
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
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		zap.L().Fatal(fmt.Sprintf("Failed to get kubernetes client. %s", err))
	}
	lr := linkerd.Reader{Client: client}
	ctx := context.Background()
	pem, err := lr.FetchTrustAnchor(ctx, "linkerd")
	if err != nil {
		zap.L().Fatal(fmt.Sprintf("Failed to retrive trustAnchorPEM. %s", err))
	}
	expiring, date, err := certificate.AboutToExpire(pem, viper.GetString("earlyexpire.anchor"))
	if err != nil {
		zap.L().Fatal(fmt.Sprintf("Failed to check trust anchor certificate. %s", err))
	}
	if expiring {
		zap.L().Warn(fmt.Sprintf("trust anchor cert about to expire. Expiring: %s", date))
		err := notification.SendSlackNotification(notification.SlackRequestBody{
			Username: viper.GetString("slack.username"),
			Channel:  viper.GetString("slack.channel"),
			Text:     fmt.Sprintf("Trust anchor cert about to expire. Expiring: %s", date),
		})
		if err != nil {
			zap.L().Fatal(fmt.Sprintf("Failed to send message to slack. %s", err))
		}
	} else {
		zap.L().Info(fmt.Sprintf("trust anchor cert not about to expire. Expiring: %s", date))
	}

	pem, err = lr.FetchIssuerCert(ctx, "linkerd")
	if err != nil {
		zap.L().Fatal(fmt.Sprintf("Failed to retrive issuerPEM. %s", err))
	}
	expiring, date, err = certificate.AboutToExpire(pem, viper.GetString("earlyexpire.issuer"))
	if err != nil {
		zap.L().Fatal(fmt.Sprintf("Failed to check issuer certificate. %s", err))
	}
	if expiring {
		zap.L().Warn(fmt.Sprintf("issuer cert about to expire. Expiring: %s", date))
		err := notification.SendSlackNotification(notification.SlackRequestBody{
			Username: viper.GetString("slack.username"),
			Channel:  viper.GetString("slack.channel"),
			Text:     fmt.Sprintf("Issuer cert about to expire. Expiring: %s", date),
		})
		if err != nil {
			zap.L().Fatal(fmt.Sprintf("Failed to send message to slack. %s", err))
		}
	} else {
		zap.L().Info(fmt.Sprintf("issuer cert not about to expire. Expiring: %s", date))
	}
}

func setupViper() {
	viper.SetDefault("development", false)
	viper.SetDefault("namespace", "linkerd")
	viper.SetDefault("earlyexpire.anchor", "1440h")
	viper.SetDefault("earlyexpire.issuer", "1440h")
	viper.SetDefault("slack.username", "linkerd-cert-notifier")
	viper.SetDefault("slack.channel", "linkerd-alerts")
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
