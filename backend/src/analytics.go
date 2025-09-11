package sentinel

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/oschwald/geoip2-golang"
	"github.com/ua-parser/uap-go/uaparser"
)

// --- EVENT TRACKING & STORAGE ---

type Event struct {
	SiteID      string `json:"siteId"`
	URL         string `json:"url"`
	Referrer    string `json:"referrer"`
	ScreenWidth int    `json:"screenWidth"`
}

type EventData struct {
	Timestamp   time.Time `json:"timestamp"`
	SiteID      string    `json:"siteId"`
	ClientIP    string    `json:"client_ip"`
	URL         string    `json:"url"`
	Referrer    string    `json:"referrer"`
	ScreenWidth int       `json:"screenWidth"`
	Browser     string    `json:"browser"`
	OS          string    `json:"os"`
	Country     string    `json:"country"`
}

type Store struct {
	mu         sync.Mutex
	fileLogger *log.Logger
}

var eventStore = NewStore("events.log")

func NewStore(logFilePath string) *Store {
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %s", err)
	}
	return &Store{
		fileLogger: log.New(file, "", 0),
	}
}

func (s *Store) AddEvent(event EventData) {
	s.mu.Lock()
	defer s.mu.Unlock()
	logEntry, _ := json.Marshal(event)
	s.fileLogger.Println(string(logEntry))
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
	TotalViews        int         `json:"totalViews"`
	UniqueVisitors    int         `json:"uniqueVisitors"`
	AvgVisitDuration  string      `json:"avgVisitDuration"`
	BounceRate        int         `json:"bounceRate"`
	TopPages          []CountStat `json:"topPages"`
	TopReferrers      []CountStat `json:"topReferrers"`
	TopBrowsers       []CountStat `json:"topBrowsers"`
	TopOS             []CountStat `json:"topOS"`
	TopCountries      []CountStat `json:"topCountries"`
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
		ScreenWidth: event.ScreenWidth,
		Browser:     browser,
		OS:          osFamily,
		Country:     country,
	}
	eventStore.AddEvent(eventData)

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

	stats := calculateStats(siteID, days)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func calculateStats(siteID string, days int) Stats {
	file, err := os.Open("events.log")
	if err != nil {
		return Stats{}
	}
	defer file.Close()

	timeCutoff := time.Now().UTC().AddDate(0, 0, -days)
	
	pageCounts := make(map[string]int)
	referrerCounts := make(map[string]int)
	browserCounts := make(map[string]int)
	osCounts := make(map[string]int)
	countryCounts := make(map[string]int)
	visitorSessions := make(map[string][]time.Time)
	totalViews := 0 // Initialize total views counter

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var event EventData
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			continue
		}
		
		if event.SiteID == siteID && event.Timestamp.After(timeCutoff) {
			totalViews++ // THIS IS THE CRITICAL FIX
			visitorSessions[event.ClientIP] = append(visitorSessions[event.ClientIP], event.Timestamp)
			pageCounts[event.URL]++
			
			if event.Referrer != "" {
				refURL, err1 := url.Parse(event.Referrer)
				pageURL, err2 := url.Parse(event.URL)
				if err1 == nil && err2 == nil && refURL.Host != pageURL.Host {
					referrerCounts[event.Referrer]++
				}
			}

			browserCounts[event.Browser]++
			osCounts[event.OS]++
			countryCounts[event.Country]++
		}
	}

	bounces := 0
	totalDuration := 0.0
	for _, timestamps := range visitorSessions {
		if len(timestamps) == 1 {
			bounces++
		}
		if len(timestamps) > 1 {
			sort.Slice(timestamps, func(i, j int) bool { return timestamps[i].Before(timestamps[j]) })
			duration := timestamps[len(timestamps)-1].Sub(timestamps[0]).Seconds()
			totalDuration += duration
		}
	}

	bounceRate := 0
	if len(visitorSessions) > 0 {
		bounceRate = (bounces * 100) / len(visitorSessions)
	}
	
	avgDuration := 0.0
	if (len(visitorSessions) - bounces) > 0 {
		avgDuration = totalDuration / float64(len(visitorSessions)-bounces)
	}

	return Stats{
		TotalViews:        totalViews, // Use the correct total views count
		UniqueVisitors:    len(visitorSessions),
		BounceRate:        bounceRate,
		AvgVisitDuration:  strconv.Itoa(int(avgDuration)) + "s",
		TopPages:          sortMap(pageCounts),
		TopReferrers:      sortMap(referrerCounts),
		TopBrowsers:       sortMap(browserCounts),
		TopOS:             sortMap(osCounts),
		TopCountries:      sortMap(countryCounts),
	}
}

func sortMap(counts map[string]int) []CountStat {
	stats := make([]CountStat, 0, len(counts))
	for value, count := range counts {
		stats = append(stats, CountStat{Value: value, Count: count})
	}
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Count > stats[j].Count
	})
	if len(stats) > 10 {
		return stats[:10]
	}
	return stats
}

