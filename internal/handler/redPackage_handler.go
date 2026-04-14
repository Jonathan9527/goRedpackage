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
	Number  int    `json:"number" binding:"required"`
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

	if req.Amount <= 0 || req.Number <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount and number must be positive"})
		return
	}

	reqAmount := decimal.NewFromInt(int64(req.Amount))
	fmt.Printf("Received request to send red package: account=%s, amount=%d, number=%d\n", req.Account, req.Amount, req.Number)

	if (req.Amount*100)/req.Number <= 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "单个红包最低金额不能小于0.01元"})
		return
	}
	redPackageList, err := h.redPackageService.CreateRedPackage(c.Request.Context(), req.Account, reqAmount, req.Number)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, redPackageList)
}

func (h *RedPackageHandler) GetRedPackage(c *gin.Context) {
	redPackageid := c.Query("redPackageId")
	if redPackageid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "redPackageId is required"})
		return
	}
	useId := c.Query("userId")
	if useId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"redPackageId": redPackageid,
		"userId":       useId,
	})
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
