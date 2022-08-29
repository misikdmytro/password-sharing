package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/misikdmitriy/password-sharing/model"
	"github.com/misikdmitriy/password-sharing/service"
)

type getLinkController struct {
	service service.PasswordService
}

func NewGetLinkController(service service.PasswordService) Controller {
	return &getLinkController{
		service: service,
	}
}

func (ctrl *getLinkController) Hander() gin.HandlerFunc {
	return func(c *gin.Context) {
		link := c.Param("link")
		if link == "" {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{
				Message: "bad request",
			})

			return
		}

		password, err := ctrl.service.GetPasswordFromLink(link)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{
				Message: "internal error",
			})

			return
		}

		c.JSON(http.StatusOK, model.PasswordResponse{
			Password: *password,
		})
	}
}

func (ctrl *getLinkController) Route() string {
	return "/pwd/:link"
}

func (ctrl *getLinkController) Method() string {
	return http.MethodGet
}
