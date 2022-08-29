package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/misikdmitriy/password-sharing/controller"
)

type Server interface {
	Run(addr ...string) error
}

type server struct {
	controllers []controller.Controller
}

func NewServer(controllers ...controller.Controller) Server {
	return &server{
		controllers: controllers,
	}
}

func (s *server) Run(addr ...string) error {
	r := gin.Default()

	for _, ctrl := range s.controllers {
		method := ctrl.Method()

		switch method {
		case http.MethodGet:
			r.GET(ctrl.Route(), ctrl.Hander())
		case http.MethodPost:
			r.POST(ctrl.Route(), ctrl.Hander())
		default:
			return fmt.Errorf("cannot create HTTP handler of method %s", method)
		}
	}

	return r.Run()
}
