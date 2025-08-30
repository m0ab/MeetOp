package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	token      string
	httpClient *http.Client
}

type Message struct {
	Channel     string       `json:"channel"`
	Text        string       `json:"text,omitempty"`
	Blocks      []Block      `json:"blocks,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

type Block struct {
	Type string      `json:"type"`
	Text *TextObject `json:"text,omitempty"`
}

type TextObject struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Attachment struct {
	Color  string  `json:"color,omitempty"`
	Blocks []Block `json:"blocks,omitempty"`
}

type SlackResponse struct {
	OK    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
	TS    string `json:"ts,omitempty"`
}

func NewClient(token string) *Client {
	return &Client{
		token:      token,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *Client) PostMessage(channel, eventTitle, eventURL, venue, eventDate, eventTime string, numSpeakers int, sponsor, sponsorURL string) error {
	message := c.buildEventMessage(channel, eventTitle, eventURL, venue, eventDate, eventTime, numSpeakers, sponsor, sponsorURL)

	reqBody, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	req, err := http.NewRequest("POST", "https://slack.com/api/chat.postMessage", bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var slackResp SlackResponse
	if err := json.Unmarshal(body, &slackResp); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !slackResp.OK {
		return fmt.Errorf("slack API error: %s", slackResp.Error)
	}

	return nil
}

func (c *Client) PostWebhookMessage(webhookURL, eventTitle, eventURL, venue, eventDate, eventTime string, numSpeakers int, sponsor, sponsorURL string) error {
	message := c.buildEventMessage("", eventTitle, eventURL, venue, eventDate, eventTime, numSpeakers, sponsor, sponsorURL)
	
	// Remove channel from message for webhook (webhooks don't use channel in payload)
	message.Channel = ""

	reqBody, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook message: %w", err)
	}

	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create webhook request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("slack webhook returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (c *Client) buildEventMessage(channel, eventTitle, eventURL, venue, eventDate, eventTime string, numSpeakers int, sponsor, sponsorURL string) Message {
	dateTime, _ := time.Parse("2006-01-02 15:04", eventDate+" "+eventTime)
	formattedDate := dateTime.Format("Monday, January 2, 2006 at 3:04 PM")
	
	speakerText := "speaker"
	if numSpeakers != 1 {
		speakerText = "speakers"
	}

	mainText := fmt.Sprintf("🎉 *New Meetup Event: %s*", eventTitle)
	
	blocks := []Block{
		{
			Type: "section",
			Text: &TextObject{
				Type: "mrkdwn",
				Text: mainText,
			},
		},
		{
			Type: "section",
			Text: &TextObject{
				Type: "mrkdwn",
				Text: fmt.Sprintf("📅 *When:* %s\n📍 *Where:* %s\n🎤 *Speakers:* %d %s", 
					formattedDate, venue, numSpeakers, speakerText),
			},
		},
	}

	if sponsor != "" {
		sponsorText := fmt.Sprintf("🏢 *Sponsored by:* %s", sponsor)
		if sponsorURL != "" {
			sponsorText = fmt.Sprintf("🏢 *Sponsored by:* <%s|%s>", sponsorURL, sponsor)
		}
		blocks = append(blocks, Block{
			Type: "section",
			Text: &TextObject{
				Type: "mrkdwn",
				Text: sponsorText,
			},
		})
	}

	blocks = append(blocks, Block{
		Type: "section",
		Text: &TextObject{
			Type: "mrkdwn",
			Text: fmt.Sprintf("🔗 *Register now:* <%s|View Event>", eventURL),
		},
	})

	return Message{
		Channel: channel,
		Blocks:  blocks,
	}
}