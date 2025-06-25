# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Architecture

This is **Lagoon** - a WhatsApp webhook integration service combining a Go backend with TypeScript frontend. The application uses Twilio's WhatsApp API to receive and send messages via webhooks.

### Technology Stack
- **Backend**: Go with Echo v4 framework
- **Database**: Turso (libSQL/SQLite-compatible) with golang-migrate
- **WhatsApp Integration**: Twilio WhatsApp API
- **Frontend Build**: esbuild + TypeScript + Tailwind CSS + HTMX
- **Module Name**: `tidebot` (internal Go module name)

### Key Architecture Patterns
- **Clean Architecture**: `cmd/` for application entry, `pkg/` for business logic
- **Interface-Based Design**: `WhatsappClient` interface for dependency injection
- **Migration-First Database**: Automatic schema migrations on startup
- **Webhook-Driven**: Receives WhatsApp messages via `/message` endpoint

## Common Development Commands

### Build Commands
```bash
# Full build (CSS + TypeScript)
npm run build

# Production build (minified)
npm run build:prod

# Individual builds
npm run build:css    # Tailwind CSS compilation
npm run build:ts     # TypeScript via esbuild
```

### Go Application
```bash
# Run the server (port 42069)
go run cmd/main.go

# Build binary
go build -o lagoon cmd/main.go
```

### Database Operations
Database migrations run automatically on startup. Migration files are in `db/migrations/`:
- `000001_create_users_table.up.sql` - Creates users table
- `000001_create_users_table.down.sql` - Rollback migration

## Environment Configuration

Required environment variables:
- `TURSO_DB_URL` - Database connection string
- `TURSO_DB_AUTH_TOKEN` - Database authentication token  
- `TWILIO_WHATSAPP_FROM` - WhatsApp sender number (format: `whatsapp:+1234567890`)

Environment files are loaded in order:
1. Base `.env` file
2. Environment-specific `.env.{environment}` file

## Core Components

### Backend Structure
- `cmd/main.go` - Server setup, middleware, route registration
- `cmd/environment.go` - Environment configuration management
- `pkg/whatsapp/apiclient.go` - Twilio API client implementation
- `pkg/whatsapp/webhookcontroller.go` - Webhook handlers for incoming messages

### Frontend Structure
- All `.ts` files are automatically discovered and bundled
- `index.d.ts` - Global TypeScript definitions for HTMX and custom types
- `assets/` - Static files served by the application
- `css/styles.css` - Tailwind CSS input file

### Database Schema
Single `users` table with:
- `id` (PRIMARY KEY)
- `phone_number` (UNIQUE, INDEXED) - WhatsApp phone numbers
- `name`, `created_at`, `updated_at`

## Development Notes

- Server runs on port 42069 (hardcoded)
- All webhook data is logged for debugging
- TypeScript uses strict mode with ES6 target
- Static assets served from `/assets` route
- Database migrations run automatically on application start
- Phone numbers stored clean in database; `whatsapp:` prefix only added for Twilio API calls