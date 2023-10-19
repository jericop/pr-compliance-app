package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jericop/pr-compliance-app/fakes"
	"github.com/jericop/pr-compliance-app/storage/postgres"
)

func NewMockedApiServer(querier *fakes.Querier) *Server {
	server := &Server{
		querier:                 querier,
		frontEndUrl:             "http://localhost:8080/approval",
		githubWebhookSecret:     "0123456789abcdef",
		jsonMarshalFunc:         json.Marshal,
		router:                  mux.NewRouter(),
		knownPullRequestActions: map[string]struct{}{},
		githubPrivateKey:        testPrivateKey,
	}

	server.githubFactory = NewMockGithubClientFactory(server).
		WithValidateWebhookRequestReturns(newPrEvent("opened").getEvent(), nil)

	server.dbTxFactory = NewMockPostgresTransactionFactory(querier) // No-op transaction mock

	return server
}

func getRouteUrlPath(t *testing.T, router *mux.Router, routeName string, pairs ...string) string {
	urlPath, err := router.Get(routeName).URL(pairs...)
	if err != nil {
		t.Fatalf("expected 'err' (%v) be nil", err)
	}
	return urlPath.String()
}

func makeHttpRequest(t *testing.T, expectedStatusCode int, httpRequestFunc func() (resp *http.Response, err error)) *http.Response {
	resp, err := httpRequestFunc() // http.Get, http.Post, etc.. functions get called here
	if err != nil {
		t.Fatalf("expected 'err' (%v) be nil", err)
	}

	if resp.StatusCode != expectedStatusCode {
		f, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("expected 'err' (%v) be nil", err)
		}
		resp.Body.Close()
		t.Fatalf("expected 'resp.StatusCode' (%v) to equal 'expectedStatusCode' (%v) resp.Body:\n%v", resp.StatusCode, expectedStatusCode, string(f))
	}

	return resp
}

type mockPostgresTransactionFactory struct {
	querier           postgres.Querier
	ExecWithTxReturns error
}

func NewMockPostgresTransactionFactory(q postgres.Querier) *mockPostgresTransactionFactory {
	return &mockPostgresTransactionFactory{querier: q}
}

func (p *mockPostgresTransactionFactory) ExecWithTx(ctx context.Context, fn func(postgres.Querier) error) error {
	return fn(p.querier)
}
