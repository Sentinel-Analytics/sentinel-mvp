package sentinel

import (
	_ "database/sql"
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"strings"
)

type FirewallRule struct {
	ID       string `json:"id"`
	SiteID   string `json:"siteId"`
	RuleType string `json:"ruleType"` // e.g., "ip", "country", "asn"
	Value    string `json:"value"`    // The actual IP, country code, or ASN
}

// FirewallApiHandler routes requests to appropriate functions based on HTTP method.
func FirewallApiHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handleListFirewallRules(w, r)
	case "POST":
		handleCreateFirewallRule(w, r)
	case "DELETE":
		handleDeleteFirewallRule(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// @Summary List firewall rules
// @Description Get a list of all firewall rules for a specific site.
// @Tags firewall
// @Produce  json
// @Param siteId query string true "Site ID"
// @Success 200 {array} FirewallRule
// @Router /api/firewall [get]
func handleListFirewallRules(w http.ResponseWriter, r *http.Request) {
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

	rows, err := db.Query("SELECT id, site_id, rule_type, value FROM firewall_rules WHERE site_id = $1 ORDER BY rule_type, value", siteID)
	if err != nil {
		http.Error(w, "Failed to fetch firewall rules", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	rules := []FirewallRule{}
	for rows.Next() {
		var rule FirewallRule
		if err := rows.Scan(&rule.ID, &rule.SiteID, &rule.RuleType, &rule.Value); err != nil {
			http.Error(w, "Failed to scan firewall rule", http.StatusInternalServerError)
			return
		}
		rules = append(rules, rule)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rules)
}

// @Summary Create a new firewall rule
// @Description Add a new firewall rule for a site.
// @Tags firewall
// @Accept  json
// @Produce  json
// @Param rule body FirewallRule true "Firewall rule to create"
// @Success 201 {object} FirewallRule
// @Router /api/firewall [post]
func handleCreateFirewallRule(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	var rule FirewallRule
	if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Verify site ownership
	var ownerID int
	err := db.QueryRow("SELECT user_id FROM sites WHERE id = $1", rule.SiteID).Scan(&ownerID)
	if err != nil || ownerID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Basic validation for rule type and value
	switch rule.RuleType {
	case "ip":
		if net.ParseIP(rule.Value) == nil && !strings.Contains(rule.Value, "/") {
			http.Error(w, "Invalid IP address or CIDR", http.StatusBadRequest)
			return
		}
	case "country":
		if len(rule.Value) != 2 {
			http.Error(w, "Country code must be 2 characters (ISO 3166-1 alpha-2)", http.StatusBadRequest)
			return
		}
	case "asn":
		// ASN values are typically numbers, but can be prefixed with AS. Simple check for now.
		if !strings.HasPrefix(strings.ToUpper(rule.Value), "AS") {
			// Attempt to parse as integer if no AS prefix
			if _, err := strconv.Atoi(rule.Value); err != nil {
				http.Error(w, "Invalid ASN value", http.StatusBadRequest)
				return
			}
		}
	default:
		http.Error(w, "Invalid rule type. Must be 'ip', 'country', or 'asn'", http.StatusBadRequest)
		return
	}

	var newRuleID string
	err = db.QueryRow("INSERT INTO firewall_rules (site_id, rule_type, value) VALUES ($1, $2, $3) RETURNING id", rule.SiteID, rule.RuleType, rule.Value).Scan(&newRuleID)
	if err != nil {
		http.Error(w, "Failed to create firewall rule", http.StatusInternalServerError)
		return
	}

	rule.ID = newRuleID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rule)
}

// @Summary Delete a firewall rule
// @Description Delete an existing firewall rule.
// @Tags firewall
// @Param id query string true "Rule ID"
// @Success 204 "No Content"
// @Router /api/firewall [delete]
func handleDeleteFirewallRule(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	ruleID := r.URL.Query().Get("id")
	if ruleID == "" {
		http.Error(w, "id query parameter is required", http.StatusBadRequest)
		return
	}

	// Verify rule ownership via site ownership
	var siteID string
	var ownerID int
	err := db.QueryRow("SELECT site_id FROM firewall_rules WHERE id = $1", ruleID).Scan(&siteID)
	if err != nil {
		http.Error(w, "Firewall rule not found", http.StatusNotFound)
		return
	}
	err = db.QueryRow("SELECT user_id FROM sites WHERE id = $1", siteID).Scan(&ownerID)
	if err != nil || ownerID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	_, err = db.Exec("DELETE FROM firewall_rules WHERE id = $1", ruleID)
	if err != nil {
		http.Error(w, "Failed to delete firewall rule", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
