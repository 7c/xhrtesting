package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-resty/resty/v2"
)

var addr = flag.String("addr", ":8080", "HTTP network address")

func PublicIP() (string, error) {
	client := resty.New()
	resp, err := client.R().Get("https://ip4.ip8.com")
	if err != nil {
		return "", err
	}
	if resp.StatusCode() == 200 {
		return string(resp.Body()), nil
	}
	return "", errors.New("Error: " + resp.Status())
}

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

	// Cookie random endpoint
	r.Get("/cookie/random", func(w http.ResponseWriter, r *http.Request) {
		// Generate a random number as the cookie value
		rand.Seed(time.Now().UnixNano())
		randomValue := fmt.Sprintf("%d", rand.Intn(1000000))

		// Set the cookie to the response
		http.SetCookie(w, &http.Cookie{
			Name:    "Random-Cookie",
			Value:   randomValue,
			Expires: time.Now().Add(24 * time.Hour),
			Path:    "/",
		})

		w.Write([]byte("Random cookie set"))
	})

	// Cookie random endpoint for setting a specified number of cookies
	r.Get("/cookie/random/{number}", func(w http.ResponseWriter, r *http.Request) {
		// Get the number from the URL parameter
		numberStr := chi.URLParam(r, "number")
		number, err := strconv.Atoi(numberStr)
		if err != nil {
			http.Error(w, "Invalid number of cookies", http.StatusBadRequest)
			return
		}

		// Generate and set the specified number of cookies
		for i := 0; i < number; i++ {
			randomValue := fmt.Sprintf("%d", rand.Intn(1000000))
			cookie := &http.Cookie{
				Name:    fmt.Sprintf("Random-Cookie-%d", i+1),
				Value:   randomValue,
				Expires: time.Now().Add(24 * time.Hour),
				Path:    "/",
			}
			http.SetCookie(w, cookie)
		}

		w.Write([]byte(fmt.Sprintf("Set %d random cookies", number)))
	})

	ip4, _ := PublicIP()
	fmt.Printf("Your external IP is : %s\n", ip4)
	fmt.Printf("Launching xhrtesting server at %s\n", *addr)
	log.Fatal(http.ListenAndServe(*addr, r))
}
