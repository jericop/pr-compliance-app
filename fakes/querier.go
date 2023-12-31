package fakes

import (
	"context"
	"sync"

	"github.com/jericop/pr-compliance-app/storage/postgres"
)

type Querier struct {
	CreateApprovalCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Ctx context.Context
			Arg postgres.CreateApprovalParams
		}
		Returns struct {
			Approval postgres.Approval
			Error    error
		}
		Stub func(context.Context, postgres.CreateApprovalParams) (postgres.Approval, error)
	}
	CreateApprovalYesAnswerCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Ctx context.Context
			Arg postgres.CreateApprovalYesAnswerParams
		}
		Returns struct {
			ApprovalYesAnswer postgres.ApprovalYesAnswer
			Error             error
		}
		Stub func(context.Context, postgres.CreateApprovalYesAnswerParams) (postgres.ApprovalYesAnswer, error)
	}
	CreateApprovalYesAnswerByUuidCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Ctx context.Context
			Arg postgres.CreateApprovalYesAnswerByUuidParams
		}
		Returns struct {
			ApprovalYesAnswer postgres.ApprovalYesAnswer
			Error             error
		}
		Stub func(context.Context, postgres.CreateApprovalYesAnswerByUuidParams) (postgres.ApprovalYesAnswer, error)
	}
	CreateGithubUserCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Ctx context.Context
			Arg postgres.CreateGithubUserParams
		}
		Returns struct {
			GhUser postgres.GhUser
			Error  error
		}
		Stub func(context.Context, postgres.CreateGithubUserParams) (postgres.GhUser, error)
	}
	CreateInstallationCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Ctx context.Context
			Id  int32
		}
		Returns struct {
			Int32 int32
			Error error
		}
		Stub func(context.Context, int32) (int32, error)
	}
	CreatePullRequestCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Ctx context.Context
			Arg postgres.CreatePullRequestParams
		}
		Returns struct {
			PullRequest postgres.PullRequest
			Error       error
		}
		Stub func(context.Context, postgres.CreatePullRequestParams) (postgres.PullRequest, error)
	}
	CreatePullRequestActionCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Ctx  context.Context
			Name string
		}
		Returns struct {
			String string
			Error  error
		}
		Stub func(context.Context, string) (string, error)
	}
	CreatePullRequestEventCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Ctx context.Context
			Arg postgres.CreatePullRequestEventParams
		}
		Returns struct {
			PullRequestEvent postgres.PullRequestEvent
			Error            error
		}
		Stub func(context.Context, postgres.CreatePullRequestEventParams) (postgres.PullRequestEvent, error)
	}
	CreateRepoCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Ctx context.Context
			Arg postgres.CreateRepoParams
		}
		Returns struct {
			Repo  postgres.Repo
			Error error
		}
		Stub func(context.Context, postgres.CreateRepoParams) (postgres.Repo, error)
	}
	DeleteApprovalCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Ctx context.Context
			Id  int32
		}
		Returns struct {
			Error error
		}
		Stub func(context.Context, int32) error
	}
	DeleteApprovalYesAnswerCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Ctx context.Context
			Arg postgres.DeleteApprovalYesAnswerParams
		}
		Returns struct {
			Error error
		}
		Stub func(context.Context, postgres.DeleteApprovalYesAnswerParams) error
	}
	DeleteApprovalYesAnswerByUuidCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Ctx context.Context
			Arg postgres.DeleteApprovalYesAnswerByUuidParams
		}
		Returns struct {
			Error error
		}
		Stub func(context.Context, postgres.DeleteApprovalYesAnswerByUuidParams) error
	}
	DeleteGithubUserCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Ctx context.Context
			Id  int32
		}
		Returns struct {
			Error error
		}
		Stub func(context.Context, int32) error
	}
	DeleteInstallationCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Ctx context.Context
			Id  int32
		}
		Returns struct {
			Error error
		}
		Stub func(context.Context, int32) error
	}
	DeletePullRequestCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Ctx context.Context
			Id  int32
		}
		Returns struct {
			Error error
		}
		Stub func(context.Context, int32) error
	}
	DeletePullRequestActionCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Ctx  context.Context
			Name string
		}
		Returns struct {
			Error error
		}
		Stub func(context.Context, string) error
	}
	DeletePullRequestEventCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Ctx context.Context
			Id  int32
		}
		Returns struct {
			Error error
		}
		Stub func(context.Context, int32) error
	}
	DeleteRepoCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Ctx context.Context
			Id  int32
		}
		Returns struct {
			Error error
		}
		Stub func(context.Context, int32) error
	}
	GetApprovalByIdCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Ctx context.Context
			Id  int32
		}
		Returns struct {
			Approval postgres.Approval
			Error    error
		}
		Stub func(context.Context, int32) (postgres.Approval, error)
	}
	GetApprovalByPrIDShaCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Ctx context.Context
			Arg postgres.GetApprovalByPrIDShaParams
		}
		Returns struct {
			Approval postgres.Approval
			Error    error
		}
		Stub func(context.Context, postgres.GetApprovalByPrIDShaParams) (postgres.Approval, error)
	}
	GetApprovalByUuidCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Ctx  context.Context
			Uuid string
		}
		Returns struct {
			Approval postgres.Approval
			Error    error
		}
		Stub func(context.Context, string) (postgres.Approval, error)
	}
	GetApprovalSchemaByIdCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Ctx context.Context
			Id  int32
		}
		Returns struct {
			ApprovalSchema postgres.ApprovalSchema
			Error          error
		}
		Stub func(context.Context, int32) (postgres.ApprovalSchema, error)
	}
	GetApprovalsCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Ctx context.Context
		}
		Returns struct {
			ApprovalSlice []postgres.Approval
			Error         error
		}
		Stub func(context.Context) ([]postgres.Approval, error)
	}
	GetCreateStatusInputsFromApprovalUuidCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Ctx  context.Context
			Uuid string
		}
		Returns struct {
			GetCreateStatusInputsFromApprovalUuidRow postgres.GetCreateStatusInputsFromApprovalUuidRow
			Error                                    error
		}
		Stub func(context.Context, string) (postgres.GetCreateStatusInputsFromApprovalUuidRow, error)
	}
	GetDefaultApprovalSchemaCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Ctx context.Context
		}
		Returns struct {
			ApprovalSchema postgres.ApprovalSchema
			Error          error
		}
		Stub func(context.Context) (postgres.ApprovalSchema, error)
	}
	GetGithubUserCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Ctx context.Context
			Id  int32
		}
		Returns struct {
			GhUser postgres.GhUser
			Error  error
		}
		Stub func(context.Context, int32) (postgres.GhUser, error)
	}
	GetGithubUsersCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Ctx context.Context
		}
		Returns struct {
			GhUserSlice []postgres.GhUser
			Error       error
		}
		Stub func(context.Context) ([]postgres.GhUser, error)
	}
	GetInstallationCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Ctx context.Context
			Id  int32
		}
		Returns struct {
			Int32 int32
			Error error
		}
		Stub func(context.Context, int32) (int32, error)
	}
	GetInstallationsCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Ctx context.Context
		}
		Returns struct {
			Int32Slice []int32
			Error      error
		}
		Stub func(context.Context) ([]int32, error)
	}
	GetPullRequestActionCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Ctx  context.Context
			Name string
		}
		Returns struct {
			String string
			Error  error
		}
		Stub func(context.Context, string) (string, error)
	}
	GetPullRequestActionsCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Ctx context.Context
		}
		Returns struct {
			StringSlice []string
			Error       error
		}
		Stub func(context.Context) ([]string, error)
	}
	GetPullRequestByIdCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Ctx context.Context
			Id  int32
		}
		Returns struct {
			PullRequest postgres.PullRequest
			Error       error
		}
		Stub func(context.Context, int32) (postgres.PullRequest, error)
	}
	GetPullRequestByRepoIdPrIdCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Ctx context.Context
			Arg postgres.GetPullRequestByRepoIdPrIdParams
		}
		Returns struct {
			PullRequest postgres.PullRequest
			Error       error
		}
		Stub func(context.Context, postgres.GetPullRequestByRepoIdPrIdParams) (postgres.PullRequest, error)
	}
	GetPullRequestEventCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Ctx context.Context
			Id  int32
		}
		Returns struct {
			PullRequestEvent postgres.PullRequestEvent
			Error            error
		}
		Stub func(context.Context, int32) (postgres.PullRequestEvent, error)
	}
	GetPullRequestEventsCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Ctx context.Context
		}
		Returns struct {
			PullRequestEventSlice []postgres.PullRequestEvent
			Error                 error
		}
		Stub func(context.Context) ([]postgres.PullRequestEvent, error)
	}
	GetPullRequestsCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Ctx context.Context
		}
		Returns struct {
			PullRequestSlice []postgres.PullRequest
			Error            error
		}
		Stub func(context.Context) ([]postgres.PullRequest, error)
	}
	GetRepoCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Ctx context.Context
			Id  int32
		}
		Returns struct {
			Repo  postgres.Repo
			Error error
		}
		Stub func(context.Context, int32) (postgres.Repo, error)
	}
	GetReposCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Ctx context.Context
		}
		Returns struct {
			RepoSlice []postgres.Repo
			Error     error
		}
		Stub func(context.Context) ([]postgres.Repo, error)
	}
	GetSortedApprovalYesNoQuestionAnswersByUuidCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Ctx  context.Context
			Uuid string
		}
		Returns struct {
			GetSortedApprovalYesNoQuestionAnswersByUuidRowSlice []postgres.GetSortedApprovalYesNoQuestionAnswersByUuidRow
			Error                                               error
		}
		Stub func(context.Context, string) ([]postgres.GetSortedApprovalYesNoQuestionAnswersByUuidRow, error)
	}
	GetSortedApprovalYesNoQuestionsBySchemaIdCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Ctx  context.Context
			Name string
		}
		Returns struct {
			GetSortedApprovalYesNoQuestionsBySchemaIdRowSlice []postgres.GetSortedApprovalYesNoQuestionsBySchemaIdRow
			Error                                             error
		}
		Stub func(context.Context, string) ([]postgres.GetSortedApprovalYesNoQuestionsBySchemaIdRow, error)
	}
	UpdateApprovalByUuidCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Ctx context.Context
			Arg postgres.UpdateApprovalByUuidParams
		}
		Returns struct {
			Error error
		}
		Stub func(context.Context, postgres.UpdateApprovalByUuidParams) error
	}
	UpdatePullRequestIsMergedCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Ctx context.Context
			Arg postgres.UpdatePullRequestIsMergedParams
		}
		Returns struct {
			PullRequest postgres.PullRequest
			Error       error
		}
		Stub func(context.Context, postgres.UpdatePullRequestIsMergedParams) (postgres.PullRequest, error)
	}
	UpdateRepoNameCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Ctx context.Context
			Arg postgres.UpdateRepoNameParams
		}
		Returns struct {
			Repo  postgres.Repo
			Error error
		}
		Stub func(context.Context, postgres.UpdateRepoNameParams) (postgres.Repo, error)
	}
	UpdateRepoOrgCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Ctx context.Context
			Arg postgres.UpdateRepoOrgParams
		}
		Returns struct {
			Repo  postgres.Repo
			Error error
		}
		Stub func(context.Context, postgres.UpdateRepoOrgParams) (postgres.Repo, error)
	}
}

