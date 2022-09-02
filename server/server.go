package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"github.com/misikdmitriy/password-sharing/config"
	"github.com/misikdmitriy/password-sharing/controller"
	"github.com/misikdmitriy/password-sharing/logger"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"go.uber.org/zap"
)

type Server interface {
	Run() error
}

type server struct {
	controllers   []controller.Controller
	loggerFactory logger.LoggerFactory
	config        *config.Config
}

func NewServer(loggerFactory logger.LoggerFactory, config *config.Config, controllers ...controller.Controller) Server {
	return &server{
		controllers:   controllers,
		config:        config,
		loggerFactory: loggerFactory,
	}
}

const serviceName = "passwordsharing"
const healthCheck = "healthcheck"

func (s *server) Run() error {
	appLogger, closeLogger, err := s.loggerFactory.NewLogger()
	if err != nil {
		return err
	}
	defer closeLogger()

	appLogger.Info("starting web server...")
	router, err := s.buildRouter(appLogger)
	if err != nil {
		return err
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.App.Port),
		Handler: router,
	}

	startupFailed := make(chan interface{}, 1)
	defer close(startupFailed)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			startupFailed <- true
		}
	}()

	deregister, err := s.registerInConsul()
	if err != nil {
		return err
	}

	quit := make(chan os.Signal, 1)
	defer close(quit)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		log.Println("shutdown web server ...")
		deregister()
		return nil
	case <-startupFailed:
		return fmt.Errorf("error on server start")
	}
}

func (s *server) buildRouter(appLogger *zap.Logger) (*gin.Engine, error) {
	router := gin.Default()

	router.Use(ginzap.Ginzap(appLogger, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(appLogger, true))

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
			return nil, fmt.Errorf("cannot create HTTP handler of method %s", method)
		}
	}

	return router, nil
}

func (s *server) registerInConsul() (func(), error) {
	client, err := api.NewClient(&api.Config{
		Address: s.config.App.ConsulAddress,
		Scheme:  "http",
	})
	if err != nil {
		return nil, err
	}

	registration := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%s-%d", serviceName, s.config.App.ServiceId),
		Name:    serviceName,
		Port:    s.config.App.Port,
		Address: s.config.App.Address,
		Checks: []*api.AgentServiceCheck{
			{
				CheckID:  fmt.Sprintf("%s-%d", healthCheck, s.config.App.ServiceId),
				Name:     healthCheck,
				Timeout:  "5s",
				HTTP:     fmt.Sprintf("http://%s:%d/health", s.config.App.Address, s.config.App.Port),
				Method:   http.MethodGet,
				Interval: "15s",
			},
		},
	}

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		return nil, err
	}

	deregister := func() {
		client.Agent().ServiceDeregister(registration.ID)
	}

	return deregister, nil
}
