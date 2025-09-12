// This file is now in the 'backend' directory, not 'src'.
package main // Note: The package is now 'main'

import (
	"log"
	"net/http"

	// This now correctly imports your 'sentinel' library package
	sentinel "sentinel.go/src"

	"github.com/rs/cors"
)

func main() {
	// All functions from your library are now prefixed with 'sentinel.'
	sentinel.InitDB()
	sentinel.InitAnalyticsEngine()

	mux := http.NewServeMux()

	// The file server now needs to look inside the 'static' folder
	// which will be created in the Docker container.
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// --- Public API Routes ---
	mux.HandleFunc("/auth/signup", sentinel.SignupHandler)
	mux.HandleFunc("/auth/login", sentinel.LoginHandler)
	mux.HandleFunc("/track", sentinel.TrackHandler)

	// --- Protected API Routes ---
	mux.HandleFunc("/logout", sentinel.AuthMiddleware(sentinel.LogoutHandler))
	mux.HandleFunc("/api/sites/", sentinel.AuthMiddleware(sentinel.SitesApiHandler))
	mux.HandleFunc("/api/dashboard", sentinel.AuthMiddleware(sentinel.DashboardApiHandler))

	// CORS Middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})
	handler := c.Handler(mux)

	log.Println("Sentinel Go server starting on :8000")
	if err := http.ListenAndServe(":8000", handler); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
