package api

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/google/go-github/v53/github"
	"github.com/google/uuid"
	"github.com/jericop/pr-compliance-app/storage/postgres"
)

const (
	statusContext   = "Pull Request Compliance"
	statusTitle     = "User Review Required"
)

func (server *Server) AddWebhookEventsRoutes() {
	server.router.HandleFunc("/webhook_events", server.PostWebhookEvent).Methods("Post").Name("PostWebhookEvent")
}

func (server *Server) PostWebhookEvent(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()

	// event, err := server.validatWebhookRequest(r)
	event, err := server.githubFactory.ValidatWebhookRequest(r)
	if err != nil {
		fmt.Printf("validatWebhookRequest err: %v\n", err)
		http.Error(w, fmt.Sprintf("webhook event validation error %v", err), http.StatusBadRequest)
		return
	}

	// Possible event types
	switch event := event.(type) {
	case *github.PullRequestEvent:
		log.Printf("Handling PullRequestEvent %s for PR %d on repo %s", event.GetAction(), event.GetNumber(), event.GetRepo().GetName())
		// // client, err := server.getInstallationGitHubClient(ctx, event.GetInstallation().GetID())
		// // if err != nil {
		// // 	http.Error(w, fmt.Sprintf("github client error %v", err), http.StatusBadRequest)
		// // 	return
		// // }
		// fmt.Printf("client (%T)", client)
		// fmt.Printf("server.githubAppId (%T) %v", server.githubAppId, server.githubAppId)
		// fmt.Printf("Installation ID (%T) %v", client, client)
		err = server.processPullRequestEvent(ctx, event)
		if err != nil {
			log.Printf("error processing github pull request event %v", err)
			http.Error(w, fmt.Sprintf("error processing github pull request event %v", err), http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

type PullRequestEventCreateParams struct {
	GithubUser       postgres.CreateGithubUserParams
	Repo             postgres.CreateRepoParams
	PullRequest      postgres.CreatePullRequestParams
	PullRequestEvent postgres.CreatePullRequestEventParams
	Approval         postgres.CreateApprovalParams
}

func getCreateParamsFromEvent(event *github.PullRequestEvent) PullRequestEventCreateParams {
	orgName := "" // blank for non-organization accounts
	if event.Organization != nil {
		log.Printf("nil event.Organization %#v", event.Organization)
		orgName = *event.Organization.Name
	}

	return PullRequestEventCreateParams{
		Repo: postgres.CreateRepoParams{
			Org:  orgName,
			Name: event.Repo.GetName(),
			ID:   int32(event.Repo.GetID()),
		},
		PullRequest: postgres.CreatePullRequestParams{
			RepoID:   int32(event.Repo.GetID()),
			PrID:     int32(event.PullRequest.GetID()),
			PrNumber: int32(event.PullRequest.GetNumber()),
			OpenedBy: int32(*event.GetSender().ID),
		},
		PullRequestEvent: postgres.CreatePullRequestEventParams{
			PrID:     int32(event.PullRequest.GetID()),
			Action:   event.GetAction(),
			Sha:      *event.GetPullRequest().GetHead().SHA,
			IsMerged: event.PullRequest.GetMerged(),
		},
		GithubUser: postgres.CreateGithubUserParams{
			ID:    int32(*event.GetSender().ID),
			Login: *event.GetSender().Login,
		},
		Approval: postgres.CreateApprovalParams{
			Uuid:       uuid.New().String(),
			PrID:       int32(event.PullRequest.GetID()),
			Sha:        *event.GetPullRequest().GetHead().SHA,
			IsApproved: false,
		},
	}
}

// Handles a PullRequestEvent
func (server *Server) processPullRequestEvent(ctx context.Context, event *github.PullRequestEvent) error {
	// Params for creating db items populated from event info
	createParams := getCreateParamsFromEvent(event)

	switch event.GetAction() {
	case "opened", "synchronize", "reopened", "closed":
		if err := server.GetOrCreatePullRequestAction(ctx, event.GetAction()); err != nil {
			return err
		}

		ghUser, err := server.GetOrCreateGithubUser(ctx, createParams.GithubUser)
		if err != nil {
			return err
		}

		repo, err := server.GetOrCreateRepo(ctx, createParams.Repo)
		if err != nil {
			return err
		}

		pr, err := server.GetOrCreatePullRequest(ctx, createParams.PullRequest)
		if err != nil {
			return err
		}

		_, err = server.querier.CreatePullRequestEvent(ctx, createParams.PullRequestEvent)
		if err != nil {
			return err
		}

		approval, err := server.GetOrCreateApproval(ctx, createParams.Approval)
		if err != nil {
			return err
		}

		log.Printf("Creating a commit status for pull request %s/%d created by %s", repo.Name, pr.PrNumber, ghUser.Login)
		log.Printf("approval %v", approval)

		// client, err = server.githubClientFactory.NewClient(ctx, event.GetInstallation().GetID())
		// if err != nil {
		// 	return fmt.Errorf("github client error %v", err)
		// }

		// client, err := server.getInstallationGitHubClient(ctx, event.GetInstallation().GetID())
		// if err != nil {
		// 	return err
		// }

		client, err := server.githubFactory.NewInstallationClient(ctx, event.GetInstallation().GetID())
		if err != nil {
			return fmt.Errorf("github client error %v", err)
		}

		fmt.Printf("client (%T)", client)
		fmt.Printf("server.githubAppId (%T) %v", server.githubAppId, server.githubAppId)
		fmt.Printf("Installation ID (%T) %v", client, client)

		failedStatus := &github.RepoStatus{
			// TODO: Get these fields from the database at startup and then use them for all requests
			Context:     github.String(statusContext),
			Description: github.String(statusTitle),
			TargetURL:   github.String(fmt.Sprintf("http://localhost:8080/approval/%s", approval.Uuid)),
			// TargetURL:   github.String(fmt.Sprintf("https://localhost:8080/approval?id=%s", approval.Uuid)),
			State: github.String("error"), // "error" or "failure" show up as a red X
		}

		status, response, err := client.Repositories.CreateStatus(
			ctx, createParams.GithubUser.Login, createParams.Repo.Name, createParams.PullRequestEvent.Sha, failedStatus,
		)

		if err == nil {
			log.Printf("created repo status (check) for pull request %s/%d", repo.Name, pr.PrNumber)
		} else {
			return fmt.Errorf("failed to create repo status (check) for pull request %s/%d: status: %v, response: %v, err: %v", repo.Name, pr.PrNumber, status, response, err)
		}

		// case "closed":
		// 	log.Printf("Pull request %s/%d created by %s has been closed.", repoName, prNum, owner)
		// 	// TODO: Add the event to an activity table with the status of whether it was merged or not. This allows for multiple closed/reopened events.
		// 	// See: https://docs.github.com/en/rest/pulls/pulls?apiVersion=2022-11-28#check-if-a-pull-request-has-been-merged
	}
	return nil
}

func (server *Server) GetOrCreateGithubUser(ctx context.Context, p postgres.CreateGithubUserParams) (postgres.GhUser, error) {
	log.Printf("GetOrCreateGithubUser(ctx, %#v)", p)

	ghUser, err := server.querier.GetGithubUser(ctx, p.ID)
	if err != nil {
		log.Printf("Error getting ghUser %#v: %v", p, err)
		log.Printf("Creating ghUser %#v", p)

		ghUser, err = server.querier.CreateGithubUser(ctx, p)
		if err != nil {
			log.Printf("Error creating ghUser %#v: %v", p, err)
			return postgres.GhUser{}, err
		}
	}

	return ghUser, nil
}

func (server *Server) GetOrCreateRepo(ctx context.Context, p postgres.CreateRepoParams) (postgres.Repo, error) {
	log.Printf("GetOrCreateRepo(ctx, %#v)", p)

	repo, err := server.querier.GetRepo(ctx, p.ID)
	if err != nil {
		log.Printf("Error getting repo %#v: %v", p, err)
		log.Printf("Creating repo %#v", p)

		repo, err = server.querier.CreateRepo(ctx, p)
		if err != nil {
			log.Printf("Error creating repo %#v: %v", p, err)
			return postgres.Repo{}, err
		}
	}

	return repo, nil
}

func (server *Server) GetOrCreatePullRequest(ctx context.Context, p postgres.CreatePullRequestParams) (postgres.PullRequest, error) {
	log.Printf("GetOrCreatePullRequest(ctx, %#v)", p)

	queryParams := postgres.GetPullRequestByRepoIdPrIdParams{
		RepoID: p.RepoID,
		PrID:   p.PrID,
	}

	pr, err := server.querier.GetPullRequestByRepoIdPrId(ctx, queryParams)
	if err != nil {
		log.Printf("Error getting pr %#v : %v", queryParams, err)
		log.Printf("Creating pr %#v", p)

		pr, err = server.querier.CreatePullRequest(ctx, p)
		if err != nil {
			log.Printf("Error creating pr %#v: %v", p, err)
			return postgres.PullRequest{}, err
		}
	}

	return pr, nil
}

func (server *Server) GetOrCreateApproval(ctx context.Context, p postgres.CreateApprovalParams) (postgres.Approval, error) {
	log.Printf("GetOrCreateApproval(ctx, %#v)", p)

	queryParams := postgres.GetApprovalByPrIDShaParams{
		PrID: p.PrID,
		Sha:  p.Sha,
	}

	approval, err := server.querier.GetApprovalByPrIDSha(ctx, queryParams)
	if err != nil {
		log.Printf("Error getting approval %#v: %v", p, err)
		log.Printf("Creating approval %#v", p)

		approval, err = server.querier.CreateApproval(ctx, p)
		if err != nil {
			log.Printf("Error creating approval %#v: %v", p, err)
			return postgres.Approval{}, err
		}
	}

	return approval, nil
}

func (server *Server) GetOrCreatePullRequestAction(ctx context.Context, name string) error {
	log.Printf("GetOrCreatePullRequestAction(ctx, %s)", name)

	_, known := server.KnownPullRequestActions[name]
	if known {
		log.Printf("known action %s", name)
		return nil
	}

	if _, err := server.querier.GetPullRequestAction(ctx, name); err != nil {
		log.Printf("action '%s' does not exist: %v", name, err)
		log.Printf("Creating pull request (event) action %s", name)

		_, err = server.querier.CreatePullRequestAction(ctx, name)
		if err != nil {
			log.Printf("Error creating pull_request_action %s: %v", name, err)
			return err
		}
		server.KnownPullRequestActions[name] = struct{}{}
	}

	return nil
}
