package main

import (
	"log"
	"strings"
	"time"

	ical "github.com/arran4/golang-ical"
	"github.com/teambition/rrule-go"
)

type EventType string

const (
	EventTypeUnknown  EventType = "unknown"
	EventTypePlanning EventType = "planning"
	EventTypeDaily    EventType = "daily"
	EventTypeGrooming EventType = "grooming"
	EventTypeReview   EventType = "review"
)

const DatetimeLayout = "20060102T150405"

type Event struct {
	Title     string
	StartTime time.Time
	EndTime   time.Time
	Type      EventType
	RR        *rrule.RRule
}

func ToEvent(ie *ical.VEvent) *Event {
	title := title(ie)
	et := toEventType(ie)
	if len(title) == 0 || isCanceledEvent(ie) || et == EventTypeUnknown {
		return nil
	}

	start, err := time.ParseInLocation(DatetimeLayout, ie.GetProperty(ical.ComponentPropertyDtStart).Value, time.Local)
	if err != nil {
		log.Printf("Error parsing start time: %s", err)
		return nil
	}

	end, err := time.ParseInLocation(DatetimeLayout, ie.GetProperty(ical.ComponentPropertyDtEnd).Value, time.Local)
	if err != nil {
		log.Printf("Error parsing end time: %s", err)
		return nil
	}

	rr := (*rrule.RRule)(nil)
	rruleProp := ie.GetProperty(ical.ComponentPropertyRrule)
	if rruleProp != nil {
		rr, err = rrule.StrToRRule(rruleProp.Value)
		if err != nil {
			log.Printf("Error parsing RRULE: %s", err)
			return nil
		}
		rr.DTStart(start)
	}

	return &Event{
		Title:     title,
		StartTime: start,
		EndTime:   end,
		Type:      et,
		RR:        rr,
	}
}

func toEventType(ie *ical.VEvent) EventType {
	titleProp := ie.GetProperty(ical.ComponentPropertySummary)
	if titleProp == nil {
		return EventTypeUnknown
	}

	title := strings.ToLower(titleProp.Value)

	switch {
	case strings.Contains(title, string(EventTypePlanning)):
		return EventTypePlanning
	case strings.Contains(title, string(EventTypeDaily)):
		return EventTypeDaily
	case strings.Contains(title, string(EventTypeGrooming)):
		return EventTypeGrooming
	case strings.Contains(title, string(EventTypeReview)):
		return EventTypeReview
	default:
		return EventTypeUnknown
	}
}

func isCanceledEvent(ie *ical.VEvent) bool {
	titleProp := ie.GetProperty(ical.ComponentPropertySummary)
	if titleProp == nil {
		return true
	}
	return strings.HasPrefix(titleProp.Value, "Canceled")
}

func title(ie *ical.VEvent) string {
	titleProp := ie.GetProperty(ical.ComponentPropertySummary)
	if titleProp == nil {
		return ""
	}
	return titleProp.Value
}
