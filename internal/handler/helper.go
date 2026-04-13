package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

func parseIntQuery(c *gin.Context, key string, defaultValue int) int {
	valueStr := c.Query(key)
	if valueStr == "" {
		return defaultValue
	}
	var value int
	_, err := fmt.Sscanf(valueStr, "%d", &value)
	if err != nil {
		return defaultValue
	}
	return value
}

func toYuanDecimal(fen int64) decimal.Decimal {
	return decimal.NewFromInt(fen).Div(decimal.NewFromInt(100))
}
