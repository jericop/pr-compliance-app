package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Widget is a thing
type Widget struct {
	Id   int64
	Name string
	Role string
}

var AllWidgets []Widget

func init() {
	AllWidgets = []Widget{
		Widget{0, "Nothing", "N/A"},
		Widget{1, "Jibjab", "Instrument"},
		Widget{2, "Pencil", "Utensil"},
		Widget{3, "Fork", "Utensil"},
		Widget{4, "Password", "Credential"},
		Widget{5, "SpanFrankisco", "Home"},
		Widget{6, "Doc", "Villain"},
		Widget{7, "Coff3e", "Hack"},
	}
}

func (server *Server) AddPullRequestRoutes() {
	server.router.HandleFunc("/pull_requests", server.GetPullRequests).Methods("GET").Name("GetPullRequests")
	// server.router.HandleFunc("/pull_requests", server.PostPullRequest).Methods("POST").Name("PostPullRequest")
	// server.router.HandleFunc("/pull_requests/{id}", server.GetPullRequest).Methods("GET").Name("GetPullRequest")
}

// GetPullRequests retrieves the collection of PullRequests
func (server *Server) GetPullRequests(w http.ResponseWriter, req *http.Request) {
	pullRequests, err := server.store.GetPullRequests(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pullRequestsJSON, err := server.jsonMarshal(pullRequests)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("pullRequestsJSON: %v\n", pullRequestsJSON)

	w.Header().Set("Content-Type", "application/json")

	fmt.Fprintf(w, string(pullRequestsJSON))
}

// GetPullRequest retrieves a single PullRequest
func (server *Server) GetPullRequest(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(req)["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if id >= int64(len(AllWidgets)) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	widgetJSON, err := json.Marshal(AllWidgets[id])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	fmt.Fprintf(w, string(widgetJSON))
}

// PostPullRequest adds a PullRequest to the collection
func (server *Server) PostPullRequest(w http.ResponseWriter, req *http.Request) {
	var widget Widget
	decoder := json.NewDecoder(req.Body)

	err := decoder.Decode(&widget)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(widget.Name) == 0 {
		err = errors.New("Widget name can't be blank.")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// not thread safe...
	widget.Id = int64(len(AllWidgets))
	AllWidgets = append(AllWidgets, widget)

	widgetJSON, err := json.Marshal(widget)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, string(widgetJSON))

}
