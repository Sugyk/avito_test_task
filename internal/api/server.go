package api

import (
	"log/slog"

	"github.com/Sugyk/avito_test_task/internal/api/handlers"
	"github.com/Sugyk/avito_test_task/pkg/database"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
}

func (s Server) Run() {}

func NewServer() (*Server, error) {
	router := gin.New()

	db, err := database.NewDbConnection()
	if err != nil {
		return nil, err
	}
	slog.Info("successful connected to db")

	handlers.Register(router, db)
	return &Server{
		router: router,
	}, nil
}
