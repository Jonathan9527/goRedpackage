package service

import (
	"context"
	"database/sql"
	"time"

	"learnGO/internal/model"
)

type HomeService struct {
	db *sql.DB
}

func NewHomeService(db *sql.DB) *HomeService {
	return &HomeService{
		db: db,
	}
}

func (s *HomeService) Greeting() model.MessageResponse {
	return model.MessageResponse{
		Message: "Hello, Gin!",
	}
}

func (s *HomeService) Health() model.HealthResponse {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	databaseStatus := "ok"
	if err := s.db.PingContext(ctx); err != nil {
		databaseStatus = "error"
	}

	return model.HealthResponse{
		Status:   "ok",
		Database: databaseStatus,
	}
}
