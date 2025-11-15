package sentinel

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/oschwald/geoip2-golang"
	"github.com/ua-parser/uap-go/uaparser"
)

// --- EVENT TRACKING ---

type Event struct {
	SiteID      string   `json:"siteId"`
	URL         string   `json:"url"`
	Referrer    string   `json:"referrer"`
	ScreenWidth int      `json:"screenWidth"`
	LCP         *float64 `json:"LCP,omitempty"`
	CLS         *float64 `json:"CLS,omitempty"`
	FID         *float64 `json:"FID,omitempty"`
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
	TrustScore  uint8
	LCP         sql.NullFloat64
	CLS         sql.NullFloat64
	FID         sql.NullFloat64
}

// --- ANALYTICS ENGINE ---

var uaParser *uaparser.Parser
var geoipDb *geoip2.Reader
var asnDb *geoip2.Reader

func InitAnalyticsEngine() {
	var err error
	uaParser = uaparser.NewFromSaved()

	geoipDb, err = geoip2.Open("GeoLite2-Country.mmdb")
	if err != nil {
		log.Printf("Warning: GeoIP database 'GeoLite2-Country.mmdb' not found. Country lookups will be disabled. Error: %v", err)
	}

	asnDb, err = geoip2.Open("GeoLite2-ASN.mmdb")
	if err != nil {
		log.Printf("Warning: ASN database 'GeoLite2-ASN.mmdb' not found. Bot detection will be less accurate. Error: %v", err)
	}
}

type Stats struct {
	TotalViews          uint64      `json:"totalViews"`
	UniqueVisitors      uint64      `json:"uniqueVisitors"`
	BounceRate          float64     `json:"bounceRate"`
	AvgVisitTime        string      `json:"avgVisitTime"`
	TrafficQualityScore float64     `json:"trafficQualityScore"`
	AvgLCP              float64     `json:"avgLcp"`
	AvgCLS              float64     `json:"avgCls"`
	AvgFID              float64     `json:"avgFid"`
	TopPages            []CountStat `json:"topPages"`
	TopReferrers        []CountStat `json:"topReferrers"`
	TopBrowsers         []CountStat `json:"topBrowsers"`
	TopOS               []CountStat `json:"topOS"`
	TopCountries        []CountStat `json:"topCountries"`
}

type CountStat struct {
	Value string `json:"value"`
	Count uint64 `json:"count"`
}

func getClientIP(r *http.Request) string {
	// Check for X-Forwarded-For header first
	forwardedFor := r.Header.Get("X-Forwarded-For")
	if forwardedFor != "" {
		// The header can contain a comma-separated list of IPs. The first one is the original client.
		ips := strings.Split(forwardedFor, ",")
		return strings.TrimSpace(ips[0])
	}

	// Fallback to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr // or handle error appropriately
	}
	return ip
}

// List of known bot user agent substrings
var botUserAgents = []string{
	"bot", "spider", "crawler", "monitor", "Go-http-client", "python-requests",
}

