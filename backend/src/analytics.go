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

	// Percentage changes
	TotalViewsChange          float64 `json:"totalViewsChange"`
	UniqueVisitorsChange      float64 `json:"uniqueVisitorsChange"`
	BounceRateChange          float64 `json:"bounceRateChange"`
	AvgVisitTimeChange        float64 `json:"avgVisitTimeChange"`
	TrafficQualityScoreChange float64 `json:"trafficQualityScoreChange"`
	AvgLCPChange              float64 `json:"avgLcpChange"`
	AvgCLSChange              float64 `json:"avgClsChange"`
	AvgFIDChange              float64 `json:"avgFidChange"`
}

type CoreStats struct {
	TotalViews          uint64
	UniqueVisitors      uint64
	BounceRate          float64
	AvgVisitTime        float64 // in seconds
	TrafficQualityScore float64
	AvgLCP              float64
	AvgCLS              float64
	AvgFID              float64
}

type CountStat struct {
	Value string `json:"value"`
	Count uint64 `json:"count"`
}

func getClientIP(r *http.Request) string {
	forwardedFor := r.Header.Get("X-Forwarded-For")
	if forwardedFor != "" {
		ips := strings.Split(forwardedFor, ",")
		return strings.TrimSpace(ips[0])
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

var botUserAgents = []string{
	"bot", "spider", "crawler", "monitor", "Go-http-client", "python-requests",
}

func calculateTrustScore(ip net.IP, userAgent string) uint8 {
	score := 100
	for _, botString := range botUserAgents {
		if strings.Contains(strings.ToLower(userAgent), botString) {
			score -= 50
			break
		}
	}
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

	ctx := context.Background()
	err := chConn.AsyncInsert(ctx, "INSERT INTO sentinel.events VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", false,
		eventData.Timestamp, eventData.SiteID, eventData.ClientIP, eventData.URL, eventData.Referrer,
		eventData.ScreenWidth, eventData.Browser, eventData.OS, eventData.Country, eventData.TrustScore,
		eventData.LCP, eventData.CLS, eventData.FID,
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
		days = 30 // Default to 30 days
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

func calculateChange(current, previous float64) float64 {
	if previous == 0 {
		if current > 0 {
			return 100.0 // Infinite growth, capped at 100%
		}
		return 0.0
	}
	return ((current - previous) / previous) * 100
}

func getCoreStats(ctx context.Context, siteID string, startDaysAgo, endDaysAgo int) (CoreStats, error) {
	var stats CoreStats
	
	// Total Views
	queryTotalViews := "SELECT count() FROM events WHERE SiteID = ? AND Timestamp BETWEEN now() - INTERVAL ? DAY AND now() - INTERVAL ? DAY"
	err := chConn.QueryRow(ctx, queryTotalViews, siteID, startDaysAgo, endDaysAgo).Scan(&stats.TotalViews)
	if err != nil && err != sql.ErrNoRows {
		return stats, err
	}

	// Unique Visitors
	queryUniqueVisitors := "SELECT uniq(ClientIP) FROM events WHERE SiteID = ? AND Timestamp BETWEEN now() - INTERVAL ? DAY AND now() - INTERVAL ? DAY"
	err = chConn.QueryRow(ctx, queryUniqueVisitors, siteID, startDaysAgo, endDaysAgo).Scan(&stats.UniqueVisitors)
	if err != nil && err != sql.ErrNoRows {
		return stats, err
	}

	// Bounce Rate
	queryBounceRate := `
		SELECT (countIf(pageviews = 1) / count()) * 100
		FROM (
			SELECT ClientIP, count() AS pageviews
			FROM events
			WHERE SiteID = ? AND Timestamp BETWEEN now() - INTERVAL ? DAY AND now() - INTERVAL ? DAY
			GROUP BY ClientIP
		)`
	err = chConn.QueryRow(ctx, queryBounceRate, siteID, startDaysAgo, endDaysAgo).Scan(&stats.BounceRate)
	if err != nil {
		stats.BounceRate = 0
	}

	// Average Visit Duration
	queryAvgVisitTime := `
		SELECT avg(duration)
		FROM (
			SELECT ClientIP, date_diff('second', min(Timestamp), max(Timestamp)) AS duration
			FROM events
			WHERE SiteID = ? AND Timestamp BETWEEN now() - INTERVAL ? DAY AND now() - INTERVAL ? DAY
			GROUP BY ClientIP
		)`
	err = chConn.QueryRow(ctx, queryAvgVisitTime, siteID, startDaysAgo, endDaysAgo).Scan(&stats.AvgVisitTime)
	if err != nil {
		stats.AvgVisitTime = 0
	}

	// Traffic Quality Score
	queryGoodTraffic := "SELECT count() FROM events WHERE SiteID = ? AND Timestamp BETWEEN now() - INTERVAL ? DAY AND now() - INTERVAL ? DAY AND TrustScore > 50"
	var goodTrafficCount uint64
	err = chConn.QueryRow(ctx, queryGoodTraffic, siteID, startDaysAgo, endDaysAgo).Scan(&goodTrafficCount)
	if err != nil || stats.TotalViews == 0 {
		stats.TrafficQualityScore = 0
	} else {
		stats.TrafficQualityScore = (float64(goodTrafficCount) / float64(stats.TotalViews)) * 100
	}

	// Web Vitals
	chConn.QueryRow(ctx, "SELECT avg(LCP) FROM events WHERE SiteID = ? AND Timestamp BETWEEN now() - INTERVAL ? DAY AND now() - INTERVAL ? DAY", siteID, startDaysAgo, endDaysAgo).Scan(&stats.AvgLCP)
	chConn.QueryRow(ctx, "SELECT avg(CLS) FROM events WHERE SiteID = ? AND Timestamp BETWEEN now() - INTERVAL ? DAY AND now() - INTERVAL ? DAY", siteID, startDaysAgo, endDaysAgo).Scan(&stats.AvgCLS)
	chConn.QueryRow(ctx, "SELECT avg(FID) FROM events WHERE SiteID = ? AND Timestamp BETWEEN now() - INTERVAL ? DAY AND now() - INTERVAL ? DAY", siteID, startDaysAgo, endDaysAgo).Scan(&stats.AvgFID)

	return stats, nil
}

func calculateStats(siteID string, days int) (Stats, error) {
	ctx := context.Background()
	var finalStats Stats

	// Get stats for the current period (e.g., last 30 days)
	currentStats, err := getCoreStats(ctx, siteID, days, 0)
	if err != nil {
		return finalStats, err
	}

	// Get stats for the previous period (e.g., 31-60 days ago)
	previousStats, err := getCoreStats(ctx, siteID, days*2, days)
	if err != nil {
		return finalStats, err
	}

	// Populate the final stats struct
	finalStats.TotalViews = currentStats.TotalViews
	finalStats.UniqueVisitors = currentStats.UniqueVisitors
	finalStats.BounceRate = currentStats.BounceRate
	d := time.Duration(currentStats.AvgVisitTime) * time.Second
	finalStats.AvgVisitTime = d.Round(time.Second).String()
	finalStats.TrafficQualityScore = currentStats.TrafficQualityScore
	finalStats.AvgLCP = currentStats.AvgLCP
	finalStats.AvgCLS = currentStats.AvgCLS
	finalStats.AvgFID = currentStats.AvgFID

	// Calculate percentage changes
	finalStats.TotalViewsChange = calculateChange(float64(currentStats.TotalViews), float64(previousStats.TotalViews))
	finalStats.UniqueVisitorsChange = calculateChange(float64(currentStats.UniqueVisitors), float64(previousStats.UniqueVisitors))
	finalStats.BounceRateChange = calculateChange(currentStats.BounceRate, previousStats.BounceRate)
	finalStats.AvgVisitTimeChange = calculateChange(currentStats.AvgVisitTime, previousStats.AvgVisitTime)
	finalStats.TrafficQualityScoreChange = calculateChange(currentStats.TrafficQualityScore, previousStats.TrafficQualityScore)
	finalStats.AvgLCPChange = calculateChange(currentStats.AvgLCP, previousStats.AvgLCP)
	finalStats.AvgCLSChange = calculateChange(currentStats.AvgCLS, previousStats.AvgCLS)
	finalStats.AvgFIDChange = calculateChange(currentStats.AvgFID, previousStats.AvgFID)

	// Top stats are still for the current period
	finalStats.TopPages, _ = queryTopStats(ctx, "URL", siteID, days)
	finalStats.TopReferrers, _ = queryTopStats(ctx, "Referrer", siteID, days)
	finalStats.TopBrowsers, _ = queryTopStats(ctx, "Browser", siteID, days)
	finalStats.TopOS, _ = queryTopStats(ctx, "OS", siteID, days)
	finalStats.TopCountries, _ = queryTopStats(ctx, "Country", siteID, days)

	return finalStats, nil
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

func nullFloat64(f *float64) sql.NullFloat64 {
	if f == nil {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{Float64: *f, Valid: true}
}

