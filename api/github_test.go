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

	api := NewMockedApiServer(&fakes.Querier{})
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
	defaultClient                 *http.Client
	server                        *Server
	ValidateWebhookRequestReturns struct {
		event interface{}
		err   error
	}
	NewAppClientReturns struct {
		httpClient *http.Client
		err        error
	}
	NewInstallationClientReturns struct {
		httpClient *http.Client
		err        error
	}
}

func NewMockGithubClientFactory(server *Server) *mockGithubFactory {
	factory := &mockGithubFactory{server: server}

	// default mocks for factory
	factory.defaultClient = mock.NewMockedHTTPClient(
		mock.WithRequestMatch(
			mock.PostReposStatusesByOwnerByRepoBySha,
			// Multiple return responses supported (and required for some tests)
			// Adding additional responses fixes: http: panic serving 127.0.0.1:59729: runtime error: index out of range [1] with length 1
			github.RepoStatus{
				ID: github.Int64(int64(333)),
			},
			github.RepoStatus{
				ID: github.Int64(int64(333)),
			},
			github.RepoStatus{
				ID: github.Int64(int64(333)),
			},
		),
		mock.WithRequestMatch(
			mock.PostAppInstallationsAccessTokensByInstallationId,
			github.InstallationToken{
				Token: github.String("something"),
			},
			github.InstallationToken{
				Token: github.String("something"),
			},
			github.InstallationToken{
				Token: github.String("something"),
			},
		),
	)

	factory.NewAppClientReturns.httpClient = factory.defaultClient
	factory.NewInstallationClientReturns.httpClient = factory.defaultClient

	return factory
}

func (f *mockGithubFactory) WithValidateWebhookRequestReturns(event interface{}, err error) *mockGithubFactory {
	f.ValidateWebhookRequestReturns.event = event
	f.ValidateWebhookRequestReturns.err = err
	return f
}

func (f *mockGithubFactory) WithNewAppClientReturns(client *http.Client, err error) *mockGithubFactory {
	f.NewAppClientReturns.httpClient = client
	f.NewAppClientReturns.err = err
	return f
}

func (f *mockGithubFactory) WithNewInstallationClientReturns(client *http.Client, err error) *mockGithubFactory {
	f.NewInstallationClientReturns.httpClient = client
	f.NewInstallationClientReturns.err = err
	return f
}

func (f *mockGithubFactory) ValidatWebhookRequest(r *http.Request) (interface{}, error) {
	return f.ValidateWebhookRequestReturns.event, f.ValidateWebhookRequestReturns.err
}

func (f *mockGithubFactory) NewInstallationClient(ctx context.Context, installationID int64) (*github.Client, error) {
	return github.NewClient(f.NewInstallationClientReturns.httpClient), f.NewInstallationClientReturns.err
}

func (f *mockGithubFactory) NewAppClient(ctx context.Context) (*github.Client, error) {
	return github.NewClient(f.NewAppClientReturns.httpClient), f.NewAppClientReturns.err
}
