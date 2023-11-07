package middleware

import (
	"net/http"
)

// CloudflareIPMiddleware retrieves the real IP from the CF-Connecting-IP header
func CloudflareIPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// litter.Dump(r.Header)

		// Get the CF-Connecting-IP header value
		cfIP := r.Header.Get("CF-Connecting-IP")
		if cfIP != "" {
			r.RemoteAddr = cfIP
		} else {
			// Fallback to the RemoteAddr if header is not set
			// log.Printf("Direct External IP: %s", r.RemoteAddr)
		}
		next.ServeHTTP(w, r)
	})
}
