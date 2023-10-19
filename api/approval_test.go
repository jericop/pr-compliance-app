package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jericop/pr-compliance-app/fakes"
	"github.com/jericop/pr-compliance-app/storage/postgres"
	"github.com/migueleliasweb/go-github-mock/src/mock"
)

var (
	question1 string = "Are there any bugs in this code"
	question2 string = "Are you sure there are no bugs"
)

func TestGetApproval(t *testing.T) {
	// Local querier and api for testing failures
	querier := &fakes.Querier{}
	api := NewMockedApiServer(querier).WithRoutes()

	// Requests to this http server will not show up in the api blueprint document.
	server := httptest.NewServer(api.router)
	defer server.Close()

	questions := []postgres.GetSortedApprovalYesNoQuestionAnswersByUuidRow{
		{ID: 1, QuestionText: question1, AnsweredYes: false},
		{ID: 2, QuestionText: question2, AnsweredYes: false},
	}

	queriers := []*fakes.Querier{fakeQuerier, querier}

	// Set up successful call return values for both queriers
	for _, q := range queriers {
		q.GetSortedApprovalYesNoQuestionAnswersByUuidCall.Returns.GetSortedApprovalYesNoQuestionAnswersByUuidRowSlice = questions
	}

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

	for _, test := range getTests {
		t.Run(fmt.Sprintf("StatusOK test2doc %s", test.name), func(t *testing.T) {

			makeHttpRequest(t, http.StatusOK, func() (resp *http.Response, err error) {
				return http.Get(test.url)
			})
		})
	}

	t.Run("StatusInternalServerError url missing query param", func(t *testing.T) {
		makeHttpRequest(t, http.StatusInternalServerError, func() (resp *http.Response, err error) {
			return http.Get(server.URL + getRouteUrlPath(t, api.router, "GetApprovalQueryParam"))
		})
	})

	t.Run("StatusInternalServerError json marshal error", func(t *testing.T) {
		api.jsonMarshalFunc = func(v interface{}) ([]byte, error) {
			return []byte{}, fmt.Errorf("Marshalling failed")
		}

		_ = makeHttpRequest(t, http.StatusInternalServerError, func() (resp *http.Response, err error) {
			return http.Get(server.URL + getRouteUrlPath(t, api.router, "GetApproval", "id", approvalId))
		})

		api.jsonMarshalFunc = json.Marshal
	})
	t.Run("StatusInternalServerError querier.GetSortedApprovalYesNoQuestionAnswersByUuid error", func(t *testing.T) {
		querier.GetSortedApprovalYesNoQuestionAnswersByUuidCall.Returns.Error = fmt.Errorf("db error")

		_ = makeHttpRequest(t, http.StatusInternalServerError, func() (resp *http.Response, err error) {
			return http.Get(server.URL + getRouteUrlPath(t, api.router, "GetApproval", "id", approvalId))
		})

		api.jsonMarshalFunc = json.Marshal
	})
}

