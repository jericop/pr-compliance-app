package api

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/hex"
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

	t.Run("StatusOK test2doc", func(t *testing.T) {
		for n, test := range tests {
			pJson, err := apiServer.jsonMarshal(test.payload)
			if err != nil {
				t.Fatalf("Marshal(%#v): %v", test.payload, err)
			}

			buf := bytes.NewBuffer(pJson)

			makeHttpRequest(t, http.StatusOK, func() (resp *http.Response, err error) {
				// Only use the first payload to keep the api documentation concise
				serverUrl := server.URL
				if n == 0 {
					serverUrl = test2docServer.URL
				}

				req, err := http.NewRequest("POST", serverUrl+urlPath, buf)
				if err != nil {
					t.Errorf("NewRequest: %v", err)
				}

				// Get payload sha1 signature using githubWebhookSecret
				mac := hmac.New(sha1.New, []byte(apiServer.githubWebhookSecret))
				mac.Write(pJson)
				expectedMAC := mac.Sum(nil)

				req.Header.Set("Content-Type", "application/json")
				req.Header.Set(github.SHA1SignatureHeader, "sha1="+hex.EncodeToString(expectedMAC))
				req.Header.Set(github.EventTypeHeader, test.messageType)

				return http.DefaultClient.Do(req)
			})
		}
	})

	// http.StatusBadRequest

	// approvalId := uuid.New().String()
	// urlPath := getRouteUrlPath(t, apiServer.router, "PostWebhookEvent")

	// p := postgres.UpdateApprovalByUuidParams{
	// 	Uuid:       approvalId,
	// 	IsApproved: true,
	// }
	// expected := p // The response may contain different json data than POST request body

	// pJson, err := json.Marshal(p)
	// if err != nil {
	// 	t.Fatalf("expected 'err' (%v) be nil", err)
	// }

	t.Run("StatusBadRequest validate webhook event error", func(t *testing.T) {
		buf := bytes.NewBufferString("not valid json")

		apiServer.githubFactory = NewMockGithubClientFactory(apiServer).
			WithValidateWebhookReturns(struct{}{}, fmt.Errorf("json marshalerror"))

		makeHttpRequest(t, http.StatusBadRequest, func() (resp *http.Response, err error) {
			return http.Post(server.URL+urlPath, "application/json", buf)
		})
	})

	// t.Run("StatusBadRequest querier error", func(t *testing.T) {
	// 	buf := bytes.NewBufferString("ignored because apiServer.githubFactory is being mocked to return a valid value")

	// 	apiServer.githubFactory = NewMockGithubClientFactory(apiServer).
	// 		WithValidateWebhookReturns(newPrEvent("opened").getEvent(), nil)

	// 	// fakeStore.GetGithubUserCall.Returns.Error = fmt.Errorf("querier error")
	// 	fakeStore.GetRepoCall.Returns.Repo = postgres.Repo{}
	// 	fakeStore.GetRepoCall.Returns.Error = fmt.Errorf("querier error")

	// 	fakeStore.GetPullRequestActionCall.Returns.Error = fmt.Errorf("querier error")
	// 	fmt.Printf("fakeStore.GetPullRequestActionCall before : %d\n", fakeStore.GetPullRequestActionCall.CallCount)

	// 	makeHttpRequest(t, http.StatusBadRequest, func() (resp *http.Response, err error) {
	// 		return http.Post(server.URL+urlPath, "application/json", buf)
	// 	})

	// 	fmt.Printf("fakeStore.GetPullRequestActionCall after : %d\n", fakeStore.GetPullRequestActionCall.CallCount)
	// })

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
