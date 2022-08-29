package controller

import "github.com/gin-gonic/gin"

type Controller interface {
	Hander() gin.HandlerFunc
	Route() string
	Method() string
}
