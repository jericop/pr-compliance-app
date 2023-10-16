package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jericop/pr-compliance-app/fakes"
	"github.com/jericop/pr-compliance-app/storage/postgres"
)

func TestGetApproval(t *testing.T) {
	t.Skip()
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
		makeHttpRequest(t, http.StatusInternalServerError, func() (resp *http.Response, err error) {
			return http.Get(server.URL + getRouteUrlPath(t, apiServer.router, "GetApprovalQueryParam"))
		})
	})

	t.Run("StatusInternalServerError json marshal error", func(t *testing.T) {
		apiServer.jsonMarshal = func(v interface{}) ([]byte, error) {
			return []byte{}, fmt.Errorf("Marshalling failed")
		}

		makeHttpRequest(t, http.StatusInternalServerError, func() (resp *http.Response, err error) {
			return http.Get(server.URL + getRouteUrlPath(t, apiServer.router, "GetApproval", "id", approvalId))
		})

		apiServer.jsonMarshal = json.Marshal
	})

	t.Run("StatusInternalServerError querier error", func(t *testing.T) {
		// Test an error from the database
		fakeStore.GetApprovalByUuidCall.Returns.Error = fmt.Errorf("db error")
		fakeStore.GetApprovalByUuidCall.Returns.Approval = postgres.Approval{}

		makeHttpRequest(t, http.StatusInternalServerError, func() (resp *http.Response, err error) {
			return http.Get(server.URL + getRouteUrlPath(t, apiServer.router, "GetApproval", "id", approvalId))
		})
	})

}

func TestUpdateApproval(t *testing.T) {

	api := getApiServer(&fakes.Querier{}).WithRoutes()
	// Requests to this http server will not show up in the api blueprint document.
	server := httptest.NewServer(api.router)
	defer server.Close()

	approvalId := uuid.New().String()
	urlPath := getRouteUrlPath(t, apiServer.router, "UpdateApproval")

	p := postgres.UpdateApprovalByUuidParams{
		Uuid:       approvalId,
		IsApproved: true,
	}

	pJson, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("expected 'err' (%v) be nil", err)
	}

	tests := []struct {
		contentType    string
		ioReader       func() io.Reader
		beforeFunc     func()
		url            string
		wantError      bool
		wantStatusCode int
	}{
		{
			contentType: "application/json",
			ioReader: func() io.Reader {
				return bytes.NewBuffer(pJson)
			},
			wantStatusCode: http.StatusCreated,
			url:            test2docServer.URL + urlPath,
		},
		{
			contentType: "application/x-www-form-urlencoded",
			ioReader: func() io.Reader {
				formData := url.Values{}
				formData.Set("uuid", approvalId)
				formData.Set("is_approved", "true")
				return strings.NewReader(formData.Encode())
			},
			wantStatusCode: http.StatusCreated,
			url:            test2docServer.URL + urlPath,
		},
	}

	fakeStore.GetApprovalByUuidCall.Returns.Error = nil
	fakeStore.GetCreateStatusInputsFromApprovalUuidCall.Returns.GetCreateStatusInputsFromApprovalUuidRow = postgres.GetCreateStatusInputsFromApprovalUuidRow{
		InstallationID: 54321,
		Login:          "some-user",
		Name:           "some-repo",                                // the name of the github repo
		Sha:            "038d718da6a1ebbc6a7780a96ed75a70cc2ad6e2", // echo testing | git hash-object --stdin -w
	}

	apiServer.githubFactory = NewMockGithubClientFactory(apiServer)
	log.Printf("apiServer.githubFactory %v", apiServer.githubFactory)

	for _, test := range tests {
		t.Run(fmt.Sprintf("StatusOK test2doc %v", test.contentType), func(t *testing.T) {
			if test.beforeFunc != nil {
				log.Printf("Calling beforeFunc")
				test.beforeFunc()
			}
			makeHttpRequest(t, test.wantStatusCode, func() (resp *http.Response, err error) {
				return http.Post(test.url, test.contentType, test.ioReader())
			})
		})
	}

	// t.Run("StatusOK test2doc", func(t *testing.T) {
	// 	fakeStore.GetApprovalByUuidCall.Returns.Error = nil
	// 	fakeStore.GetCreateStatusInputsFromApprovalUuidCall.Returns.GetCreateStatusInputsFromApprovalUuidRow = postgres.GetCreateStatusInputsFromApprovalUuidRow{
	// 		InstallationID: 54321,
	// 		Login:          "some-user",
	// 		Name:           "some-repo",                                // the name of the github repo
	// 		Sha:            "038d718da6a1ebbc6a7780a96ed75a70cc2ad6e2", // echo testing | git hash-object --stdin -w
	// 	}

	// 	apiServer.githubFactory = NewMockGithubClientFactory(apiServer)
	// 	log.Printf("apiServer.githubFactory %v", apiServer.githubFactory)

	// 	buf := bytes.NewBuffer(pJson)

	// 	resp := makeHttpRequest(t, http.StatusCreated, func() (resp *http.Response, err error) {
	// 		return http.Post(test2docServer.URL+urlPath, "application/json", buf)
	// 	})

	// 	decoder := json.NewDecoder(resp.Body)
	// 	defer resp.Body.Close()

	// 	var result postgres.UpdateApprovalByUuidParams
	// 	if err := decoder.Decode(&result); err != nil {
	// 		t.Fatalf("expected 'err' (%v) be nil", err)
	// 	}

	// 	if !reflect.DeepEqual(result, expected) {
	// 		t.Fatalf("expected 'result' (%v) to equal 'expected' (%v)", result, expected)
	// 	}

	// })

	t.Run("StatusBadRequest body json decode error", func(t *testing.T) {
		buf := bytes.NewBuffer([]byte("not valid json"))

		makeHttpRequest(t, http.StatusBadRequest, func() (resp *http.Response, err error) {
			return http.Post(server.URL+urlPath, "application/json", buf)
		})

	})

	t.Run("StatusInternalServerError json marshal error", func(t *testing.T) {
		buf := bytes.NewBuffer(pJson)

		apiServer.jsonMarshal = func(v interface{}) ([]byte, error) {
			return []byte{}, fmt.Errorf("Marshalling failed")
		}

		makeHttpRequest(t, http.StatusInternalServerError, func() (resp *http.Response, err error) {
			return http.Post(server.URL+urlPath, "application/json", buf)
		})

		apiServer.jsonMarshal = json.Marshal
	})

	t.Run("StatusInternalServerError querier error", func(t *testing.T) {
		buf := bytes.NewBuffer(pJson)

		fakeStore.UpdateApprovalByUuidCall.Returns.Error = fmt.Errorf("db error")

		makeHttpRequest(t, http.StatusInternalServerError, func() (resp *http.Response, err error) {
			return http.Post(server.URL+urlPath, "application/json", buf)
		})
	})

}
