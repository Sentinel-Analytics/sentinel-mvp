package sentinel

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/oschwald/geoip2-golang"
	"github.com/ua-parser/uap-go/uaparser"
)

// --- EVENT TRACKING ---

type Event struct {
	SiteID      string `json:"siteId"`
	URL         string `json:"url"`
	Referrer    string `json:"referrer"`
	ScreenWidth int    `json:"screenWidth"`
}

type EventData struct {
	Timestamp   time.Time
	SiteID      string
	ClientIP    string
	URL         string
	Referrer    string
	ScreenWidth uint16
	Browser     string
	OS          string
	Country     string
}

// --- ANALYTICS ENGINE ---

var uaParser *uaparser.Parser
var geoipDb *geoip2.Reader

func InitAnalyticsEngine() {
	var err error
	uaParser = uaparser.NewFromSaved()

	geoipDb, err = geoip2.Open("GeoLite2-Country.mmdb")
	if err != nil {
		log.Printf("Warning: GeoIP database 'GeoLite2-Country.mmdb' not found. Country lookups will be disabled. Error: %v", err)
	}
}

type Stats struct {
	TotalViews     int         `json:"totalViews"`
	UniqueVisitors int         `json:"uniqueVisitors"`
	TopPages       []CountStat `json:"topPages"`
	TopReferrers   []CountStat `json:"topReferrers"`
	TopBrowsers    []CountStat `json:"topBrowsers"`
	TopOS          []CountStat `json:"topOS"`
	TopCountries   []CountStat `json:"topCountries"`
}

type CountStat struct {
	Value string `json:"value"`
	Count int    `json:"count"`
}

func TrackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var event Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	userAgent := r.UserAgent()
	client := uaParser.Parse(userAgent)

	ipStr, _, _ := net.SplitHostPort(r.RemoteAddr)
	ip := net.ParseIP(ipStr)

	country := "Unknown"
	if geoipDb != nil && ip != nil {
		record, err := geoipDb.Country(ip)
		if err == nil && record.Country.IsoCode != "" {
			country = record.Country.IsoCode
		}
	}

	browser := client.UserAgent.Family
	if browser == "Other" {
		browser = "Unknown"
	}
	osFamily := client.Os.Family
	if osFamily == "Other" {
		osFamily = "Unknown"
	}

	eventData := EventData{
		Timestamp:   time.Now().UTC(),
		SiteID:      event.SiteID,
		ClientIP:    r.RemoteAddr,
		URL:         event.URL,
		Referrer:    event.Referrer,
		ScreenWidth: uint16(event.ScreenWidth),
		Browser:     browser,
		OS:          osFamily,
		Country:     country,
	}

	// Insert into ClickHouse
	ctx := context.Background()
	err := chConn.AsyncInsert(ctx, "INSERT INTO events VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", false,
		eventData.Timestamp,
		eventData.SiteID,
		eventData.ClientIP,
		eventData.URL,
		eventData.Referrer,
		eventData.ScreenWidth,
		eventData.Browser,
		eventData.OS,
		eventData.Country,
	)
	if err != nil {
		log.Printf("Error inserting event into ClickHouse: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func DashboardApiHandler(w http.ResponseWriter, r *http.Request) {
	siteID := r.URL.Query().Get("siteId")
	if siteID == "" {
		http.Error(w, "siteId query parameter is required", http.StatusBadRequest)
		return
	}

	daysStr := r.URL.Query().Get("days")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		days = 1
	}

	stats, err := calculateStats(siteID, days)
	if err != nil {
		log.Printf("Error calculating stats: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func calculateStats(siteID string, days int) (Stats, error) {
	ctx := context.Background()
	var stats Stats

	// Total Views
	queryTotalViews := "SELECT count() FROM events WHERE SiteID = ? AND Timestamp >= now() - INTERVAL ? DAY"
	err := chConn.QueryRow(ctx, queryTotalViews, siteID, days).Scan(&stats.TotalViews)
	if err != nil {
		return stats, err
	}

	// Unique Visitors
	queryUniqueVisitors := "SELECT uniq(ClientIP) FROM events WHERE SiteID = ? AND Timestamp >= now() - INTERVAL ? DAY"
	err = chConn.QueryRow(ctx, queryUniqueVisitors, siteID, days).Scan(&stats.UniqueVisitors)
	if err != nil {
		return stats, err
	}

	// Top Pages
	stats.TopPages, err = queryTopStats(ctx, "URL", siteID, days)
	if err != nil {
		return stats, err
	}

	// Top Referrers
	stats.TopReferrers, err = queryTopStats(ctx, "Referrer", siteID, days)
	if err != nil {
		return stats, err
	}

	// Top Browsers
	stats.TopBrowsers, err = queryTopStats(ctx, "Browser", siteID, days)
	if err != nil {
		return stats, err
	}

	// Top OS
	stats.TopOS, err = queryTopStats(ctx, "OS", siteID, days)
	if err != nil {
		return stats, err
	}

	// Top Countries
	stats.TopCountries, err = queryTopStats(ctx, "Country", siteID, days)
	if err != nil {
		return stats, err
	}

	return stats, nil
}

func queryTopStats(ctx context.Context, column, siteID string, days int) ([]CountStat, error) {
	query := "SELECT " + column + ", count() AS c FROM events WHERE SiteID = ? AND Timestamp >= now() - INTERVAL ? DAY GROUP BY " + column + " ORDER BY c DESC LIMIT 10"
	rows, err := chConn.Query(ctx, query, siteID, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []CountStat
	for rows.Next() {
		var stat CountStat
		if err := rows.Scan(&stat.Value, &stat.Count); err != nil {
			return nil, err
		}
		result = append(result, stat)
	}

	return result, nil
}

