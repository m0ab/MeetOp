package meetup

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	apiKey       string
	groupURLName string
	httpClient   *http.Client
}

type Event struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	EventURL    string `json:"event_url"`
	DateTime    int64  `json:"time"`
	Venue       *Venue `json:"venue,omitempty"`
}

type Venue struct {
	Name      string  `json:"name"`
	Address1  string  `json:"address_1"`
	City      string  `json:"city"`
	State     string  `json:"state"`
	Zip       string  `json:"zip"`
	Country   string  `json:"country"`
	Latitude  float64 `json:"lat,omitempty"`
	Longitude float64 `json:"lon,omitempty"`
}

type CreateEventRequest struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Time         int64  `json:"time"`
	Duration     int    `json:"duration"` // in milliseconds
	VenueName    string `json:"venue_name,omitempty"`
	VenueAddress string `json:"venue_address,omitempty"`
	HowToFindUs  string `json:"how_to_find_us,omitempty"`
	GuestLimit   int    `json:"guest_limit,omitempty"`
}

type EventResponse struct {
	ID       string `json:"id"`
	Link     string `json:"link"`
	Name     string `json:"name"`
	Status   string `json:"status"`
	Time     int64  `json:"time"`
	EventURL string `json:"event_url"`
}

func NewClient(apiKey, groupURLName string) *Client {
	return &Client{
		apiKey:       apiKey,
		groupURLName: groupURLName,
		httpClient:   &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *Client) CreateEvent(title, description, eventDate, eventTime, venue, venueAddress string) (*EventResponse, error) {
	// Parse date and time
	dateTime, err := time.Parse("2006-01-02 15:04", eventDate+" "+eventTime)
	if err != nil {
		return nil, fmt.Errorf("failed to parse date/time: %w", err)
	}

	// Convert to Unix timestamp in milliseconds
	timestamp := dateTime.Unix() * 1000

	req := CreateEventRequest{
		Name:         title,
		Description:  description,
		Time:         timestamp,
		Duration:     7200000, // 2 hours in milliseconds
		VenueName:    venue,
		VenueAddress: venueAddress,
		HowToFindUs:  fmt.Sprintf("Event will be held at %s, %s", venue, venueAddress),
		GuestLimit:   100, // Default guest limit
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("https://api.meetup.com/%s/events", c.groupURLName)
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("meetup API returned status %d: %s", resp.StatusCode, string(body))
	}

	var eventResp EventResponse
	if err := json.Unmarshal(body, &eventResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Construct event URL if not provided
	if eventResp.EventURL == "" && eventResp.ID != "" {
		eventResp.EventURL = fmt.Sprintf("https://www.meetup.com/%s/events/%s/", c.groupURLName, eventResp.ID)
	}

	return &eventResp, nil
}

func (c *Client) GetEvent(eventID string) (*Event, error) {
	url := fmt.Sprintf("https://api.meetup.com/%s/events/%s", c.groupURLName, eventID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("meetup API returned status %d: %s", resp.StatusCode, string(body))
	}

	var event Event
	if err := json.NewDecoder(resp.Body).Decode(&event); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &event, nil
}
