package server

import (
	"github.com/gin-gonic/gin"
	"github.com/misikdmitriy/password-sharing/controller"
)

type Server interface {
	Run(addr ...string) error
}

type server struct {
	pwdController controller.PasswordController
}

func NewServer(pwdController controller.PasswordController) Server {
	return &server{
		pwdController: pwdController,
	}
}

func (s *server) Run(addr ...string) error {
	r := gin.Default()
	r.POST("/link", s.pwdController.CreateLinkFromPassword)

	return r.Run()
}
