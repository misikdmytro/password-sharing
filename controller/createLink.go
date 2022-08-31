package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	pserror "github.com/misikdmitriy/password-sharing/error"
	"github.com/misikdmitriy/password-sharing/model"
	"github.com/misikdmitriy/password-sharing/service"
)

type createLinkController struct {
	service service.PasswordService
}

func NewCreateLinkController(service service.PasswordService) Controller {
	return &createLinkController{
		service: service,
	}
}

func (ctrl *createLinkController) Hander() gin.HandlerFunc {
	type Body struct {
		Password string `json:"password"`
	}

	return func(c *gin.Context) {
		body := &Body{}
		err := c.BindJSON(body)
		if err != nil {
			c.JSON(pserror.BadRequestError())

			return
		}

		link, err := ctrl.service.CreateLinkFromPassword(c, body.Password)
		if err != nil {
			psError := pserror.AsPasswordSharingError(err)
			c.JSON(psError.ToResponse())

			return
		}

		c.JSON(http.StatusCreated, model.LinkResponse{
			Link: link,
		})
	}
}

func (ctrl *createLinkController) Route() string {
	return "/link"
}

func (ctrl *createLinkController) Method() string {
	return http.MethodPost
}
