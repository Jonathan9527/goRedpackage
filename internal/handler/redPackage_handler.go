package handler

import (
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"

	"learnGO/internal/service"

	"github.com/gin-gonic/gin"
)

type RedPackageHandler struct {
	userService       *service.UserService
	redPackageService *service.RedPackageService
}

type sendRedPackageRequest struct {
	Account string `json:"account" binding:"required"`
	Amount  int    `json:"amount" binding:"required"`
}

func NewRedPackageHandler(redPackageService *service.RedPackageService) *RedPackageHandler {
	return &RedPackageHandler{
		redPackageService: redPackageService,
	}
}

func (h *RedPackageHandler) SendRedPackage(c *gin.Context) {
	var req sendRedPackageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	reqAmount := decimal.NewFromInt(int64(req.Amount))
	fmt.Printf("Received request to send red package: account=%s, amount=%d\n", req.Account, req.Amount)

	redPackageList, err := h.redPackageService.CreateRedPackage(c.Request.Context(), req.Account, reqAmount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, redPackageList)
}

func (h *RedPackageHandler) List(c *gin.Context) {
	limit := parseIntQuery(c, "limit", 20)
	offset := parseIntQuery(c, "offset", 0)

	users, err := h.userService.List(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query users failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   users,
		"limit":  limit,
		"offset": offset,
	})
}
