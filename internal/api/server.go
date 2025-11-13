package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Sugyk/avito_test_task/internal/api/handlers"
)

var (
	teamGetPostfix          = "/team/get"
	usersSetIsActivePostfix = "/users/setIsActive"
	prCreatePostfix         = "/pullRequest/create"
	prMergePostfix          = "/pullRequest/merge"
	prReassignPostfix       = "/pullRequest/reassign"
	usersGetReviewPostfix   = "/users/getReview"
)

type Router struct {
	server   *http.Server
	handlers *handlers.Handler
}

func NewRouter(port string, handler *handlers.Handler) *Router {
	mux := http.NewServeMux()

	// mux.HandleFunc("POST /team/add", nil)

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
