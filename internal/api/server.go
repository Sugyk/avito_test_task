package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Sugyk/avito_test_task/internal/api/handlers"
)

type Router struct {
	server   *http.Server
	handlers *handlers.Handler
}

func NewRouter(port string, handler *handlers.Handler) *Router {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /team/add", handler.TeamAdd)
	mux.HandleFunc("GET /team/get", handler.TeamGet)
	mux.HandleFunc("POST /users/setIsActive", handler.UsersSetIsActive)
	mux.HandleFunc("POST /pullRequest/create", handler.PullRequestCreate)
	mux.HandleFunc("POST /pullRequest/merge", handler.PullRequestMerge)
	mux.HandleFunc("POST /pullRequest/reassign", handler.PullRequestReassign)
	mux.HandleFunc("GET /users/getReview", handler.UsersGetReview)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	return &Router{
		server:   server,
		handlers: handler,
	}
}

func (r *Router) Start() error {
	return r.server.ListenAndServe()
}

func (r *Router) Shutdown(ctx context.Context) error {
	if err := r.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}
	return nil
}
