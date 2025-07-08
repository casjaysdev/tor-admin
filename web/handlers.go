// File: web/handlers.go
// Purpose: Handlers for UI and API routes (login, setup, API endpoints)

package web

import (
	"html/template"
	"net/http"
	"os"

	"tor-admin/internal/auth"
	"tor-admin/internal/bandwidth"
)

func IndexHandler(tfs templateFS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFS(tfs, "index.html"))
		_ = tmpl.Execute(w, map[string]any{
			"Username": auth.GetUser(r),
			"Theme":    "dracula", // TODO: load from cookie
		})
	}
}

func LoginHandler(tfs templateFS, cfg *auth.UserConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			r.ParseForm()
			user := r.FormValue("username")
			pass := r.FormValue("password")

			if user == cfg.Username && auth.CheckPassword(pass, cfg.PasswordHash) {
				_ = auth.SetSession(w, r, user)
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		tmpl := template.Must(template.ParseFS(tfs, "login.html"))
		_ = tmpl.Execute(w, nil)
	}
}

func SetupHandler(tfs templateFS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := auth.LoadOrInitConfig(); err == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		if r.Method == http.MethodPost {
			r.ParseForm()
			user := r.FormValue("username")
			pass := r.FormValue("password")

			if len(user) < 3 || len(pass) < 6 {
				http.Error(w, "Username or password too short", http.StatusBadRequest)
				return
			}
			hashed, _ := auth.HashPassword(pass)
			token := os.Getenv("API_TOKEN")
			if token == "" {
				token = auth.GenerateRandomToken()
			}

			_ = auth.SaveConfig(&auth.UserConfig{
				Username:     user,
				PasswordHash: hashed,
				APIToken:     token,
			})
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		tmpl := template.Must(template.ParseFS(tfs, "setup.html"))
		_ = tmpl.Execute(w, nil)
	}
}

func LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth.ClearSession(w)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

// ==== API ====

func StatusAPIHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"ok":true}`))
	}
}

func BandwidthAPIHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Example: convert 5 GB/month to bytes/sec
		bps, _ := bandwidth.ToBytesPerSecond(5, bandwidth.GB, bandwidth.Monthly)
		w.Write([]byte(`{"bps":` + bandwidth.PrettyPrintBytes(bps) + `}`))
	}
}

func HiddenServicesAPIHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"onions":["abc123.onion","xyz456.onion"]}`))
	}
}

func TorrcUpdateAPIHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// For now, just pretend we updated it
		w.Write([]byte(`{"saved":true}`))
	}
}