func TestUpdateApproval(t *testing.T) {
	// Local querier and api for testing failures
	querier := &fakes.Querier{}
	api := NewMockedApiServer(querier).WithRoutes()

	// Requests to this http server will not show up in the api blueprint document.
	server := httptest.NewServer(api.router)
	defer server.Close()

	urlPath := getRouteUrlPath(t, apiServer.router, "UpdateApproval")

	questions := []postgres.GetSortedApprovalYesNoQuestionAnswersByUuidRow{
		{ID: 1, QuestionText: question1, AnsweredYes: false},
		{ID: 2, QuestionText: question2, AnsweredYes: false},
	}

	queriers := []*fakes.Querier{fakeQuerier, querier}

	// Set up successful call return values for both queriers
	for _, q := range queriers {
		q.GetSortedApprovalYesNoQuestionAnswersByUuidCall.Returns.GetSortedApprovalYesNoQuestionAnswersByUuidRowSlice = questions

		q.GetCreateStatusInputsFromApprovalUuidCall.Returns.
			GetCreateStatusInputsFromApprovalUuidRow = postgres.GetCreateStatusInputsFromApprovalUuidRow{
			InstallationID: 54321,
			Login:          "some-user",
			Name:           "some-repo",                                // the name of the github repo
			Sha:            "038d718da6a1ebbc6a7780a96ed75a70cc2ad6e2", // echo testing | git hash-object --stdin -w
		}
	}

	apiServer.githubFactory = NewMockGithubClientFactory(apiServer)
	api.githubFactory = NewMockGithubClientFactory(apiServer)

	p := ApprovalYesNoQuestionAnswersResponse{
		Uuid:      testUuid,
		Questions: questions,
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
			name:        "StatusInternalServerError marshal error",
			contentType: "application/json",
			beforeFunc: func() {
				api.jsonMarshalFunc = func(v interface{}) ([]byte, error) {
					return []byte{}, fmt.Errorf("marshal error")
				}
			},
			ioReader: func() io.Reader {
				return bytes.NewBuffer(pJson)
			},
			afterFunc: func() {
				api.jsonMarshalFunc = json.Marshal
			},
			wantStatusCode: http.StatusInternalServerError,
			url:            server.URL + urlPath,
		},
		{
			name:        "StatusCreated querier.DeleteApprovalYesAnswerByUuid called",
			contentType: "application/json",
			beforeFunc: func() {
				querier.GetSortedApprovalYesNoQuestionAnswersByUuidCall.Returns.
					GetSortedApprovalYesNoQuestionAnswersByUuidRowSlice = []postgres.GetSortedApprovalYesNoQuestionAnswersByUuidRow{
					{ID: 1, QuestionText: question1, AnsweredYes: true},
				}
			},
			ioReader: func() io.Reader {
				questionsPayload := ApprovalYesNoQuestionAnswersResponse{
					Uuid: testUuid,
					Questions: []postgres.GetSortedApprovalYesNoQuestionAnswersByUuidRow{
						{ID: 1, QuestionText: question1, AnsweredYes: false},
					},
				}
				questionsJSON, err := json.Marshal(questionsPayload)
				if err != nil {
					t.Fatalf("expected 'err' (%v) be nil", err)
				}

				return bytes.NewBuffer(questionsJSON)
			},
			wantStatusCode: http.StatusCreated,
			url:            server.URL + urlPath,
		},
		{
			name:        "StatusInternalServerError querier.DeleteApprovalYesAnswerByUuid error",
			contentType: "application/json",
			beforeFunc: func() {
				querier.DeleteApprovalYesAnswerByUuidCall.Returns.Error = fmt.Errorf("db error")
			},
			ioReader: func() io.Reader {
				questionsPayload := ApprovalYesNoQuestionAnswersResponse{
					Uuid: testUuid,
					Questions: []postgres.GetSortedApprovalYesNoQuestionAnswersByUuidRow{
						{ID: 1, QuestionText: question1, AnsweredYes: false},
					},
				}
				questionsJSON, err := json.Marshal(questionsPayload)
				if err != nil {
					t.Fatalf("expected 'err' (%v) be nil", err)
				}

				return bytes.NewBuffer(questionsJSON)
			},
			afterFunc: func() {
				querier.DeleteApprovalYesAnswerByUuidCall.Returns.Error = nil
			},
			wantStatusCode: http.StatusInternalServerError,
			url:            server.URL + urlPath,
		},

		{
			name:        "StatusCreated querier.CreateApprovalYesAnswerByUuid called",
			contentType: "application/json",
			beforeFunc: func() {
				querier.GetSortedApprovalYesNoQuestionAnswersByUuidCall.Returns.
					GetSortedApprovalYesNoQuestionAnswersByUuidRowSlice = []postgres.GetSortedApprovalYesNoQuestionAnswersByUuidRow{
					{ID: 1, QuestionText: question1, AnsweredYes: false},
				}
			},
			ioReader: func() io.Reader {
				questionsPayload := ApprovalYesNoQuestionAnswersResponse{
					Uuid: testUuid,
					Questions: []postgres.GetSortedApprovalYesNoQuestionAnswersByUuidRow{
						{ID: 1, QuestionText: question1, AnsweredYes: true},
					},
				}
				questionsJSON, err := json.Marshal(questionsPayload)
				if err != nil {
					t.Fatalf("expected 'err' (%v) be nil", err)
				}

				return bytes.NewBuffer(questionsJSON)
			},
			wantStatusCode: http.StatusCreated,
			url:            server.URL + urlPath,
		},
		{
			name:        "StatusInternalServerError querier.CreateApprovalYesAnswerByUuid error",
			contentType: "application/json",
			beforeFunc: func() {
				querier.CreateApprovalYesAnswerByUuidCall.Returns.Error = fmt.Errorf("db error")
			},
			ioReader: func() io.Reader {
				questionsPayload := ApprovalYesNoQuestionAnswersResponse{
					Uuid: testUuid,
					Questions: []postgres.GetSortedApprovalYesNoQuestionAnswersByUuidRow{
						{ID: 1, QuestionText: question1, AnsweredYes: true},
					},
				}
				questionsJSON, err := json.Marshal(questionsPayload)
				if err != nil {
					t.Fatalf("expected 'err' (%v) be nil", err)
				}

				return bytes.NewBuffer(questionsJSON)
			},
			afterFunc: func() {
				querier.CreateApprovalYesAnswerByUuidCall.Returns.Error = nil
				querier.GetSortedApprovalYesNoQuestionAnswersByUuidCall.Returns.GetSortedApprovalYesNoQuestionAnswersByUuidRowSlice = questions
			},
			wantStatusCode: http.StatusInternalServerError,
			url:            server.URL + urlPath,
		},

		{
			name:        "StatusBadRequest test2doc",
			contentType: "application/x-www-form-urlencoded",
			ioReader: func() io.Reader {
				return strings.NewReader("")
			},
			wantStatusCode: http.StatusBadRequest,
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

				api.githubFactory = NewMockGithubClientFactory(apiServer).
					WithNewInstallationClientReturns(badClient, nil)
			},
			ioReader: func() io.Reader {
				return bytes.NewBuffer(pJson)
			},
			wantStatusCode: http.StatusInternalServerError,
			url:            server.URL + urlPath,
		},
		{
			name:        "StatusInternalServerError get installation client error",
			contentType: "application/json",
			beforeFunc: func() {
				api.githubFactory = NewMockGithubClientFactory(apiServer).
					WithNewInstallationClientReturns(&http.Client{}, fmt.Errorf("github client error"))
			},
			ioReader: func() io.Reader {
				return bytes.NewBuffer(pJson)
			},
			wantStatusCode: http.StatusInternalServerError,
			url:            server.URL + urlPath,
		},
		{
			name:        "StatusInternalServerError querier.GetCreateStatusInputsFromApprovalUuid error",
			contentType: "application/json",
			beforeFunc: func() {
				querier.GetCreateStatusInputsFromApprovalUuidCall.Returns.Error = fmt.Errorf("db error")
			},
			ioReader: func() io.Reader {
				return bytes.NewBuffer(pJson)
			},
			afterFunc: func() {
				querier.GetCreateStatusInputsFromApprovalUuidCall.Returns.Error = nil
			},
			wantStatusCode: http.StatusInternalServerError,
			url:            server.URL + urlPath,
		},
		{
			name:        "StatusInternalServerError querier.UpdateApprovalByUuid error",
			contentType: "application/json",
			beforeFunc: func() {
				querier.UpdateApprovalByUuidCall.Returns.Error = fmt.Errorf("db error")
			},
			ioReader: func() io.Reader {
				return bytes.NewBuffer(pJson)
			},
			afterFunc: func() {
				querier.UpdateApprovalByUuidCall.Returns.Error = nil
			},
			wantStatusCode: http.StatusInternalServerError,
			url:            server.URL + urlPath,
		},
		{
			name:        "StatusInternalServerError querier.GetSortedApprovalYesNoQuestionAnswersByUuid error",
			contentType: "application/json",
			beforeFunc: func() {
				querier.GetSortedApprovalYesNoQuestionAnswersByUuidCall.Returns.Error = fmt.Errorf("db error")
			},
			ioReader: func() io.Reader {
				return bytes.NewBuffer(pJson)
			},
			afterFunc: func() {
				querier.GetSortedApprovalYesNoQuestionAnswersByUuidCall.Returns.Error = nil
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
