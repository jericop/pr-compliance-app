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
}

func NewServer(connPool *pgxpool.Pool, querier postgres.Querier) (*Server, error) {
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

	actions, err := querier.GetPullRequestActions(context.Background())
	for _, action := range actions {
		server.knownPullRequestActions[action] = struct{}{}
	}

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

// ExecWithTx executes a function within a database transaction
func (server *Server) ExecWithTx(ctx context.Context, fn func(postgres.Querier) error) error {
	tx, err := server.connPool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("tx begin err: %v", err)
	}

	q := postgres.New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	err = tx.Commit(ctx)
	return err
}
