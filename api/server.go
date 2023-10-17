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
	jsonMarshal             func(v any) ([]byte, error) // Allows json.Marshal to be mocked
	githubAppId             string
	githubWebhookSecret     string
	githubPrivateKey        *rsa.PrivateKey
	githubFactory           githubFactoryInterface // Allows github operations to be mocked
	router                  *mux.Router
	KnownPullRequestActions map[string]struct{}
}

// NewServer creates a new HTTP server and sets up routing.
func NewServer(connPool *pgxpool.Pool, querier postgres.Querier) (*Server, error) {
	server := &Server{
		connPool:                connPool,
		querier:                 querier,
		jsonMarshal:             json.Marshal,
		githubAppId:             os.Getenv("GITHUB_APP_IDENTIFIER"),
		githubWebhookSecret:     os.Getenv("GITHUB_WEBHOOK_SECRET"),
		router:                  mux.NewRouter(),
		KnownPullRequestActions: make(map[string]struct{}),
	}

	server.githubFactory = NewGithubFactory(server)

	server.AddAllRoutes()

	privateKeyData := os.Getenv("GITHUB_PRIVATE_KEY")
	parsedPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKeyData))
	if err != nil {
		return &Server{}, fmt.Errorf("Failed to parse private key: %v", err)
	}
	server.githubPrivateKey = parsedPrivateKey

	actions, err := querier.GetPullRequestActions(context.Background())
	for _, action := range actions {
		server.KnownPullRequestActions[action] = struct{}{}
	}

	return server, nil
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

// ExecTx executes a function within a database transaction
func (server *Server) execTx(ctx context.Context, fn func(postgres.Querier) error) error {
	tx, err := server.connPool.Begin(ctx)
	if err != nil {
		return err
	}

	q := postgres.New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
