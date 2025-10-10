package main

import (
	"context"
	"fmt"

	"github.com/sethvargo/go-envconfig"
)

var m map[EventType]string = nil

type Config struct {
	JiraEmail    string `env:"JIRA_EMAIL"`
	JiraAPIToken string `env:"JIRA_API_TOKEN"`
	JiraURL      string `env:"JIRA_URL"`
	ICalURL      string `env:"ICAL_URL"`
	GroomingID   string `env:"GROOMING_ID"`
	PlanningID   string `env:"PLANNING_ID"`
	ReviewID     string `env:"REVIEW_ID"`
	DailyID      string `env:"DAILY_ID"`
}

func loadConfig() (*Config, error) {
	ctx := context.Background()
	var c Config
	if err := envconfig.Process(ctx, &c); err != nil {
		return nil, fmt.Errorf("error processing env config: %w", err)
	}

	return &c, nil
}

func (c *Config) EventTypeToTicketID() map[EventType]string {
	if m == nil {
		m = map[EventType]string{
			EventTypeDaily:    c.DailyID,
			EventTypePlanning: c.PlanningID,
			EventTypeReview:   c.ReviewID,
			EventTypeGrooming: c.GroomingID,
		}
	}
	return m
}
