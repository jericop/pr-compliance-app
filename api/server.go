package api

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jericop/pr-compliance-app/storage/postgres"
)

// Server serves HTTP requests for the api and uses store to interact with the Querier interface created by sqlc.
type Server struct {
	connPool                *pgxpool.Pool
	querier                 postgres.Querier
	jsonMarshalFunc         func(v any) ([]byte, error) // Allows json.Marshal to be mocked
	frontEndUrl             string
	githubAppId             string
	githubWebhookSecret     string
	githubPrivateKey        *rsa.PrivateKey
	githubFactory           githubFactoryInterface // Allows github operations to be mocked
	router                  *mux.Router
	dbTxFactory             postgres.DatabaseTransactionFactory // Allows db transactions to be mocked
	knownPullRequestActions map[string]struct{}
	schema                  postgres.ApprovalSchema
}

func NewServer(connPool *pgxpool.Pool, querier postgres.Querier) (*Server, error) {
	ctx := context.Background()
	server := &Server{
		connPool:                connPool,
		querier:                 querier,
		jsonMarshalFunc:         json.Marshal,
		frontEndUrl:             os.Getenv("APP_FRONTEND_URL"),
		githubAppId:             os.Getenv("GITHUB_APP_IDENTIFIER"),
		githubWebhookSecret:     os.Getenv("GITHUB_WEBHOOK_SECRET"),
		router:                  mux.NewRouter(),
		knownPullRequestActions: make(map[string]struct{}),
	}

	server.githubFactory = NewGithubFactory(server)
	server.dbTxFactory = postgres.NewPostgresTransactionFactory(connPool)

	server.AddAllRoutes()

	privateKeyData := os.Getenv("GITHUB_PRIVATE_KEY")
	parsedPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKeyData))
	if err != nil {
		return &Server{}, fmt.Errorf("Failed to parse private key: %v", err)
	}
	server.githubPrivateKey = parsedPrivateKey

	actions, err := querier.GetPullRequestActions(ctx)
	if err != nil {
		return &Server{}, err
	}
	for _, action := range actions {
		server.knownPullRequestActions[action] = struct{}{}
	}

	schema, err := querier.GetDefaultApprovalSchema(ctx)
	if err != nil {
		return &Server{}, err
	}
	server.schema = schema

	return server, nil
}

func (server *Server) GetQuerier() *postgres.Queries {
	return postgres.New(server.connPool)
}

func (server *Server) GetRouter() *mux.Router {
	return server.router
}

// Add all routes here
func (server *Server) AddAllRoutes() {
	server.AddHealthRoutes()
	server.AddWebhookEventsRoutes()
	server.AddApprovalRoutes()
}

func (server *Server) WithRoutes() *Server {
	server.AddAllRoutes()
	return server
}

func (server *Server) WithPrivateKey(key *rsa.PrivateKey) *Server {
	server.githubPrivateKey = key
	return server
}
