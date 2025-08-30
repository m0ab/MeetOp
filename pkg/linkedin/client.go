package linkedin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	accessToken string
	personURN   string
	httpClient  *http.Client
}

type ShareRequest struct {
	Author      string      `json:"author"`
	Commentary  string      `json:"commentary"`
	Visibility  Visibility  `json:"visibility"`
	Distribution Distribution `json:"distribution"`
}

type Visibility struct {
	MemberNetworkVisibility string `json:"com.linkedin.ugc.MemberNetworkVisibility"`
}

type Distribution struct {
	FeedDistribution string `json:"com.linkedin.ugc.FeedDistribution"`
}

type ShareResponse struct {
	ID       string `json:"id"`
	Activity string `json:"activity"`
}

func NewClient(accessToken, personURN string) *Client {
	return &Client{
		accessToken: accessToken,
		personURN:   personURN,
		httpClient:  &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *Client) ShareEvent(eventTitle, eventURL, venue, eventDate, eventTime string, numSpeakers int, sponsor, sponsorURL string) error {
	commentary := c.buildEventPost(eventTitle, eventURL, venue, eventDate, eventTime, numSpeakers, sponsor, sponsorURL)

	shareReq := ShareRequest{
		Author:     c.personURN,
		Commentary: commentary,
		Visibility: Visibility{
			MemberNetworkVisibility: "PUBLIC",
		},
		Distribution: Distribution{
			FeedDistribution: "MAIN_FEED",
		},
	}

	reqBody, err := json.Marshal(shareReq)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.linkedin.com/v2/ugcPosts", bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Restli-Protocol-Version", "2.0.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("linkedin API returned status %d: %s", resp.StatusCode, string(body))
	}

	var shareResp ShareResponse
	if err := json.Unmarshal(body, &shareResp); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}

func (c *Client) buildEventPost(eventTitle, eventURL, venue, eventDate, eventTime string, numSpeakers int, sponsor, sponsorURL string) string {
	dateTime, _ := time.Parse("2006-01-02 15:04", eventDate+" "+eventTime)
	formattedDate := dateTime.Format("Monday, January 2, 2006 at 3:04 PM")

	speakerText := "speaker"
	if numSpeakers != 1 {
		speakerText = "speakers"
	}

	post := fmt.Sprintf("🎉 Excited to announce our upcoming meetup: %s\n\n", eventTitle)
	post += fmt.Sprintf("📅 Date & Time: %s\n", formattedDate)
	post += fmt.Sprintf("📍 Location: %s\n", venue)
	post += fmt.Sprintf("🎤 Featuring %d amazing %s\n\n", numSpeakers, speakerText)

	if sponsor != "" {
		if sponsorURL != "" {
			post += fmt.Sprintf("Special thanks to our sponsor: %s (%s)\n\n", sponsor, sponsorURL)
		} else {
			post += fmt.Sprintf("Special thanks to our sponsor: %s\n\n", sponsor)
		}
	}

	post += fmt.Sprintf("Don't miss out! Register now: %s\n\n", eventURL)
	post += "#meetup #tech #community #networking"

	return post
}