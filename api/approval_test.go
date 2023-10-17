package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/jericop/pr-compliance-app/storage/postgres"
	"github.com/migueleliasweb/go-github-mock/src/mock"
)

func TestGetApproval(t *testing.T) {
	// Requests to this http server will not show up in the api blueprint document.
	server := httptest.NewServer(apiServer.router)
	defer server.Close()

	approvalId := testUuid
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

	fakeQuerier.GetApprovalByUuidCall.Returns.Approval = postgres.Approval{
		ID:   1,
		Uuid: approvalId,
		PrID: 555,
		Sha:  "038d718da6a1ebbc6a7780a96ed75a70cc2ad6e2", // echo testing | git hash-object --stdin -w
	}
	fakeQuerier.GetApprovalByUuidCall.Returns.Error = nil

	for _, test := range getTests {
		t.Run(fmt.Sprintf("StatusOK test2doc %s", test.name), func(t *testing.T) {

			makeHttpRequest(t, http.StatusOK, func() (resp *http.Response, err error) {
				return http.Get(test.url)
			})
		})
	}

	t.Run("StatusInternalServerError url missing query param", func(t *testing.T) {
		makeHttpRequest(t, http.StatusInternalServerError, func() (resp *http.Response, err error) {
			return http.Get(server.URL + getRouteUrlPath(t, apiServer.router, "GetApprovalQueryParam"))
		})
	})
}

func TestUpdateApproval(t *testing.T) {
	// Requests to this http server will not show up in the api blueprint document.
	server := httptest.NewServer(apiServer.router)
	defer server.Close()

	approvalId := testUuid
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
		name           string
		contentType    string
		ioReader       func() io.Reader
		beforeFunc     func()
		afterFunc      func()
		url            string
		wantStatusCode int
	}{
		{
			name:        "StatusCreated test2doc",
			contentType: "application/json",
			ioReader: func() io.Reader {
				return bytes.NewBuffer(pJson)
			},
			wantStatusCode: http.StatusCreated,
			url:            test2docServer.URL + urlPath,
		},
		{
			name:        "StatusCreated test2doc",
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
		{
			name:        "StatusInternalServerError marshal error",
			contentType: "application/json",
			beforeFunc: func() {
				apiServer.jsonMarshal = func(v interface{}) ([]byte, error) {
					return []byte{}, fmt.Errorf("marshal error")
				}
			},
			ioReader: func() io.Reader {
				return bytes.NewBuffer(pJson)
			},
			afterFunc: func() {
				apiServer.jsonMarshal = json.Marshal
			},
			wantStatusCode: http.StatusInternalServerError,
			url:            server.URL + urlPath,
		},

		{
			name:        "StatusBadRequest ParseForm error",
			contentType: "application/x-www-form-urlencoded",
			ioReader: func() io.Reader {
				return nil
			},
			wantStatusCode: http.StatusBadRequest,
			url:            server.URL + urlPath,
		},
		{
			name:        "StatusBadRequest invalid bool",
			contentType: "application/x-www-form-urlencoded",
			ioReader: func() io.Reader {
				formData := url.Values{}
				formData.Set("uuid", approvalId)
				formData.Set("is_approved", "invalidBool")
				return strings.NewReader(formData.Encode())
			},
			wantStatusCode: http.StatusBadRequest,
			url:            server.URL + urlPath,
		},
		{
			name:        "StatusInternalServerError querier.UpdateApprovalByUuid error",
			contentType: "application/json",
			beforeFunc: func() {
				fakeQuerier.UpdateApprovalByUuidCall.Returns.Error = fmt.Errorf("db error")
			},
			ioReader: func() io.Reader {
				return bytes.NewBuffer(pJson)
			},
			afterFunc: func() {
				fakeQuerier.UpdateApprovalByUuidCall.Returns.Error = nil
			},
			wantStatusCode: http.StatusInternalServerError,
			url:            server.URL + urlPath,
		},
		{
			name:        "StatusInternalServerError querier.GetCreateStatusInputsFromApprovalUuid error",
			contentType: "application/json",
			beforeFunc: func() {
				fakeQuerier.GetCreateStatusInputsFromApprovalUuidCall.Returns.Error = fmt.Errorf("db error")
			},
			ioReader: func() io.Reader {
				return bytes.NewBuffer(pJson)
			},
			afterFunc: func() {
				fakeQuerier.GetCreateStatusInputsFromApprovalUuidCall.Returns.Error = nil
			},
			wantStatusCode: http.StatusInternalServerError,
			url:            server.URL + urlPath,
		},
		{
			name:        "StatusInternalServerError get installation client error",
			contentType: "application/json",
			beforeFunc: func() {
				apiServer.githubFactory = NewMockGithubClientFactory(apiServer).
					WithNewInstallationClientReturns(&http.Client{}, fmt.Errorf("github client error"))
			},
			ioReader: func() io.Reader {
				return bytes.NewBuffer(pJson)
			},
			wantStatusCode: http.StatusInternalServerError,
			url:            server.URL + urlPath,
		},
		{
			name:        "StatusInternalServerError create status error",
			contentType: "application/json",
			beforeFunc: func() {
				badClient := mock.NewMockedHTTPClient(
					mock.WithRequestMatchHandler(
						mock.PostReposStatusesByOwnerByRepoBySha,
						http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
							mock.WriteError(
								w,
								http.StatusInternalServerError,
								"github went belly up or something",
							)
						}),
					),
				)

				apiServer.githubFactory = NewMockGithubClientFactory(apiServer).
					WithNewInstallationClientReturns(badClient, nil)
			},
			ioReader: func() io.Reader {
				return bytes.NewBuffer(pJson)
			},
			wantStatusCode: http.StatusInternalServerError,
			url:            server.URL + urlPath,
		},
		{
			name:        "StatusBadRequest invalid json",
			contentType: "application/json",
			ioReader: func() io.Reader {
				return bytes.NewBuffer([]byte("not valid json"))
			},
			wantStatusCode: http.StatusBadRequest,
			url:            server.URL + urlPath,
		},
	}

	fakeQuerier.GetApprovalByUuidCall.Returns.Error = nil
	fakeQuerier.GetCreateStatusInputsFromApprovalUuidCall.Returns.GetCreateStatusInputsFromApprovalUuidRow = postgres.GetCreateStatusInputsFromApprovalUuidRow{
		InstallationID: 54321,
		Login:          "some-user",
		Name:           "some-repo",                                // the name of the github repo
		Sha:            "038d718da6a1ebbc6a7780a96ed75a70cc2ad6e2", // echo testing | git hash-object --stdin -w
	}

	apiServer.githubFactory = NewMockGithubClientFactory(apiServer)

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v %v", test.name, test.contentType), func(t *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			makeHttpRequest(t, test.wantStatusCode, func() (resp *http.Response, err error) {
				return http.Post(test.url, test.contentType, test.ioReader())
			})

			if test.afterFunc != nil {
				test.afterFunc()
			}
		})
	}
}

type badReader struct{}

func (m *badReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("failure is my middle name")
}
