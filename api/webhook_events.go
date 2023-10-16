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
	server.Router.HandleFunc("/webhook_events", server.PostWebhookEvent).Methods("Post").Name("PostWebhookEvent")
}

func (server *Server) PostWebhookEvent(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()

	event, err := server.validatWebhookRequest(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("webhook event validation error %v", err), http.StatusBadRequest)
		return
	}

	// Possible event types
	switch event := event.(type) {
	case *github.PullRequestEvent:
		log.Printf("Handling PullRequestEvent %s for PR %d on repo %s", event.GetAction(), event.GetNumber(), event.GetRepo().GetName())
		client, err := server.getInstallationGitHubClient(ctx, event.GetInstallation().GetID())
		if err != nil {
			http.Error(w, fmt.Sprintf("github client error %v", err), http.StatusBadRequest)
			return
		}
		fmt.Printf("client (%T)", client)
		fmt.Printf("server.githubAppId (%T) %v", server.githubAppId, server.githubAppId)
		fmt.Printf("Installation ID (%T) %v", client, client)
		err = server.processPullRequestEvent(ctx, client, event)
		if err != nil {
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
			PrID:     int32(event.Repo.GetID()),
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
			PrID:       int32(event.Repo.GetID()),
			Sha:        *event.GetPullRequest().GetHead().SHA,
			IsApproved: false,
		},
	}
}

// Handles a PullRequestEvent
func (server *Server) processPullRequestEvent(ctx context.Context, client *github.Client, event *github.PullRequestEvent) error {
	// Params for creating db items populated from event info
	createParams := getCreateParamsFromEvent(event)

	switch event.GetAction() {
	case "opened", "synchronize", "reopened", "closed":
		// tx, err := server.querier.(*postgres.Queries).
		// err := server.execTx(ctx, func(postgres.Querier) error {
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

		// _, err := server.GetOrCreateGithubUser(ctx, createParams.GithubUser)
		// if err != nil {
		// 	log.Printf("Error getting ghUser: %v", err)
		// 	return
		// }

		// repo, err := server.GetOrCreateRepo(ctx, createParams.Repo)
		// if err != nil {
		// 	log.Printf("Error getting repo: %v", err)
		// 	return
		// }

		// pr, err := server.GetOrCreatePullRequest(ctx, createParams.PullRequest)
		// if err != nil {
		// 	log.Printf("Error getting pr: %v", err)
		// 	return
		// }

		// prEvent, err := server.querier.CreatePullRequestEvent(ctx, createParams.PullRequestEvent)
		// if err != nil {
		// 	log.Printf("Error creating prEvent: %v", err)
		// }

		// approval, err := server.GetOrCreateApproval(ctx, createParams.Approval)
		// if err != nil {
		// 	log.Printf("Error creating prEvent: %v", err)
		// }

		// log.Printf("Creating a commit status for pull request %s/%d created by %s", repoName, prNum, owner)

		// failedStatus := &github.RepoStatus{
		// 	// TODO: Get these fields from the database at startup and then use them for all requests
		// 	Context:     github.String(statusContext),
		// 	Description: github.String(statusTitle),
		// 	// TargetURL: github.String("https://github.com"),
		// 	State: github.String("error"), // "error" or "failure" show up as a red X
		// }

		// status, response, err := client.Repositories.CreateStatus(ctx, owner, repo.Name, sha, failedStatus)

		// if err == nil {
		// 	log.Printf("Successfully created commit status for pull request %s/%d", repoName, prNum)
		// } else {
		// 	log.Printf("Failed to create commit status for pull request %s/%d", repoName, prNum)
		// 	log.Printf("status: %v", status)
		// 	log.Printf("response: %v", response) // You can keep track of rate limit usage from this
		// 	log.Printf("err: %v", err)
		// 	return
		// }

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
