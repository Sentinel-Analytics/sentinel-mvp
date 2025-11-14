package main

import (
	"log"
	"net/http"

	_ "sentinel-backend/docs"
	sentinel "sentinel-backend/src"

	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Sentinel API
// @version 1.0
// @description This is a sample server for a web analytics platform.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:6060
// @BasePath /
func main() {
	// All functions from your library are now prefixed with 'sentinel.'
	sentinel.InitDB()
	sentinel.InitAnalyticsEngine()
	sentinel.InitClickHouse()

	mux := http.NewServeMux()

	// The file server now needs to look inside the 'static' folder
	// which will be created in the Docker container.
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// --- CORS Policies ---
	// Permissive CORS for the tracking endpoint
	trackCors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"POST", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	})

	// Strict CORS for the dashboard and API
	apiCors := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://sentinel-mvp.getmusterup.com", "https://sentinel.getmusterup.com", "http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// --- Public API Routes ---
	mux.Handle("/auth/signup", apiCors.Handler(http.HandlerFunc(sentinel.SignupHandler)))
	mux.Handle("/auth/login", apiCors.Handler(http.HandlerFunc(sentinel.LoginHandler)))
	mux.Handle("/track", trackCors.Handler(http.HandlerFunc(sentinel.TrackHandler)))
	mux.Handle("/session", trackCors.Handler(http.HandlerFunc(sentinel.SessionHandler)))
	mux.Handle("/api/session", trackCors.Handler(http.HandlerFunc(sentinel.SessionHandler)))

	// --- Protected API Routes ---
	mux.Handle("/logout", apiCors.Handler(sentinel.AuthMiddleware(sentinel.LogoutHandler)))
	mux.Handle("/api/sites/", apiCors.Handler(sentinel.AuthMiddleware(sentinel.SitesApiHandler)))
	mux.Handle("/api/dashboard", apiCors.Handler(sentinel.AuthMiddleware(sentinel.DashboardApiHandler)))
	mux.Handle("/api/firewall", apiCors.Handler(sentinel.AuthMiddleware(sentinel.FirewallApiHandler)))
	mux.Handle("/api/session/events", apiCors.Handler(sentinel.AuthMiddleware(sentinel.GetSessionEventsHandler)))
	mux.Handle("/api/sessions", apiCors.Handler(sentinel.AuthMiddleware(sentinel.ListSessionsHandler)))
	mux.Handle("/api/funnels/", apiCors.Handler(sentinel.AuthMiddleware(sentinel.FunnelsApiHandler)))

	// Swagger documentation
	mux.HandleFunc("/docs/", httpSwagger.WrapHandler)

	log.Println("Sentinel Go server starting on :6060")
	if err := http.ListenAndServe(":6060", mux); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}

