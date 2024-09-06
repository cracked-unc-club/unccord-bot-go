
# UNCCORD BOT GO
[Join Cracked Unc Club](https://discord.gg/3jfKWTwbeM)

This is a Golang-based Discord bot using the [disgo](https://github.com/disgoorg/disgo) library. The project is configured with a CI pipeline using GitHub Actions for building and testing the application. All secrets (such as the Discord bot token) are stored in GitHub Secrets.

## Library

[![image](https://github.com/user-attachments/assets/b382c075-b992-401b-8565-d46224345b44)](https://github.com/disgoorg/disgo)

I trust him.

## Roadmap

### Planned Features for Version 0.1:
- [ ] Be able to replace the Color-Chan and Unpaid Intern (Carlbot) bots utilizing native discord features. (Anything related to role assignment currently in [#react-roles](https://discord.com/channels/1276883668559724544/1277649676698386585)).
- [ ] Build Voice-Master Functionality
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
├── go.mod                    # Go module file
└── go.sum                    # Go dependencies
```

## Getting Started

### Prerequisites

- [Go 1.20+](https://golang.org/dl/)
- TODO: Describe the process for obtaining and setting bot token locally

### Environment Variables

Set the required environment variables:

```bash
export DISCORD_BOT_TOKEN="your_discord_bot_token"
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

3. Run the bot:

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

- `DISCORD_BOT_TOKEN`: Your Discord bot token.

## Contributing

Feel free to fork this repository, make your changes in a new branch, and submit a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
