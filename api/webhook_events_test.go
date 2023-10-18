package api

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-github/v53/github"
	"github.com/jericop/pr-compliance-app/fakes"
	"github.com/jericop/pr-compliance-app/storage/postgres"
	"github.com/migueleliasweb/go-github-mock/src/mock"
)

const (
	testUuid           = "5fcae4d5-fbb6-417d-8267-70f9b0f6d28f"
	testInstallationID = 8675309
)

func TestPostWebhookEvent(t *testing.T) {

	// Requests to this http server will not show up in the api blueprint document.
	server := httptest.NewServer(apiServer.router)
	defer server.Close()

	tests := []struct {
		name        string
		payload     interface{}
		result      interface{}
		messageType string
	}{
		{
			name:        "pull_request opened no-org",
			payload:     newPrEvent("opened").getEvent(),
			result:      &github.PullRequestEvent{},
			messageType: "pull_request",
		},
		{
			name:        "pull_request opened with org",
			payload:     newPrEvent("opened").withOrg("some-org").getEvent(),
			result:      &github.PullRequestEvent{},
			messageType: "pull_request",
		},
	}

	urlPath := getRouteUrlPath(t, apiServer.router, "PostWebhookEvent")

	updateQuerierWithSuccessPostWebhookEventReturns(fakeQuerier)

	apiServer.githubFactory = NewMockGithubClientFactory(apiServer).
		WithValidateWebhookRequestReturns(newPrEvent("opened").getEvent(), nil)

	for n, test := range tests {
		testNamePrefix := ""
		serverURL := server.URL

		// Only use test2doc server for the first payload to keep the api documentation concise
		if n == 0 {
			testNamePrefix = "test2doc"
			serverURL = test2docServer.URL
		}

		t.Run(fmt.Sprintf("StatusOK %v %v", testNamePrefix, test.name), func(t *testing.T) {
			apiServer.githubFactory = NewMockGithubClientFactory(apiServer).
				WithValidateWebhookRequestReturns(newPrEvent("opened").getEvent(), nil)

			makeHttpRequest(t, http.StatusOK, func() (resp *http.Response, err error) {
				return http.Post(serverURL+urlPath, "application/json", bytes.NewBufferString("{}"))
			})
		})
	}

	t.Run("StatusBadRequest validate webhook event error", func(t *testing.T) {
		apiServer.githubFactory = NewMockGithubClientFactory(apiServer).
			WithValidateWebhookRequestReturns(struct{}{}, fmt.Errorf("json marshal error"))

		makeHttpRequest(t, http.StatusBadRequest, func() (resp *http.Response, err error) {
			return http.Post(server.URL+urlPath, "application/json", bytes.NewBufferString("{}"))
		})
	})

	t.Run("StatusBadRequest create commit status error", func(t *testing.T) {
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
			WithValidateWebhookRequestReturns(newPrEvent("opened").getEvent(), nil).
			WithNewInstallationClientReturns(badClient, nil)

		makeHttpRequest(t, http.StatusBadRequest, func() (resp *http.Response, err error) {
			return http.Post(server.URL+urlPath, "application/json", bytes.NewBufferString("{}"))
		})
	})

	t.Run("StatusBadRequest new installation client error", func(t *testing.T) {
		apiServer.githubFactory = NewMockGithubClientFactory(apiServer).
			WithValidateWebhookRequestReturns(newPrEvent("opened").getEvent(), nil).
			WithNewInstallationClientReturns(&http.Client{}, fmt.Errorf("github client error"))

		makeHttpRequest(t, http.StatusBadRequest, func() (resp *http.Response, err error) {
			return http.Post(server.URL+urlPath, "application/json", bytes.NewBufferString("{}"))
		})
	})

	t.Run("StatusBadRequest querier error UpdateApprovalByUuid", func(t *testing.T) {
		fakeQuerier.UpdateApprovalByUuidCall.Returns.Error = fmt.Errorf("querier error")

		makeHttpRequest(t, http.StatusBadRequest, func() (resp *http.Response, err error) {
			return http.Post(server.URL+urlPath, "application/json", bytes.NewBufferString("{}"))
		})
	})

	t.Run("StatusBadRequest querier error GetOrCreateApproval", func(t *testing.T) {
		fakeQuerier.GetApprovalByPrIDShaCall.Returns.Error = fmt.Errorf("querier error")
		fakeQuerier.CreateApprovalCall.Returns.Error = fmt.Errorf("querier error")

		makeHttpRequest(t, http.StatusBadRequest, func() (resp *http.Response, err error) {
			return http.Post(server.URL+urlPath, "application/json", bytes.NewBufferString("{}"))
		})
	})

	t.Run("StatusBadRequest querier error CreatePullRequestEvent", func(t *testing.T) {
		fakeQuerier.CreatePullRequestEventCall.Returns.Error = fmt.Errorf("querier error")

		makeHttpRequest(t, http.StatusBadRequest, func() (resp *http.Response, err error) {
			return http.Post(server.URL+urlPath, "application/json", bytes.NewBufferString("{}"))
		})
	})

	t.Run("StatusBadRequest querier error GetOrCreatePullRequest", func(t *testing.T) {
		fakeQuerier.GetPullRequestByRepoIdPrIdCall.Returns.Error = fmt.Errorf("querier error")
		fakeQuerier.CreatePullRequestCall.Returns.Error = fmt.Errorf("querier error")

		makeHttpRequest(t, http.StatusBadRequest, func() (resp *http.Response, err error) {
			return http.Post(server.URL+urlPath, "application/json", bytes.NewBufferString("{}"))
		})
	})

	t.Run("StatusBadRequest querier error GetOrCreateRepo", func(t *testing.T) {
		fakeQuerier.GetRepoCall.Returns.Error = fmt.Errorf("querier error")
		fakeQuerier.CreateRepoCall.Returns.Error = fmt.Errorf("querier error")

		makeHttpRequest(t, http.StatusBadRequest, func() (resp *http.Response, err error) {
			return http.Post(server.URL+urlPath, "application/json", bytes.NewBufferString("{}"))
		})
	})

	t.Run("StatusBadRequest querier error GetOrCreateGithubUser", func(t *testing.T) {
		fakeQuerier.GetGithubUserCall.Returns.Error = fmt.Errorf("querier error")
		fakeQuerier.CreateGithubUserCall.Returns.Error = fmt.Errorf("querier error")

		makeHttpRequest(t, http.StatusBadRequest, func() (resp *http.Response, err error) {
			return http.Post(server.URL+urlPath, "application/json", bytes.NewBufferString("{}"))
		})
	})

	t.Run("StatusBadRequest querier error GetOrCreatePullRequestAction", func(t *testing.T) {
		fakeQuerier.GetPullRequestActionCall.Returns.Error = fmt.Errorf("querier error")
		fakeQuerier.CreatePullRequestActionCall.Returns.Error = fmt.Errorf("querier error")

		makeHttpRequest(t, http.StatusBadRequest, func() (resp *http.Response, err error) {
			return http.Post(server.URL+urlPath, "application/json", bytes.NewBufferString("{}"))
		})
	})

	t.Run("StatusBadRequest querier GetOrCreatePullRequestAction map update", func(t *testing.T) {
		fakeQuerier.GetPullRequestActionCall.Returns.Error = fmt.Errorf("querier error")
		fakeQuerier.CreatePullRequestActionCall.Returns.Error = nil

		makeHttpRequest(t, http.StatusBadRequest, func() (resp *http.Response, err error) {
			return http.Post(server.URL+urlPath, "application/json", bytes.NewBufferString("{}"))
		})

		if _, known := apiServer.knownPullRequestActions["opened"]; !known {
			t.Errorf("Expected 'opened' to be in map %#v", apiServer.knownPullRequestActions)
		}

		// If it's not in the map then it will return an error
		fakeQuerier.CreatePullRequestActionCall.Returns.Error = fmt.Errorf("querier error")

		makeHttpRequest(t, http.StatusBadRequest, func() (resp *http.Response, err error) {
			return http.Post(server.URL+urlPath, "application/json", bytes.NewBufferString("{}"))
		})
	})

	t.Run("StatusBadRequest querier error GetOrInstallation", func(t *testing.T) {
		fakeQuerier.GetInstallationCall.Returns.Error = fmt.Errorf("querier error")
		fakeQuerier.CreateInstallationCall.Returns.Error = fmt.Errorf("querier error")

		makeHttpRequest(t, http.StatusBadRequest, func() (resp *http.Response, err error) {
			return http.Post(server.URL+urlPath, "application/json", bytes.NewBufferString("{}"))
		})
	})

}

