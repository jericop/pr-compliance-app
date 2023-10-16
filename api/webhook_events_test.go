package api

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-github/v53/github"
)

func TestValidatWebhookRequest(t *testing.T) {
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
				fmt.Printf("setting conttentType: %v \n", test.contentType)
				req.Header.Set("Content-Type", test.contentType)
			}

			if test.xGithubEvent != "" {
				fmt.Printf("setting eventType: %v \n", test.xGithubEvent)
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

func TestValidatWebhookRequest1(t *testing.T) {
	t.Skip()

	payload := `{"yo":true}`
	// signature := "sha1=3374ef144403e8035423b23b02e2c9d7a4c50368"
	// signature := "sha256=b1f8020f5b4cd42042f807dd939015c4a418bc1ff7f604dd55b0a19b5d953d9b"
	signature := "sha256=3182e6c2ee3cd5cbdb0360ead9b9d1f14e62d813dfc65fb2f350a5b9b8d34e37"
	// "encoding/hex"

	// sigParts := strings.SplitN(signature, "=", 2)
	// if len(sigParts) != 2 {
	// 	t.Errorf("error parsing signature %q", signature)
	// }
	// buf, err := hex.DecodeString(sigParts[1])
	// if err != nil {
	// 	t.Errorf("error decoding signature %q: %v", signature, err)
	// }

	// mac := hmac.New(sha256.New, key)
	// mac.Write(message)

	// fmt.Printf("len(buf) %d\n", len(buf))

	// form := url.Values{}
	// form.Add("payload", payload)
	buf := bytes.NewBufferString(payload)
	req, err := http.NewRequest("POST", "http://localhost/event", buf)
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}
	// req.PostForm = form
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// req.Header.Set(github.SHA1SignatureHeader, signature)
	req.Header.Set(github.SHA256SignatureHeader, signature)

	event, err := apiServer.validatWebhookRequest(req)
	if err != nil {
		t.Errorf("got: err = %v, want nil", err)
	}

	fmt.Printf("event (%T) %v\n", event, event)

	// got, err := github.ValidatePayload(req, secretKey)
	// if err != nil {
	// 	t.Errorf("ValidatePayload(%#v): err = %v, want nil", payload, err)
	// }
	// if string(got) != payload {
	// 	t.Errorf("ValidatePayload = %q, want %q", got, payload)
	// }

	// // check that if payload is invalid we get error
	// req.Header.Set(github.SHA1SignatureHeader, "invalid signature")
	// if _, err = github.ValidatePayload(req, []byte{0}); err == nil {
	// 	t.Error("ValidatePayload = nil, want err")
	// }

	// tests := []struct {
	// 	contentType string
	// }{
	// 	{contentType: "application/json"},
	// 	{contentType: "application/x-www-form-urlencoded"},
	// }

	// for i, tt := range tests {
	// 	t.Run(fmt.Sprintf("test #%v", i), func(t *testing.T) {
	// 		req := &http.Request{
	// 			Header: http.Header{"Content-Type": []string{"application/json"}},
	// 			Body:   &badReader{},
	// 		}
	// 		if _, err := validatWebhookRequest(req, nil); err == nil {
	// 			t.Fatal("ValidatePayload returned nil; want error")
	// 		}
	// 	})
	// }

	// event, err := getAndValidatePayload(r)
	// if err != nil {
	// 	log.Printf("Failed to parse payload: %v", err)
	// 	http.Error(w, "Invalid payload", http.StatusBadRequest)
	// 	return
	// }
}

// alwaysErrorReader satisfies io.Reader but always returns an error.
type alwaysErrorReader struct{}

func (m *alwaysErrorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}

func (b *alwaysErrorReader) Close() error { return errors.New("close error") }

// // TODO: See below url for testing github.{ValidatePayload,ValidateSignature,ParseWebHook}
// // https://github.com/google/go-github/blob/v53.2.0/github/messages_test.go
// func (server *Server) validatWebhookRequest(r *http.Request) (interface{}, error) {
// 	// Validate and parse payload from request
// 	payload, err := github.ValidatePayload(r, []byte(server.webhookSecret))
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to validate and parse payload: %v", err)
// 	}

