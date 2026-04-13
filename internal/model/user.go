package model

import "github.com/shopspring/decimal"

type User struct {
	ID        int64           `json:"id"`
	Account   string          `json:"account"`
	Username  string          `json:"username"`
	Balance   decimal.Decimal `json:"balance"`
	CreatedAt string          `json:"created_at"`
	UpdatedAt string          `json:"updated_at"`
}
