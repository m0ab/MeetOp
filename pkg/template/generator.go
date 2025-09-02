package template

import (
	"fmt"
	"strings"
	"time"

	"github.com/m0ab/meetop/pkg/config"
)

type MeetupTemplate struct {
	Title       string
	Description string
	DateTime    string
	Venue       string
	Address     string
	EventURL    string
}

type SlackTemplate struct {
	Content string
}

type LinkedInTemplate struct {
	Content string
}

type Generator struct {
	groupURLName string
}

func NewGenerator(groupURLName string) *Generator {
	return &Generator{
		groupURLName: groupURLName,
	}
}

func (g *Generator) GenerateMeetupTemplate(cfg *config.Config) *MeetupTemplate {
	// Parse and format date/time
	dateTime, err := time.Parse("2006-01-02 15:04", cfg.EventDate+" "+cfg.EventTime)
	formattedDateTime := cfg.EventDate + " " + cfg.EventTime
	if err == nil {
		formattedDateTime = dateTime.Format("Monday, January 2, 2006 at 3:04 PM")
	}

	// Generate the expected Meetup URL (for copying into templates)
	eventURL := fmt.Sprintf("https://www.meetup.com/%s/events/[EVENT_ID]/", g.groupURLName)

	return &MeetupTemplate{
		Title:       cfg.EventTitle,
		Description: g.buildEventDescription(cfg),
		DateTime:    formattedDateTime,
		Venue:       cfg.Venue,
		Address:     cfg.VenueAddress,
		EventURL:    eventURL,
	}
}

func (g *Generator) GenerateSlackTemplate(cfg *config.Config, eventURL string) *SlackTemplate {
	dateTime, err := time.Parse("2006-01-02 15:04", cfg.EventDate+" "+cfg.EventTime)
	formattedDate := cfg.EventDate + " " + cfg.EventTime
	if err == nil {
		formattedDate = dateTime.Format("Monday, January 2, 2006 at 3:04 PM")
	}

	var content strings.Builder
	if cfg.EventType == config.EventTypeSocial {
		content.WriteString(fmt.Sprintf("🍻 *%s*\n\n", cfg.EventTitle))
		content.WriteString(fmt.Sprintf("📅 *When:* %s\n", formattedDate))
		content.WriteString(fmt.Sprintf("📍 *Where:* %s\n\n", cfg.Venue))
		content.WriteString("🤝 Join us for networking, great conversations, and community building!\n\n")
	} else {
		speakerText := "speaker"
		if cfg.NumSpeakers != 1 {
			speakerText = "speakers"
		}
		content.WriteString(fmt.Sprintf("🎉 *%s*\n\n", cfg.EventTitle))
		content.WriteString(fmt.Sprintf("📅 *When:* %s\n", formattedDate))
		content.WriteString(fmt.Sprintf("📍 *Where:* %s\n", cfg.Venue))
		content.WriteString(fmt.Sprintf("🎤 *Featuring:* %d amazing %s\n\n", cfg.NumSpeakers, speakerText))
		
		if cfg.Sponsor != "" {
			if cfg.SponsorURL != "" {
				content.WriteString(fmt.Sprintf("💝 Special thanks to our sponsor: <%s|%s>\n\n", cfg.SponsorURL, cfg.Sponsor))
			} else {
				content.WriteString(fmt.Sprintf("💝 Special thanks to our sponsor: %s\n\n", cfg.Sponsor))
			}
		}
	}

	content.WriteString(fmt.Sprintf("🎫 *Register here:* <%s|Join us!>\n\n", eventURL))
	content.WriteString("See you there! 👋")

	return &SlackTemplate{Content: content.String()}
}

