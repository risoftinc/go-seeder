package goseeder

import (
	"flag"
	"log"
	"os"
	"strings"
)

// CLI handles command line interface for seeder operations
type CLI struct {
	manager *SeederManager
	appName string // Application name for usage display
}

// NewCLI creates a new CLI instance
func NewCLI(manager *SeederManager) *CLI {
	return &CLI{
		manager: manager,
		appName: "seeder", // Default app name
	}
}

// NewCLIWithAppName creates a new CLI instance with custom app name
func NewCLIWithAppName(manager *SeederManager, appName string) *CLI {
	return &CLI{
		manager: manager,
		appName: appName,
	}
}

// Run executes the seeder based on command line arguments
func (cli *CLI) Run() error {
	// Parse command line flags
	seedType := flag.String("type", "", "Type of seeder to run (all, or specific seeder name)")
	flag.Parse()

	// If no type specified, show usage and available seeders
	if *seedType == "" {
		cli.Usage()
		return nil
	}

	log.Printf("Starting seeder with type: %s", *seedType)

	switch *seedType {
	case "all":
		return cli.manager.RunAllSeeders()
	default:
		// Check if it's a specific seeder name
		if cli.manager.IsSeederRegistered(*seedType) {
			return cli.manager.RunSeederByName(*seedType)
		} else {
			log.Printf("Unknown seeder type: %s", *seedType)
			log.Printf("Available seeders: %v", cli.manager.GetRegisteredSeeders())
			cli.Usage()
			os.Exit(1)
		}
	}

	return nil
}

// Usage prints the usage information for the seeder
func (cli *CLI) Usage() {
	log.Println("=" + strings.Repeat("=", 60))
	log.Printf("DATABASE SEEDER - %s", strings.ToUpper(cli.appName))
	log.Println("=" + strings.Repeat("=", 60))
	log.Println("")
	log.Println("Usage:")
	log.Printf("  %s -type=all     # Run all seeders", cli.appName)
	log.Printf("  %s -type=<name>  # Run specific seeder", cli.appName)
	log.Printf("  %s               # Show this help", cli.appName)
	log.Println("")

	// Get registered seeders
	seeders := cli.manager.GetRegisteredSeeders()

	if len(seeders) == 0 {
		log.Println("No seeders registered yet.")
		return
	}

	log.Println("Available seeders (in execution order):")
	log.Println("-" + strings.Repeat("-", 40))

	// Show seeders with numbering
	for i, name := range seeders {
		log.Printf("  %d. %s", i+1, name)
		log.Printf("     Command: %s -type=%s", cli.appName, name)
		log.Println("")
	}

	log.Println("Quick commands:")
	log.Printf("  %s -type=all     # Run all seeders", cli.appName)
	log.Println("=" + strings.Repeat("=", 60))
}
