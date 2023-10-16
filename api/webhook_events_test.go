package api

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-github/v53/github"
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

			_, err = apiServer.validatWebhookRequest(req)
			// fmt.Printf("err: %v\n", err)

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
