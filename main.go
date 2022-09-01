package main

import (
	"fmt"

	"github.com/misikdmitriy/password-sharing/config"
	"github.com/misikdmitriy/password-sharing/controller"
	"github.com/misikdmitriy/password-sharing/database"
	"github.com/misikdmitriy/password-sharing/health"
	"github.com/misikdmitriy/password-sharing/helper"
	"github.com/misikdmitriy/password-sharing/logger"
	"github.com/misikdmitriy/password-sharing/server"
	"github.com/misikdmitriy/password-sharing/service"
	"go.uber.org/zap"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	log, close, err := logger.NewLogger(conf)
	if err != nil {
		panic(err)
	}
	defer close()

	dbf := database.NewFactory(conf, log)
	rf := helper.NewRandomFactory()
	service := service.NewPasswordService(dbf, conf, rf, log)

	pgHealthCheck := health.NewPgHealthCheck(dbf, func(err error) {
		log.Error("postgres health check failed",
			zap.Error(err))
	})

	server := server.NewServer(
		log,
		controller.NewCreateLinkController(service),
		controller.NewGetLinkController(service),
		controller.NewHealthController(pgHealthCheck),
	)

	if err := server.Run(fmt.Sprintf(":%d", conf.App.Port)); err != nil {
		panic(err)
	}
}
