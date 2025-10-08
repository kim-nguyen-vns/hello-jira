package main

import (
	"flag"
	"log"
	"time"

	ical "github.com/arran4/golang-ical"
)

func main() {
	envPath := flag.String("env", "", "Path to .env file (required)")
	date := flag.String("date", "", "Date to log works for (format: YYYY-MM-DD) (required)")
	flag.Parse()

	if len(*envPath) == 0 {
		log.Fatalln("Error: env path is required")
	}

	if len(*date) == 0 {
		log.Fatalln("Error: date is required")
	}

	pDate, err := time.Parse(time.DateOnly, *date)
	if err != nil {
		log.Fatalf("Error parsing date: %s", err)
	}

	log.Println("Log works started...")

	cfg, err := loadConfig(*envPath)
	if err != nil {
		log.Fatalf("Error loading config: %s", err)
	}

	jiraClient, err := NewJiraClient(cfg)
	if err != nil {
		log.Fatalf("Error creating Jira client: %s", err)
	}

	// Download and parse the calendar
	cal, err := ical.ParseCalendarFromUrl(cfg.ICalURL)
	if err != nil {
		log.Fatalf("Error parsing calendar: %s", err)
	}

	startOfDay := time.Date(pDate.Year(), pDate.Month(), pDate.Day(), 0, 0, 0, 0, time.Local)
	endOfDay := startOfDay.Add(24 * time.Hour)

	handledEvents := make(map[string]bool) // to avoid logging the same event multiple times

	for _, ie := range cal.Events() {
		event := ToEvent(ie)
		if event == nil { // invalid event
			continue
		}

		if (event.RR == nil && !isSameDate(event.StartTime, pDate)) ||
			(event.RR != nil && len(event.RR.Between(startOfDay, endOfDay, true)) == 0) { // event not occurring today
			continue
		}

		if handledEvents[event.Title] || jiraClient.HasLogged(event) { // already logged this event
			continue
		}

		if err := jiraClient.LogWork(event); err != nil {
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
