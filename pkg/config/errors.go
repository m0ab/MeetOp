package config

import "errors"

var (
	ErrMissingMeetupAPIKey       = errors.New("MEETUP_API_KEY is required")
	ErrMissingMeetupGroupURLName = errors.New("MEETUP_GROUP_URLNAME is required")
	ErrMissingEventTitle         = errors.New("EVENT_TITLE is required")
	ErrMissingEventDescription   = errors.New("EVENT_DESCRIPTION is required")
	ErrMissingEventDate          = errors.New("EVENT_DATE is required")
	ErrMissingEventTime          = errors.New("EVENT_TIME is required")
	ErrMissingVenue              = errors.New("VENUE is required")
	ErrMissingVenueAddress       = errors.New("VENUE_ADDRESS is required")
	ErrInvalidDateFormat         = errors.New("EVENT_DATE must be in YYYY-MM-DD format")
	ErrInvalidTimeFormat         = errors.New("EVENT_TIME must be in HH:MM format")
)