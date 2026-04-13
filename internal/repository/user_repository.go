package repository

import (
	"context"
	"database/sql"
	"time"

	"learnGO/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) FindByAccount(ctx context.Context, account string) (model.User, error) {
	var user model.User
	var createdAt time.Time
	var updatedAt time.Time

	err := r.db.QueryRowContext(ctx, `
		SELECT id, account, username, balance::text, created_at, updated_at
		FROM users
		WHERE account = $1
	`, account).Scan(
		&user.ID,
		&user.Account,
		&user.Username,
		&user.Balance,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return model.User{}, err
	}

	user.CreatedAt = createdAt.Format(time.RFC3339)
	user.UpdatedAt = updatedAt.Format(time.RFC3339)
	return user, nil
}

func (r *UserRepository) List(ctx context.Context, limit int, offset int) ([]model.User, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, account, username, balance::text, created_at, updated_at
		FROM users
		ORDER BY id
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]model.User, 0, limit)
	for rows.Next() {
		var user model.User
		var createdAt time.Time
		var updatedAt time.Time

		if err := rows.Scan(
			&user.ID,
			&user.Account,
			&user.Username,
			&user.Balance,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, err
		}

		user.CreatedAt = createdAt.Format(time.RFC3339)
		user.UpdatedAt = updatedAt.Format(time.RFC3339)
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
