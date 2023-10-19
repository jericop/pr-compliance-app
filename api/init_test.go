package api

import (
	"crypto/rand"
	"crypto/rsa"
	"os"
	"testing"

	"github.com/jericop/pr-compliance-app/fakes"

	"github.com/s-mang/test2doc/test"
	"github.com/s-mang/test2doc/vars"
)

var test2docServer *test.Server
var apiServer *Server
var fakeQuerier *fakes.Querier
var testPrivateKey *rsa.PrivateKey

func TestMain(m *testing.M) {
	var err error
	fakeQuerier = &fakes.Querier{}

	// Generate RSA key.
	testPrivateKey, err = rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err.Error())
	}

	apiServer = NewMockedApiServer(fakeQuerier).
		WithRoutes().
		WithPrivateKey(testPrivateKey)

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
