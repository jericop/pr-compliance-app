package api

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jericop/pr-compliance-app/storage/postgres"
)

// Server serves HTTP requests for the api and uses store to interact with the Querier interface created by sqlc.
type Server struct {
	store         postgres.Querier
	jsonMarshal   func(v any) ([]byte, error)
	webhookSecret string
	Router        *mux.Router
}

// NewServer creates a new HTTP server and sets up routing.
func NewServer(store postgres.Querier) *Server {
	server := &Server{
		store:         store,
		jsonMarshal:   json.Marshal,
		webhookSecret: os.Getenv("GITHUB_WEBHOOK_SECRET"),
		Router:        mux.NewRouter(),
	}

	server.AddAllRoutes()

	return server
}

// Add all routes here
func (server *Server) AddAllRoutes() {
	server.AddPullRequestRoutes()
}

func (server *Server) Start(address string) error {
	return http.ListenAndServe(address, server.Router)
}