func calculateTrustScore(ip net.IP, userAgent string) uint8 {
	score := 100 // Start with a perfect score

	// Penalty for known bot user agents
	for _, botString := range botUserAgents {
		if strings.Contains(strings.ToLower(userAgent), botString) {
			score -= 50
			break // Apply penalty once
		}
	}

	// Penalty for being from a data center (ASN lookup)
	if asnDb != nil && ip != nil {
		_, err := asnDb.ASN(ip)
		if err == nil {
			score -= 40
		}
	}

	if score < 0 {
		return 0
	}
	return uint8(score)
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

	ipStr := getClientIP(r)
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

	trustScore := calculateTrustScore(ip, userAgent)

	// --- Firewall Blocking Logic ---
	// Check if the request should be blocked by firewall rules
    var asn string
    if asnDb != nil && ip != nil {
        record, err := asnDb.ASN(ip)
        if err == nil {
            asn = record.AutonomousSystemOrganization
        }
    }
	if isBlocked(event.SiteID, ipStr, country, asn) {
		http.Error(w, "Forbidden by firewall", http.StatusForbidden)
		return
	}

	eventData := EventData{
		Timestamp:   time.Now().UTC(),
		SiteID:      event.SiteID,
		ClientIP:    ipStr,
		URL:         event.URL,
		Referrer:    event.Referrer,
		ScreenWidth: uint16(event.ScreenWidth),
		Browser:     browser,
		OS:          osFamily,
		Country:     country,
		TrustScore:  trustScore,
		LCP:         nullFloat64(event.LCP),
		CLS:         nullFloat64(event.CLS),
		FID:         nullFloat64(event.FID),
	}

	// Insert into ClickHouse
	ctx := context.Background()
	err := chConn.AsyncInsert(ctx, "INSERT INTO events VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", false,
		eventData.Timestamp,
		eventData.SiteID,
		eventData.ClientIP,
		eventData.URL,
		eventData.Referrer,
		eventData.ScreenWidth,
		eventData.Browser,
		eventData.OS,
		eventData.Country,
		eventData.TrustScore,
		eventData.LCP,
		eventData.CLS,
		eventData.FID,
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

func isBlocked(siteID, ip, country, asn string) bool {
	rows, err := db.Query("SELECT rule_type, value FROM firewall_rules WHERE site_id = $1", siteID)
	if err != nil {
		log.Printf("Error querying firewall rules: %v", err)
		return false
	}
	defer rows.Close()

	for rows.Next() {
		var ruleType, value string
		if err := rows.Scan(&ruleType, &value); err != nil {
			log.Printf("Error scanning firewall rule: %v", err)
			continue
		}

		switch ruleType {
		case "ip":
			// Handle CIDR ranges as well
			if strings.Contains(value, "/") {
				_, ipNet, err := net.ParseCIDR(value)
				if err == nil && ipNet.Contains(net.ParseIP(ip)) {
					return true
				}
			} else if value == ip {
				return true
			}
		case "country":
			if value == country {
				return true
			}
		case "asn":
			// This is a simplified check. A more robust solution would involve ASN lookup for the incoming IP.
			if value == asn {
				return true
			}
		}
	}
	return false
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

	// Bounce Rate
	// This query first calculates the number of pageviews for each visitor (session).
	// Then it counts how many of those visitors had only one pageview (bounces).
	// Finally, it divides the number of bounces by the total number of visitors.
	queryBounceRate := `
		SELECT
			(countIf(pageviews = 1) / count()) * 100
		FROM (
			SELECT
				ClientIP,
				count() AS pageviews
			FROM events
			WHERE SiteID = ? AND Timestamp >= now() - INTERVAL ? DAY
			GROUP BY ClientIP
		)
	`
	err = chConn.QueryRow(ctx, queryBounceRate, siteID, days).Scan(&stats.BounceRate)
	if err != nil {
		// It's possible there's no data, which can cause an error. Default to 0.
		stats.BounceRate = 0
	}

	// Average Visit Duration
	// This query calculates the duration of each session (time between max and min timestamp for each visitor).
	// Then it calculates the average of these durations.
	queryAvgVisitTime := `
		SELECT
			avg(duration)
		FROM (
			SELECT
				ClientIP,
				date_diff('second', min(Timestamp), max(Timestamp)) AS duration
			FROM events
			WHERE SiteID = ? AND Timestamp >= now() - INTERVAL ? DAY
			GROUP BY ClientIP
		)
	`
	var avgSeconds float64
	err = chConn.QueryRow(ctx, queryAvgVisitTime, siteID, days).Scan(&avgSeconds)
	if err != nil {
		stats.AvgVisitTime = "0s"
	} else {
		// Format the duration into a more readable string (e.g., 1m 23s)
		d := time.Duration(avgSeconds) * time.Second
		stats.AvgVisitTime = d.Round(time.Second).String()
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

	// Traffic Quality Score
	queryGoodTraffic := "SELECT count() FROM events WHERE SiteID = ? AND Timestamp >= now() - INTERVAL ? DAY AND TrustScore > 50"
	var goodTrafficCount uint64
	err = chConn.QueryRow(ctx, queryGoodTraffic, siteID, days).Scan(&goodTrafficCount)
	if err != nil {
		log.Printf("Error querying good traffic count: %v", err)
		stats.TrafficQualityScore = 0
	} else if stats.TotalViews > 0 {
		stats.TrafficQualityScore = (float64(goodTrafficCount) / float64(stats.TotalViews)) * 100
	} else {
		stats.TrafficQualityScore = 0
	}

	// Average Web Vitals
	queryAvgLCP := "SELECT avg(LCP) FROM events WHERE SiteID = ? AND Timestamp >= now() - INTERVAL ? DAY"
	err = chConn.QueryRow(ctx, queryAvgLCP, siteID, days).Scan(&stats.AvgLCP)
	if err != nil {
		stats.AvgLCP = 0
	}

	queryAvgCLS := "SELECT avg(CLS) FROM events WHERE SiteID = ? AND Timestamp >= now() - INTERVAL ? DAY"
	err = chConn.QueryRow(ctx, queryAvgCLS, siteID, days).Scan(&stats.AvgCLS)
	if err != nil {
		stats.AvgCLS = 0
	}

	queryAvgFID := "SELECT avg(FID) FROM events WHERE SiteID = ? AND Timestamp >= now() - INTERVAL ? DAY"
	err = chConn.QueryRow(ctx, queryAvgFID, siteID, days).Scan(&stats.AvgFID)
	if err != nil {
		stats.AvgFID = 0
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

// Helper function to convert *float64 to sql.NullFloat64
func nullFloat64(f *float64) sql.NullFloat64 {
	if f == nil {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{Float64: *f, Valid: true}
}

