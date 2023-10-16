// endpoints.go
package api

import (
	"net/http"
	"testing"
)

func TestGetHealth(t *testing.T) {
	makeHttpRequest(t, http.StatusOK, func() (resp *http.Response, err error) {
		return http.Get(test2docServer.URL + getRouteUrlPath(t, apiServer.router, "GetHealth"))
	})
}
