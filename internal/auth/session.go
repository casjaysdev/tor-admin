// File: internal/auth/session.go
// Purpose: Web session management using secure cookies

package auth

import (
	"context"
	"net/http"

	"github.com/gorilla/securecookie"
)

var (
	// Key generation should happen once and be stored securely
	hashKey  = []byte("super-secret-hash-key-change-me")
	blockKey = []byte("super-secret-block-key-32b-123456789012") // 32 bytes

	sc = securecookie.New(hashKey, blockKey)
)

const sessionCookieName = "tor_admin_session"
const contextUserKey = "authenticatedUser"

// SetSession writes the session cookie after successful login
func SetSession(w http.ResponseWriter, r *http.Request, username string) error {
	encoded, err := sc.Encode(sessionCookieName, map[string]string{
		"user": username,
	})
	if err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    encoded,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true if behind HTTPS
	})
	return nil
}

// ClearSession clears the user's session cookie
func ClearSession(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
}

// WithSession wraps an http.Handler and injects session info into context
func WithSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(sessionCookieName)
		if err == nil {
			val := make(map[string]string)
			if err := sc.Decode(sessionCookieName, cookie.Value, &val); err == nil {
				if user := val["user"]; user != "" {
					ctx := context.WithValue(r.Context(), contextUserKey, user)
					r = r.WithContext(ctx)
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}

// GetUser returns the currently logged-in username (or empty string)
func GetUser(r *http.Request) string {
	if user, ok := r.Context().Value(contextUserKey).(string); ok {
		return user
	}
	return ""
}
