# MeetOp

Automate meetup event creation and social media communication across Meetup.com, Slack, and LinkedIn.

## Overview

MeetOp streamlines the process of creating meetup events by:
1. Creating events on Meetup.com via their API
2. Automatically sharing event announcements to Slack
3. Publishing event posts to LinkedIn
4. Running everything through GitHub Actions with manual triggers

## Features

- **Automated Event Creation**: Create Meetup.com events with venue details, date/time, and speaker information
- **Multi-platform Sharing**: Share to Slack and LinkedIn simultaneously
- **Smart Sponsor Attribution**: LinkedIn gets full sponsor URLs, Slack gets hyperlinked sponsor names
- **Flexible Slack Integration**: Support for both webhooks and bot tokens
- **GitHub Actions Integration**: Run via workflow_dispatch with custom inputs
- **Secure Credential Management**: All API tokens stored as GitHub Secrets
- **Comprehensive Logging**: Full audit trail of all operations

## GitHub Actions Workflow

The project uses GitHub Actions with `workflow_dispatch` to allow manual triggering with the following inputs:

- **Event Title**: The name of your meetup event
- **Event Description**: Detailed description of the event
- **Event Date**: Date in YYYY-MM-DD format
- **Event Time**: Time in HH:MM (24-hour) format
- **Venue**: Name of the venue
- **Venue Address**: Full address of the venue
- **Number of Speakers**: How many speakers will present
- **Sponsor**: Optional sponsor name
- **Sponsor URL**: Optional sponsor website URL for proper attribution
- **Share to Slack**: Toggle Slack sharing on/off
- **Share to LinkedIn**: Toggle LinkedIn sharing on/off

## Required GitHub Secrets

Configure these secrets in your GitHub repository settings:

### Meetup.com API
- `MEETUP_API_KEY`: Your Meetup.com API key
- `MEETUP_GROUP_URLNAME`: Your meetup group's URL name

### Slack API (Choose one method)
- `SLACK_BOT_TOKEN`: Your Slack bot token (xoxb-...) + `SLACK_CHANNEL`: Target channel
- `SLACK_WEBHOOK_URL`: Incoming webhook URL (preferred method)

### LinkedIn API
- `LINKEDIN_ACCESS_TOKEN`: Your LinkedIn API access token
- `LINKEDIN_PERSON_URN`: Your LinkedIn person URN

## API Setup Instructions

### Meetup.com API
1. Go to [Meetup.com API](https://www.meetup.com/api/)
2. Create an API key for your account
3. Note your group's URL name from your group page

### Slack API (Two Options)

**Option 1: Incoming Webhook (Recommended)**
1. Go to your Slack workspace settings
2. Navigate to "Apps" and search for "Incoming Webhooks"
3. Add to Slack and choose your target channel
4. Copy the webhook URL

**Option 2: Bot Token**
1. Create a new app at [Slack API](https://api.slack.com/apps)
2. Add the `chat:write` OAuth scope
3. Install the app to your workspace
4. Copy the Bot User OAuth Token

### LinkedIn API
1. Create an app at [LinkedIn Developers](https://www.linkedin.com/developers/)
2. Request access to the LinkedIn Share API
3. Generate an access token with the required scopes
4. Get your person URN from the LinkedIn API

## Usage

1. Navigate to the "Actions" tab in your GitHub repository
2. Select "Create Meetup Event and Share" workflow
3. Click "Run workflow"
4. Fill in all the required event details
5. Choose which platforms to share to
6. Click "Run workflow" to start the automation

## Local Development

```bash
# Clone the repository
git clone https://github.com/m0ab/meetop.git
cd meetop

# Install dependencies
go mod tidy

# Create .env file with your credentials
cp .env.example .env

# Build the application
go build -o meetop ./cmd/meetop

# Run locally (make sure environment variables are set)
./meetop
```

## Project Structure

```
├── cmd/meetop/          # Main application entry point
├── pkg/
│   ├── config/          # Configuration management
│   ├── meetup/          # Meetup.com API client
│   ├── slack/           # Slack API client
│   └── linkedin/        # LinkedIn API client
├── .github/workflows/   # GitHub Actions workflows
└── README.md           # This file
```

## Error Handling

The application includes comprehensive error handling and logging:
- All operations are logged to `meetop.log`
- Failed operations won't prevent other integrations from running
- GitHub Actions will upload logs as artifacts for debugging

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request