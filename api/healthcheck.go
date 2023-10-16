// endpoints.go
package api

import (
	"io"
	"net/http"
)

func (server *Server) AddHealthRoutes() {
	server.router.HandleFunc("/health", server.GetHealth).Methods("GET").Name("GetHealth")
}

func (server *Server) GetHealth(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// We could report some meaningful status here if needed.
	io.WriteString(w, `{"alive": true}`)
}
