package api

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-github/v53/github"
	"github.com/jericop/pr-compliance-app/fakes"
	"github.com/migueleliasweb/go-github-mock/src/mock"
)

func TestValidatWebhookRequest(t *testing.T) {
	// Tests are based on: https://github.com/google/go-github/blob/v53.2.0/github/messages_test.go
	body := `{"yo":true}`
	contentType := "application/json"
	signature := "sha1=126f2c800419c60137ce748d7672e77b65cf16d6"

	tests := []struct {
		name         string
		wantError    bool
		contentType  string
		xGithubEvent string
	}{
		{
			name:      "github.ValidatePayload error no media type",
			wantError: true,
		},
		{
			name:        "github.ParseWebHook error unknown X-Github-Event in message",
			wantError:   true,
			contentType: contentType,
		},
		{
			name:         "success",
			wantError:    false,
			contentType:  contentType,
			xGithubEvent: "pull_request",
		},
	}

	api := getApiServer(&fakes.Querier{})
	api.githubFactory = NewGithubFactory(api)

	for n, test := range tests {
		t.Run(fmt.Sprintf("%d %s", n, test.name), func(t *testing.T) {
			buf := bytes.NewBufferString(body)
			req, err := http.NewRequest("POST", "http://localhost/event", buf)
			if err != nil {
				t.Errorf("NewRequest: %v", err)
			}
			req.Header.Set(github.SHA1SignatureHeader, signature)

			if test.contentType != "" {
				req.Header.Set("Content-Type", test.contentType)
			}

			if test.xGithubEvent != "" {
				req.Header.Set(github.EventTypeHeader, test.xGithubEvent)
			}

			_, err = api.githubFactory.ValidatWebhookRequest(req)

			if test.wantError {
				if err == nil {
					t.Errorf("expected an error but got nil")
				}
			} else if err != nil {
				t.Errorf("got: err = %v, want nil", err)
			}
		})
	}
}

type mockGithubFactory struct {
	server                 *Server
	validateWebhookReturns struct {
		event interface{}
		err   error
	}
}

func NewMockGithubClientFactory(server *Server) *mockGithubFactory {
	return &mockGithubFactory{server: server}
}

func (f *mockGithubFactory) WithValidateWebhookReturns(event interface{}, err error) *mockGithubFactory {
	f.validateWebhookReturns.event = event
	f.validateWebhookReturns.err = err
	return f
}

func (f *mockGithubFactory) ValidatWebhookRequest(r *http.Request) (interface{}, error) {
	return f.validateWebhookReturns.event, f.validateWebhookReturns.err
}

func (f *mockGithubFactory) NewInstallationClient(ctx context.Context, installationID int64) (*github.Client, error) {
	mockedHTTPClient := mock.NewMockedHTTPClient(
		mock.WithRequestMatch(
			// https://github.com/google/go-github/blob/16e695dadf7afb7983193499816a009ae2227a61/github/repos_statuses.go#L75
			// https://github.com/migueleliasweb/go-github-mock/blob/1784f27b54c9c95f12160449d7e16344dd512e88/src/mock/endpointpattern.go#L3470
			mock.PostReposStatusesByOwnerByRepoBySha,
			github.RepoStatus{
				ID: github.Int64(int64(333)),
			},
		),
		mock.WithRequestMatch(
			// https://github.com/google/go-github/blob/16e695dadf7afb7983193499816a009ae2227a61/github/apps.go#L310
			// https://github.com/migueleliasweb/go-github-mock/blob/1784f27b54c9c95f12160449d7e16344dd512e88/src/mock/endpointpattern.go#L75C5-L75C53
			mock.PostAppInstallationsAccessTokensByInstallationId,
			github.InstallationToken{
				Token: github.String("something"),
			},
		),
	)
	return github.NewClient(mockedHTTPClient), nil
}
