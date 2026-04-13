package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type UserTransactionRecord struct {
	ID            int64           `gorm:"primaryKey" json:"id"`
	UserID        int64           `gorm:"not null;index" json:"user_id"`
	Type          string          `gorm:"size:64;not null" json:"type"`
	Amount        decimal.Decimal `gorm:"type:numeric(12,2);default:0" json:"amount"`
	BeforeBalance decimal.Decimal `gorm:"type:numeric(12,2);default:0" json:"before_balance"`
	AfterBalance  decimal.Decimal `gorm:"type:numeric(12,2);default:0" json:"after_balance"`
	InOrOut       int             `gorm:"size:16;not null" json:"in_or_out"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

func (UserTransactionRecord) TableName() string {
	return "user_transactions"
}
