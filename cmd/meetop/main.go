package main

import (
	"fmt"
	"log"
	"os"

	"github.com/m0ab/meetop/pkg/config"
	"github.com/m0ab/meetop/pkg/linkedin"
	"github.com/m0ab/meetop/pkg/meetup"
	"github.com/m0ab/meetop/pkg/slack"
)

func main() {
	// Set up logging
	logFile, err := os.OpenFile("meetop.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		os.Exit(1)
	}
	defer logFile.Close()

	logger := log.New(logFile, "", log.LstdFlags|log.Lshortfile)
	logger.Println("Starting MeetOp automation...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Printf("Failed to load configuration: %v", err)
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		logger.Printf("Configuration validation failed: %v", err)
		fmt.Printf("Configuration validation failed: %v\n", err)
		os.Exit(1)
	}

	logger.Printf("Configuration loaded successfully for event: %s", cfg.EventTitle)

	// Create Meetup event
	logger.Println("Creating Meetup event...")
	meetupClient := meetup.NewClient(cfg.MeetupAPIKey, cfg.MeetupGroupURLName)
	
	event, err := meetupClient.CreateEvent(
		cfg.EventTitle,
		cfg.EventDescription,
		cfg.EventDate,
		cfg.EventTime,
		cfg.Venue,
		cfg.VenueAddress,
	)
	if err != nil {
		logger.Printf("Failed to create Meetup event: %v", err)
		fmt.Printf("Failed to create Meetup event: %v\n", err)
		os.Exit(1)
	}

	logger.Printf("Successfully created Meetup event with ID: %s, URL: %s", event.ID, event.EventURL)
	fmt.Printf("✅ Meetup event created successfully!\n")
	fmt.Printf("Event ID: %s\n", event.ID)
	fmt.Printf("Event URL: %s\n", event.EventURL)

	// Share to Slack if enabled
	if cfg.ShareSlack {
		logger.Println("Sharing event to Slack...")
		slackClient := slack.NewClient(cfg.SlackBotToken)
		
		var err error
		
		// Use webhook if configured, otherwise use bot token
		if cfg.SlackWebhookURL != "" {
			logger.Println("Using Slack webhook for message posting")
			err = slackClient.PostWebhookMessage(
				cfg.SlackWebhookURL,
				cfg.EventTitle,
				event.EventURL,
				cfg.Venue,
				cfg.EventDate,
				cfg.EventTime,
				cfg.NumSpeakers,
				cfg.Sponsor,
				cfg.SponsorURL,
			)
		} else if cfg.SlackBotToken != "" && cfg.SlackChannel != "" {
			logger.Println("Using Slack bot token for message posting")
			err = slackClient.PostMessage(
				cfg.SlackChannel,
				cfg.EventTitle,
				event.EventURL,
				cfg.Venue,
				cfg.EventDate,
				cfg.EventTime,
				cfg.NumSpeakers,
				cfg.Sponsor,
				cfg.SponsorURL,
			)
		} else {
			logger.Println("Skipping Slack sharing (missing webhook URL or bot token/channel)")
			fmt.Printf("⏭️  Slack sharing skipped (missing configuration)\n")
		}
		
		if err != nil {
			logger.Printf("Failed to share to Slack: %v", err)
			fmt.Printf("❌ Failed to share to Slack: %v\n", err)
		} else if cfg.SlackWebhookURL != "" || (cfg.SlackBotToken != "" && cfg.SlackChannel != "") {
			logger.Println("Successfully shared event to Slack")
			fmt.Printf("✅ Event shared to Slack successfully!\n")
		}
	} else {
		logger.Println("Skipping Slack sharing (disabled)")
		fmt.Printf("⏭️  Slack sharing skipped\n")
	}

	// Share to LinkedIn if enabled
	if cfg.ShareLinkedIn && cfg.LinkedInAccessToken != "" && cfg.LinkedInPersonURN != "" {
		logger.Println("Sharing event to LinkedIn...")
		linkedinClient := linkedin.NewClient(cfg.LinkedInAccessToken, cfg.LinkedInPersonURN)
		
		err := linkedinClient.ShareEvent(
			cfg.EventTitle,
			event.EventURL,
			cfg.Venue,
			cfg.EventDate,
			cfg.EventTime,
			cfg.NumSpeakers,
			cfg.Sponsor,
			cfg.SponsorURL,
		)
		if err != nil {
			logger.Printf("Failed to share to LinkedIn: %v", err)
			fmt.Printf("❌ Failed to share to LinkedIn: %v\n", err)
		} else {
			logger.Println("Successfully shared event to LinkedIn")
			fmt.Printf("✅ Event shared to LinkedIn successfully!\n")
		}
	} else {
		logger.Println("Skipping LinkedIn sharing (disabled or missing configuration)")
		fmt.Printf("⏭️  LinkedIn sharing skipped\n")
	}

	logger.Println("MeetOp automation completed successfully")
	fmt.Printf("\n🎉 All done! Your meetup event has been created and shared according to your settings.\n")
}