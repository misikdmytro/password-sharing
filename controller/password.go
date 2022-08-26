package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/misikdmitriy/password-sharing/model"
	"github.com/misikdmitriy/password-sharing/service"
)

type PasswordController interface {
	CreateLinkFromPassword(*gin.Context)
}

type passwordController struct {
	service service.PasswordService
}

func NewPasswordController(service service.PasswordService) PasswordController {
	return &passwordController{
		service: service,
	}
}

func (ctrl *passwordController) CreateLinkFromPassword(c *gin.Context) {
	type Body struct {
		Password string `json:"password"`
	}

	body := &Body{}
	err := c.BindJSON(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Message: "bad request",
		})

		return
	}

	link, err := ctrl.service.CreateLinkFromPassword(body.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Message: "internal error",
		})

		return
	}

	c.JSON(http.StatusCreated, model.LinkResponse{
		Link: link,
	})
}
