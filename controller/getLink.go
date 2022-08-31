package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	pserror "github.com/misikdmitriy/password-sharing/error"
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
			c.JSON(pserror.BadRequestError())

			return
		}

		password, err := ctrl.service.GetPasswordFromLink(c, link)
		if err != nil {
			psError := pserror.AsPasswordSharingError(err)
			c.JSON(psError.ToResponse())

			return
		}

		c.JSON(http.StatusOK, model.PasswordResponse{
			Password: password,
		})
	}
}

func (ctrl *getLinkController) Route() string {
	return "/pwd/:link"
}

func (ctrl *getLinkController) Method() string {
	return http.MethodGet
}
