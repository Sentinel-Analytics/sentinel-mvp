package sentinel

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

// ... (InitDB and createTables are unchanged) ...
func InitDB() {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://sentinel:password@db:5432/sentinel?sslmode=disable"
		log.Println("DATABASE_URL not found, using default Docker connection string.")
	}
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}
	for i := 0; i < 5; i++ {
		err = db.Ping()
		if err == nil {
			break
		}
		log.Printf("Database not ready, retrying in 2 seconds... (%v)", err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("Could not connect to the database after several retries: %q", err)
	}
	log.Println("Successfully connected to the database.")
	createTables()
}

func createTables() {
	enableExtension := `CREATE EXTENSION IF NOT EXISTS "pgcrypto";`
	if _, err := db.Exec(enableExtension); err != nil {
		log.Fatalf("Could not enable pgcrypto extension: %v", err)
	}
	createUsersTable := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        email TEXT NOT NULL UNIQUE,
        password_hash TEXT NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );`
	if _, err := db.Exec(createUsersTable); err != nil {
		log.Fatalf("Could not create users table: %v", err)
	}
	createSitesTable := `
    CREATE TABLE IF NOT EXISTS sites (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        name TEXT NOT NULL,
        domain TEXT,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );`
	if _, err := db.Exec(createSitesTable); err != nil {
		log.Fatalf("Could not create sites table: %v", err)
	}
	log.Println("Database tables are set up.")
}


// --- AUTHENTICATION & PAGE HANDLERS ---

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("sentinel_session")
		if err != nil {
			http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
			return
		}
		userID, err := strconv.Atoi(cookie.Value)
		if err != nil || userID == 0 {
			http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}


func SignupPageHandler(w http.ResponseWriter, r *http.Request) {
	// This will now be handled by the React frontend router
	// This function can be removed if you don't need a direct server-side route
}

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
    // This will now be handled by the React frontend router
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	email := creds.Email
	password := creds.Password

	if email == "" || password == "" {
		http.Error(w, `{"error": "Email and password cannot be empty"}`, http.StatusBadRequest)
		return
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}
	var userID int
	err = db.QueryRow("INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id", email, hashedPassword).Scan(&userID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			http.Error(w, `{"error": "Could not create user (email might be taken)"}`, http.StatusBadRequest)
		} else {
			log.Printf("Error creating user: %v", err)
			http.Error(w, `{"error": "An unexpected error occurred"}`, http.StatusInternalServerError)
		}
		return
	}

	// Set session cookie upon successful signup
	http.SetCookie(w, &http.Cookie{
		Name:     "sentinel_session",
		Value:    strconv.Itoa(userID),
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true, // Important for cross-domain
		SameSite: http.SameSiteNoneMode,
		Domain:   ".getmusterup.com", // Set to the parent domain
	})

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}

// CORRECTED LoginHandler
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	email := creds.Email
	password := creds.Password

	var storedHash string
	var userID int
	err := db.QueryRow("SELECT id, password_hash FROM users WHERE email = $1", email).Scan(&userID, &storedHash)
	if err != nil {
		http.Error(w, `{"error": "Invalid credentials"}`, http.StatusUnauthorized)
		return
	}
	if !checkPasswordHash(password, storedHash) {
		http.Error(w, `{"error": "Invalid credentials"}`, http.StatusUnauthorized)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "sentinel_session",
		Value:    strconv.Itoa(userID),
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true, // Important for cross-domain
		SameSite: http.SameSiteNoneMode,
		Domain:   ".getmusterup.com", // Set to the parent domain
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "sentinel_session",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Domain:   ".getmusterup.com",
	})
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out"})
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// This will be handled by the React app's routing
}

