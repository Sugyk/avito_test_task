package handlers

import (
	"database/sql"

	"github.com/Sugyk/avito_test_task/internal/service"
)

type Service interface {
}

type Handler struct {
	service Service
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		service: service.NewService(db),
	}
}
