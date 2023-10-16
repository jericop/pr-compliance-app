package api

import (
	"context"
	"fmt"
	"net/http"
)

func (server *Server) AddPullRequestRoutes() {
	server.router.HandleFunc("/pull_requests", server.GetPullRequests).Methods("GET").Name("GetPullRequests")
}

// GetPullRequests retrieves the collection of PullRequests
func (server *Server) GetPullRequests(w http.ResponseWriter, req *http.Request) {
	pullRequests, err := server.querier.GetPullRequests(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pullRequestsJSON, err := server.jsonMarshal(pullRequests)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	fmt.Fprintf(w, string(pullRequestsJSON))
}
