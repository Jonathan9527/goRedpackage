package database

import (
	"context"
	"fmt"
	"time"

	"learnGO/internal/model"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.WithContext(ctx).AutoMigrate(&model.UserRecord{}, &model.UserTransactionRecord{}); err != nil {
		return fmt.Errorf("create users table: %w", err)
	}

	if err := db.WithContext(ctx).Exec(`
		INSERT INTO users (account, username, balance)
		SELECT
			'user' || LPAD(n::text, 4, '0'),
			'用户' || LPAD(n::text, 4, '0'),
			100.00
		FROM generate_series(1, 1000) AS n
		ON CONFLICT (account) DO NOTHING;
	`).Error; err != nil {
		return fmt.Errorf("seed users: %w", err)
	}

	return nil
}
