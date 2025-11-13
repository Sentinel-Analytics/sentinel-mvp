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

	// --- Public API Routes ---
	mux.HandleFunc("/auth/signup", sentinel.SignupHandler)
	mux.HandleFunc("/auth/login", sentinel.LoginHandler)
	mux.HandleFunc("/track", sentinel.TrackHandler)

	// --- Protected API Routes ---
	mux.HandleFunc("/logout", sentinel.AuthMiddleware(sentinel.LogoutHandler))
	mux.HandleFunc("/api/sites/", sentinel.AuthMiddleware(sentinel.SitesApiHandler))
	mux.HandleFunc("/api/dashboard", sentinel.AuthMiddleware(sentinel.DashboardApiHandler))

	// Swagger documentation
	mux.HandleFunc("/docs/", httpSwagger.WrapHandler)

	// CORS Middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://sentinel-mvp.getmusterup.com", "https://sentinel.getmusterup.com"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})
	handler := c.Handler(mux)

	log.Println("Sentinel Go server starting on :6060")
	if err := http.ListenAndServe(":6060", handler); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}

