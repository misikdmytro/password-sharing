package server

import (
	"fmt"
	"net/http"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/misikdmitriy/password-sharing/controller"
	"go.uber.org/zap"
)

type Server interface {
	Run(addr ...string) error
}

type server struct {
	controllers []controller.Controller
	logger      *zap.Logger
}

func NewServer(logger *zap.Logger, controllers ...controller.Controller) Server {
	return &server{
		controllers: controllers,
		logger:      logger,
	}
}

func (s *server) Run(addr ...string) error {
	r := gin.Default()

	r.Use(ginzap.Ginzap(s.logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(s.logger, true))

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

	return r.Run(addr...)
}
