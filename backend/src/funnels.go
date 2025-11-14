package sentinel

import (
	"encoding/json"
	"net/http"
)

type Funnel struct {
	ID     string   `json:"id"`
	SiteID string   `json:"siteId"`
	Name   string   `json:"name"`
	Steps  []string `json:"steps"`
}

// FunnelsApiHandler routes requests to appropriate functions based on HTTP method.
func FunnelsApiHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handleListFunnels(w, r)
	case "POST":
		handleCreateFunnel(w, r)
	case "PUT":
		handleUpdateFunnel(w, r)
	case "DELETE":
		handleDeleteFunnel(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleListFunnels(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	siteID := r.URL.Query().Get("siteId")
	if siteID == "" {
		http.Error(w, "siteId query parameter is required", http.StatusBadRequest)
		return
	}

	// Verify site ownership
	var ownerID int
	err := db.QueryRow("SELECT user_id FROM sites WHERE id = $1", siteID).Scan(&ownerID)
	if err != nil || ownerID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	rows, err := db.Query("SELECT id, site_id, name, steps FROM funnels WHERE site_id = $1 ORDER BY name", siteID)
	if err != nil {
		http.Error(w, "Failed to fetch funnels", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	funnels := []Funnel{}
	for rows.Next() {
		var funnel Funnel
		var stepsJSON []byte
		if err := rows.Scan(&funnel.ID, &funnel.SiteID, &funnel.Name, &stepsJSON); err != nil {
			http.Error(w, "Failed to scan funnel", http.StatusInternalServerError)
			return
		}
		if err := json.Unmarshal(stepsJSON, &funnel.Steps); err != nil {
			http.Error(w, "Failed to parse funnel steps", http.StatusInternalServerError)
			return
		}
		funnels = append(funnels, funnel)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(funnels)
}

func handleCreateFunnel(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	var funnel Funnel
	if err := json.NewDecoder(r.Body).Decode(&funnel); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Verify site ownership
	var ownerID int
	err := db.QueryRow("SELECT user_id FROM sites WHERE id = $1", funnel.SiteID).Scan(&ownerID)
	if err != nil || ownerID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	stepsJSON, err := json.Marshal(funnel.Steps)
	if err != nil {
		http.Error(w, "Failed to serialize funnel steps", http.StatusInternalServerError)
		return
	}

	var newFunnelID string
	err = db.QueryRow("INSERT INTO funnels (site_id, name, steps) VALUES ($1, $2, $3) RETURNING id", funnel.SiteID, funnel.Name, stepsJSON).Scan(&newFunnelID)
	if err != nil {
		http.Error(w, "Failed to create funnel", http.StatusInternalServerError)
		return
	}

	funnel.ID = newFunnelID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(funnel)
}

func handleUpdateFunnel(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	var funnel Funnel
	if err := json.NewDecoder(r.Body).Decode(&funnel); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Verify funnel ownership via site ownership
	var ownerID int
	err := db.QueryRow("SELECT s.user_id FROM sites s JOIN funnels f ON s.id = f.site_id WHERE f.id = $1", funnel.ID).Scan(&ownerID)
	if err != nil || ownerID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	stepsJSON, err := json.Marshal(funnel.Steps)
	if err != nil {
		http.Error(w, "Failed to serialize funnel steps", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("UPDATE funnels SET name = $1, steps = $2 WHERE id = $3", funnel.Name, stepsJSON, funnel.ID)
	if err != nil {
		http.Error(w, "Failed to update funnel", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(funnel)
}

func handleDeleteFunnel(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	funnelID := r.URL.Query().Get("id")
	if funnelID == "" {
		http.Error(w, "id query parameter is required", http.StatusBadRequest)
		return
	}

	// Verify funnel ownership via site ownership
	var ownerID int
	err := db.QueryRow("SELECT s.user_id FROM sites s JOIN funnels f ON s.id = f.site_id WHERE f.id = $1", funnelID).Scan(&ownerID)
	if err != nil || ownerID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	_, err = db.Exec("DELETE FROM funnels WHERE id = $1", funnelID)
	if err != nil {
		http.Error(w, "Failed to delete funnel", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
