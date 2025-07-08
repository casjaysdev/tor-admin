// File: internal/auth/middleware.go
// Purpose: Auth enforcement middleware for both UI and API access

package auth

import (
	"encoding/base64"
	"net/http"
	"os"
	"strings"
)

var disabledAuth = os.Getenv("DISABLE_AUTH") == "true"

// ProtectAPI ensures API token is present and correct
func ProtectAPI(token string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if disabledAuth {
			next.ServeHTTP(w, r)
			return
		}
		authHeader := r.Header.Get("Authorization")
		if isValidToken(authHeader, token) {
			next.ServeHTTP(w, r)
			return
		}
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

// RequireLogin ensures a valid session is active for web UI
func RequireLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if disabledAuth {
			next.ServeHTTP(w, r)
			return
		}
		// Skip auth for login/setup/static
		path := r.URL.Path
		if path == "/login" || path == "/setup" || path == "/logout" || strings.HasPrefix(path, "/static/") {
			next.ServeHTTP(w, r)
			return
		}
		user := GetUser(r)
		if user == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Helper: check if Authorization token matches expected
func isValidToken(header, expectedToken string) bool {
	if header == "" || expectedToken == "" {
		return false
	}
	header = strings.TrimSpace(header)

	// Bearer <token>
	if strings.HasPrefix(header, "Bearer ") {
		return strings.TrimPrefix(header, "Bearer ") == expectedToken
	}
	// Token <token>
	if strings.HasPrefix(header, "Token ") {
		return strings.TrimPrefix(header, "Token ") == expectedToken
	}
	// Basic <base64(token:)>
	if strings.HasPrefix(header, "Basic ") {
		decoded, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(header, "Basic "))
		if err == nil && strings.HasSuffix(string(decoded), ":") {
			return strings.TrimSuffix(string(decoded), ":") == expectedToken
		}
	}
	return false
}
