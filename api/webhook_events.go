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

func (server *Server) AddWebhookEventsRoutes() {
	server.router.HandleFunc("/webhook_events", server.PostWebhookEvent).Methods("Post").Name("PostWebhookEvent")
}

func (server *Server) PostWebhookEvent(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	event, err := server.githubFactory.ValidatWebhookRequest(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("webhook event validation error %v", err), http.StatusBadRequest)
		return
	}

	// Possible event types
	switch event := event.(type) {
	case *github.PullRequestEvent:
		log.Printf("Handling PullRequestEvent %s for PR %d on repo %s", event.GetAction(), event.GetNumber(), event.GetRepo().GetName())

		err := server.dbTxFactory.ExecWithTx(ctx, func(q postgres.Querier) error {
			return server.processPullRequestEvent(ctx, q, event)
		})

		if err != nil {
			log.Printf("error processing github pull request event %v", err)
			http.Error(w, fmt.Sprintf("error processing github pull request event %v", err), http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (server *Server) processPullRequestEvent(ctx context.Context, querier postgres.Querier, event *github.PullRequestEvent) error {
	// Params for creating db items populated from event info and schema id
	p := getCreateParamsFromEvent(event, server.schema.ID)

	switch event.GetAction() {
	case "opened", "synchronize", "reopened":
		if err := server.GetOrCreateInstallation(ctx, querier, p.InstallationID); err != nil {
			return err
		}

		if err := server.GetOrCreatePullRequestAction(ctx, querier, event.GetAction()); err != nil {
			return err
		}

		ghUser, err := server.GetOrCreateGithubUser(ctx, querier, p.GithubUser)
		if err != nil {
			return err
		}

		repo, err := server.GetOrCreateRepo(ctx, querier, p.Repo)
		if err != nil {
			return err
		}

		pr, err := server.GetOrCreatePullRequest(ctx, querier, p.PullRequest)
		if err != nil {
			return err
		}

		_, err = querier.CreatePullRequestEvent(ctx, p.PullRequestEvent)
		if err != nil {
			return err
		}

		approval, err := server.GetOrCreateApproval(ctx, querier, p.Approval)
		if err != nil {
			return err
		}

		if approval.IsApproved {
			log.Printf("Removing approval for pull request %s/%d with sha %s", p.Repo.Name, p.PullRequest.PrNumber, p.PullRequestEvent.Sha)

			p := postgres.UpdateApprovalByUuidParams{
				Uuid:       approval.Uuid,
				IsApproved: false,
			}

			if err := querier.UpdateApprovalByUuid(ctx, p); err != nil {
				return err
			}
		}

		client, err := server.githubFactory.NewInstallationClient(ctx, event.GetInstallation().GetID())
		if err != nil {
			return fmt.Errorf("github client error %v", err)
		}

		s := &github.RepoStatus{
			Context:     github.String(server.schema.StatusContext),
			Description: github.String(server.schema.StatusTitle),
			TargetURL:   github.String(fmt.Sprintf("%s/%s", server.frontEndUrl, approval.Uuid)),
			// TargetURL: github.String(fmt.Sprintf("%s?id=%s", server.frontEndUrl, approval.Uuid)),
			State: github.String("error"),
		}

		log.Printf("Creating an error commit status for pull request %s/%d created by %s", repo.Name, pr.PrNumber, ghUser.Login)

		_, _, err = client.Repositories.CreateStatus(ctx, p.GithubUser.Login, p.Repo.Name, p.PullRequestEvent.Sha, s)
		if err != nil {
			return err
		}
	}
	return nil
}

type PullRequestEventCreateParams struct {
	InstallationID   int32
	GithubUser       postgres.CreateGithubUserParams
	Repo             postgres.CreateRepoParams
	PullRequest      postgres.CreatePullRequestParams
	PullRequestEvent postgres.CreatePullRequestEventParams
	Approval         postgres.CreateApprovalParams
}

func getCreateParamsFromEvent(event *github.PullRequestEvent, schemaId int32) PullRequestEventCreateParams {
	orgName := "" // blank for non-organization accounts
	if event.Organization != nil {
		orgName = *event.Organization.Name
	}

	return PullRequestEventCreateParams{
		InstallationID: int32(event.GetInstallation().GetID()),
		GithubUser: postgres.CreateGithubUserParams{
			ID:    int32(*event.GetSender().ID),
			Login: *event.GetSender().Login,
		},
		Repo: postgres.CreateRepoParams{
			Org:  orgName,
			Name: event.Repo.GetName(),
			ID:   int32(event.Repo.GetID()),
		},
		PullRequest: postgres.CreatePullRequestParams{
			RepoID:         int32(event.Repo.GetID()),
			PrID:           int32(event.PullRequest.GetID()),
			PrNumber:       int32(event.PullRequest.GetNumber()),
			OpenedBy:       int32(*event.GetSender().ID),
			InstallationID: int32(event.GetInstallation().GetID()),
		},
		PullRequestEvent: postgres.CreatePullRequestEventParams{
			PrID:     int32(event.PullRequest.GetID()),
			Action:   event.GetAction(),
			Sha:      *event.GetPullRequest().GetHead().SHA,
			IsMerged: event.PullRequest.GetMerged(),
		},
		Approval: postgres.CreateApprovalParams{
			SchemaID:   int32(schemaId),
			Uuid:       uuid.New().String(),
			PrID:       int32(event.PullRequest.GetID()),
			Sha:        *event.GetPullRequest().GetHead().SHA,
			IsApproved: false,
		},
	}
}

func (server *Server) GetOrCreateInstallation(ctx context.Context, querier postgres.Querier, id int32) error {
	if _, err := querier.GetInstallation(ctx, id); err != nil {
		log.Printf("Creating installation %d", id)

		_, err = querier.CreateInstallation(ctx, id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (server *Server) GetOrCreateGithubUser(ctx context.Context, querier postgres.Querier, p postgres.CreateGithubUserParams) (postgres.GhUser, error) {
	ghUser, err := querier.GetGithubUser(ctx, p.ID)
	if err != nil {
		log.Printf("Creating ghUser %#v", p)

		ghUser, err = querier.CreateGithubUser(ctx, p)
		if err != nil {
			return postgres.GhUser{}, err
		}
	}

	return ghUser, nil
}

func (server *Server) GetOrCreateRepo(ctx context.Context, querier postgres.Querier, p postgres.CreateRepoParams) (postgres.Repo, error) {
	repo, err := querier.GetRepo(ctx, p.ID)
	if err != nil {
		log.Printf("Creating repo %#v", p)

		repo, err = querier.CreateRepo(ctx, p)
		if err != nil {
			return postgres.Repo{}, err
		}
	}

	return repo, nil
}

func (server *Server) GetOrCreatePullRequest(ctx context.Context, querier postgres.Querier, p postgres.CreatePullRequestParams) (postgres.PullRequest, error) {
	queryParams := postgres.GetPullRequestByRepoIdPrIdParams{
		RepoID: p.RepoID,
		PrID:   p.PrID,
	}

	pr, err := querier.GetPullRequestByRepoIdPrId(ctx, queryParams)
	if err != nil {
		log.Printf("Creating pr %#v", p)

		pr, err = querier.CreatePullRequest(ctx, p)
		if err != nil {
			return postgres.PullRequest{}, err
		}
	}

	return pr, nil
}

func (server *Server) GetOrCreateApproval(ctx context.Context, querier postgres.Querier, p postgres.CreateApprovalParams) (postgres.Approval, error) {
	queryParams := postgres.GetApprovalByPrIDShaParams{
		PrID: p.PrID,
		Sha:  p.Sha,
	}

	approval, err := querier.GetApprovalByPrIDSha(ctx, queryParams)
	if err != nil {
		log.Printf("Creating approval %#v", p)

		approval, err = querier.CreateApproval(ctx, p)
		if err != nil {
			return postgres.Approval{}, err
		}
	}

	return approval, nil
}

func (server *Server) GetOrCreatePullRequestAction(ctx context.Context, querier postgres.Querier, name string) error {
	_, known := server.knownPullRequestActions[name]
	if known {
		return nil
	}

	if _, err := querier.GetPullRequestAction(ctx, name); err != nil {
		log.Printf("Creating pull request (event) action %s", name)

		_, err = querier.CreatePullRequestAction(ctx, name)
		if err != nil {
			return err
		}
		server.knownPullRequestActions[name] = struct{}{}
	}

	return nil
}
