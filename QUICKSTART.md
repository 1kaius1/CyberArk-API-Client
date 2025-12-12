# Quick Start Guide

## Setup

1. **Ensure Go is installed**
   ```bash
   go version  # Should show 1.21 or higher
   ```

2. **Clone and setup the project**
   ```bash
   git clone https://github.com/yourusername/cyberark-api-harness.git
   cd cyberark-api-harness
   
   # Initialize Go module (if needed)
   go mod init github.com/yourusername/cyberark-api-harness
   go mod tidy
   ```

3. **Create configuration file**
   ```bash
   cp .cyberark_api.sample ~/.cyberark_api
   chmod 600 ~/.cyberark_api
   
   # Edit with your actual credentials
   nano ~/.cyberark_api
   ```

4. **Build the project**
   ```bash
   # Using make (recommended)
   make build
   
   # Or directly with go
   go build -o cyberark
   ```

5. **Test the installation**
   ```bash
   ./cyberark --help
   ./cyberark verify
   ```

## Project Structure

```
cyberark-api-harness/
├── main.go              # Entry point, config, workflow registry
├── api_client.go        # HTTP client for API calls
├── list_accounts.go     # Example workflow
├── go.mod              # Module definition (like requirements.txt)
├── .cyberark_api        # Your config (not committed)
├── .cyberark_api.sample # Sample config (committed)
├── Makefile            # Build automation
└── README.md           # Documentation
```

## How It Works

### 1. The Main Program (main.go)

The main program:
- Parses command-line arguments
- Loads configuration from `~/.cyberark_api`
- Looks up the requested workflow in the registry
- Executes the workflow

### 2. Workflow Registration

Each workflow file has an `init()` function that registers itself:

```go
// In list_accounts.go
func init() {
    RegisterWorkflow("list-accounts", &ListAccountsWorkflow{})
}
```

Go automatically runs all `init()` functions before `main()`. This means:
- **Just create a new `.go` file** with a workflow
- **It automatically becomes available** when you rebuild
- **No need to modify main.go** or maintain a list of workflows

### 3. The Workflow Interface

Every workflow must implement:

```go
type Workflow interface {
    Execute(config *Config, args []string) error
    Help() string
}
```

## Creating Your First Workflow

Let's create a workflow to retrieve a specific account:

```bash
# Create new file
nano get_account.go
```

```go
package main

import (
    "encoding/json"
    "flag"
    "fmt"
)

// Account represents a CyberArk account
type Account struct {
    ID       string `json:"id"`
    Name     string `json:"name"`
    Address  string `json:"address"`
    UserName string `json:"userName"`
    SafeName string `json:"safeName"`
}

type GetAccountWorkflow struct{}

func (w *GetAccountWorkflow) Execute(config *Config, args []string) error {
    fs := flag.NewFlagSet("get-account", flag.ExitOnError)
    help := fs.Bool("help", false, "Show help")
    fs.BoolVar(help, "h", false, "Show help")
    accountID := fs.String("id", "", "Account ID (required)")
    
    fs.Parse(args)
    
    if *help {
        fmt.Println("Get Account Workflow")
        fmt.Println("\nUsage:")
        fmt.Println("  cyberark get-account --id <account-id>")
        return nil
    }
    
    if *accountID == "" {
        return fmt.Errorf("--id is required")
    }
    
    // Create API client
    client := NewAPIClient(config)
    
    // Make API request
    endpoint := fmt.Sprintf("accounts/%s", *accountID)
    data, err := client.Get(endpoint)
    if err != nil {
        return err
    }
    
    // Parse response
    var account Account
    if err := json.Unmarshal(data, &account); err != nil {
        return fmt.Errorf("failed to parse response: %w", err)
    }
    
    // Display results
    fmt.Printf("Account ID: %s\n", account.ID)
    fmt.Printf("Name: %s\n", account.Name)
    fmt.Printf("Address: %s\n", account.Address)
    fmt.Printf("Username: %s\n", account.UserName)
    fmt.Printf("Safe: %s\n", account.SafeName)
    
    return nil
}

func (w *GetAccountWorkflow) Help() string {
    return "Retrieve a specific account by ID"
}

func init() {
    RegisterWorkflow("get-account", &GetAccountWorkflow{})
}
```

