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
	// Meetup Group (for template generation)
	MeetupGroupURLName string

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

	// Parse combined datetime and venue info
	eventDate, eventTime := parseEventDateTime()
	venue, venueAddress := parseVenueInfo()

	return &Config{
		MeetupGroupURLName: os.Getenv("MEETUP_GROUP_URLNAME"),
		EventType:           eventType,
		EventTitle:          os.Getenv("EVENT_TITLE"),
		EventDescription:    os.Getenv("EVENT_DESCRIPTION"),
		EventDate:           eventDate,
		EventTime:           eventTime,
		Venue:               venue,
		VenueAddress:        venueAddress,
		NumSpeakers:         numSpeakers,
		Sponsor:             os.Getenv("SPONSOR"),
		SponsorURL:          os.Getenv("SPONSOR_URL"),
		ShareSlack:          shareSlack,
		ShareLinkedIn:       shareLinkedIn,
	}, nil
}

func (c *Config) Validate() error {
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

// parseEventDateTime parses EVENT_DATETIME or falls back to EVENT_DATE and EVENT_TIME
func parseEventDateTime() (date, time string) {
	// Try combined datetime first
	if eventDateTime := os.Getenv("EVENT_DATETIME"); eventDateTime != "" {
		parts := strings.Split(eventDateTime, " ")
		if len(parts) == 2 {
			return parts[0], parts[1]
		}
	}
	
	// Fallback to separate date and time
	return os.Getenv("EVENT_DATE"), os.Getenv("EVENT_TIME")
}

// parseVenueInfo parses VENUE_INFO or falls back to VENUE and VENUE_ADDRESS
func parseVenueInfo() (venue, address string) {
	// Try combined venue info first
	if venueInfo := os.Getenv("VENUE_INFO"); venueInfo != "" {
		parts := strings.Split(venueInfo, "|")
		if len(parts) == 2 {
			return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		}
		// If no separator, treat entire string as venue
		return strings.TrimSpace(venueInfo), ""
	}
	
	// Fallback to separate venue and address
	return os.Getenv("VENUE"), os.Getenv("VENUE_ADDRESS")
}
