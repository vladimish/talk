# Talk - Telegram AI Bot

A Go-based Telegram Bot / LLM Wrapper made for 2025 T3 Chat Cloneathon


Try it out at [@ookaa_chat_bot](https://t.me/ookaa_chat_bot)


[Video overview](https://youtu.be/pyXcZAZu_to)

## Features

- AI-powered conversations using OpenRouter's API
- Message history tracking with PostgreSQL
- Streaming responses for better user experience
- Markdown formatting support via Telegramify
- Clean architecture using Domain-Driven Design (DDD) with Hexagonal/Ports & Adapters pattern
- Automatic database migrations on startup
- Graceful shutdown handling

## Prerequisites

- Go 1.24 or higher
- PostgreSQL database
- Docker and Docker Compose (for running PostgreSQL)
- Telegram Bot Token (from [@BotFather](https://t.me/botfather))
- OpenAI API Key
- Telegramify service (optional, for enhanced markdown formatting)

## Quick Start

Just run the unified startup script that handles everything (don't forget to set TG_TOKEN with your bot token and OPENAI_API_KEY with your OpenRouter key):

```bash
git clone https://github.com/vladimish/talk.git
cd talk
./start_docker.sh
```

The script will:
- Load existing environment variables from `.env` file (if it exists)
- Prompt for any missing required variables (TG_TOKEN, OPENAI_API_KEY)
- Set up sensible defaults for all other variables
- Start all dependencies (PostgreSQL, Redis, MinIO, Telegramify)
- Wait for services to be healthy

That's it! The bot will automatically apply any pending database migrations on startup.

## Architecture

The project follows Hexagonal Architecture (Ports & Adapters) with clear separation of concerns:

```
internal/
├── adapter/         # External system integrations
│   ├── in/         # Inbound adapters
│   │   └── tg/     # Telegram bot handler
│   └── out/        # Outbound adapters
│       ├── openai/ # OpenAI API client
│       ├── pg/     # PostgreSQL storage
│       ├── telegramify/ # Markdown formatter
│       └── tg/     # Telegram message sender
├── domain/         # Core business entities
├── port/          # Interface definitions
└── service/       # Business logic orchestration
```

## Development

### Database Management

The project uses:
- **Goose** for database migrations (automatic on startup)
- **SQLC** for type-safe SQL code generation

To create a new migration:
```bash
cd db/migrations
GOOSE_DBSTRING=postgresql://postgres:postgres@localhost:5432/postgres \
  go tool goose postgres create your_migration_name sql
```

After modifying SQL queries or schema:
```bash
make generate
```

### Building the Project

```bash
go build -o talk cmd/bot/main.go
```

### Running with Docker

Build and run the bot using Docker:
```bash
docker build -t talk .
docker run --env-file .env talk
```

## The Startup Script

The project includes a single unified startup script `./start.sh` that handles all scenarios:

- **Interactive setup**: Prompts for missing environment variables
- **File-based config**: Loads from `.env` file if it exists  
- **Docker management**: Creates networks, starts services, waits for health
- **Flexible execution**: Choose between local Go execution or Docker container
- **User-friendly**: Provides clear status updates and helpful commands

## Environment Variables

| Variable | Required | Description | Default |
|----------|----------|-------------|---------|
| `DATABASE_URL` | Yes | PostgreSQL connection string | - |
| `TG_TOKEN` | Yes | Telegram bot token | - |
| `OPENAI_API_KEY` | Yes | OpenAI API key | - |
| `TELEGRAMIFY_URL` | No | Telegramify service URL | `http://localhost:8000` |

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.