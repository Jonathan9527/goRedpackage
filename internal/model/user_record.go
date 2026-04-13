package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type UserRecord struct {
	ID        int64           `gorm:"primaryKey" json:"id"`
	Account   string          `gorm:"size:64;uniqueIndex;not null" json:"account"`
	Username  string          `gorm:"size:64;not null" json:"username"`
	Balance   decimal.Decimal `gorm:"type:numeric(12,2);not null;default:0" json:"balance"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

func (UserRecord) TableName() string {
	return "users"
}
