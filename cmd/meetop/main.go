package main

import (
	"fmt"
	"log"
	"os"

	"github.com/m0ab/meetop/pkg/config"
	"github.com/m0ab/meetop/pkg/template"
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

	// Generate templates
	logger.Println("Generating event templates...")
	templateGenerator := template.NewGenerator(cfg.MeetupGroupURLName)

	// Always show Meetup template first
	meetupTemplate := templateGenerator.GenerateMeetupTemplate(cfg)
	templateGenerator.PrintMeetupTemplate(meetupTemplate)

	// Get event URL from environment or prompt for input
	eventURL := os.Getenv("EVENT_URL")
	if eventURL == "" {
		fmt.Print("\nAfter creating the event on Meetup.com, enter the event URL (or press Enter to skip social templates): ")
		fmt.Scanln(&eventURL)
	}

	if eventURL == "" {
		logger.Println("No event URL provided, social media templates skipped")
		fmt.Println("⏭️  Meetup template generated. Social media templates skipped.")
		fmt.Println("\nTo generate social templates later, run with EVENT_URL environment variable set.")
		return
	}

	logger.Printf("Generating social media templates with URL: %s", eventURL)
	
	// Generate platform-specific templates
	if cfg.ShareSlack {
		slackTemplate := templateGenerator.GenerateSlackTemplate(cfg, eventURL)
		templateGenerator.PrintSlackTemplate(slackTemplate)
	}
	
	if cfg.ShareLinkedIn {
		linkedinTemplate := templateGenerator.GenerateLinkedInTemplate(cfg, eventURL)
		templateGenerator.PrintLinkedInTemplate(linkedinTemplate)
	}

	fmt.Printf("\n✅ All templates generated! Copy and paste to your platforms.\n")

	logger.Println("Template generation completed successfully")
	fmt.Printf("\n🎉 Templates ready! Create your Meetup event manually, then copy/paste the social media content.\n")
}