type prEvent struct {
	event *github.PullRequestEvent
}

func newPrEvent(action string) *prEvent {
	return &prEvent{
		event: &github.PullRequestEvent{
			Action: github.String(action),
			Number: github.Int(1),
			Repo: &github.Repository{
				ID:   github.Int64(int64(80)),
				Name: github.String("hello-world-app"),
			},
			PullRequest: &github.PullRequest{
				ID:     github.Int64(int64(9991)),
				Number: github.Int(1),
				Merged: github.Bool(false),
				Head: &github.PullRequestBranch{
					SHA: github.String("038d718da6a1ebbc6a7780a96ed75a70cc2ad6e2"), // echo testing | git hash-object --stdin -w
				},
			},
			Sender: &github.User{
				Login: github.String("frodo"),
				ID:    github.Int64(int64(1)),
			},
			Installation: &github.Installation{
				ID: github.Int64(testInstallationID),
			},
			Organization: nil,
		},
	}
}

func (e *prEvent) withSender(name string) *prEvent {
	e.event.Sender = &github.User{
		ID:   github.Int64(1),
		Name: github.String(name),
	}
	return e
}

func (e *prEvent) withOrg(name string) *prEvent {
	e.event.Organization = &github.Organization{
		Name: github.String(name),
	}
	return e
}

