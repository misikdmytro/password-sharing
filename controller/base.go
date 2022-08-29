package controller

import "github.com/gin-gonic/gin"

type HandlerFunc func(c *gin.Context)

type Controller interface {
	Hander() gin.HandlerFunc
	Route() string
	Method() string
}
