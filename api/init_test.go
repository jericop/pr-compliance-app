package api

import (
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jericop/pr-compliance-app/fakes"

	"github.com/s-mang/test2doc/test"
	"github.com/s-mang/test2doc/vars"
)

var router *mux.Router
var server *test.Server
var apiServer *Server
var fakeStore *fakes.Querier

func TestMain(m *testing.M) {
	var err error

	fakeStore = &fakes.Querier{}
	apiServer = NewServer(fakeStore)
	router = apiServer.router
	test.RegisterURLVarExtractor(vars.MakeGorillaMuxExtractor(router))

	server, err = test.NewServer(router)
	if err != nil {
		panic(err.Error())
	}
	code := m.Run()
	server.Finish()
	os.Exit(code)
}

func TestNothing(t *testing.T) {
	var result bool = true

	if !result {
		t.Fatal("something has gone wrong")
	}
}
