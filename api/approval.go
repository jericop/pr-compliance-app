package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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
	w.Header().Set("Content-Type", "text/html")
	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
    <head>
    <title>PR Compliance App</title>    
		
    </head>
    <body>
		<p>The following form must be submitted before the PR can be merged</p>
		<br>
		<h1>Approval</h1>
		<form action="%s" method="POST" novalidate>
			<input type="hidden" id="uuid" name="uuid" value="%s">
			<input type="hidden" id="is_approved" name="is_approved" value="true">	
				<input type="submit" value="Approve">
			</div>
		</form>
    </body>
</html>
	`, "http://localhost:8080/approval", uuid)
	fmt.Fprintf(w, html)
}

func (server *Server) UpdateApproval(w http.ResponseWriter, req *http.Request) {
	var p postgres.UpdateApprovalByUuidParams

	if req.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
		// untested - hard to force an error because it always seems to returns nil
		err := req.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		isApproved, err := strconv.ParseBool(req.Form.Get("is_approved"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		p = postgres.UpdateApprovalByUuidParams{
			Uuid:       req.Form.Get("uuid"),
			IsApproved: isApproved,
		}
	} else {
		decoder := json.NewDecoder(req.Body)

		err := decoder.Decode(&p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	err := server.querier.UpdateApprovalByUuid(context.Background(), p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := context.Background()

	inputs, err := server.querier.GetCreateStatusInputsFromApprovalUuid(ctx, p.Uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client, err := server.githubFactory.NewInstallationClient(ctx, int64(inputs.InstallationID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	repoStatus := &github.RepoStatus{
		// TODO: Get these fields from the database at startup and then use them for all requests
		Context:     github.String(statusContext),
		Description: github.String(statusTitle),
		TargetURL:   github.String(fmt.Sprintf("http://localhost:8080/approval/%s", p.Uuid)),
		// TargetURL:   github.String(fmt.Sprintf("https://localhost:8080/approval?id=%s", p.Uuid)),
		State: github.String("success"),
	}

	_, _, err = client.Repositories.CreateStatus(ctx, inputs.Login, inputs.Name, inputs.Sha, repoStatus)
	if err != nil {
		http.Error(w, fmt.Sprintf("installations error %v", err), http.StatusInternalServerError)
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
