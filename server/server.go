package server

import (
	"fmt"
	"net/http"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/misikdmitriy/password-sharing/controller"
	"github.com/penglongli/gin-metrics/ginmetrics"
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
	router := gin.Default()

	router.Use(ginzap.Ginzap(s.logger, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(s.logger, true))

	metrics := ginmetrics.GetMonitor()

	metrics.SetMetricPath("/metrics")
	metrics.Use(router)

	for _, ctrl := range s.controllers {
		method := ctrl.Method()

		switch method {
		case http.MethodGet:
			router.GET(ctrl.Route(), ctrl.Hander())
		case http.MethodPost:
			router.POST(ctrl.Route(), ctrl.Hander())
		default:
			return fmt.Errorf("cannot create HTTP handler of method %s", method)
		}
	}

	return router.Run(addr...)
}
