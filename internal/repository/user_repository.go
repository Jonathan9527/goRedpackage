package repository

import (
	"context"
	"errors"
	"time"

	"learnGO/internal/model"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) FindByAccount(ctx context.Context, account string) (model.User, error) {
	var record model.UserRecord
	if err := r.db.WithContext(ctx).
		Where("account = ?", account).
		First(&record).Error; err != nil {
		return model.User{}, err
	}

	return toUserModel(record), nil
}

func (r *UserRepository) List(ctx context.Context, limit int, offset int) ([]model.User, error) {
	var records []model.UserRecord
	if err := r.db.WithContext(ctx).
		Order("id ASC").
		Limit(limit).
		Offset(offset).
		Find(&records).Error; err != nil {
		return nil, err
	}

	users := make([]model.User, 0, limit)
	for _, record := range records {
		users = append(users, toUserModel(record))
	}

	return users, nil
}

func (r *UserRepository) UpdateBalance(ctx context.Context, u model.User, redPackageAmount decimal.Decimal) error {
	newBalance := u.Balance.Sub(redPackageAmount)
	tx := r.db.WithContext(ctx).Begin()
	res := tx.Model(&model.UserRecord{}).
		Where(&model.UserRecord{Account: u.Account, Balance: u.Balance}).
		Updates(map[string]interface{}{
			"balance":    newBalance,
			"updated_at": time.Now(),
		})
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}
	if res.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("balance changed, update failed")
	}
	err := tx.Create(&model.UserTransactionRecord{
		UserID:        u.ID,
		Type:          "red_package",
		Amount:        redPackageAmount,
		BeforeBalance: u.Balance,
		AfterBalance:  newBalance,
		InOrOut:       0,
	}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func toUserModel(record model.UserRecord) model.User {
	return model.User{
		ID:        record.ID,
		Account:   record.Account,
		Username:  record.Username,
		Balance:   record.Balance,
		CreatedAt: record.CreatedAt.Format(time.RFC3339),
		UpdatedAt: record.UpdatedAt.Format(time.RFC3339),
	}
}
