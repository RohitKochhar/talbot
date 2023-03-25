package main

import (
	"net/http"
)

// Replies with a 200 status indicating server is up and running
func (a *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	replyTextContent(w, r, http.StatusOK, "OK")
}

// replyTextContent wraps text content in a HTTP response and sends it
func replyTextContent(w http.ResponseWriter, r *http.Request, status int, content string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
	w.Write([]byte(content + "\n"))
}

