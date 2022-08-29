package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

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
}

func (ctrl *createLinkController) Route() string {
	return "/link"
}

func (ctrl *createLinkController) Method() string {
	return http.MethodPost
}
