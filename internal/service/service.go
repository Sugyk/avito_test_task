package service

import (
	"database/sql"

	"github.com/Sugyk/avito_test_task/internal/repository"
)

type Repository interface{}

type Service struct {
	repo Repository
}

func NewService(db *sql.DB) *Service {
	return &Service{
		repo: repository.NewRepository(db),
	}
}
