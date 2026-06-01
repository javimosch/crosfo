package main

import (
	"ffaf/pkg/database"
	"ffaf/pkg/handlers"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func startServer(port int) {
	dbPath := database.GetDefaultDBPath()
	if err := database.Init(dbPath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	mux := http.NewServeMux()

	// Static files
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/", handleHome)
	mux.HandleFunc("/c/", handleCommunityOrUser) // handles /c/{community} and /c/{community}/user/{nickname}
	mux.HandleFunc("/stats", handleStats)

	// 404 handler
	mux.HandleFunc("/404", handle404)

	// API endpoints
	mux.HandleFunc("/api/community-admins/", handlers.HandleCommunityAdmins) // Admin management (must be before /api/community/)
	mux.HandleFunc("/api/community/update", handlers.HandleUpdateCommunity) // Community editing
	mux.HandleFunc("/api/communities", handlers.HandleCommunities)
	mux.HandleFunc("/api/community/", handlers.HandleCommunity)
	mux.HandleFunc("/api/entries/", handlers.HandleEntries)
	mux.HandleFunc("/api/thumb-up", handlers.HandleThumbUp)
	mux.HandleFunc("/api/thumbs-up", handlers.HandleThumbsUp)
	mux.HandleFunc("/api/user-entries", handlers.HandleUserEntries)
	mux.HandleFunc("/api/auth", handlers.HandleAuth)
	mux.HandleFunc("/api/entry/update", handlers.HandleUpdateEntry)
	mux.HandleFunc("/api/entry/delete", handlers.HandleDeleteEntry)
	mux.HandleFunc("/api/thumbs-breakdown/", handlers.HandleThumbsBreakdown)
	mux.HandleFunc("/api/user-profile", handlers.HandleUserProfile)
	mux.HandleFunc("/api/stats", handlers.HandleStats)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	log.Printf("Crosfo server starting on http://localhost:%d", port)
	log.Printf("Database: %s", dbPath)
	log.Printf("Press Ctrl+C to stop")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	html, err := os.ReadFile("templates/index.html")
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(html)
}

// handleCommunityOrUser serves community.html for /c/{community}
// and user.html for /c/{community}/user/{nickname}
func handleCommunityOrUser(w http.ResponseWriter, r *http.Request) {
	// Check if path matches /c/{community}/user/{nickname}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) >= 4 && parts[2] == "user" {
		html, err := os.ReadFile("templates/user.html")
		if err != nil {
			http.Error(w, "Template not found", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write(html)
		return
	}

	html, err := os.ReadFile("templates/community.html")
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(html)
}

func handleStats(w http.ResponseWriter, r *http.Request) {
	html, err := os.ReadFile("templates/stats.html")
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(html)
}

func handle404(w http.ResponseWriter, r *http.Request) {
	html, err := os.ReadFile("templates/404.html")
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(html)
}