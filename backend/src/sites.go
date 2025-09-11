package sentinel

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Site struct represents a website being tracked in the database.
type Site struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Domain string `json:"domain,omitempty"`
}

// SitesApiHandler now routes to different functions based on the request.
// This is a more robust way to handle RESTful routing.
func SitesApiHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/sites")
	path = strings.TrimSuffix(path, "/")

	// If the path is empty, it's a request for the whole collection (GET list, POST create)
	if path == "" {
		switch r.Method {
		case "GET":
			handleListSites(w, r)
		case "POST":
			handleCreateSite(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	// If the path is not empty, it should be an ID for a specific site
	// (PUT update, DELETE remove)
	siteID := path[1:] // remove the leading "/"
	switch r.Method {
	case "PUT":
		handleUpdateSite(w, r, siteID)
	case "DELETE":
		handleDeleteSite(w, r, siteID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleListSites handles GET /api/sites
func handleListSites(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)

	rows, err := db.Query("SELECT id, name, domain FROM sites WHERE user_id = $1 ORDER BY created_at DESC", userID)
	if err != nil {
		http.Error(w, "Failed to fetch sites", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	sites := []Site{}
	for rows.Next() {
		var s Site
		if err := rows.Scan(&s.ID, &s.Name, &s.Domain); err != nil {
			http.Error(w, "Failed to scan site", http.StatusInternalServerError)
			return
		}
		sites = append(sites, s)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sites)
}

// handleCreateSite handles POST /api/sites
func handleCreateSite(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	var site Site
	if err := json.NewDecoder(r.Body).Decode(&site); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var newSiteID string
	err := db.QueryRow("INSERT INTO sites (user_id, name, domain) VALUES ($1, $2, $3) RETURNING id", userID, site.Name, site.Domain).Scan(&newSiteID)
	if err != nil {
		http.Error(w, "Failed to create site", http.StatusInternalServerError)
		return
	}

	site.ID = newSiteID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(site)
}

// handleUpdateSite handles PUT /api/sites/{id}
func handleUpdateSite(w http.ResponseWriter, r *http.Request, siteID string) {
	userID := r.Context().Value("userID").(int)
	var site Site
	if err := json.NewDecoder(r.Body).Decode(&site); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check for ownership before updating
	var ownerID int
	err := db.QueryRow("SELECT user_id FROM sites WHERE id = $1", siteID).Scan(&ownerID)
	if err != nil || ownerID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	_, err = db.Exec("UPDATE sites SET name = $1, domain = $2 WHERE id = $3 AND user_id = $4", site.Name, site.Domain, siteID, userID)
	if err != nil {
		http.Error(w, "Failed to update site", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(site)
}

// handleDeleteSite handles DELETE /api/sites/{id}
func handleDeleteSite(w http.ResponseWriter, r *http.Request, siteID string) {
	userID := r.Context().Value("userID").(int)

	// Check for ownership before deleting
	var ownerID int
	err := db.QueryRow("SELECT user_id FROM sites WHERE id = $1", siteID).Scan(&ownerID)
	if err != nil || ownerID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	_, err = db.Exec("DELETE FROM sites WHERE id = $1 AND user_id = $2", siteID, userID)
	if err != nil {
		http.Error(w, "Failed to delete site", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

