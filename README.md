# CyberArk API Command Harness

A modular command-line tool for interacting with the CyberArk API, written in Go. This project demonstrates a plugin-style architecture where each workflow is contained in its own source file and automatically registered at runtime.

## Features

- **Modular Workflow System**: Each workflow is a separate `.go` file that implements the `Workflow` interface
- **Automatic Registration**: Workflows self-register using Go's `init()` function
- **Secure Configuration**: Config files must have 600 permissions for security
- **Flexible Command Line**: Global options and workflow-specific options
- **Built-in Help**: Context-sensitive help for each workflow

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/cyberark-api-harness.git
cd cyberark-api-harness

# Build the program
go build -o cyberark

# Or build and install to $GOPATH/bin
go install
```

## Configuration

Create a configuration file at `~/.cyberark_api` with the following JSON structure:

```json
{
  "api_secret": "your-api-secret-here",
  "base_url": "https://your-cyberark-instance.example.com/PasswordVault/api",
  "username": "api_user",
  "timeout": 30
}
```

**Important**: Set proper permissions on the config file:
```bash
chmod 600 ~/.cyberark_api
```

You can also use a custom config file location:
```bash
cyberark --config /path/to/config.json workflow_name
```

## Usage

```bash
# General syntax
cyberark [--global-options] workflow_name [--workflow-options]

# Get help
cyberark --help
cyberark workflow_name --help

# Verify API connectivity
cyberark verify

# Using a custom config file
cyberark --config /path/to/config.json verify
```

## Built-in Workflows

### verify
Verifies API connectivity and configuration validity.

```bash
cyberark verify
```

## Available Workflows

### list-accounts
Lists accounts from CyberArk with optional filtering.

```bash
# List all accounts (default limit: 50)
cyberark list-accounts

# Filter by safe name
cyberark list-accounts --safe ProductionSafe

# Set custom limit
cyberark list-accounts --limit 100
```

## Adding New Workflows

To add a new workflow:

1. Create a new `.go` file in the project directory (e.g., `my_workflow.go`)
2. Implement the `Workflow` interface:
   ```go
   type Workflow interface {
       Execute(config *Config, args []string) error
       Help() string
   }
   ```
3. Register your workflow in an `init()` function:
   ```go
   func init() {
       RegisterWorkflow("my-workflow", &MyWorkflow{})
   }
   ```
4. Rebuild the program: `go build`

See `list_accounts.go` for a complete example.

## Project Structure

```
.
├── main.go              # Main program, config handling, workflow registry
├── list_accounts.go     # Example workflow module
├── .cyberark_api.sample # Sample configuration file
└── README.md           # This file
```

## Go Language Concepts Demonstrated

This project is designed as a learning tool for Go and demonstrates:

- **Interfaces**: Implicit interface satisfaction
- **Struct tags**: JSON marshaling/unmarshaling
- **Error handling**: Explicit error returns and checking
- **Flag parsing**: Command-line argument handling
- **init() functions**: Package initialization
- **Pointers**: Pass by reference for efficiency
- **Methods**: Attaching functions to types
- **Maps**: The workflow registry
- **Slices**: Dynamic arrays for arguments

## Security Notes

- Config files must have 600 permissions (owner read/write only)
- API secrets are never logged or displayed
- Sensitive data should never be passed via command-line arguments

## Development

```bash
# Format code (Go standard)
go fmt ./...

# Run with race detector
go run -race .

# Build for different platforms
GOOS=linux GOARCH=amd64 go build
GOOS=windows GOARCH=amd64 go build
GOOS=darwin GOARCH=amd64 go build
```

## License

[Your License Here]

## Contributing

Contributions are welcome! Please ensure:
- Code is properly formatted (`go fmt`)
- All workflows implement the `Workflow` interface
- Workflows include comprehensive help text
- New workflows are documented in this README