func (g *Generator) GenerateLinkedInTemplate(cfg *config.Config, eventURL string) *LinkedInTemplate {
	dateTime, err := time.Parse("2006-01-02 15:04", cfg.EventDate+" "+cfg.EventTime)
	formattedDate := cfg.EventDate + " " + cfg.EventTime
	if err == nil {
		formattedDate = dateTime.Format("Monday, January 2, 2006 at 3:04 PM")
	}

	var content strings.Builder
	if cfg.EventType == config.EventTypeSocial {
		content.WriteString(fmt.Sprintf("🍻 Excited to invite everyone to our upcoming social meetup: %s\n\n", cfg.EventTitle))
		content.WriteString(fmt.Sprintf("📅 Date & Time: %s\n", formattedDate))
		content.WriteString(fmt.Sprintf("📍 Location: %s\n\n", cfg.Venue))
		content.WriteString("🤝 Perfect opportunity for networking, meaningful conversations, and connecting with like-minded professionals in our community!\n\n")
	} else {
		speakerText := "speaker"
		if cfg.NumSpeakers != 1 {
			speakerText = "speakers"
		}
		content.WriteString(fmt.Sprintf("🎉 Thrilled to announce our upcoming meetup: %s\n\n", cfg.EventTitle))
		content.WriteString(fmt.Sprintf("📅 Date & Time: %s\n", formattedDate))
		content.WriteString(fmt.Sprintf("📍 Location: %s\n", cfg.Venue))
		content.WriteString(fmt.Sprintf("🎤 We'll be featuring %d incredible %s who will share valuable insights and expertise!\n\n", cfg.NumSpeakers, speakerText))
		
		if cfg.Sponsor != "" {
			if cfg.SponsorURL != "" {
				content.WriteString(fmt.Sprintf("A big thank you to our amazing sponsor %s (%s) for making this event possible!\n\n", cfg.Sponsor, cfg.SponsorURL))
			} else {
				content.WriteString(fmt.Sprintf("A big thank you to our amazing sponsor %s for making this event possible!\n\n", cfg.Sponsor))
			}
		}
	}

	content.WriteString(fmt.Sprintf("Don't miss this opportunity to learn, network, and grow! Register now: %s\n\n", eventURL))
	
	if cfg.EventType == config.EventTypeSocial {
		content.WriteString("#meetup #social #community #networking #socializing #professional")
	} else {
		content.WriteString("#meetup #tech #community #networking #speakers #learning #professional")
	}

	return &LinkedInTemplate{Content: content.String()}
}

func (g *Generator) buildEventDescription(cfg *config.Config) string {
	var description strings.Builder
	
	// Add the main description
	description.WriteString(cfg.EventDescription)
	description.WriteString("\n\n")

	// Add event type specific content
	if cfg.EventType == config.EventTypeSpeaker {
		speakerText := "speaker"
		if cfg.NumSpeakers != 1 {
			speakerText = "speakers"
		}
		
		description.WriteString(fmt.Sprintf("🎤 This event will feature %d amazing %s sharing their expertise!\n\n", cfg.NumSpeakers, speakerText))
		
		// Add sponsor information for speaker events
		if cfg.Sponsor != "" {
			if cfg.SponsorURL != "" {
				description.WriteString(fmt.Sprintf("Special thanks to our sponsor: %s (%s)\n\n", cfg.Sponsor, cfg.SponsorURL))
			} else {
				description.WriteString(fmt.Sprintf("Special thanks to our sponsor: %s\n\n", cfg.Sponsor))
			}
		}
	} else {
		description.WriteString("🤝 Join us for networking, great conversations, and community building!\n\n")
	}

	description.WriteString("Looking forward to seeing you there! 🎉")
	
	return description.String()
}

func (g *Generator) PrintMeetupTemplate(template *MeetupTemplate) {
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║          MEETUP EVENT TEMPLATE           ║")
	fmt.Println("╚══════════════════════════════════════════╝")
	fmt.Println()
	fmt.Printf("📋 TITLE:\n%s\n\n", template.Title)
	fmt.Printf("📅 DATE & TIME:\n%s\n\n", template.DateTime)
	fmt.Printf("📍 VENUE:\n%s\n%s\n\n", template.Venue, template.Address)
	fmt.Printf("📝 DESCRIPTION:\n%s\n\n", template.Description)
	fmt.Printf("🔗 EVENT URL (after creation):\n%s\n\n", template.EventURL)
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║  Copy the above details to create your  ║")
	fmt.Println("║  event manually on meetup.com           ║")
	fmt.Println("╚══════════════════════════════════════════╝")
}

func (g *Generator) PrintSlackTemplate(template *SlackTemplate) {
	fmt.Println("\n╔══════════════════════════════════════════╗")
	fmt.Println("║            SLACK MESSAGE                 ║")
	fmt.Println("╚══════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println(template.Content)
	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║  Copy and paste the above to Slack      ║")
	fmt.Println("╚══════════════════════════════════════════╝")
}

func (g *Generator) PrintLinkedInTemplate(template *LinkedInTemplate) {
	fmt.Println("\n╔══════════════════════════════════════════╗")
	fmt.Println("║           LINKEDIN POST                  ║")
	fmt.Println("╚══════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println(template.Content)
	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║  Copy and paste the above to LinkedIn   ║")
	fmt.Println("╚══════════════════════════════════════════╝")
}