// 	// Validate signature of payload
// 	err = github.ValidateSignature(r.Header.Get("X-Hub-Signature"), payload, []byte(server.webhookSecret))
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to validate payload signature: %v", err)
// 	}

// 	// Parse event from payload
// 	event, err := github.ParseWebHook(github.WebHookType(r), payload)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to parse webhook into event: %v", err)
// 	}

// 	return event, nil
// }

// var validPullRequests []postgres.PullRequest

// func TestGetPullRequests(t *testing.T) {
// 	fakeStore, apiServer, urlPath := getQuerierServerRouteUrl(t, "GetPullRequests")

// 	// This http server does not record requests and requests to it will not show up in the api blueprint document.
// 	server := httptest.NewServer(apiServer.Router)
// 	defer server.Close()

// 	expected := []postgres.PullRequest{
// 		{ID: 1, PrID: 991, RepoID: 444, PrNumber: 1, OpenedBy: 651, IsMerged: false},
// 	}

// 	// Setup and test the happy path first
// 	fakeStore.GetPullRequestsCall.Returns.Error = nil
// 	fakeStore.GetPullRequestsCall.Returns.PullRequestSlice = expected

// 	resp := makeHttpRequest(t, http.StatusOK, func() (resp *http.Response, err error) {
// 		return http.Get(test2docServer.URL + urlPath)
// 	})

// 	decoder := json.NewDecoder(resp.Body)
// 	defer resp.Body.Close()

// 	result := []postgres.PullRequest{}
// 	if err := decoder.Decode(&result); err != nil {
// 		t.Fatalf("expected 'err' (%v) be nil", err)
// 	}

// 	if !reflect.DeepEqual(result, expected) {
// 		t.Fatalf("expected 'result' (%v) to equal 'expected' (%v)", result, expected)
// 	}

// 	if fakeStore.GetPullRequestsCall.CallCount != 1 {
// 		t.Error("unexpected call count")
// 	}

// 	// Next we simulate an error decoding result from the database
// 	// Set apiServer.jsonMarshal to an anonymous function with the same signature as json.Marshal in order to force an error.
// 	apiServer.jsonMarshal = func(v interface{}) ([]byte, error) {
// 		return []byte{}, fmt.Errorf("Marshalling failed")
// 	}

// 	_ = makeHttpRequest(t, http.StatusInternalServerError, func() (resp *http.Response, err error) {
// 		return http.Get(server.URL + urlPath)
// 	})

// 	if fakeStore.GetPullRequestsCall.CallCount != 2 {
// 		t.Error("unexpected call count")
// 	}

// 	apiServer.jsonMarshal = json.Marshal

// 	// Test an error from the database
// 	fakeStore.GetPullRequestsCall.Returns.Error = fmt.Errorf("db error")
// 	fakeStore.GetPullRequestsCall.Returns.PullRequestSlice = []postgres.PullRequest{}

// 	_ = makeHttpRequest(t, http.StatusInternalServerError, func() (resp *http.Response, err error) {
// 		return http.Get(server.URL + urlPath)
// 	})

// 	if fakeStore.GetPullRequestsCall.CallCount != 3 {
// 		t.Error("unexpected call count")
// 	}
// }

// func TestGetPullRequestsOriginal(t *testing.T) {
// 	t.Skip()

// 	urlPath, err := router.Get("GetPullRequests").URL()
// 	if err != nil {
// 		t.Fatalf("expected 'err' (%v) be nil", err)
// 	}

// 	result := []postgres.PullRequest{}
// 	expected := []postgres.PullRequest{
// 		{ID: 1, PrID: 991, RepoID: 444, PrNumber: 1, OpenedBy: 651, IsMerged: false},
// 	}

