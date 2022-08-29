package main

import (
	"github.com/misikdmitriy/password-sharing/config"
	"github.com/misikdmitriy/password-sharing/controller"
	"github.com/misikdmitriy/password-sharing/database"
	"github.com/misikdmitriy/password-sharing/helper"
	"github.com/misikdmitriy/password-sharing/logger"
	"github.com/misikdmitriy/password-sharing/server"
	"github.com/misikdmitriy/password-sharing/service"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	log, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}
	defer log.Close()

	dbf := database.NewFactory(conf, log)
	rf := helper.NewRandomFactory()
	service := service.NewPasswordService(dbf, conf, rf, log)

	server := server.NewServer(
		controller.NewCreateLinkController(service),
		controller.NewGetLinkController(service),
	)

	if err := server.Run(); err != nil {
		panic(err)
	}
}
