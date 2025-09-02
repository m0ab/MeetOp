# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

MeetOp is a Go-based template generator for managing meetup events. It generates formatted templates for manual Meetup.com event creation and produces ready-to-use copy-paste content for social platforms (Slack, LinkedIn) with platform-specific messaging. The project runs through GitHub Actions or locally.

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

### Template Generation
- **Meetup.com**: Generates formatted templates for manual event creation (no API integration due to Pro plan requirements)
- **Slack**: Produces formatted messages with Slack markup, optimized for team communication
- **LinkedIn**: Creates professional posts with appropriate hashtags and LinkedIn-style messaging

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
- `MEETUP_GROUP_URLNAME` for template generation (used to create expected event URLs)
- `EVENT_URL` for social media template generation (optional, can be provided interactively)
- `SPONSOR` and `SPONSOR_URL` for sponsor attribution (speaker events only)
- `SHARE_SLACK` and `SHARE_LINKEDIN` to control which platform templates to generate

## GitHub Actions Workflow

The main workflow (`.github/workflows/create-event.yml`) accepts these inputs:
- `event_type`: Choice between `speaker` and `social` event types
- Event details: title, description, date, time, venue, address
- Speaker count (ignored for social events) and sponsor information (speaker events only)
- Boolean toggles for Slack and LinkedIn sharing

Note: The workflow now generates templates for manual Meetup event creation rather than creating events automatically.

## Workflow Process

1. **Meetup Template**: Generates formatted event details for manual Meetup.com creation
2. **Manual Event Creation**: User creates the event on Meetup.com using the generated template
3. **Social Template Generation**: Creates platform-specific copy-paste content for Slack and LinkedIn
4. **Manual Sharing**: User copies and pastes the generated content to their social platforms
5. **Platform Optimization**: Each platform gets tailored messaging and formatting for maximum engagement

## Development Notes

- Uses Go modules for dependency management (`go mod tidy` handles dependencies correctly)
- Follows standard Go project layout conventions
- All secrets must be stored as GitHub repository secrets
- Local development requires `.env` file (copy from `.env.example`)
- Build system properly configured with correct module dependencies