func (e *prEvent) getEvent() *github.PullRequestEvent {
	return e.event
}

func TestGetCreateParamsFromEvent(t *testing.T) {
	tests := []struct {
		name         string
		event        *github.PullRequestEvent
		wantEmtpyOrg bool
	}{
		{
			name:         "no org",
			event:        newPrEvent("opened").getEvent(),
			wantEmtpyOrg: true,
		},
		{
			name:         "with org",
			event:        newPrEvent("opened").withOrg("some-org").getEvent(),
			wantEmtpyOrg: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := getCreateParamsFromEvent(test.event)
			apiServer.githubFactory = NewMockGithubClientFactory(apiServer).
				WithValidateWebhookRequestReturns(newPrEvent("opened").getEvent(), nil)

			if test.wantEmtpyOrg && p.Repo.Org != "" {
				t.Error("Expected empty org")
			}
			if !test.wantEmtpyOrg && p.Repo.Org == "" {
				t.Error("expected non-empty org")
			}
		})
	}

}

func updateQuerierWithSuccessPostWebhookEventReturns(querier *fakes.Querier) *fakes.Querier {

	e := newPrEvent("opened").getEvent()

	querier.GetInstallationCall.Returns.Int32 = testInstallationID

	querier.GetPullRequestActionCall.Returns.String = "opened"

	querier.GetGithubUserCall.Returns.GhUser = postgres.GhUser{
		ID:    int32(*e.Sender.ID),
		Login: *e.Sender.Login,
	}

	querier.GetRepoCall.Returns.Repo = postgres.Repo{
		ID:   int32(*e.Repo.ID),
		Org:  "",
		Name: *e.Repo.Name,
	}

	querier.GetPullRequestByRepoIdPrIdCall.Returns.PullRequest = postgres.PullRequest{
		ID:             123,
		RepoID:         int32(*e.Repo.ID),
		PrID:           int32(*e.PullRequest.ID),
		PrNumber:       int32(*e.PullRequest.Number),
		OpenedBy:       int32(*e.Sender.ID),
		InstallationID: testInstallationID,
		IsMerged:       false,
	}

	querier.GetApprovalByPrIDShaCall.Returns.Approval = postgres.Approval{
		ID:         321,
		Uuid:       testUuid,
		PrID:       int32(*e.PullRequest.ID),
		Sha:        *e.PullRequest.Head.SHA,
		IsApproved: true,
	}

	return querier
}
