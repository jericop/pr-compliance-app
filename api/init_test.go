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

func TestMain(m *testing.M) {
	var err error

	db := &fakes.Querier{}
	apiServer := NewServer(db)
	test.RegisterURLVarExtractor(vars.MakeGorillaMuxExtractor(apiServer.router))

	server, err = test.NewServer(apiServer.router)
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
