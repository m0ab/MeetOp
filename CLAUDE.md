# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

MeetOp is a Go-based automation tool for creating meetup events and sharing them across multiple platforms (Meetup.com, Slack, LinkedIn). The project runs primarily through GitHub Actions with manual triggers.

## Commands

### Build and Run
```bash
# Build the application
go build -o meetop ./cmd/meetop

# Run locally (requires environment variables)
./meetop

# Install dependencies
go mod tidy

# Download new dependencies
go mod download
```

### Development
```bash
# Run with auto-restart during development
go run ./cmd/meetop

# Format code
go fmt ./...

# Run tests (when implemented)
go test ./...

# Check for race conditions
go run -race ./cmd/meetop
```

## Architecture

### Project Structure
- `cmd/meetop/`: Main application entry point
- `pkg/config/`: Configuration management with environment variables
- `pkg/meetup/`: Meetup.com API client implementation
- `pkg/slack/`: Slack API client with block-based messaging
- `pkg/linkedin/`: LinkedIn API client for social posts
- `.github/workflows/`: GitHub Actions workflow definitions

### Key Components
- **Configuration System**: Environment-based config with validation in `pkg/config/`
- **API Clients**: Each platform has its own package with specialized methods
- **Error Handling**: Comprehensive logging to `meetop.log` file
- **GitHub Actions**: Uses `workflow_dispatch` for manual triggering with custom inputs

### API Integration Details
- **Meetup.com**: Uses OAuth2 Bearer token authentication, creates events with venue details
- **Slack**: Supports both webhooks (preferred) and bot tokens, sends structured block messages with customized content based on event type
- **LinkedIn**: Uses personal access token, posts to user's feed with UGC API, includes customized messaging and hashtags based on event type

### Event Types
The application supports two types of meetup events:

- **Speaker Events** (`EVENT_TYPE=speaker`): Traditional meetups with speakers and presentations
  - Requires speaker count (`NUM_SPEAKERS` > 0)
  - Supports sponsor information (`SPONSOR` and `SPONSOR_URL`)
  - Messages emphasize speakers and technical content
  - Uses hashtags: #meetup #tech #community #networking #speakers

- **Social Events** (`EVENT_TYPE=social`): Networking-focused meetups without formal presentations
  - Speaker count defaults to 0 and is ignored
  - Sponsor information is not displayed (social events are sponsor-free)
  - Messages emphasize networking and community building
  - Uses hashtags: #meetup #social #community #networking #socializing

## Environment Variables

All configuration is handled through environment variables. See `.env.example` for the complete list. Critical variables include:
- `EVENT_TYPE` for event type (`speaker` or `social`, defaults to `speaker`)
- `MEETUP_API_KEY` and `MEETUP_GROUP_URLNAME` for event creation
- Slack integration (choose one):
  - `SLACK_WEBHOOK_URL` for webhook posting (preferred)
  - `SLACK_BOT_TOKEN` and `SLACK_CHANNEL` for bot token posting
- `LINKEDIN_ACCESS_TOKEN` and `LINKEDIN_PERSON_URN` for LinkedIn sharing
- `SPONSOR` and `SPONSOR_URL` for sponsor attribution (speaker events only)

## GitHub Actions Workflow

The main workflow (`.github/workflows/create-event.yml`) accepts these inputs:
- `event_type`: Choice between `speaker` and `social` event types
- Event details: title, description, date, time, venue, address
- Speaker count (ignored for social events) and sponsor information (speaker events only)
- Boolean toggles for Slack and LinkedIn sharing

## Error Handling

- Each API integration has independent error handling
- Failures in one platform don't prevent others from running
- All operations are logged with timestamps
- GitHub Actions uploads logs as artifacts for debugging

## Development Notes

- Uses Go modules for dependency management (`go mod tidy` handles dependencies correctly)
- Follows standard Go project layout conventions
- All secrets must be stored as GitHub repository secrets
- Local development requires `.env` file (copy from `.env.example`)
- Build system properly configured with correct module dependencies