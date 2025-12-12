package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

// Config represents the structure of our configuration file
// In Go, struct tags (like `json:"api_secret"`) tell the JSON encoder/decoder
// how to map JSON fields to struct fields
type Config struct {
	APISecret string `json:"api_secret"` // Required field
	BaseURL   string `json:"base_url"`   // CyberArk API base URL
	Username  string `json:"username"`   // Optional: API username
	Timeout   int    `json:"timeout"`    // Optional: request timeout in seconds
}

// Workflow is an interface that all workflow modules must implement
// Interfaces in Go are implicit - any type that has these methods
// automatically satisfies this interface (unlike Python's explicit inheritance)
type Workflow interface {
	// Execute runs the workflow with the provided config and arguments
	Execute(config *Config, args []string) error

	// Help returns usage information for this workflow
	Help() string
}

// WorkflowRegistry maps workflow names to their implementations
// This is similar to a Python dictionary with string keys
var WorkflowRegistry = make(map[string]Workflow)

// RegisterWorkflow adds a workflow to the registry
// This will be called by each workflow module's init() function
func RegisterWorkflow(name string, workflow Workflow) {
	WorkflowRegistry[name] = workflow
}

// loadConfig reads and parses the configuration file
// Go functions can return multiple values - here we return both
// the config and an error (idiomatic Go error handling)
func loadConfig(path string) (*Config, error) {
	// Expand ~ to home directory if present
	// Unlike Python, Go doesn't automatically expand ~ in paths
	if len(path) > 0 && path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		path = filepath.Join(home, path[1:])
	}

	// Check file permissions (must be 600 for security)
	// In Go, we use the os.Stat function to get file info
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to stat config file: %w", err)
	}

	// Get file permissions - this is Unix-specific
	// The & 0777 masks out everything except permission bits
	mode := info.Mode().Perm()
	if mode != 0600 {
		return nil, fmt.Errorf("config file must have 0600 permissions, has %o", mode)
	}

	// Read the entire file into memory
	// In Go, we explicitly handle the byte slice; no automatic string conversion
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse JSON into our Config struct
	// The & operator gets the address of config (pointer)
	// This is necessary because Unmarshal modifies the struct
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %w", err)
	}

	// Validate required fields
	if config.APISecret == "" {
		return nil, fmt.Errorf("api_secret is required in config file")
	}
	if config.BaseURL == "" {
		return nil, fmt.Errorf("base_url is required in config file")
	}

	// Return pointer to config and nil error (success)
	return &config, nil
}

// printUsage displays general program usage information
func printUsage() {
	fmt.Println("CyberArk API Command Harness")
	fmt.Println("\nUsage:")
	fmt.Println("  cyberark [--global-options] workflow_name [--workflow-options]")
	fmt.Println("\nGlobal Options:")
	fmt.Println("  -c, --config PATH    Path to configuration file (default: ~/.cyberark_api)")
	fmt.Println("  -h, --help           Show this help message")
	fmt.Println("\nBuilt-in Workflows:")
	fmt.Println("  verify               Verify API connectivity")
	fmt.Println("\nRegistered Workflows:")

	// Range is Go's way of iterating over maps, slices, arrays, etc.
	// It's similar to Python's for...in but gives you both key and value
	for name := range WorkflowRegistry {
		fmt.Printf("  %s\n", name)
	}

	fmt.Println("\nFor workflow-specific help:")
	fmt.Println("  cyberark workflow_name --help")
}

// verifyConnectivity is the built-in verify workflow
func verifyConnectivity(config *Config, args []string) error {
	// Create a new FlagSet for this workflow's arguments
	// FlagSet is like Python's argparse, but more manual
	fs := flag.NewFlagSet("verify", flag.ExitOnError)
	help := fs.Bool("help", false, "Show help for verify workflow")
	fs.BoolVar(help, "h", false, "Show help (shorthand)")

	// Parse the workflow-specific arguments
	fs.Parse(args)

	if *help {
		fmt.Println("Verify Workflow - Test CyberArk API connectivity")
		fmt.Println("\nUsage:")
		fmt.Println("  cyberark verify [options]")
		fmt.Println("\nOptions:")
		fmt.Println("  -h, --help    Show this help message")
		return nil
	}

	fmt.Println("Verifying CyberArk API connectivity...")
	fmt.Printf("Base URL: %s\n", config.BaseURL)
	fmt.Println("API Secret: [REDACTED]")

	// TODO: Implement actual API call here
	// For now, just verify we have the configuration
	fmt.Println("\n✓ Configuration loaded successfully")
	fmt.Println("✓ API credentials present")
	fmt.Println("\nNote: Actual API connectivity test not yet implemented")

	return nil
}

func main() {
	// Define global flags
	// In Go, flags must be defined before parsing
	var configPath string
	var showHelp bool

	// flag.StringVar binds a flag to an existing variable
	// This is different from Python where you typically get a namespace object
	flag.StringVar(&configPath, "config", "~/.cyberark_api", "Path to config file")
	flag.StringVar(&configPath, "c", "~/.cyberark_api", "Path to config file (shorthand)")
	flag.BoolVar(&showHelp, "help", false, "Show help")
	flag.BoolVar(&showHelp, "h", false, "Show help (shorthand)")

	// Custom usage function
	flag.Usage = printUsage

	// Parse only the global flags
	// We need to manually handle the positional argument (workflow name)
	flag.Parse()

	// Check for help flag
	if showHelp {
		printUsage()
		os.Exit(0)
	}

	// Get remaining arguments (after flags)
	// flag.Args() returns a slice of strings - similar to Python list
	args := flag.Args()

	// Check if workflow name was provided
	if len(args) == 0 {
		fmt.Println("Error: workflow name required")
		printUsage()
		os.Exit(1)
	}

	workflowName := args[0]
	workflowArgs := args[1:] // Slice syntax: from index 1 to end

	// Load configuration
	config, err := loadConfig(configPath)
	if err != nil {
		// In Go, we explicitly check errors after each operation
		// This is more verbose than Python's try/except but more explicit
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Handle built-in verify workflow
	if workflowName == "verify" {
		if err := verifyConnectivity(config, workflowArgs); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Look up workflow in registry
	// The "comma ok" idiom checks if a key exists in a map
	// Similar to dict.get() in Python, but built into the language
	workflow, ok := WorkflowRegistry[workflowName]
	if !ok {
		fmt.Fprintf(os.Stderr, "Error: unknown workflow '%s'\n", workflowName)
		printUsage()
		os.Exit(1)
	}

	// Execute the workflow
	if err := workflow.Execute(config, workflowArgs); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing workflow: %v\n", err)
		os.Exit(1)
	}
}

// Note: In Go, you can also check permissions more portably using:
func checkFilePermissions(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	// For Unix systems, check exact permissions
	if info.Mode().Perm() != 0600 {
		// On Windows, this check may need to be different
		// Windows uses a different permission model (ACLs)
		if runtime.GOOS != "windows" {
			return fmt.Errorf("file must have 0600 permissions")
		}
	}
	return nil
}
