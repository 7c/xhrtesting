package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
	"xhrtesting/shared"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-resty/resty/v2"
)

var addr = flag.String("addr", ":8080", "HTTP network address")
var tlsDomain = flag.String("tlsdomain", "", "Letsencrypt TLS Domain")

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

	// Setup CORS options
	corsMiddleware := cors.New(cors.Options{
		// AllowedOrigins is a list of origins a cross-domain request can be executed from.
		// Use "*" to allow all origins. To allow specific origins, use a list like the following:
		// []string{"https://foo.com", "https://bar.com"}
		AllowedOrigins: []string{"*"},
		// AllowedMethods is a list of methods the client is allowed to use with cross-domain requests.
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		// AllowedHeaders is list of non simple headers the client is allowed to use with cross-domain requests.
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		// ExposedHeaders indicates which headers are safe to expose to the API of a CORS API specification
		ExposedHeaders: []string{},
		// AllowCredentials indicates whether the request can include user credentials like
		// cookies, HTTP authentication or client side SSL certificates.
		AllowCredentials: false,
		// MaxAge indicates how long (in seconds) the results of a preflight request can be cached.
		MaxAge: 300, // Maximum value not to exceed 600 seconds
	})

	// Using the logger middleware
	r.Use(middleware.Logger)

	// Use the CORS middleware
	r.Use(corsMiddleware.Handler)

	// Ping endpoint
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	r.HandleFunc("/to/301", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("301"))
	})

	r.HandleFunc("/to/302", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("302"))
	})

	r.HandleFunc("/status/200", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.HandleFunc("/status/201", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Created"))
	})

	r.HandleFunc("/status/204", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		// No content to return
	})

	r.HandleFunc("/status/301", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/to/301", http.StatusMovedPermanently)
	})

	r.HandleFunc("/status/302", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/to/302", http.StatusFound)
	})

	//304 Not Modified
	r.HandleFunc("/status/304", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotModified)
		// No content to return
	})

	r.HandleFunc("/status/400", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	})

	// long body response
	r.HandleFunc("/long/body/{number}", func(w http.ResponseWriter, r *http.Request) {
		numberStr := chi.URLParam(r, "number")
		number, err := strconv.Atoi(numberStr)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// Ensure the number is non-negative
		if number < 0 {
			http.Error(w, "Number must be non-negative", http.StatusBadRequest)
			return
		}

		// Create a string that is 'number' characters long
		longString := strings.Repeat("x", number)

		w.Write([]byte(longString))
	})

	// Random JSON endpoint
	r.HandleFunc("/json/random", func(w http.ResponseWriter, r *http.Request) {
		// Generate random values of different types
		rand.Seed(time.Now().UnixNano())
		response := map[string]interface{}{
			"randomInt":    rand.Intn(100),                                           // Random integer
			"randomFloat":  rand.Float64(),                                           // Random float
			"randomBool":   rand.Intn(2) == 1,                                        // Random boolean
			"randomString": shared.RandomString(10),                                  // Random string
			"randomArray":  shared.RandomArray(5),                                    // Random array
			"randomObject": map[string]int{"a": rand.Intn(100), "b": rand.Intn(100)}, // Random object/map
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	r.HandleFunc("/status/401", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted area"`)
		http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
	})

	r.HandleFunc("/status/403", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Forbidden", http.StatusForbidden)
	})

	// 404 Not found
	r.HandleFunc("/status/404", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	r.HandleFunc("/status/408", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Request Timeout", http.StatusRequestTimeout)
	})

	r.HandleFunc("/status/500", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	})

	r.HandleFunc("/status/501", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Implemented", http.StatusNotImplemented)
	})

	r.HandleFunc("/status/503", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
	})

	// Cookie random endpoint
	r.HandleFunc("/cookie/random", func(w http.ResponseWriter, r *http.Request) {
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
	r.HandleFunc("/cookie/random/{number}", func(w http.ResponseWriter, r *http.Request) {
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

	if *tlsDomain != "" {
		fmt.Printf("HTTPS Mode\n")
		fmt.Printf("TLS Domain '%s' enabled\n", *tlsDomain)
		cert, keyFile := shared.LECerts(*tlsDomain)
		if cert == "" || keyFile == "" {
			log.Fatalln("Could not find cert,keyfiles from TLS Domain")
		}
		fmt.Printf("Listening on https://%s\n", *addr)
		log.Fatal(http.ListenAndServeTLS(*addr, cert, keyFile, r))
	}
	fmt.Printf("HTTP Mode\n")
	fmt.Printf("Listening on http://%s\n", *addr)
	log.Fatal(http.ListenAndServe(*addr, r))

}