Now rebuild and use it:

```bash
make build
./cyberark get-account --id 12345
```

## Common Patterns

### 1. Making API Calls

```go
client := NewAPIClient(config)

// GET request
data, err := client.Get("accounts")
if err != nil {
    return err
}

// POST request
payload := map[string]interface{}{
    "name": "New Account",
    "address": "192.168.1.1",
}
data, err := client.Post("accounts", payload)

// PUT request
data, err := client.Put("accounts/123", payload)

// DELETE request
err := client.Delete("accounts/123")
```

### 2. Parsing JSON Responses

```go
// Define struct matching API response
type Response struct {
    Items []Account `json:"items"`
    Total int       `json:"total"`
}

// Parse
var response Response
if err := json.Unmarshal(data, &response); err != nil {
    return fmt.Errorf("parse error: %w", err)
}
```

### 3. Command-Line Flags

```go
fs := flag.NewFlagSet("my-workflow", flag.ExitOnError)

// Boolean flag
verbose := fs.Bool("verbose", false, "Verbose output")

// String flag with default
format := fs.String("format", "json", "Output format")

// Integer flag
limit := fs.Int("limit", 50, "Max results")

fs.Parse(args)

// Use flags (must dereference with *)
if *verbose {
    fmt.Println("Verbose mode enabled")
}
```

### 4. Error Handling

```go
// Check errors immediately
data, err := client.Get("endpoint")
if err != nil {
    // Wrap error with context
    return fmt.Errorf("failed to get data: %w", err)
}

// Multiple operations
if err := step1(); err != nil {
    return err
}
if err := step2(); err != nil {
    return err
}
```

### 5. Formatted Output

```go
// Simple
fmt.Println("Hello")
fmt.Printf("Name: %s, Age: %d\n", name, age)

// To stderr
fmt.Fprintf(os.Stderr, "Error: %v\n", err)

// JSON output
jsonData, _ := json.MarshalIndent(data, "", "  ")
fmt.Println(string(jsonData))
```

## Development Workflow

```bash
# Make changes to code
nano my_workflow.go

# Format code (Go standard)
go fmt ./...

# Build
make build

# Test
./cyberark my-workflow --help
./cyberark my-workflow --option value

# When ready, commit
git add .
git commit -m "Add my-workflow"
git push
```

## Debugging Tips

### 1. Add Debug Output
```go
fmt.Printf("DEBUG: value = %+v\n", myStruct)  // %+v shows field names
```

### 2. Check Error Details
```go
if err != nil {
    fmt.Fprintf(os.Stderr, "Detailed error: %+v\n", err)
    return err
}
```

### 3. Use Go's Built-in Tools
```bash
# Check for common mistakes
go vet

# Format code
go fmt ./...

# See what would be built
go list -f '{{.GoFiles}}'
```

## Next Steps

1. **Read the language comparison**: See `GO_VS_PYTHON.md` for detailed Python vs Go comparisons
2. **Study the example workflow**: `list_accounts.go` shows all the patterns
3. **Create your workflows**: Add new `.go` files for each CyberArk operation
4. **Refer to CyberArk API docs**: Adapt the examples to actual API endpoints

## Common Gotchas

1. **Forgetting to check errors**: Every function that returns `error` must be checked
2. **Not using pointers**: Use `*` when you need to modify structs
3. **Capitalization matters**: `MyFunction` is public, `myFunction` is private
4. **Nil vs zero values**: Empty string is `""`, not `nil`
5. **Slices vs arrays**: Use slices `[]int` (dynamic), not arrays `[10]int` (fixed)
6. **Must rebuild after changes**: Unlike Python, Go must be recompiled

## Resources

- [Official Go Tour](https://go.dev/tour/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go by Example](https://gobyexample.com/)
- [CyberArk REST API Documentation](https://docs.cyberark.com/)