func (f *Querier) CreateApproval(param1 context.Context, param2 postgres.CreateApprovalParams) (postgres.Approval, error) {
	f.CreateApprovalCall.mutex.Lock()
	defer f.CreateApprovalCall.mutex.Unlock()
	f.CreateApprovalCall.CallCount++
	f.CreateApprovalCall.Receives.Ctx = param1
	f.CreateApprovalCall.Receives.Arg = param2
	if f.CreateApprovalCall.Stub != nil {
		return f.CreateApprovalCall.Stub(param1, param2)
	}
	return f.CreateApprovalCall.Returns.Approval, f.CreateApprovalCall.Returns.Error
}
func (f *Querier) CreateApprovalYesAnswer(param1 context.Context, param2 postgres.CreateApprovalYesAnswerParams) (postgres.ApprovalYesAnswer, error) {
	f.CreateApprovalYesAnswerCall.mutex.Lock()
	defer f.CreateApprovalYesAnswerCall.mutex.Unlock()
	f.CreateApprovalYesAnswerCall.CallCount++
	f.CreateApprovalYesAnswerCall.Receives.Ctx = param1
	f.CreateApprovalYesAnswerCall.Receives.Arg = param2
	if f.CreateApprovalYesAnswerCall.Stub != nil {
		return f.CreateApprovalYesAnswerCall.Stub(param1, param2)
	}
	return f.CreateApprovalYesAnswerCall.Returns.ApprovalYesAnswer, f.CreateApprovalYesAnswerCall.Returns.Error
}
func (f *Querier) CreateApprovalYesAnswerByUuid(param1 context.Context, param2 postgres.CreateApprovalYesAnswerByUuidParams) (postgres.ApprovalYesAnswer, error) {
	f.CreateApprovalYesAnswerByUuidCall.mutex.Lock()
	defer f.CreateApprovalYesAnswerByUuidCall.mutex.Unlock()
	f.CreateApprovalYesAnswerByUuidCall.CallCount++
	f.CreateApprovalYesAnswerByUuidCall.Receives.Ctx = param1
	f.CreateApprovalYesAnswerByUuidCall.Receives.Arg = param2
	if f.CreateApprovalYesAnswerByUuidCall.Stub != nil {
		return f.CreateApprovalYesAnswerByUuidCall.Stub(param1, param2)
	}
	return f.CreateApprovalYesAnswerByUuidCall.Returns.ApprovalYesAnswer, f.CreateApprovalYesAnswerByUuidCall.Returns.Error
}
func (f *Querier) CreateGithubUser(param1 context.Context, param2 postgres.CreateGithubUserParams) (postgres.GhUser, error) {
	f.CreateGithubUserCall.mutex.Lock()
	defer f.CreateGithubUserCall.mutex.Unlock()
	f.CreateGithubUserCall.CallCount++
	f.CreateGithubUserCall.Receives.Ctx = param1
	f.CreateGithubUserCall.Receives.Arg = param2
	if f.CreateGithubUserCall.Stub != nil {
		return f.CreateGithubUserCall.Stub(param1, param2)
	}
	return f.CreateGithubUserCall.Returns.GhUser, f.CreateGithubUserCall.Returns.Error
}
func (f *Querier) CreateInstallation(param1 context.Context, param2 int32) (int32, error) {
	f.CreateInstallationCall.mutex.Lock()
	defer f.CreateInstallationCall.mutex.Unlock()
	f.CreateInstallationCall.CallCount++
	f.CreateInstallationCall.Receives.Ctx = param1
	f.CreateInstallationCall.Receives.Id = param2
	if f.CreateInstallationCall.Stub != nil {
		return f.CreateInstallationCall.Stub(param1, param2)
	}
	return f.CreateInstallationCall.Returns.Int32, f.CreateInstallationCall.Returns.Error
}
func (f *Querier) CreatePullRequest(param1 context.Context, param2 postgres.CreatePullRequestParams) (postgres.PullRequest, error) {
	f.CreatePullRequestCall.mutex.Lock()
	defer f.CreatePullRequestCall.mutex.Unlock()
	f.CreatePullRequestCall.CallCount++
	f.CreatePullRequestCall.Receives.Ctx = param1
	f.CreatePullRequestCall.Receives.Arg = param2
	if f.CreatePullRequestCall.Stub != nil {
		return f.CreatePullRequestCall.Stub(param1, param2)
	}
	return f.CreatePullRequestCall.Returns.PullRequest, f.CreatePullRequestCall.Returns.Error
}
func (f *Querier) CreatePullRequestAction(param1 context.Context, param2 string) (string, error) {
	f.CreatePullRequestActionCall.mutex.Lock()
	defer f.CreatePullRequestActionCall.mutex.Unlock()
	f.CreatePullRequestActionCall.CallCount++
	f.CreatePullRequestActionCall.Receives.Ctx = param1
	f.CreatePullRequestActionCall.Receives.Name = param2
	if f.CreatePullRequestActionCall.Stub != nil {
		return f.CreatePullRequestActionCall.Stub(param1, param2)
	}
	return f.CreatePullRequestActionCall.Returns.String, f.CreatePullRequestActionCall.Returns.Error
}
func (f *Querier) CreatePullRequestEvent(param1 context.Context, param2 postgres.CreatePullRequestEventParams) (postgres.PullRequestEvent, error) {
	f.CreatePullRequestEventCall.mutex.Lock()
	defer f.CreatePullRequestEventCall.mutex.Unlock()
	f.CreatePullRequestEventCall.CallCount++
	f.CreatePullRequestEventCall.Receives.Ctx = param1
	f.CreatePullRequestEventCall.Receives.Arg = param2
	if f.CreatePullRequestEventCall.Stub != nil {
		return f.CreatePullRequestEventCall.Stub(param1, param2)
	}
	return f.CreatePullRequestEventCall.Returns.PullRequestEvent, f.CreatePullRequestEventCall.Returns.Error
}
func (f *Querier) CreateRepo(param1 context.Context, param2 postgres.CreateRepoParams) (postgres.Repo, error) {
	f.CreateRepoCall.mutex.Lock()
	defer f.CreateRepoCall.mutex.Unlock()
	f.CreateRepoCall.CallCount++
	f.CreateRepoCall.Receives.Ctx = param1
	f.CreateRepoCall.Receives.Arg = param2
	if f.CreateRepoCall.Stub != nil {
		return f.CreateRepoCall.Stub(param1, param2)
	}
	return f.CreateRepoCall.Returns.Repo, f.CreateRepoCall.Returns.Error
}
func (f *Querier) DeleteApproval(param1 context.Context, param2 int32) error {
	f.DeleteApprovalCall.mutex.Lock()
	defer f.DeleteApprovalCall.mutex.Unlock()
	f.DeleteApprovalCall.CallCount++
	f.DeleteApprovalCall.Receives.Ctx = param1
	f.DeleteApprovalCall.Receives.Id = param2
	if f.DeleteApprovalCall.Stub != nil {
		return f.DeleteApprovalCall.Stub(param1, param2)
	}
	return f.DeleteApprovalCall.Returns.Error
}
func (f *Querier) DeleteApprovalYesAnswer(param1 context.Context, param2 postgres.DeleteApprovalYesAnswerParams) error {
	f.DeleteApprovalYesAnswerCall.mutex.Lock()
	defer f.DeleteApprovalYesAnswerCall.mutex.Unlock()
	f.DeleteApprovalYesAnswerCall.CallCount++
	f.DeleteApprovalYesAnswerCall.Receives.Ctx = param1
	f.DeleteApprovalYesAnswerCall.Receives.Arg = param2
	if f.DeleteApprovalYesAnswerCall.Stub != nil {
		return f.DeleteApprovalYesAnswerCall.Stub(param1, param2)
	}
	return f.DeleteApprovalYesAnswerCall.Returns.Error
}
func (f *Querier) DeleteApprovalYesAnswerByUuid(param1 context.Context, param2 postgres.DeleteApprovalYesAnswerByUuidParams) error {
	f.DeleteApprovalYesAnswerByUuidCall.mutex.Lock()
	defer f.DeleteApprovalYesAnswerByUuidCall.mutex.Unlock()
	f.DeleteApprovalYesAnswerByUuidCall.CallCount++
	f.DeleteApprovalYesAnswerByUuidCall.Receives.Ctx = param1
	f.DeleteApprovalYesAnswerByUuidCall.Receives.Arg = param2
	if f.DeleteApprovalYesAnswerByUuidCall.Stub != nil {
		return f.DeleteApprovalYesAnswerByUuidCall.Stub(param1, param2)
	}
	return f.DeleteApprovalYesAnswerByUuidCall.Returns.Error
}
func (f *Querier) DeleteGithubUser(param1 context.Context, param2 int32) error {
	f.DeleteGithubUserCall.mutex.Lock()
	defer f.DeleteGithubUserCall.mutex.Unlock()
	f.DeleteGithubUserCall.CallCount++
	f.DeleteGithubUserCall.Receives.Ctx = param1
	f.DeleteGithubUserCall.Receives.Id = param2
	if f.DeleteGithubUserCall.Stub != nil {
		return f.DeleteGithubUserCall.Stub(param1, param2)
	}
	return f.DeleteGithubUserCall.Returns.Error
}
func (f *Querier) DeleteInstallation(param1 context.Context, param2 int32) error {
	f.DeleteInstallationCall.mutex.Lock()
	defer f.DeleteInstallationCall.mutex.Unlock()
	f.DeleteInstallationCall.CallCount++
	f.DeleteInstallationCall.Receives.Ctx = param1
	f.DeleteInstallationCall.Receives.Id = param2
	if f.DeleteInstallationCall.Stub != nil {
		return f.DeleteInstallationCall.Stub(param1, param2)
	}
	return f.DeleteInstallationCall.Returns.Error
}
func (f *Querier) DeletePullRequest(param1 context.Context, param2 int32) error {
	f.DeletePullRequestCall.mutex.Lock()
	defer f.DeletePullRequestCall.mutex.Unlock()
	f.DeletePullRequestCall.CallCount++
	f.DeletePullRequestCall.Receives.Ctx = param1
	f.DeletePullRequestCall.Receives.Id = param2
	if f.DeletePullRequestCall.Stub != nil {
		return f.DeletePullRequestCall.Stub(param1, param2)
	}
	return f.DeletePullRequestCall.Returns.Error
}
func (f *Querier) DeletePullRequestAction(param1 context.Context, param2 string) error {
	f.DeletePullRequestActionCall.mutex.Lock()
	defer f.DeletePullRequestActionCall.mutex.Unlock()
	f.DeletePullRequestActionCall.CallCount++
	f.DeletePullRequestActionCall.Receives.Ctx = param1
	f.DeletePullRequestActionCall.Receives.Name = param2
	if f.DeletePullRequestActionCall.Stub != nil {
		return f.DeletePullRequestActionCall.Stub(param1, param2)
	}
	return f.DeletePullRequestActionCall.Returns.Error
}
func (f *Querier) DeletePullRequestEvent(param1 context.Context, param2 int32) error {
	f.DeletePullRequestEventCall.mutex.Lock()
	defer f.DeletePullRequestEventCall.mutex.Unlock()
	f.DeletePullRequestEventCall.CallCount++
	f.DeletePullRequestEventCall.Receives.Ctx = param1
	f.DeletePullRequestEventCall.Receives.Id = param2
	if f.DeletePullRequestEventCall.Stub != nil {
		return f.DeletePullRequestEventCall.Stub(param1, param2)
	}
	return f.DeletePullRequestEventCall.Returns.Error
}
func (f *Querier) DeleteRepo(param1 context.Context, param2 int32) error {
	f.DeleteRepoCall.mutex.Lock()
	defer f.DeleteRepoCall.mutex.Unlock()
	f.DeleteRepoCall.CallCount++
	f.DeleteRepoCall.Receives.Ctx = param1
	f.DeleteRepoCall.Receives.Id = param2
	if f.DeleteRepoCall.Stub != nil {
		return f.DeleteRepoCall.Stub(param1, param2)
	}
	return f.DeleteRepoCall.Returns.Error
}
func (f *Querier) GetApprovalById(param1 context.Context, param2 int32) (postgres.Approval, error) {
	f.GetApprovalByIdCall.mutex.Lock()
	defer f.GetApprovalByIdCall.mutex.Unlock()
	f.GetApprovalByIdCall.CallCount++
	f.GetApprovalByIdCall.Receives.Ctx = param1
	f.GetApprovalByIdCall.Receives.Id = param2
	if f.GetApprovalByIdCall.Stub != nil {
		return f.GetApprovalByIdCall.Stub(param1, param2)
	}
	return f.GetApprovalByIdCall.Returns.Approval, f.GetApprovalByIdCall.Returns.Error
}
func (f *Querier) GetApprovalByPrIDSha(param1 context.Context, param2 postgres.GetApprovalByPrIDShaParams) (postgres.Approval, error) {
	f.GetApprovalByPrIDShaCall.mutex.Lock()
	defer f.GetApprovalByPrIDShaCall.mutex.Unlock()
	f.GetApprovalByPrIDShaCall.CallCount++
	f.GetApprovalByPrIDShaCall.Receives.Ctx = param1
	f.GetApprovalByPrIDShaCall.Receives.Arg = param2
	if f.GetApprovalByPrIDShaCall.Stub != nil {
		return f.GetApprovalByPrIDShaCall.Stub(param1, param2)
	}
	return f.GetApprovalByPrIDShaCall.Returns.Approval, f.GetApprovalByPrIDShaCall.Returns.Error
}
func (f *Querier) GetApprovalByUuid(param1 context.Context, param2 string) (postgres.Approval, error) {
	f.GetApprovalByUuidCall.mutex.Lock()
	defer f.GetApprovalByUuidCall.mutex.Unlock()
	f.GetApprovalByUuidCall.CallCount++
	f.GetApprovalByUuidCall.Receives.Ctx = param1
	f.GetApprovalByUuidCall.Receives.Uuid = param2
	if f.GetApprovalByUuidCall.Stub != nil {
		return f.GetApprovalByUuidCall.Stub(param1, param2)
	}
	return f.GetApprovalByUuidCall.Returns.Approval, f.GetApprovalByUuidCall.Returns.Error
}
func (f *Querier) GetApprovalSchemaById(param1 context.Context, param2 int32) (postgres.ApprovalSchema, error) {
	f.GetApprovalSchemaByIdCall.mutex.Lock()
	defer f.GetApprovalSchemaByIdCall.mutex.Unlock()
	f.GetApprovalSchemaByIdCall.CallCount++
	f.GetApprovalSchemaByIdCall.Receives.Ctx = param1
	f.GetApprovalSchemaByIdCall.Receives.Id = param2
	if f.GetApprovalSchemaByIdCall.Stub != nil {
		return f.GetApprovalSchemaByIdCall.Stub(param1, param2)
	}
	return f.GetApprovalSchemaByIdCall.Returns.ApprovalSchema, f.GetApprovalSchemaByIdCall.Returns.Error
}
func (f *Querier) GetApprovals(param1 context.Context) ([]postgres.Approval, error) {
	f.GetApprovalsCall.mutex.Lock()
	defer f.GetApprovalsCall.mutex.Unlock()
	f.GetApprovalsCall.CallCount++
	f.GetApprovalsCall.Receives.Ctx = param1
	if f.GetApprovalsCall.Stub != nil {
		return f.GetApprovalsCall.Stub(param1)
	}
	return f.GetApprovalsCall.Returns.ApprovalSlice, f.GetApprovalsCall.Returns.Error
}
func (f *Querier) GetCreateStatusInputsFromApprovalUuid(param1 context.Context, param2 string) (postgres.GetCreateStatusInputsFromApprovalUuidRow, error) {
	f.GetCreateStatusInputsFromApprovalUuidCall.mutex.Lock()
	defer f.GetCreateStatusInputsFromApprovalUuidCall.mutex.Unlock()
	f.GetCreateStatusInputsFromApprovalUuidCall.CallCount++
	f.GetCreateStatusInputsFromApprovalUuidCall.Receives.Ctx = param1
	f.GetCreateStatusInputsFromApprovalUuidCall.Receives.Uuid = param2
	if f.GetCreateStatusInputsFromApprovalUuidCall.Stub != nil {
		return f.GetCreateStatusInputsFromApprovalUuidCall.Stub(param1, param2)
	}
	return f.GetCreateStatusInputsFromApprovalUuidCall.Returns.GetCreateStatusInputsFromApprovalUuidRow, f.GetCreateStatusInputsFromApprovalUuidCall.Returns.Error
}
func (f *Querier) GetDefaultApprovalSchema(param1 context.Context) (postgres.ApprovalSchema, error) {
	f.GetDefaultApprovalSchemaCall.mutex.Lock()
	defer f.GetDefaultApprovalSchemaCall.mutex.Unlock()
	f.GetDefaultApprovalSchemaCall.CallCount++
	f.GetDefaultApprovalSchemaCall.Receives.Ctx = param1
	if f.GetDefaultApprovalSchemaCall.Stub != nil {
		return f.GetDefaultApprovalSchemaCall.Stub(param1)
	}
	return f.GetDefaultApprovalSchemaCall.Returns.ApprovalSchema, f.GetDefaultApprovalSchemaCall.Returns.Error
}
func (f *Querier) GetGithubUser(param1 context.Context, param2 int32) (postgres.GhUser, error) {
	f.GetGithubUserCall.mutex.Lock()
	defer f.GetGithubUserCall.mutex.Unlock()
	f.GetGithubUserCall.CallCount++
	f.GetGithubUserCall.Receives.Ctx = param1
	f.GetGithubUserCall.Receives.Id = param2
	if f.GetGithubUserCall.Stub != nil {
		return f.GetGithubUserCall.Stub(param1, param2)
	}
	return f.GetGithubUserCall.Returns.GhUser, f.GetGithubUserCall.Returns.Error
}
func (f *Querier) GetGithubUsers(param1 context.Context) ([]postgres.GhUser, error) {
	f.GetGithubUsersCall.mutex.Lock()
	defer f.GetGithubUsersCall.mutex.Unlock()
	f.GetGithubUsersCall.CallCount++
	f.GetGithubUsersCall.Receives.Ctx = param1
	if f.GetGithubUsersCall.Stub != nil {
		return f.GetGithubUsersCall.Stub(param1)
	}
	return f.GetGithubUsersCall.Returns.GhUserSlice, f.GetGithubUsersCall.Returns.Error
}
func (f *Querier) GetInstallation(param1 context.Context, param2 int32) (int32, error) {
	f.GetInstallationCall.mutex.Lock()
	defer f.GetInstallationCall.mutex.Unlock()
	f.GetInstallationCall.CallCount++
	f.GetInstallationCall.Receives.Ctx = param1
	f.GetInstallationCall.Receives.Id = param2
	if f.GetInstallationCall.Stub != nil {
		return f.GetInstallationCall.Stub(param1, param2)
	}
	return f.GetInstallationCall.Returns.Int32, f.GetInstallationCall.Returns.Error
}
func (f *Querier) GetInstallations(param1 context.Context) ([]int32, error) {
	f.GetInstallationsCall.mutex.Lock()
	defer f.GetInstallationsCall.mutex.Unlock()
	f.GetInstallationsCall.CallCount++
	f.GetInstallationsCall.Receives.Ctx = param1
	if f.GetInstallationsCall.Stub != nil {
		return f.GetInstallationsCall.Stub(param1)
	}
	return f.GetInstallationsCall.Returns.Int32Slice, f.GetInstallationsCall.Returns.Error
}
func (f *Querier) GetPullRequestAction(param1 context.Context, param2 string) (string, error) {
	f.GetPullRequestActionCall.mutex.Lock()
	defer f.GetPullRequestActionCall.mutex.Unlock()
	f.GetPullRequestActionCall.CallCount++
	f.GetPullRequestActionCall.Receives.Ctx = param1
	f.GetPullRequestActionCall.Receives.Name = param2
	if f.GetPullRequestActionCall.Stub != nil {
		return f.GetPullRequestActionCall.Stub(param1, param2)
	}
	return f.GetPullRequestActionCall.Returns.String, f.GetPullRequestActionCall.Returns.Error
}
func (f *Querier) GetPullRequestActions(param1 context.Context) ([]string, error) {
	f.GetPullRequestActionsCall.mutex.Lock()
	defer f.GetPullRequestActionsCall.mutex.Unlock()
	f.GetPullRequestActionsCall.CallCount++
	f.GetPullRequestActionsCall.Receives.Ctx = param1
	if f.GetPullRequestActionsCall.Stub != nil {
		return f.GetPullRequestActionsCall.Stub(param1)
	}
	return f.GetPullRequestActionsCall.Returns.StringSlice, f.GetPullRequestActionsCall.Returns.Error
}
func (f *Querier) GetPullRequestById(param1 context.Context, param2 int32) (postgres.PullRequest, error) {
	f.GetPullRequestByIdCall.mutex.Lock()
	defer f.GetPullRequestByIdCall.mutex.Unlock()
	f.GetPullRequestByIdCall.CallCount++
	f.GetPullRequestByIdCall.Receives.Ctx = param1
	f.GetPullRequestByIdCall.Receives.Id = param2
	if f.GetPullRequestByIdCall.Stub != nil {
		return f.GetPullRequestByIdCall.Stub(param1, param2)
	}
	return f.GetPullRequestByIdCall.Returns.PullRequest, f.GetPullRequestByIdCall.Returns.Error
}
func (f *Querier) GetPullRequestByRepoIdPrId(param1 context.Context, param2 postgres.GetPullRequestByRepoIdPrIdParams) (postgres.PullRequest, error) {
	f.GetPullRequestByRepoIdPrIdCall.mutex.Lock()
	defer f.GetPullRequestByRepoIdPrIdCall.mutex.Unlock()
	f.GetPullRequestByRepoIdPrIdCall.CallCount++
	f.GetPullRequestByRepoIdPrIdCall.Receives.Ctx = param1
	f.GetPullRequestByRepoIdPrIdCall.Receives.Arg = param2
	if f.GetPullRequestByRepoIdPrIdCall.Stub != nil {
		return f.GetPullRequestByRepoIdPrIdCall.Stub(param1, param2)
	}
	return f.GetPullRequestByRepoIdPrIdCall.Returns.PullRequest, f.GetPullRequestByRepoIdPrIdCall.Returns.Error
}
func (f *Querier) GetPullRequestEvent(param1 context.Context, param2 int32) (postgres.PullRequestEvent, error) {
	f.GetPullRequestEventCall.mutex.Lock()
	defer f.GetPullRequestEventCall.mutex.Unlock()
	f.GetPullRequestEventCall.CallCount++
	f.GetPullRequestEventCall.Receives.Ctx = param1
	f.GetPullRequestEventCall.Receives.Id = param2
	if f.GetPullRequestEventCall.Stub != nil {
		return f.GetPullRequestEventCall.Stub(param1, param2)
	}
	return f.GetPullRequestEventCall.Returns.PullRequestEvent, f.GetPullRequestEventCall.Returns.Error
}
func (f *Querier) GetPullRequestEvents(param1 context.Context) ([]postgres.PullRequestEvent, error) {
	f.GetPullRequestEventsCall.mutex.Lock()
	defer f.GetPullRequestEventsCall.mutex.Unlock()
	f.GetPullRequestEventsCall.CallCount++
	f.GetPullRequestEventsCall.Receives.Ctx = param1
	if f.GetPullRequestEventsCall.Stub != nil {
		return f.GetPullRequestEventsCall.Stub(param1)
	}
	return f.GetPullRequestEventsCall.Returns.PullRequestEventSlice, f.GetPullRequestEventsCall.Returns.Error
}
func (f *Querier) GetPullRequests(param1 context.Context) ([]postgres.PullRequest, error) {
	f.GetPullRequestsCall.mutex.Lock()
	defer f.GetPullRequestsCall.mutex.Unlock()
	f.GetPullRequestsCall.CallCount++
	f.GetPullRequestsCall.Receives.Ctx = param1
	if f.GetPullRequestsCall.Stub != nil {
		return f.GetPullRequestsCall.Stub(param1)
	}
	return f.GetPullRequestsCall.Returns.PullRequestSlice, f.GetPullRequestsCall.Returns.Error
}
func (f *Querier) GetRepo(param1 context.Context, param2 int32) (postgres.Repo, error) {
	f.GetRepoCall.mutex.Lock()
	defer f.GetRepoCall.mutex.Unlock()
	f.GetRepoCall.CallCount++
	f.GetRepoCall.Receives.Ctx = param1
	f.GetRepoCall.Receives.Id = param2
	if f.GetRepoCall.Stub != nil {
		return f.GetRepoCall.Stub(param1, param2)
	}
	return f.GetRepoCall.Returns.Repo, f.GetRepoCall.Returns.Error
}
func (f *Querier) GetRepos(param1 context.Context) ([]postgres.Repo, error) {
	f.GetReposCall.mutex.Lock()
	defer f.GetReposCall.mutex.Unlock()
	f.GetReposCall.CallCount++
	f.GetReposCall.Receives.Ctx = param1
	if f.GetReposCall.Stub != nil {
		return f.GetReposCall.Stub(param1)
	}
	return f.GetReposCall.Returns.RepoSlice, f.GetReposCall.Returns.Error
}
func (f *Querier) GetSortedApprovalYesNoQuestionAnswersByUuid(param1 context.Context, param2 string) ([]postgres.GetSortedApprovalYesNoQuestionAnswersByUuidRow, error) {
	f.GetSortedApprovalYesNoQuestionAnswersByUuidCall.mutex.Lock()
	defer f.GetSortedApprovalYesNoQuestionAnswersByUuidCall.mutex.Unlock()
	f.GetSortedApprovalYesNoQuestionAnswersByUuidCall.CallCount++
	f.GetSortedApprovalYesNoQuestionAnswersByUuidCall.Receives.Ctx = param1
	f.GetSortedApprovalYesNoQuestionAnswersByUuidCall.Receives.Uuid = param2
	if f.GetSortedApprovalYesNoQuestionAnswersByUuidCall.Stub != nil {
		return f.GetSortedApprovalYesNoQuestionAnswersByUuidCall.Stub(param1, param2)
	}
	return f.GetSortedApprovalYesNoQuestionAnswersByUuidCall.Returns.GetSortedApprovalYesNoQuestionAnswersByUuidRowSlice, f.GetSortedApprovalYesNoQuestionAnswersByUuidCall.Returns.Error
}
func (f *Querier) GetSortedApprovalYesNoQuestionsBySchemaId(param1 context.Context, param2 string) ([]postgres.GetSortedApprovalYesNoQuestionsBySchemaIdRow, error) {
	f.GetSortedApprovalYesNoQuestionsBySchemaIdCall.mutex.Lock()
	defer f.GetSortedApprovalYesNoQuestionsBySchemaIdCall.mutex.Unlock()
	f.GetSortedApprovalYesNoQuestionsBySchemaIdCall.CallCount++
	f.GetSortedApprovalYesNoQuestionsBySchemaIdCall.Receives.Ctx = param1
	f.GetSortedApprovalYesNoQuestionsBySchemaIdCall.Receives.Name = param2
	if f.GetSortedApprovalYesNoQuestionsBySchemaIdCall.Stub != nil {
		return f.GetSortedApprovalYesNoQuestionsBySchemaIdCall.Stub(param1, param2)
	}
	return f.GetSortedApprovalYesNoQuestionsBySchemaIdCall.Returns.GetSortedApprovalYesNoQuestionsBySchemaIdRowSlice, f.GetSortedApprovalYesNoQuestionsBySchemaIdCall.Returns.Error
}
func (f *Querier) UpdateApprovalByUuid(param1 context.Context, param2 postgres.UpdateApprovalByUuidParams) error {
	f.UpdateApprovalByUuidCall.mutex.Lock()
	defer f.UpdateApprovalByUuidCall.mutex.Unlock()
	f.UpdateApprovalByUuidCall.CallCount++
	f.UpdateApprovalByUuidCall.Receives.Ctx = param1
	f.UpdateApprovalByUuidCall.Receives.Arg = param2
	if f.UpdateApprovalByUuidCall.Stub != nil {
		return f.UpdateApprovalByUuidCall.Stub(param1, param2)
	}
	return f.UpdateApprovalByUuidCall.Returns.Error
}
func (f *Querier) UpdatePullRequestIsMerged(param1 context.Context, param2 postgres.UpdatePullRequestIsMergedParams) (postgres.PullRequest, error) {
	f.UpdatePullRequestIsMergedCall.mutex.Lock()
	defer f.UpdatePullRequestIsMergedCall.mutex.Unlock()
	f.UpdatePullRequestIsMergedCall.CallCount++
	f.UpdatePullRequestIsMergedCall.Receives.Ctx = param1
	f.UpdatePullRequestIsMergedCall.Receives.Arg = param2
	if f.UpdatePullRequestIsMergedCall.Stub != nil {
		return f.UpdatePullRequestIsMergedCall.Stub(param1, param2)
	}
	return f.UpdatePullRequestIsMergedCall.Returns.PullRequest, f.UpdatePullRequestIsMergedCall.Returns.Error
}
func (f *Querier) UpdateRepoName(param1 context.Context, param2 postgres.UpdateRepoNameParams) (postgres.Repo, error) {
	f.UpdateRepoNameCall.mutex.Lock()
	defer f.UpdateRepoNameCall.mutex.Unlock()
	f.UpdateRepoNameCall.CallCount++
	f.UpdateRepoNameCall.Receives.Ctx = param1
	f.UpdateRepoNameCall.Receives.Arg = param2
	if f.UpdateRepoNameCall.Stub != nil {
		return f.UpdateRepoNameCall.Stub(param1, param2)
	}
	return f.UpdateRepoNameCall.Returns.Repo, f.UpdateRepoNameCall.Returns.Error
}
func (f *Querier) UpdateRepoOrg(param1 context.Context, param2 postgres.UpdateRepoOrgParams) (postgres.Repo, error) {
	f.UpdateRepoOrgCall.mutex.Lock()
	defer f.UpdateRepoOrgCall.mutex.Unlock()
	f.UpdateRepoOrgCall.CallCount++
	f.UpdateRepoOrgCall.Receives.Ctx = param1
	f.UpdateRepoOrgCall.Receives.Arg = param2
	if f.UpdateRepoOrgCall.Stub != nil {
		return f.UpdateRepoOrgCall.Stub(param1, param2)
	}
	return f.UpdateRepoOrgCall.Returns.Repo, f.UpdateRepoOrgCall.Returns.Error
}
