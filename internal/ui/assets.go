// File: internal/ui/assets.go
// Purpose: Embed static assets and serve them via /static and templates

package ui

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//go:embed static/* templates/*
var content embed.FS

// RegisterStatic sets up the embedded file handler for static assets
func RegisterStatic(mux *http.ServeMux) {
	staticFiles, err := fs.Sub(content, "static")
	if err != nil {
		log.Fatalf("failed to load static assets: %v", err)
	}
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFiles))))
}

// GetTemplateFS exposes the template subdirectory for parsing
func GetTemplateFS() fs.FS {
	tmpls, err := fs.Sub(content, "templates")
	if err != nil {
		log.Fatalf("failed to load templates: %v", err)
	}
	return tmpls
}
