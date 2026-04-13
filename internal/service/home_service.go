package service

import (
	"context"
	"time"

	"learnGO/internal/model"

	"gorm.io/gorm"
)

type HomeService struct {
	db *gorm.DB
}

func NewHomeService(db *gorm.DB) *HomeService {
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
	sqlDB, err := s.db.DB()
	if err != nil {
		databaseStatus = "error"
	} else if err := sqlDB.PingContext(ctx); err != nil {
		databaseStatus = "error"
	}

	return model.HealthResponse{
		Status:   "ok",
		Database: databaseStatus,
	}
}
