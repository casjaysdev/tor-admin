// File: internal/web/routes.go
// Purpose: Define and register web and API routes, applying auth middleware.

package web

import (
	"net/http"

	"tor-admin/internal/auth"
	"tor-admin/internal/ui"
	"tor-admin/internal/web/handlers"
)

// RegisterRoutes sets up all HTTP routes for the web UI and API.
func RegisterRoutes(mux *http.ServeMux) {
	templateFS := ui.GetTemplateFS()

	// Static file handler (served from /static/)
	ui.RegisterStatic(mux)

	// Public routes (setup, login)
	mux.HandleFunc("/setup", handlers.SetupHandler(templateFS))
	mux.HandleFunc("/login", handlers.LoginHandler(templateFS))

	// Auth-protected routes
	mux.Handle("/", auth.RequireLogin(handlers.IndexHandler(templateFS)))
	mux.Handle("/config", auth.RequireLogin(handlers.ConfigHandler(templateFS)))
	mux.Handle("/logout", auth.RequireLogin(http.HandlerFunc(handlers.LogoutHandler)))
	mux.Handle("/api/hidden", auth.RequireLogin(http.HandlerFunc(handlers.HiddenServicesHandler)))
	mux.Handle("/api/bandwidth", auth.RequireLogin(http.HandlerFunc(handlers.BandwidthHandler)))

	// Optional: health check
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
}
