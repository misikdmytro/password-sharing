package main

import (
	"github.com/misikdmitriy/password-sharing/config"
	"github.com/misikdmitriy/password-sharing/controller"
	"github.com/misikdmitriy/password-sharing/database"
	"github.com/misikdmitriy/password-sharing/helper"
	"github.com/misikdmitriy/password-sharing/server"
	"github.com/misikdmitriy/password-sharing/service"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	dbf := database.NewFactory(conf)
	rf := helper.NewRandomFactory()
	service := service.NewPasswordService(dbf, conf, rf)
	passwordController := controller.NewPasswordController(service)

	server := server.NewServer(passwordController)

	if err := server.Run(); err != nil {
		panic(err)
	}
}
