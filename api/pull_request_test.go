package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/jericop/pr-compliance-app/storage/postgres"
)

var validPullRequests []postgres.PullRequest

func TestGetPullRequests(t *testing.T) {
	urlPath := getRouteUrlPath(t, apiServer.Router, "GetPullRequests")

	// Requests to this http server will not show up in the api blueprint document.
	server := httptest.NewServer(apiServer.Router)
	defer server.Close()

	t.Run("StatusOK test2doc", func(t *testing.T) {
		expected := []postgres.PullRequest{
			{ID: 1, PrID: 991, RepoID: 444, PrNumber: 1, OpenedBy: 651, IsMerged: false},
		}

		fakeStore.GetPullRequestsCall.Returns.Error = nil
		fakeStore.GetPullRequestsCall.Returns.PullRequestSlice = expected

		resp := makeHttpRequest(t, http.StatusOK, func() (resp *http.Response, err error) {
			return http.Get(test2docServer.URL + urlPath)
		})

		decoder := json.NewDecoder(resp.Body)
		defer resp.Body.Close()

		result := []postgres.PullRequest{}
		if err := decoder.Decode(&result); err != nil {
			t.Fatalf("expected 'err' (%v) be nil", err)
		}

		if !reflect.DeepEqual(result, expected) {
			t.Fatalf("expected 'result' (%v) to equal 'expected' (%v)", result, expected)
		}

		if fakeStore.GetPullRequestsCall.CallCount != 1 {
			t.Errorf("unexpected call count: %d\n", fakeStore.GetPullRequestsCall.CallCount)
		}
	})

	t.Run("StatusInternalServerError json marshal error", func(t *testing.T) {
		apiServer.jsonMarshal = func(v interface{}) ([]byte, error) {
			return []byte{}, fmt.Errorf("Marshalling failed")
		}

		_ = makeHttpRequest(t, http.StatusInternalServerError, func() (resp *http.Response, err error) {
			return http.Get(server.URL + urlPath)
		})

		apiServer.jsonMarshal = json.Marshal

		if fakeStore.GetPullRequestsCall.CallCount != 2 {
			t.Errorf("unexpected call count: %d\n", fakeStore.GetPullRequestsCall.CallCount)
		}
	})

	t.Run("StatusInternalServerError store", func(t *testing.T) {
		// Test an error from the database
		fakeStore.GetPullRequestsCall.Returns.Error = fmt.Errorf("db error")
		fakeStore.GetPullRequestsCall.Returns.PullRequestSlice = []postgres.PullRequest{}

		_ = makeHttpRequest(t, http.StatusInternalServerError, func() (resp *http.Response, err error) {
			return http.Get(server.URL + urlPath)
		})

		if fakeStore.GetPullRequestsCall.CallCount != 3 {
			t.Errorf("unexpected call count: %d\n", fakeStore.GetPullRequestsCall.CallCount)
		}
	})
}
