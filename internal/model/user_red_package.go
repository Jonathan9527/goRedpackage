package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type UserRedPackages struct {
	ID           int64           `gorm:"primaryKey" json:"id"`
	UserID       int64           `gorm:"not null;index" json:"user_id"`
	Type         string          `gorm:"size:64;not null" json:"type"`
	Amount       decimal.Decimal `gorm:"type:numeric(12,2);default:0" json:"amount"`
	Number       int32           `gorm:"not null;index" json:"number"`
	AfterBalance decimal.Decimal `gorm:"type:numeric(12,2);default:0" json:"after_balance"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

func (UserRedPackages) TableName() string {
	return "user_red_packages"
}
