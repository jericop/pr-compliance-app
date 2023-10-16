package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jericop/pr-compliance-app/storage/postgres"
)

// Server serves HTTP requests for the api and uses store to interact with the Querier interface created by sqlc.
type Server struct {
	store       postgres.Querier
	router      *mux.Router
	jsonMarshal func(v any) ([]byte, error)
}

// NewServer creates a new HTTP server and sets up routing.
func NewServer(store postgres.Querier) *Server {
	server := &Server{
		store:       store,
		router:      mux.NewRouter(),
		jsonMarshal: json.Marshal,
	}

	// Add all routes here
	server.AddPullRequestRoutes()

	return server
}

func (server *Server) Start(address string) error {
	return http.ListenAndServe(address, server.router)
}
