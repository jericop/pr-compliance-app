package api

import (
	"encoding/json"
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
var fakeStore *fakes.Querier

func TestMain(m *testing.M) {
	var err error
	fakeStore = &fakes.Querier{}
	apiServer = getApiServer(fakeStore)
	apiServer.AddAllRoutes()
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
	return &Server{
		querier:             querier,
		githubWebhookSecret: "0123456789abcdef",
		jsonMarshal:         json.Marshal,
		router:              mux.NewRouter(),
	}
}

func getQuerierServer() (*fakes.Querier, *Server) {
	querier := &fakes.Querier{}
	apiServer := getApiServer(querier)
	return querier, apiServer
}

func getQuerierServerWithRoutes() (*fakes.Querier, *Server) {
	querier, apiServer := getQuerierServer()
	apiServer.AddAllRoutes()
	return querier, apiServer
}

func getQuerierServerRouteUrl(t *testing.T, routeName string) (*fakes.Querier, *Server, string) {
	querier, apiServer := getQuerierServerWithRoutes()
	urlPath, err := apiServer.router.Get(routeName).URL()
	if err != nil {
		t.Fatalf("expected 'err' (%v) be nil", err)
	}
	return querier, apiServer, urlPath.String()
}

func getRouteUrlPath(t *testing.T, router *mux.Router, routeName string) string {
	urlPath, err := router.Get(routeName).URL()
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
		t.Fatalf("expected 'resp.StatusCode' (%v) to equal 'expectedStatusCode' (%v)", resp.StatusCode, http.StatusOK)
	}
	return resp
}
