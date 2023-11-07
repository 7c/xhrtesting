package handlers

import "net/http"

// PingHandler responds to requests with "pong".
// This handler is used to check if the service is running and is reachable.
func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

// Status301Handler responds with a 301 status code.
// It is used to demonstrate a permanent redirect response.
func Status301Handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("301"))
}

// Status302Handler responds with a 302 status code.
// It is used to demonstrate a temporary redirect response.
func Status302Handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("302"))
}

// Status200Handler responds with a 200 status code.
// It indicates that the request has succeeded.
func Status200Handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// Status201Handler responds with a 201 status code.
// It indicates that the request has been fulfilled and has resulted in a new resource being created.
func Status201Handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Created"))
}

// Status204Handler responds with a 204 status code.
// It indicates that the server successfully processed the request, and is not returning any content.
func Status204Handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

// Status301RedirectHandler redirects to the /to/301 endpoint with a 301 status code.
// This handler is used to demonstrate a permanent redirection.
func Status301RedirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/to/301", http.StatusMovedPermanently)
}

// Status302RedirectHandler redirects to the /to/302 endpoint with a 302 status code.
// This handler is used to demonstrate a temporary redirection.
func Status302RedirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/to/302", http.StatusFound)
}

// Status304Handler responds with a 304 status code.
// It is used to tell the client that the response has not been modified,
// so the client can continue to use the same cached version of the response.
func Status304Handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotModified)
}

// Status400Handler responds with a 400 status code.
// It indicates that the server cannot or will not process the request due to an apparent client error.
func Status400Handler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Bad Request", http.StatusBadRequest)
}

func GithubPage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://github.com/7c/xhrtesting", http.StatusMovedPermanently)
}
