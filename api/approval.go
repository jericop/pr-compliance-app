package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

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
	approval, err := server.querier.GetApprovalByUuid(context.Background(), uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	approvalJSON, err := server.jsonMarshal(approval)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	fmt.Fprintf(w, string(approvalJSON))
}

func (server *Server) UpdateApproval(w http.ResponseWriter, req *http.Request) {
	var p postgres.UpdateApprovalByUuidParams

	decoder := json.NewDecoder(req.Body)

	err := decoder.Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = server.querier.UpdateApprovalByUuid(context.Background(), p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pJSON, err := server.jsonMarshal(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, string(pJSON))

}
