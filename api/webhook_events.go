package api

import (
	"fmt"
	"net/http"

	"github.com/google/go-github/v53/github"
)

// TODO: See below url for testing github.{ValidatePayload,ValidateSignature,ParseWebHook}
// https://github.com/google/go-github/blob/v53.2.0/github/messages_test.go
func (server *Server) validatWebhookRequest(r *http.Request) (interface{}, error) {
	// Validate and parse payload from request and validate signature if set
	payload, err := github.ValidatePayload(r, []byte(server.webhookSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to validate and parse payload: %v", err)
	}

	// Parse event from payload
	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		return nil, fmt.Errorf("failed to parse webhook into event: %v", err)
	}

	return event, nil
}

// func (server *Server) AddWebhookEventsRoutes() {
// 	server.router.HandleFunc("/webhook_events", server.PostWebhookEvent).Methods("Post").Name("PostWebhookEvent")
// }

// func eventHandler(w http.ResponseWriter, r *http.Request) {
// 	ctx := context.TODO()

// 	event, err := getAndValidatePayload(r)
// 	if err != nil {
// 		log.Printf("Failed to parse payload: %v", err)
// 		http.Error(w, "Invalid payload", http.StatusBadRequest)
// 		return
// 	}

// 	// Possible event types
// 	switch event := event.(type) {
// 	case *github.CheckSuiteEvent:
// 		log.Printf("NOT Handling CheckSuiteEvent %s for check suite: %d", event.GetAction(), event.GetCheckSuite().GetID())
// 		// client := authenticateGitHub(ctx, event.GetInstallation().GetID())
// 		// processCheckSuiteEvent(client, event)
// 	case *github.CheckRunEvent:
// 		log.Printf("NOT Handling CheckRunEvent %s for check run: %d", event.GetAction(), event.GetCheckRun().GetID())
// 		// client := authenticateGitHub(ctx, event.GetInstallation().GetID())
// 		// processCheckRunEvent(client, event)
// 	case *github.IssueCommentEvent:
// 		log.Printf("NOT Handling IssueCommentEvent %s on PR #%d for comment: %d", event.GetAction(), event.GetIssue().GetNumber(), event.GetComment().GetID())
// 		// client := authenticateGitHub(ctx, event.GetInstallation().GetID())
// 		// processIssueCommentEvent(client, event)
// 	case *github.PullRequestEvent:
// 		log.Printf("Handling PullRequestEvent %s for PR %d on repo %s", event.GetAction(), event.GetNumber(), event.GetRepo().GetName())
// 		client := authenticateGitHub(ctx, event.GetInstallation().GetID())
// 		processPullRequestEvent(ctx, client, event)
// 	}

// 	w.WriteHeader(http.StatusOK)
// }

// // PostPullRequest adds a PullRequest to the collection
// func (server *Server) PostWebhookEvent(w http.ResponseWriter, req *http.Request) {
// 	var widget Widget
// 	decoder := json.NewDecoder(req.Body)

// 	err := decoder.Decode(&widget)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	if len(widget.Name) == 0 {
// 		err = errors.New("Widget name can't be blank.")
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	// not thread safe...
// 	widget.Id = int64(len(AllWidgets))
// 	AllWidgets = append(AllWidgets, widget)

// 	widgetJSON, err := json.Marshal(widget)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")

// 	w.WriteHeader(http.StatusCreated)
// 	fmt.Fprintf(w, string(widgetJSON))

// }

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
