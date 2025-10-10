package main

import (
	"flag"
	"log"
	"time"

	ical "github.com/arran4/golang-ical"
	"github.com/joho/godotenv"
)

func main() {
	envPath := flag.String("env", "", "Path to .env file. If not set, will use environment variables")
	dateParam := flag.String("date", "", "Date to log works for (format: YYYY-MM-DD). If not set, will use current date")
	flag.Parse()

	// Load env file if provided
	if len(*envPath) != 0 {
		if err := godotenv.Load(*envPath); err != nil {
			log.Fatalf("error loading .env file: %v", err)
		}
	}

	// Load configuration
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %s", err)
	}

	// Determine date to log works for
	date := time.Now()
	if len(*dateParam) != 0 {
		pDate, err := time.Parse(time.DateOnly, *dateParam)
		if err != nil {
			log.Fatalf("Error parsing date: %s", err)
		}
		date = pDate
	}

	log.Printf("Log works for %s started...", date.Format(time.DateOnly))

	jiraClient, err := NewJiraClient(cfg)
	if err != nil {
		log.Fatalf("Error creating Jira client: %s", err)
	}

	// Download and parse the calendar
	cal, err := ical.ParseCalendarFromUrl(cfg.ICalURL)
	if err != nil {
		log.Fatalf("Error parsing calendar: %s", err)
	}

	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
	endOfDay := startOfDay.Add(24 * time.Hour)

	handledEvents := make(map[string]bool) // to avoid logging the same event multiple times

	for _, ie := range cal.Events() {
		event := ToEvent(ie, date)
		if event == nil { // invalid event
			continue
		}

		if (event.RR == nil && !isSameDate(event.StartTime, date)) ||
			(event.RR != nil && len(event.RR.Between(startOfDay, endOfDay, true)) == 0) { // event not occurring today
			continue
		}

		if handledEvents[event.Title] || jiraClient.HasLogged(event) { // already logged this event
			continue
		}

		if err := jiraClient.LogWork(event, date); err != nil {
			log.Printf("Error logging work for event %s: %s", event.Title, err)
			continue
		}

		log.Printf(
			"Logged work for event: %s | Ticket ID: %s | Duration: %vm",
			event.Title,
			cfg.EventTypeToTicketID()[event.Type],
			event.EndTime.Sub(event.StartTime).Minutes(),
		)
		handledEvents[event.Title] = true
	}

	log.Printf("Done with %v works logged", len(handledEvents))
}
