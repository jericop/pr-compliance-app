// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package postgres

import (
	"time"
)

type Approval struct {
	ID          int32     `json:"id"`
	SchemaID    int32     `json:"schema_id"`
	Uuid        string    `json:"uuid"`
	PrID        int32     `json:"pr_id"`
	Sha         string    `json:"sha"`
	IsApproved  bool      `json:"is_approved"`
	LastUpdated time.Time `json:"last_updated"`
}

type ApprovalSchema struct {
	ID            int32  `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	StatusContext string `json:"status_context"`
	StatusTitle   string `json:"status_title"`
}

type ApprovalYesAnswer struct {
	ApprovalID int32 `json:"approval_id"`
	QuestionID int32 `json:"question_id"`
}

type ApprovalYesnoQuestion struct {
	ID           int32  `json:"id"`
	SchemaID     int32  `json:"schema_id"`
	QuestionText string `json:"question_text"`
}

type DefaultApprovalSchema struct {
	ID       int32 `json:"id"`
	SchemaID int32 `json:"schema_id"`
}

type GhUser struct {
	ID    int32  `json:"id"`
	Login string `json:"login"`
}

type Installation struct {
	ID int32 `json:"id"`
}

type PullRequest struct {
	ID             int32 `json:"id"`
	RepoID         int32 `json:"repo_id"`
	PrID           int32 `json:"pr_id"`
	PrNumber       int32 `json:"pr_number"`
	OpenedBy       int32 `json:"opened_by"`
	InstallationID int32 `json:"installation_id"`
	IsMerged       bool  `json:"is_merged"`
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
