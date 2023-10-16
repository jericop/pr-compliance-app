package api

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-github/v53/github"
	"github.com/google/uuid"
	"github.com/jericop/pr-compliance-app/storage/postgres"
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

	e := newPrEvent("opened").getEvent()

	fakeStore.GetPullRequestActionCall.Returns.String = "opened"

	fakeStore.GetGithubUserCall.Returns.GhUser = postgres.GhUser{
		ID:    int32(*e.Sender.ID),
		Login: *e.Sender.Login,
	}

	fakeStore.GetRepoCall.Returns.Repo = postgres.Repo{
		ID:   int32(*e.Repo.ID),
		Org:  "",
		Name: *e.Repo.Name,
	}

	fakeStore.GetPullRequestByRepoIdPrIdCall.Returns.PullRequest = postgres.PullRequest{
		ID:       123,
		RepoID:   int32(*e.Repo.ID),
		PrID:     int32(*e.PullRequest.ID),
		PrNumber: int32(*e.PullRequest.Number),
		OpenedBy: int32(*e.Sender.ID),
		IsMerged: false,
	}

	fakeStore.GetApprovalByPrIDShaCall.Returns.Approval = postgres.Approval{
		ID:         321,
		Uuid:       uuid.New().String(),
		PrID:       int32(*e.PullRequest.ID),
		Sha:        *e.PullRequest.Head.SHA,
		IsApproved: false,
	}

	apiServer.githubFactory = NewMockGithubClientFactory(apiServer).
		WithValidateWebhookReturns(newPrEvent("opened").getEvent(), nil)

	// Generate RSA key.
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		t.Errorf("rsa.GenerateKey expected 'err' (%v) be nil", err)
	}

	apiServer.githubPrivateKey = key

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
				WithValidateWebhookReturns(newPrEvent("opened").getEvent(), nil)

			makeHttpRequest(t, http.StatusOK, func() (resp *http.Response, err error) {
				return http.Post(serverURL+urlPath, "application/json", bytes.NewBufferString("{}"))
			})
		})
	}

	t.Run("StatusBadRequest validate webhook event error", func(t *testing.T) {
		apiServer.githubFactory = NewMockGithubClientFactory(apiServer).
			WithValidateWebhookReturns(struct{}{}, fmt.Errorf("json marshal error"))

		makeHttpRequest(t, http.StatusBadRequest, func() (resp *http.Response, err error) {
			return http.Post(server.URL+urlPath, "application/json", bytes.NewBufferString("{}"))
		})
	})

	t.Run("StatusBadRequest new installation client error", func(t *testing.T) {
		apiServer.githubFactory = NewMockGithubClientFactory(apiServer).
			WithValidateWebhookReturns(newPrEvent("opened").getEvent(), nil).
			WithNewInstallationClientReturns(&http.Client{}, fmt.Errorf("github client error"))

		makeHttpRequest(t, http.StatusBadRequest, func() (resp *http.Response, err error) {
			return http.Post(server.URL+urlPath, "application/json", bytes.NewBufferString("{}"))
		})
	})

	t.Run("StatusBadRequest querier error GetOrCreateApproval", func(t *testing.T) {
		apiServer.githubFactory = NewMockGithubClientFactory(apiServer).
			WithValidateWebhookReturns(newPrEvent("opened").getEvent(), nil)

		fakeStore.GetApprovalByPrIDShaCall.Returns.Error = fmt.Errorf("querier error")
		fakeStore.CreateApprovalCall.Returns.Error = fmt.Errorf("querier error")

		makeHttpRequest(t, http.StatusBadRequest, func() (resp *http.Response, err error) {
			return http.Post(server.URL+urlPath, "application/json", bytes.NewBufferString("{}"))
		})
	})

	t.Run("StatusBadRequest querier error CreatePullRequestEvent", func(t *testing.T) {
		apiServer.githubFactory = NewMockGithubClientFactory(apiServer).
			WithValidateWebhookReturns(newPrEvent("opened").getEvent(), nil)

		fakeStore.CreatePullRequestEventCall.Returns.Error = fmt.Errorf("querier error")

		makeHttpRequest(t, http.StatusBadRequest, func() (resp *http.Response, err error) {
			return http.Post(server.URL+urlPath, "application/json", bytes.NewBufferString("{}"))
		})
	})

	t.Run("StatusBadRequest querier error GetOrCreatePullRequest", func(t *testing.T) {
		apiServer.githubFactory = NewMockGithubClientFactory(apiServer).
			WithValidateWebhookReturns(newPrEvent("opened").getEvent(), nil)

		fakeStore.GetPullRequestByRepoIdPrIdCall.Returns.Error = fmt.Errorf("querier error")
		fakeStore.CreatePullRequestCall.Returns.Error = fmt.Errorf("querier error")

		makeHttpRequest(t, http.StatusBadRequest, func() (resp *http.Response, err error) {
			return http.Post(server.URL+urlPath, "application/json", bytes.NewBufferString("{}"))
		})
	})

	t.Run("StatusBadRequest querier error GetOrCreateRepo", func(t *testing.T) {
		apiServer.githubFactory = NewMockGithubClientFactory(apiServer).
			WithValidateWebhookReturns(newPrEvent("opened").getEvent(), nil)

		fakeStore.GetRepoCall.Returns.Error = fmt.Errorf("querier error")
		fakeStore.CreateRepoCall.Returns.Error = fmt.Errorf("querier error")

		makeHttpRequest(t, http.StatusBadRequest, func() (resp *http.Response, err error) {
			return http.Post(server.URL+urlPath, "application/json", bytes.NewBufferString("{}"))
		})
	})

	t.Run("StatusBadRequest querier error GetOrCreateGithubUser", func(t *testing.T) {
		apiServer.githubFactory = NewMockGithubClientFactory(apiServer).
			WithValidateWebhookReturns(newPrEvent("opened").getEvent(), nil)

		fakeStore.GetGithubUserCall.Returns.Error = fmt.Errorf("querier error")
		fakeStore.CreateGithubUserCall.Returns.Error = fmt.Errorf("querier error")

		makeHttpRequest(t, http.StatusBadRequest, func() (resp *http.Response, err error) {
			return http.Post(server.URL+urlPath, "application/json", bytes.NewBufferString("{}"))
		})
	})

	t.Run("StatusBadRequest querier error GetOrCreatePullRequestAction", func(t *testing.T) {
		apiServer.githubFactory = NewMockGithubClientFactory(apiServer).
			WithValidateWebhookReturns(newPrEvent("opened").getEvent(), nil)

		fakeStore.GetPullRequestActionCall.Returns.Error = fmt.Errorf("querier error")
		fakeStore.CreatePullRequestActionCall.Returns.Error = fmt.Errorf("querier error")

		makeHttpRequest(t, http.StatusBadRequest, func() (resp *http.Response, err error) {
			return http.Post(server.URL+urlPath, "application/json", bytes.NewBufferString("{}"))
		})
	})

	t.Run("StatusBadRequest querier GetOrCreatePullRequestAction map update", func(t *testing.T) {
		apiServer.githubFactory = NewMockGithubClientFactory(apiServer).
			WithValidateWebhookReturns(newPrEvent("opened").getEvent(), nil)

		fakeStore.GetPullRequestActionCall.Returns.Error = fmt.Errorf("querier error")
		fakeStore.CreatePullRequestActionCall.Returns.Error = nil

		makeHttpRequest(t, http.StatusBadRequest, func() (resp *http.Response, err error) {
			return http.Post(server.URL+urlPath, "application/json", bytes.NewBufferString("{}"))
		})

		if _, known := apiServer.KnownPullRequestActions["opened"]; !known {
			t.Errorf("Expected 'opened' to be in map %#v", apiServer.KnownPullRequestActions)
		}

		// If it's not in the map then it will return an error
		fakeStore.CreatePullRequestActionCall.Returns.Error = fmt.Errorf("querier error")

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
				ID: github.Int64(int64(8675309)),
			},
		},
	}
}

func (e *prEvent) withSender(name string) *prEvent {
	e.event.Organization = &github.Organization{
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
