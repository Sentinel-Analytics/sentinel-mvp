package main

import (
	"log"
	"net/http"

	_ "sentinel-backend/docs"
	sentinel "sentinel-backend/src"

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
	// Initialize DB and engines
	sentinel.InitDB()
	sentinel.InitAnalyticsEngine()
	sentinel.InitClickHouse()

	mux := http.NewServeMux()

	// Serve static files
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

	log.Println("Sentinel Go server starting on :6060")
	if err := http.ListenAndServe(":6060", CORSMiddleware(mux)); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}

// CORSMiddleware adds CORS headers to all responses
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only allow requests from your frontend domain
		frontendOrigin := "https://sentinel-mvp.getmusterup.com"
		w.Header().Set("Access-Control-Allow-Origin", frontendOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true") // Important for cookies / credentials

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

