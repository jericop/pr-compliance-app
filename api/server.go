package api

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jericop/pr-compliance-app/storage/postgres"
)

// Server serves HTTP requests for the api and uses store to interact with the Querier interface created by sqlc.
type Server struct {
	connPool            *pgxpool.Pool
	querier             postgres.Querier
	jsonMarshal         func(v any) ([]byte, error)
	githubAppId         string
	githubWebhookSecret string
	githubPrivateKey    *rsa.PrivateKey
	Router              *mux.Router
}

// NewServer creates a new HTTP server and sets up routing.
func NewServer(connPool *pgxpool.Pool, querier postgres.Querier) (*Server, error) {
	server := &Server{
		connPool:            connPool,
		querier:             querier,
		jsonMarshal:         json.Marshal,
		githubAppId:         os.Getenv("GITHUB_APP_IDENTIFIER"),
		githubWebhookSecret: os.Getenv("GITHUB_WEBHOOK_SECRET"),
		Router:              mux.NewRouter(),
	}

	server.AddAllRoutes()

	privateKeyData := os.Getenv("GITHUB_PRIVATE_KEY")
	parsedPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKeyData))
	if err != nil {
		return &Server{}, fmt.Errorf("Failed to parse private key: %v", err)
	}
	server.githubPrivateKey = parsedPrivateKey

	return server, nil
}

// Add all routes here
func (server *Server) AddAllRoutes() {
	server.AddWebhookEventsRoutes()
	server.AddPullRequestRoutes()
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

func (server *Server) Start(address string) error {
	return http.ListenAndServe(address, server.Router)
}
