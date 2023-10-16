package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/jericop/pr-compliance-app/storage/postgres"
)

func TestGetApproval(t *testing.T) {
	// Requests to this http server will not show up in the api blueprint document.
	server := httptest.NewServer(apiServer.router)
	defer server.Close()

	approvalId := uuid.New().String()
	getTests := []struct {
		name string
		url  string
	}{
		{
			name: "GetApproval",
			url:  test2docServer.URL + getRouteUrlPath(t, apiServer.router, "GetApproval", "id", approvalId),
		},
		{
			name: "GetApprovalQueryParam",
			url:  test2docServer.URL + getRouteUrlPath(t, apiServer.router, "GetApprovalQueryParam") + fmt.Sprintf("?id=%s", approvalId),
		},
	}

	expected := postgres.Approval{
		ID:   1,
		Uuid: approvalId,
		PrID: 555,
		Sha:  "038d718da6a1ebbc6a7780a96ed75a70cc2ad6e2", // echo testing | git hash-object --stdin -w
	}
	fakeStore.GetApprovalByUuidCall.Returns.Approval = expected
	fakeStore.GetApprovalByUuidCall.Returns.Error = nil

	for _, test := range getTests {
		t.Run(fmt.Sprintf("StatusOK test2doc %s", test.name), func(t *testing.T) {

			resp := makeHttpRequest(t, http.StatusOK, func() (resp *http.Response, err error) {
				return http.Get(test.url)
			})

			decoder := json.NewDecoder(resp.Body)
			defer resp.Body.Close()

			var result postgres.Approval
			if err := decoder.Decode(&result); err != nil {
				t.Fatalf("expected 'err' (%v) be nil", err)
			}

			if !reflect.DeepEqual(result, expected) {
				t.Fatalf("expected 'result' (%v) to equal 'expected' (%v)", result, expected)
			}
		})
	}

	t.Run("StatusInternalServerError url missing query param", func(t *testing.T) {
		_ = makeHttpRequest(t, http.StatusInternalServerError, func() (resp *http.Response, err error) {
			return http.Get(server.URL + getRouteUrlPath(t, apiServer.router, "GetApprovalQueryParam"))
		})
	})

	t.Run("StatusInternalServerError json marshal error", func(t *testing.T) {
		apiServer.jsonMarshal = func(v interface{}) ([]byte, error) {
			return []byte{}, fmt.Errorf("Marshalling failed")
		}

		_ = makeHttpRequest(t, http.StatusInternalServerError, func() (resp *http.Response, err error) {
			return http.Get(server.URL + getRouteUrlPath(t, apiServer.router, "GetApproval", "id", approvalId))
		})

		apiServer.jsonMarshal = json.Marshal
	})

	t.Run("StatusInternalServerError querier error", func(t *testing.T) {
		// Test an error from the database
		fakeStore.GetApprovalByUuidCall.Returns.Error = fmt.Errorf("db error")
		fakeStore.GetApprovalByUuidCall.Returns.Approval = postgres.Approval{}

		_ = makeHttpRequest(t, http.StatusInternalServerError, func() (resp *http.Response, err error) {
			return http.Get(server.URL + getRouteUrlPath(t, apiServer.router, "GetApproval", "id", approvalId))
		})
	})

}

func TestUpdateApproval(t *testing.T) {
	// Requests to this http server will not show up in the api blueprint document.
	server := httptest.NewServer(apiServer.router)
	defer server.Close()

	approvalId := uuid.New().String()
	urlPath := getRouteUrlPath(t, apiServer.router, "UpdateApproval")

	p := postgres.UpdateApprovalByUuidParams{
		Uuid:       approvalId,
		IsApproved: true,
	}
	expected := p // The response may contain different json data than POST request body

	pJson, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("expected 'err' (%v) be nil", err)
	}

	t.Run("StatusOK test2doc", func(t *testing.T) {
		fakeStore.GetApprovalByUuidCall.Returns.Error = nil

		buf := bytes.NewBuffer(pJson)

		resp := makeHttpRequest(t, http.StatusCreated, func() (resp *http.Response, err error) {
			return http.Post(test2docServer.URL+urlPath, "application/json", buf)
		})

		decoder := json.NewDecoder(resp.Body)
		defer resp.Body.Close()

		var result postgres.UpdateApprovalByUuidParams
		if err := decoder.Decode(&result); err != nil {
			t.Fatalf("expected 'err' (%v) be nil", err)
		}

		if !reflect.DeepEqual(result, expected) {
			t.Fatalf("expected 'result' (%v) to equal 'expected' (%v)", result, expected)
		}

	})

	t.Run("StatusBadRequest body json decode error", func(t *testing.T) {
		buf := bytes.NewBuffer([]byte("not valid json"))

		_ = makeHttpRequest(t, http.StatusBadRequest, func() (resp *http.Response, err error) {
			return http.Post(server.URL+urlPath, "application/json", buf)
		})

	})

	t.Run("StatusInternalServerError json marshal error", func(t *testing.T) {
		buf := bytes.NewBuffer(pJson)

		apiServer.jsonMarshal = func(v interface{}) ([]byte, error) {
			return []byte{}, fmt.Errorf("Marshalling failed")
		}

		_ = makeHttpRequest(t, http.StatusInternalServerError, func() (resp *http.Response, err error) {
			return http.Post(server.URL+urlPath, "application/json", buf)
		})

		apiServer.jsonMarshal = json.Marshal
	})

	t.Run("StatusInternalServerError querier error", func(t *testing.T) {
		buf := bytes.NewBuffer(pJson)

		fakeStore.UpdateApprovalByUuidCall.Returns.Error = fmt.Errorf("db error")

		_ = makeHttpRequest(t, http.StatusInternalServerError, func() (resp *http.Response, err error) {
			return http.Post(server.URL+urlPath, "application/json", buf)
		})
	})

}
