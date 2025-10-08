package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/andygrunwald/go-jira"
)

type JiraClient struct {
	cli    *jira.Client
	cfg    *Config
	logged map[string]map[string]bool // issueID -> map[eventTitle]bool
}

func NewJiraClient(cfg *Config) (*JiraClient, error) {
	tp := jira.BasicAuthTransport{
		Username: cfg.JiraEmail,
		Password: cfg.JiraAPIToken,
	}

	client, err := jira.NewClient(tp.Client(), cfg.JiraURL)
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}

	logged := map[string]map[string]bool{}
	for _, issueID := range cfg.EventTypeToTicketID() {
		wl, err := preFetchWorklogs(client, issueID, cfg.JiraEmail)
		if err != nil {
			return nil, fmt.Errorf("error pre-fetching work logs for %s: %w", issueID, err)
		}
		logged[issueID] = wl
	}

	return &JiraClient{cli: client, cfg: cfg, logged: logged}, nil
}

func (j *JiraClient) LogWork(event *Event, date time.Time) error {
	time := replaceTime(date, event.StartTime)
	_, _, err := j.cli.Issue.AddWorklogRecord(j.cfg.EventTypeToTicketID()[event.Type], &jira.WorklogRecord{
		Comment:   event.Title,
		TimeSpent: fmt.Sprintf("%vm", event.EndTime.Sub(event.StartTime).Minutes()),
		Started:   (*jira.Time)(&time),
	})

	return err
}

func (j *JiraClient) HasLogged(event *Event) bool {
	return j.logged[j.cfg.EventTypeToTicketID()[event.Type]][event.Title]
}

func preFetchWorklogs(cli *jira.Client, issueID string, email string) (map[string]bool, error) {
	wl, _, err := cli.Issue.GetWorklogs(issueID, func(r *http.Request) error { return nil })
	if err != nil {
		return nil, fmt.Errorf("error fetching worklogs: %w", err)
	}

	handled := make(map[string]bool)
	for _, worklog := range wl.Worklogs {
		if worklog.Author.EmailAddress != email {
			continue
		}
		handled[worklog.Comment] = true
	}

	return handled, nil
}
