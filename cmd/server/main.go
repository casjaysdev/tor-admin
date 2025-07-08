// File: cmd/server/main.go
// Purpose: Main entrypoint for the tor-admin application

package main

import (
	"log"
	"net/http"
	"os"

	"tor-admin/internal/auth"
	"tor-admin/internal/ui"
	"tor-admin/web"
)

func main() {
	// Load configuration and user data
	cfg, err := auth.LoadOrInitConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Set default listen port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create base router
	mux := http.NewServeMux()

	// Register static files and embedded assets
	ui.RegisterStatic(mux)

	// Register web routes (handlers + templates)
	web.RegisterRoutes(mux, cfg)

	// Wrap with top-level middleware
	handler := auth.WithSession(mux)

	log.Printf("🚀 tor-admin is running on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
