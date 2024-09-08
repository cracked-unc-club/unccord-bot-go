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
uncord-bot-go/
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
├── .env.example              # Example environment variable file
├── go.mod                    # Go module file
└── go.sum                    # Go dependencies
```

## Getting Started

### Prerequisites

- [Go 1.20+](https://golang.org/dl/)
- [Lavalink](https://github.com/lavalink-devs/Lavalink) server
- TODO: Describe the process of creating a Discord application and obtaining a bot token

### Environment Variables

1. Copy the `.env.example` file to `.env`:

   ```bash
   cp .env.example .env
   ```

2. Edit the `.env` file and fill in your actual values:

   ```env
   DISCORD_TOKEN=your_discord_bot_token_here
   LAVALINK_HOST=localhost
   LAVALINK_PORT=2333
   LAVALINK_PASSWORD=youshallnotpass
   LAVALINK_SECURE=false
   ```

3. Load the environment variables:

   ```bash
   source .env
   ```

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/uncord-bot-go.git
   cd uncord-bot-go
   ```

2. Install dependencies:

   ```bash
   go mod download
   ```

### Running Lavalink

1. Download the latest Lavalink.jar from the [Lavalink releases page](https://github.com/lavalink-devs/Lavalink/releases).

2. Create an `application.yml` file in the same directory as the Lavalink.jar with the following content:

```yaml
server: # REST and WS server
  port: 2333
  address: 0.0.0.0
  http2:
    enabled: false
plugins:
  youtube:
    enabled: true
lavalink:
  plugins:
    - dependency: "dev.lavalink.youtube:youtube-plugin:1.7.2"
      snapshot: false
  server:
    password: "youshallnotpass"
    sources:
      # The default Youtube source is now deprecated and won't receive further updates. Please use https://github.com/lavalink-devs/youtube-source#plugin instead.
      youtube: false
      bandcamp: true
      soundcloud: true
      twitch: true
      vimeo: true
      nico: true
      http: true # warning: keeping HTTP enabled without a proxy configured could expose your server's IP address.
      local: false
    filters: # All filters are enabled by default
      volume: true
      equalizer: true
      karaoke: true
      timescale: true
      tremolo: true
      vibrato: true
      distortion: true
      rotation: true
      channelMix: true
      lowPass: true
    bufferDurationMs: 400
    frameBufferDurationMs: 5000
    opusEncodingQuality: 10
    resamplingQuality: LOW
    trackStuckThresholdMs: 10000
    useSeekGhosting: true
    youtubePlaylistLoadLimit: 6 # Number of pages at 100 each
    playerUpdateInterval: 5
    youtubeSearchEnabled: true
    soundcloudSearchEnabled: true
    gc-warnings: true

metrics:
  prometheus:
    enabled: false
    endpoint: /metrics

sentry:
  dsn: ""
  environment: ""

logging:
  file:
    path: ./logs/

  level:
    root: INFO
    lavalink: DEBUG
    lavalink.server.io.SocketContext: TRACE
    com.sedmelluq.discord.lavaplayer.tools.ExceptionTools: DEBUG

  request:
    enabled: true
    includeClientInfo: true
    includeHeaders: false
    includeQueryString: true
    includePayload: true
    maxPayloadLength: 10000


  logback:
    rollingpolicy:
      max-file-size: 1GB
      max-history: 30
   ```

3. Run Lavalink:

   ```bash
   java -jar Lavalink.jar
   ```

### Running the Bot

Run the bot:

```bash
go run cmd/main.go
```

## GitHub Actions CI Pipeline

### Overview

The CI pipeline is set up using GitHub Actions. It runs on every push or pull request to the `main` branch, and it performs the following steps:
- **Build**: Compiles the Golang bot.
- **Test**: Runs tests for the project.

### Workflow File (`.github/workflows/ci.yml`)

The GitHub Actions workflow is defined in `.github/workflows/ci.yml`:

```yaml
name: uncord-bot-go CI

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
      run: go build -o uncord-bot-go cmd/main.go
```

### Secrets

Make sure the following secret is added in your GitHub repository:

- `DISCORD_TOKEN`: Your Discord bot token.

## Contributing

Feel free to fork this repository, make your changes in a new branch, and submit a pull request.

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.
