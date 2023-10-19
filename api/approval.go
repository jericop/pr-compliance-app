package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/go-github/v53/github"
	"github.com/gorilla/mux"
	"github.com/jericop/pr-compliance-app/storage/postgres"
)

func (server *Server) AddApprovalRoutes() {
	server.router.HandleFunc("/approval/{id}", server.GetApproval).Methods("GET").Name("GetApproval")
	server.router.HandleFunc("/approval", server.UpdateApproval).Methods("POST").Name("UpdateApproval")
	server.router.HandleFunc("/approval", server.GetApprovalQueryParam).Methods("GET").Name("GetApprovalQueryParam")
}

func (server *Server) GetApproval(w http.ResponseWriter, req *http.Request) {
	server.getApproval(w, req, mux.Vars(req)["id"])
}

func (server *Server) GetApprovalQueryParam(w http.ResponseWriter, req *http.Request) {
	uuid := req.URL.Query().Get("id")
	if uuid == "" {
		http.Error(w, "query paramter 'id' is needed", http.StatusInternalServerError)
		return
	}
	server.getApproval(w, req, uuid)
}

func (server *Server) getApproval(w http.ResponseWriter, req *http.Request, uuid string) {
	ctx := context.Background()

	questions, err := server.querier.GetSortedApprovalYesNoQuestionAnswersByUuid(ctx, uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	approvalQuestions := &ApprovalYesNoQuestionAnswersResponse{
		Uuid:      uuid,
		Questions: questions,
	}

	approvalQuestionsJSON, err := server.jsonMarshalFunc(approvalQuestions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(approvalQuestionsJSON))
}

func (server *Server) UpdateApproval(w http.ResponseWriter, req *http.Request) {
	var p ApprovalYesNoQuestionAnswersResponse

	ctx := context.Background()

	switch req.Header.Get("Content-Type") {
	case "application/json":
		decoder := json.NewDecoder(req.Body)

		err := decoder.Decode(&p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, fmt.Sprintf("invalid content-type %s", req.Header.Get("Content-Type")), http.StatusBadRequest)
		return
	}

	c, err := server.querier.GetSortedApprovalYesNoQuestionAnswersByUuid(ctx, p.Uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	currentYes := getYesAnswerMap(ApprovalYesNoQuestionAnswersResponse{Uuid: p.Uuid, Questions: c})
	postedYes := getYesAnswerMap(p)

	err = server.dbTxFactory.ExecWithTx(ctx, func(q postgres.Querier) error {
		// Create entries for questions that were answered yes
		for question_id, _ := range postedYes {
			if _, known := currentYes[question_id]; !known {
				_, err := q.CreateApprovalYesAnswerByUuid(ctx, postgres.CreateApprovalYesAnswerByUuidParams{Uuid: p.Uuid, QuestionID: question_id})
				if err != nil {
					return err
				}
			}
		}

		// Delete entries for questions that were previously answered yes, but now are not
		for question_id, _ := range currentYes {
			if _, known := postedYes[question_id]; !known {
				err := q.DeleteApprovalYesAnswerByUuid(ctx, postgres.DeleteApprovalYesAnswerByUuidParams{Uuid: p.Uuid, QuestionID: question_id})
				if err != nil {
					return err
				}
			}
		}

		err := q.UpdateApprovalByUuid(ctx, postgres.UpdateApprovalByUuidParams{Uuid: p.Uuid, IsApproved: true})
		if err != nil {
			return err
		}

		inputs, err := q.GetCreateStatusInputsFromApprovalUuid(ctx, p.Uuid)
		if err != nil {
			return err
		}

		client, err := server.githubFactory.NewInstallationClient(ctx, int64(inputs.InstallationID))
		if err != nil {
			return err
		}

		repoStatus := &github.RepoStatus{
			Context:     github.String(server.schema.StatusContext),
			Description: github.String(server.schema.StatusTitle),
			TargetURL:   github.String(fmt.Sprintf("%s/%s", server.frontEndUrl, p.Uuid)),
			// TargetURL: github.String(fmt.Sprintf("%s?id=%s", server.frontEndUrl, p.Uuid)),
			State: github.String("success"),
		}

		log.Printf("Creating a success commit status for approval uuid %s", p.Uuid)

		_, _, err = client.Repositories.CreateStatus(ctx, inputs.Login, inputs.Name, inputs.Sha, repoStatus)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pJSON, err := server.jsonMarshalFunc(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, string(pJSON))
}

func getYesAnswerMap(r ApprovalYesNoQuestionAnswersResponse) map[int32]struct{} {
	m := make(map[int32]struct{})

	for _, question := range r.Questions {
		if question.AnsweredYes {
			m[question.ID] = struct{}{}
		}
	}
	return m
}

type UrlEncodedFormAnswers struct {
	Data string `json:"data"`
}

type ApprovalYesNoQuestionAnswersResponse struct {
	Uuid      string                                                    `json:"uuid"`
	Questions []postgres.GetSortedApprovalYesNoQuestionAnswersByUuidRow `json:"questions"`
}
