package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jericop/pr-compliance-app/storage/postgres"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	store  postgres.Querier
	router *mux.Router
}

// NewServer creates a new HTTP server and sets up routing.
func NewServer(store postgres.Querier) *Server {
	server := &Server{
		store:  store,
		router: mux.NewRouter(),
	}

	server.AddPullRequestRoutes()

	return server
}

func (server *Server) Start(address string) error {
	return http.ListenAndServe(address, server.router)
}
