package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
	"xhrtesting/handlers"
	mw "xhrtesting/middleware"
	"xhrtesting/shared"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var addr = flag.String("addr", ":8080", "HTTP network address")
var tlsDomain = flag.String("tlsdomain", "", "Letsencrypt TLS Domain")

func main() {
	flag.Parse() // Parse the command line flags

	r := chi.NewRouter()

	r.Use(mw.CloudflareIPMiddleware)

	// Using the logger middleware
	r.Use(middleware.Logger)

	// Use the CORS middleware
	r.Use(mw.CorsMiddleware().Handler)

	// proxy to godocs
	// r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
	// 	backendURL, _ := url.Parse("http://127.0.0.1:6060")
	// 	pr := httputil.NewSingleHostReverseProxy(backendURL)
	// 	pr.Director = func(req *http.Request) {
	// 		// Update the request to point to the backend server
	// 		req.URL.Scheme = backendURL.Scheme
	// 		req.URL.Host = backendURL.Host
	// 		req.Header.Set("X-Forwarded-Host", req.Host)
	// 	}
	// 	pr.ServeHTTP(w, r)
	// })

	r.HandleFunc("/", handlers.GithubPage)

	r.HandleFunc("/ping", handlers.PingHandler)

	r.HandleFunc("/status/200", handlers.Status200Handler)
	r.HandleFunc("/status/201", handlers.Status201Handler)
	r.HandleFunc("/status/204", handlers.Status204Handler)
	r.HandleFunc("/status/301", handlers.Status301Handler)
	r.HandleFunc("/status/302", handlers.Status302Handler)
	//304 Not Modified
	r.HandleFunc("/status/304", handlers.Status304Handler)

	r.HandleFunc("/to/301", handlers.Status200Handler)
	r.HandleFunc("/to/302", handlers.Status200Handler)

	r.HandleFunc("/status/400", handlers.Status400Handler)

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

	ip4, _ := shared.PublicIP()
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
