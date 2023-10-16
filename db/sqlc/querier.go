// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package db

import (
	"context"
)

type Querier interface {
	CreateApproval(ctx context.Context, arg CreateApprovalParams) (Approval, error)
	CreateGithubUser(ctx context.Context, arg CreateGithubUserParams) (GhUser, error)
	CreatePullRequest(ctx context.Context, arg CreatePullRequestParams) (PullRequest, error)
	CreatePullRequestAction(ctx context.Context, name string) (string, error)
	CreatePullRequestEvent(ctx context.Context, arg CreatePullRequestEventParams) (PullRequestEvent, error)
	CreateRepo(ctx context.Context, arg CreateRepoParams) (Repo, error)
	DeleteApproval(ctx context.Context, id int32) error
	DeleteGithubUser(ctx context.Context, id int32) error
	DeletePullRequest(ctx context.Context, id int32) error
	DeletePullRequestAction(ctx context.Context, name string) error
	DeletePullRequestEvent(ctx context.Context, id int32) error
	DeleteRepo(ctx context.Context, id int32) error
	GetApprovalById(ctx context.Context, id int32) (Approval, error)
	GetApprovalByUuid(ctx context.Context, uuid string) (Approval, error)
	GetApprovals(ctx context.Context) ([]Approval, error)
	GetGithubUser(ctx context.Context, id int32) (GhUser, error)
	GetGithubUsers(ctx context.Context) ([]GhUser, error)
	GetPullRequest(ctx context.Context, id int32) (PullRequest, error)
	GetPullRequestAction(ctx context.Context, name string) (string, error)
	GetPullRequestActions(ctx context.Context) ([]string, error)
	GetPullRequestEvent(ctx context.Context, id int32) (PullRequestEvent, error)
	GetPullRequestEvents(ctx context.Context) ([]PullRequestEvent, error)
	GetPullRequestForUpdate(ctx context.Context, id int32) (PullRequest, error)
	GetPullRequests(ctx context.Context) ([]PullRequest, error)
	GetRepo(ctx context.Context, id int32) (Repo, error)
	GetRepoForUpdate(ctx context.Context, id int32) (Repo, error)
	GetRepos(ctx context.Context) ([]Repo, error)
	UpdatePullRequestIsMerged(ctx context.Context, arg UpdatePullRequestIsMergedParams) (PullRequest, error)
	UpdateRepoName(ctx context.Context, arg UpdateRepoNameParams) (Repo, error)
	UpdateRepoOrg(ctx context.Context, arg UpdateRepoOrgParams) (Repo, error)
}

var _ Querier = (*Queries)(nil)
