package handler

import (
	"net/http"

	"learnGO/internal/service"

	"github.com/gin-gonic/gin"
)

type HomeHandler struct {
	homeService *service.HomeService
}

func NewHomeHandler(homeService *service.HomeService) *HomeHandler {
	return &HomeHandler{
		homeService: homeService,
	}
}

func (h *HomeHandler) Greeting(c *gin.Context) {
	c.JSON(http.StatusOK, h.homeService.Greeting())
}

func (h *HomeHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, h.homeService.Health())
}
