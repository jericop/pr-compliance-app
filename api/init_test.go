package api

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jericop/pr-compliance-app/fakes"

	"github.com/s-mang/test2doc/test"
	"github.com/s-mang/test2doc/vars"
)

var test2docServer *test.Server
var apiServer *Server
var fakeQuerier *fakes.Querier

func TestMain(m *testing.M) {
	var err error
	fakeQuerier = &fakes.Querier{}

	// Generate RSA key.
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err.Error())
	}

	apiServer = getApiServer(fakeQuerier).WithRoutes().WithPrivateKey(key)

	test.RegisterURLVarExtractor(vars.MakeGorillaMuxExtractor(apiServer.router))

	// Requests to this http server will show up in the api blueprint document.
	test2docServer, err = test.NewServer(apiServer.router)
	if err != nil {
		panic(err.Error())
	}

	code := m.Run()
	test2docServer.Finish()
	os.Exit(code)
}

func getApiServer(querier *fakes.Querier) *Server {
	// Generate RSA key.
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err.Error())
	}

	server := &Server{
		querier:                 querier,
		githubWebhookSecret:     "0123456789abcdef",
		jsonMarshal:             json.Marshal,
		router:                  mux.NewRouter(),
		KnownPullRequestActions: map[string]struct{}{},
		githubPrivateKey:        key,
	}

	server.githubFactory = NewMockGithubClientFactory(apiServer)

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
