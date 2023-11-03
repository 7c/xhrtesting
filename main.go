package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var addr = flag.String("addr", ":8080", "HTTP network address")

func main() {
	flag.Parse() // Parse the command line flags

	r := chi.NewRouter()

	// Using the logger middleware
	r.Use(middleware.Logger)

	// Ping endpoint
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	r.Get("/to/301", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("301"))
	})

	r.Get("/to/302", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("302"))
	})

	r.Get("/status/200", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Get("/status/201", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Created"))
	})

	r.Get("/status/204", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		// No content to return
	})

	r.Get("/status/301", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/to/301", http.StatusMovedPermanently)
	})

	r.Get("/status/302", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/to/302", http.StatusFound)
	})

	//304 Not Modified
	r.Get("/status/304", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotModified)
		// No content to return
	})

	r.Get("/status/400", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	})

	r.Get("/status/401", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted area"`)
		http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
	})

	r.Get("/status/403", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Forbidden", http.StatusForbidden)
	})

	// 404 Not found
	r.Get("/status/404", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	r.Get("/status/408", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Request Timeout", http.StatusRequestTimeout)
	})

	r.Get("/status/500", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	})

	r.Get("/status/501", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Implemented", http.StatusNotImplemented)
	})

	r.Get("/status/503", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
	})

	log.Fatal(http.ListenAndServe(*addr, r))
}
