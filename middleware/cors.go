package middleware

import (
	"github.com/go-chi/cors"
)

// Setup CORS options
func CorsMiddleware() *cors.Cors {
	return cors.New(cors.Options{
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
}
