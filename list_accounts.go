package main

import (
	"flag"
	"fmt"
)

// ListAccountsWorkflow implements the Workflow interface
// In Go, you define a type (often an empty struct) and attach methods to it
// This is different from Python classes - there's no __init__ or self
type ListAccountsWorkflow struct{}

// Execute runs the list accounts workflow
// The receiver (w *ListAccountsWorkflow) is like Python's self
// The * means this is a pointer receiver (can modify the struct)
func (w *ListAccountsWorkflow) Execute(config *Config, args []string) error {
	// Create a FlagSet for this workflow's specific arguments
	fs := flag.NewFlagSet("list-accounts", flag.ExitOnError)

	// Define workflow-specific flags
	help := fs.Bool("help", false, "Show help for list-accounts workflow")
	fs.BoolVar(help, "h", false, "Show help (shorthand)")

	// Example: additional workflow-specific options
	safe := fs.String("safe", "", "Filter by safe name")
	limit := fs.Int("limit", 50, "Maximum number of accounts to return")

	// Parse the arguments passed to this workflow
	fs.Parse(args)

	// Handle help request
	if *help {
		w.printHelp()
		return nil
	}

	// Your workflow logic goes here
	fmt.Println("Listing CyberArk accounts...")
	fmt.Printf("Base URL: %s\n", config.BaseURL)

	// Access flag values using * (dereference pointer)
	if *safe != "" {
		fmt.Printf("Filtering by safe: %s\n", *safe)
	}
	fmt.Printf("Limit: %d\n", *limit)

	// TODO: Implement actual API call
	fmt.Println("\n[This would make an API call to list accounts]")

	return nil
}

// Help returns usage information for this workflow
func (w *ListAccountsWorkflow) Help() string {
	return "List accounts from CyberArk"
}

// printHelp displays detailed usage information
func (w *ListAccountsWorkflow) printHelp() {
	fmt.Println("List Accounts Workflow - Retrieve accounts from CyberArk")
	fmt.Println("\nUsage:")
	fmt.Println("  cyberark list-accounts [options]")
	fmt.Println("\nOptions:")
	fmt.Println("  -h, --help       Show this help message")
	fmt.Println("  --safe NAME      Filter accounts by safe name")
	fmt.Println("  --limit N        Maximum number of accounts to return (default: 50)")
	fmt.Println("\nExamples:")
	fmt.Println("  cyberark list-accounts")
	fmt.Println("  cyberark list-accounts --safe ProductionSafe")
	fmt.Println("  cyberark list-accounts --safe DevSafe --limit 100")
}

// init is a special function that runs automatically when the package is imported
// This is how we register our workflow - similar to Python's module-level code
// Go's init() functions run before main()
func init() {
	// Register this workflow with the name "list-accounts"
	// &ListAccountsWorkflow{} creates a pointer to a new instance
	RegisterWorkflow("list-accounts", &ListAccountsWorkflow{})
}

// Key Go Concepts Demonstrated Here:
//
// 1. Methods: Functions attached to types using receivers
//    - (w *ListAccountsWorkflow) is the receiver
//    - * means pointer receiver (can modify, more efficient)
//
// 2. Interfaces: Implicit satisfaction
//    - ListAccountsWorkflow satisfies Workflow automatically
//    - No "implements" keyword needed
//
// 3. init() function: Runs before main()
//    - Used for package initialization
//    - Perfect for registering plugins/workflows
//
// 4. Flag parsing: More manual than Python's argparse
//    - Must create FlagSet explicitly
//    - Flags are pointers (need * to access values)
//
// 5. Error handling: Explicit return values
//    - No exceptions, return error or nil
//    - Caller checks errors explicitly