// 	// This http server does not record requests
// 	server := httptest.NewServer(apiServer.Router)
// 	defer server.Close()
// 	fmt.Printf("server.URL: %v\n", server.URL)

// 	// Test the happy path first
// 	fakeStore.GetPullRequestsCall.Returns.PullRequestSlice = expected
// 	fakeStore.GetPullRequestsCall.Returns.Error = nil

// 	resp, err := http.Get(test2docServer.URL + urlPath.String())
// 	if err != nil {
// 		t.Fatalf("expected 'err' (%v) be nil", err)
// 	}

// 	if resp.StatusCode != http.StatusOK {
// 		t.Fatalf("expected 'resp.StatusCode' (%v) to equal 'http.StatusOK' (%v)", resp.StatusCode, http.StatusOK)
// 	}

// 	decoder := json.NewDecoder(resp.Body)
// 	defer resp.Body.Close()

// 	err = decoder.Decode(&result)
// 	if err != nil {
// 		t.Fatalf("expected 'err' (%v) be nil", err)
// 	}

// 	if !reflect.DeepEqual(result, expected) {
// 		t.Fatalf("expected 'result' (%v) to equal 'expected' (%v)", result, expected)
// 	}

// 	if fakeStore.GetPullRequestsCall.CallCount != 1 {
// 		t.Error("unexpected call count")
// 	}

// 	// Set apiServer.jsonMarshal to an anonymous function with the same signature as json.Marshal in order to force an error.
// 	apiServer.jsonMarshal = func(v interface{}) ([]byte, error) {
// 		return []byte{}, fmt.Errorf("Marshalling failed")
// 	}

// 	resp, err = http.Get(server.URL + urlPath.String())
// 	if err != nil {
// 		t.Fatalf("expected 'err' (%v) to be nil", err)
// 	}

// 	if resp.StatusCode != http.StatusInternalServerError {
// 		t.Fatalf("expected 'resp.StatusCode' (%v) to equal 'http.StatusInternalServerError' (%v)", resp.StatusCode, http.StatusInternalServerError)
// 	}

// 	if fakeStore.GetPullRequestsCall.CallCount != 2 {
// 		t.Error("unexpected call count")
// 	}

// 	// Set the
// 	fakeStore.GetPullRequestsCall.Returns.PullRequestSlice = []postgres.PullRequest{}
// 	fakeStore.GetPullRequestsCall.Returns.Error = fmt.Errorf("db error")

// 	resp, err = http.Get(server.URL + urlPath.String())
// 	if err != nil {
// 		t.Fatalf("expected 'err' (%v) to be nil", err)
// 	}

// 	if resp.StatusCode != http.StatusInternalServerError {
// 		t.Fatalf("expected 'resp.StatusCode' (%v) to equal 'http.StatusInternalServerError' (%v)", resp.StatusCode, http.StatusInternalServerError)
// 	}

// 	if fakeStore.GetPullRequestsCall.CallCount != 3 {
// 		t.Error("unexpected call count")
// 	}
// }

// func getQuerierServerRouteUrl(t *testing.T, routeName string) (*fakes.Querier, *Server, string) {
// 	store := &fakes.Querier{}
// 	apiServer := NewServer(store)
// 	urlPath, err := apiServer.Router.Get(routeName).URL()
// 	if err != nil {
// 		t.Fatalf("expected 'err' (%v) be nil", err)
// 	}
// 	return store, apiServer, urlPath.String()
// }

// func makeHttpRequest(t *testing.T, expectedStatusCode int, httpRequestFunc func() (resp *http.Response, err error)) *http.Response {
// 	resp, err := httpRequestFunc() // http.Get or http.Post functions get called here
// 	if err != nil {
// 		t.Fatalf("expected 'err' (%v) be nil", err)
// 	}

// 	if resp.StatusCode != expectedStatusCode {
// 		t.Fatalf("expected 'resp.StatusCode' (%v) to equal 'expectedStatusCode' (%v)", resp.StatusCode, http.StatusOK)
// 	}
// 	return resp
// }
