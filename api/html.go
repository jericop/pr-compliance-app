package api

import (
	"fmt"
	"net/http"
)

func (server *Server) AddHtmlRoutes() {
	server.router.HandleFunc("/html_test", server.GetHtml).Methods("GET").Name("GetHtml")
}

func (server *Server) GetHtml(w http.ResponseWriter, req *http.Request) {
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
			<div>
				<p><label>Your email:</label></p>
				<p><input type="email" name="email"></p>
			</div>
			<div>
				<p><label>Your message:</label></p>
				<p><textarea name="content"></textarea></p>
			</div>
			<div>
				<input type="submit" value="Send message">
			</div>
		</form>
    </body>
</html>
	`, "http://localhost:8080/approval", "7829491b-9bdc-4167-87e3-73a334fb5916")
	fmt.Fprintf(w, html)
}
