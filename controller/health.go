package controller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/misikdmitriy/password-sharing/health"
	"github.com/misikdmitriy/password-sharing/model"
)

type healthController struct {
	healthChecks []health.HealthCheck
}

func NewHealthController(healthChecks ...health.HealthCheck) Controller {
	return &healthController{
		healthChecks: healthChecks,
	}
}

func (ctrl *healthController) Hander() gin.HandlerFunc {
	return func(c *gin.Context) {
		totallyHealthy := true
		reason := ""

		for _, hc := range ctrl.healthChecks {
			healthy, err := hc.Check(c)
			if !healthy {
				totallyHealthy = false
				reason += fmt.Sprintf("%v\n", err)
			}
		}

		response := model.HealthResponse{
			Healthy: totallyHealthy,
			Reason:  strings.TrimRight(reason, "\n"),
		}

		if !totallyHealthy {
			c.JSON(500, response)
		} else {
			c.JSON(200, response)
		}
	}
}

func (c *healthController) Route() string {
	return "/health"
}

func (c *healthController) Method() string {
	return http.MethodGet
}
