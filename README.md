# TideBot 🌊

A WhatsApp bot that provides tide information for Fuerteventura, Canary Islands using the WorldTides API.

## Features

- **User Registration**: Send "overpowered" to register for tide notifications
- **Tide Notifications**: Automated daily tide extremes via WhatsApp
- **REST API**: Manual job triggering and webhook handling
- **Real-time Data**: Uses WorldTides API for accurate tide predictions

## Quick Start

### Prerequisites

- Go 1.23+
- Node.js (for frontend build)
- Turso database
- Twilio WhatsApp API
- WorldTides API key
- Ngrok (for local development webhook tunneling)

### Environment Setup

Copy `.env.example` to `.env` and fill in your credentials:

```bash
TURSO_DB_URL=your_turso_database_url
TURSO_DB_AUTH_TOKEN=your_turso_auth_token
TWILIO_ACCOUNT_SID=your_twilio_sid
TWILIO_AUTH_TOKEN=your_twilio_token
TWILIO_WHATSAPP_FROM=whatsapp:+1234567890
WORLDTIDES_API_KEY=your_worldtides_api_key
```

### Build & Run

```bash
# Build frontend assets
npm run build

# Build and run the application
go build -o tidebot ./cmd
./tidebot --env development

# Or use Air for development with hot reload
air
```

The server starts on port `42069`.

### Development Setup with Ngrok

For local development, you'll need to expose your local server to receive WhatsApp webhooks from Twilio:

```bash
# Start ngrok tunnel (in a separate terminal)
ngrok http http://localhost:42069
```

Copy the ngrok URL (e.g., `https://5aad-85-254-47-205.ngrok-free.app`) and paste it in your Twilio WhatsApp Sandbox console as the webhook URL with the `/message` endpoint: `https://5aad-85-254-47-205.ngrok-free.app/message`

## API Endpoints

### WhatsApp Webhook
- `POST /message` - Receives WhatsApp messages from Twilio

### Jobs
- `POST /jobs/send-tide-extremes` - Send tide extremes to all registered users

## Usage

### Register for Notifications
Send a WhatsApp message with the text "overpowered" to your bot number.

### Trigger Manual Notifications
```bash
curl -X POST http://localhost:42069/jobs/send-tide-extremes
```

## Architecture

- **Backend**: Go with Echo framework
- **Database**: Turso (libSQL)
- **WhatsApp**: Twilio API
- **Tides**: WorldTides API
- **Frontend**: TypeScript + Tailwind CSS

## Development

### Package Structure
```
pkg/
├── environment/     # Environment configuration
├── jobs/           # Job scheduling and execution
├── users/          # User management (models, repositories, services)
├── whatsapp/       # WhatsApp integration and messaging
└── worldtides/     # WorldTides API client
```

### Development Server
Use Air for hot reload during development:
```bash
air
```

This watches Go, TypeScript, and template files and rebuilds automatically.

## License

MIT License