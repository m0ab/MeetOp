package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type EventType string

const (
	EventTypeSocial  EventType = "social"
	EventTypeSpeaker EventType = "speaker"
)

type Config struct {
	// Meetup API
	MeetupAPIKey       string
	MeetupGroupURLName string

	// Slack API
	SlackBotToken   string
	SlackChannel    string
	SlackWebhookURL string

	// LinkedIn API
	LinkedInAccessToken string
	LinkedInPersonURN   string

	// Event details
	EventType        EventType
	EventTitle       string
	EventDescription string
	EventDate        string
	EventTime        string
	Venue            string
	VenueAddress     string
	NumSpeakers      int
	Sponsor          string
	SponsorURL       string

	// Sharing options
	ShareSlack    bool
	ShareLinkedIn bool
}

func Load() (*Config, error) {
	// Load .env file if it exists (for local development)
	_ = godotenv.Load()

	// Parse event type
	eventTypeStr := strings.ToLower(getEnv("EVENT_TYPE", "speaker"))
	var eventType EventType
	switch eventTypeStr {
	case "social":
		eventType = EventTypeSocial
	case "speaker":
		eventType = EventTypeSpeaker
	default:
		eventType = EventTypeSpeaker // default to speaker events
	}

	// For social events, default speakers to 0, for speaker events default to 1
	defaultSpeakers := "1"
	if eventType == EventTypeSocial {
		defaultSpeakers = "0"
	}
	numSpeakers, err := strconv.Atoi(getEnv("NUM_SPEAKERS", defaultSpeakers))
	if err != nil {
		numSpeakers = 1
		if eventType == EventTypeSocial {
			numSpeakers = 0
		}
	}

	shareSlack, _ := strconv.ParseBool(getEnv("SHARE_SLACK", "true"))
	shareLinkedIn, _ := strconv.ParseBool(getEnv("SHARE_LINKEDIN", "true"))

	return &Config{
		MeetupAPIKey:        os.Getenv("MEETUP_API_KEY"),
		MeetupGroupURLName:  os.Getenv("MEETUP_GROUP_URLNAME"),
		SlackBotToken:       os.Getenv("SLACK_BOT_TOKEN"),
		SlackChannel:        os.Getenv("SLACK_CHANNEL"),
		SlackWebhookURL:     os.Getenv("SLACK_WEBHOOK_URL"),
		LinkedInAccessToken: os.Getenv("LINKEDIN_ACCESS_TOKEN"),
		LinkedInPersonURN:   os.Getenv("LINKEDIN_PERSON_URN"),
		EventType:           eventType,
		EventTitle:          os.Getenv("EVENT_TITLE"),
		EventDescription:    os.Getenv("EVENT_DESCRIPTION"),
		EventDate:           os.Getenv("EVENT_DATE"),
		EventTime:           os.Getenv("EVENT_TIME"),
		Venue:               os.Getenv("VENUE"),
		VenueAddress:        os.Getenv("VENUE_ADDRESS"),
		NumSpeakers:         numSpeakers,
		Sponsor:             os.Getenv("SPONSOR"),
		SponsorURL:          os.Getenv("SPONSOR_URL"),
		ShareSlack:          shareSlack,
		ShareLinkedIn:       shareLinkedIn,
	}, nil
}

func (c *Config) Validate() error {
	if c.MeetupAPIKey == "" {
		return ErrMissingMeetupAPIKey
	}
	if c.MeetupGroupURLName == "" {
		return ErrMissingMeetupGroupURLName
	}
	if c.EventTitle == "" {
		return ErrMissingEventTitle
	}
	if c.EventDescription == "" {
		return ErrMissingEventDescription
	}
	if c.EventDate == "" {
		return ErrMissingEventDate
	}
	if c.EventTime == "" {
		return ErrMissingEventTime
	}
	if c.Venue == "" {
		return ErrMissingVenue
	}
	if c.VenueAddress == "" {
		return ErrMissingVenueAddress
	}

	// Validate date format
	if _, err := time.Parse("2006-01-02", c.EventDate); err != nil {
		return ErrInvalidDateFormat
	}

	// Validate time format
	if _, err := time.Parse("15:04", c.EventTime); err != nil {
		return ErrInvalidTimeFormat
	}

	// Validate event type specific requirements
	if c.EventType == EventTypeSpeaker {
		if c.NumSpeakers <= 0 {
			return ErrInvalidSpeakerCount
		}
	} else if c.EventType == EventTypeSocial {
		// For social events, speakers and sponsors are optional
		// No additional validation needed
	}

	return nil
}

// IsSocialEvent returns true if this is a social event
func (c *Config) IsSocialEvent() bool {
	return c.EventType == EventTypeSocial
}

// IsSpeakerEvent returns true if this is a speaker event
func (c *Config) IsSpeakerEvent() bool {
	return c.EventType == EventTypeSpeaker
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
