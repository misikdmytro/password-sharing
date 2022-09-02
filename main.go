package main

import (
	"github.com/misikdmitriy/password-sharing/config"
	"github.com/misikdmitriy/password-sharing/controller"
	"github.com/misikdmitriy/password-sharing/database"
	"github.com/misikdmitriy/password-sharing/health"
	"github.com/misikdmitriy/password-sharing/helper"
	"github.com/misikdmitriy/password-sharing/logger"
	"github.com/misikdmitriy/password-sharing/server"
	"github.com/misikdmitriy/password-sharing/service"
)

func main() {
	appConfiguration, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	appLogger := logger.NewLoggerFactory(appConfiguration)

	encoder := helper.NewEncoder(appConfiguration)
	databaseFactory := database.NewFactory(appConfiguration, appLogger)
	randomFactory := helper.NewRandomFactory()
	service := service.NewPasswordService(databaseFactory, appConfiguration, randomFactory, appLogger, encoder)

	pgHealthCheck := health.NewPgHealthCheck(databaseFactory, appLogger)

	server := server.NewServer(
		appLogger,
		appConfiguration,
		controller.NewCreateLinkController(service, appConfiguration),
		controller.NewGetLinkController(service),
		controller.NewHealthController(pgHealthCheck),
	)

	if err = server.Run(); err != nil {
		panic(err)
	}
}
