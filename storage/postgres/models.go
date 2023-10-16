// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package postgres

import (
	"time"
)

type Approval struct {
	ID         int32     `json:"id"`
	Uuid       string    `json:"uuid"`
	PrID       int32     `json:"pr_id"`
	Sha        string    `json:"sha"`
	ApprovedOn time.Time `json:"approved_on"`
}

type GhUser struct {
	ID    int32  `json:"id"`
	Login string `json:"login"`
}

type PullRequest struct {
	ID       int32 `json:"id"`
	RepoID   int32 `json:"repo_id"`
	PrID     int32 `json:"pr_id"`
	PrNumber int32 `json:"pr_number"`
	OpenedBy int32 `json:"opened_by"`
	IsMerged bool  `json:"is_merged"`
}

type PullRequestAction struct {
	Name string `json:"name"`
}

type PullRequestEvent struct {
	ID          int32     `json:"id"`
	PrID        int32     `json:"pr_id"`
	Action      string    `json:"action"`
	Sha         string    `json:"sha"`
	IsMerged    bool      `json:"is_merged"`
	LastUpdated time.Time `json:"last_updated"`
}

type Repo struct {
	ID   int32  `json:"id"`
	Org  string `json:"org"`
	Name string `json:"name"`
}
