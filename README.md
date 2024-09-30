
# UNCCORD BOT GO
[Join Cracked Unc Club](https://discord.gg/3jfKWTwbeM)

This is a Golang-based Discord bot using the [disgo](https://github.com/disgoorg/disgo) library. The project is configured with a CI pipeline using GitHub Actions for building and testing the application. All secrets (such as the Discord bot token) are stored in GitHub Secrets.

## Library

[![image](https://github.com/user-attachments/assets/b382c075-b992-401b-8565-d46224345b44)](https://github.com/disgoorg/disgo)

I trust him.

## Roadmap

### Planned Features for Version 0.1:
- [ ] Build Voice-Master Functionality.
- [ ] Basic Discord message handling using 'disgo'.
- [ ] CI Pipeline setup using GitHub Actions for build and test.
- [ ] CD Pipeline setup (Where are we hosting? What tool are we going to use?).
- [ ] Full unit testing suite for basic features. 
- [ ] Develop a modular progress structure for easy scalability.

### To version 1.0:
Replace all third-party bots.
llm integration?
Twitter integration?
Automatic cracked unc event organization in your local metropolitian area? 


### Features


## Project Structure

```
unccord-bot-go/
│
├── cmd/
│   └── main.go               # Main entry point for the bot
│
├── config/
│   └── config.go             # Configuration loader (e.g., environment variables)
│
├── handlers/
│   ├── message_handler.go    # Handles messages
│   └── interaction_handler.go # Handles slash commands & interactions
│
├── services/
│   └── discord_service.go    # Business logic for interacting with Discord API using disgo
│
├── internal/
│   └── util.go               # Utility functions (logging, parsing, etc.)
│
├── .github/
│   └── workflows/
│       └── ci.yml            # CI pipeline configuration for GitHub Actions
│
├── go.mod                    # Go module file
└── go.sum                    # Go dependencies
```

## Getting Started

### Setting Up Your Discord Bot

1. Go to the [Discord Developer Portal](https://discord.com/developers/applications).
2. Click "New Application" and name your bot.
3. Navigate to the "Bot" tab and click "Add Bot".
4. Under "Token", click "Copy" to copy your bot token. Keep this secure!
5. Enable these "Privileged Gateway Intents":
   - Presence Intent
   - Server Members Intent
   - Message Content Intent
6. Save your changes.

### Bot Permissions

Ensure your bot has these permissions:
- Read Messages/View Channels
- Send Messages
- Embed Links
- Attach Files
- Read Message History
- Add Reactions
- Connect (to voice channels)
- Speak (in voice channels)

### Inviting the Bot to Your Server

1. In the Developer Portal, go to "OAuth2" > "URL Generator".
2. Select scopes: "bot" and "applications.commands".
3. Choose the permissions listed above.
4. Copy the generated URL and open it to invite the bot to your server.

## Project Setup

### Prerequisites

- [Go 1.20+](https://golang.org/dl/)
- Docker and Docker Compose
- Your Discord bot token

### Configuration

1. Clone the repository:
   ```bash
   https://github.com/cracked-unc-club/unccord-bot-go.git
   cd unccord-bot-go
   ```

2. Update the `.env` file in the project root:
   ```
   #DB config
   DB_HOST=db
   DB_PORT=5432
   DB_USER=potclean
   DB_PASSWORD=yourpass  # Change this to a secure password
   DB_NAME=potclean

   #Starboard config
   STARBOARD_CHANNEL_ID=1282793245289484420  # Update with your channel ID
   STAR_THRESHOLD=1

   #JoinToCreate config
   JOIN_TO_CREATE_CHANNEL_ID=1286835730705813574  # Update with your channel ID

   #Discord config
   DISCORD_TOKEN=yourtoken  # Replace with your actual bot token

   # Lavalink Configuration
   *JAVA*OPTIONS=-Xmx6G                                 
   SERVER_PORT=2333
   SERVER_ADDRESS=lavalink
   LAVALINK_SERVER_PASSWORD=yourpass  # Change this to a secure password
   ```
   Replace `yourpass`, `yourtoken`, and the channel IDs with your actual values.

### Building and Running with Docker

1. Ensure Docker and Docker Compose are installed on your system.

2. Build the Docker images:
   ```bash
   docker-compose build
   ```

3. Start the services:
   ```bash
   docker-compose up -d
   ```

4. Check the logs to ensure everything is running correctly:
   ```bash
   docker-compose logs -f
   ```

5. The first time you run the bot, you'll need to authorize the YouTube integration:
   - Look for a log message from the Lavalink container with a URL and code.
   - Go to the provided URL (usually https://www.google.com/device) and enter the code.
   - Use a burner Google account for this, not your main account.

### Troubleshooting

- If the bot doesn't connect, check your `DISCORD_TOKEN` in the `.env` file.
- For database issues, ensure the `DB_PASSWORD` is correct and the PostgreSQL container is running.
- If Lavalink fails to connect, verify the `LAVALINK_SERVER_PASSWORD` matches in both the bot and Lavalink configurations.

### Maintenance

- To stop the bot: `docker-compose down`
- To update: Pull the latest changes, rebuild, and restart the containers.
- Monitor logs regularly: `docker-compose logs -f`

Remember to keep your `.env` file and bot token secure. Never commit them to public repositories.


## GitHub Actions CI Pipeline

### Overview

The CI pipeline is set up using GitHub Actions. It runs on every push or pull request to the `main` branch, and it performs the following steps:
- **Build**: Compiles the Golang bot.
- **Test**: Runs tests for the project.

### Workflow File (`.github/workflows/ci.yml`)

The GitHub Actions workflow is defined in `.github/workflows/ci.yml`:

```yaml
name: unccord-bot-go CI

on:
  push:
    branches:
      - main
  pull_request:
    branches:
      - main

jobs:
  build:
    runs-on: ubuntu-latest

    steps:
    - name: Checkout code
      uses: actions/checkout@v2

    - name: Set up Go
      uses: actions/setup-go@v2
      with:
        go-version: 1.20

    - name: Install dependencies
      run: go mod download

    - name: Run tests
      run: go test ./...

    - name: Build the bot
      run: go build -o unccord-bot-go cmd/main.go
```

### Secrets

Make sure the following secret is added in your GitHub repository:

- `DISCORD_BOT_TOKEN`: Your Discord bot token.

## Contributing

Feel free to fork this repository, make your changes in a new branch, and submit a pull request.
