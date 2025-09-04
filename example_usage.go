package goseeder

import (
	"log"
)

// ExampleBasicUsage demonstrates basic seeder usage
func ExampleBasicUsage() {
	// Create seeder manager
	manager := NewSeederManager()

	// Register individual seeder
	manager.RegisterSeeder("custom_seeder", func() error {
		log.Println("Running custom seeder...")
		// Your custom seeding logic here
		return nil
	})

	// Run specific seeder
	err := manager.RunSeederByName("custom_seeder")
	if err != nil {
		log.Printf("Error: %v", err)
	}
}

// ExampleVariadicUsage demonstrates variadic seeder registration
func ExampleVariadicUsage() {
	// Create seeder manager
	manager := NewSeederManager()

	// Register multiple seeders at once
	seeders := []SeederItem{
		{Name: "users", Function: func() error {
			log.Println("Seeding users...")
			return nil
		}},
		{Name: "departments", Function: func() error {
			log.Println("Seeding departments...")
			return nil
		}},
	}

	// Register all at once
	if err := manager.RegisterSeeders(seeders...); err != nil {
		log.Printf("Registration error: %v", err)
		return
	}

	// Run all seeders
	if err := manager.RunAllSeeders(); err != nil {
		log.Printf("Execution error: %v", err)
	}
}

// ExampleCLIUsage demonstrates CLI usage with custom app name
func ExampleCLIUsage() {
	// Create seeder manager
	manager := NewSeederManager()

	// Register some seeders
	manager.RegisterSeeder("test_seeder", func() error {
		log.Println("Test seeder executed")
		return nil
	})

	// Create CLI with custom app name
	cli := NewCLIWithAppName(manager, "my-app seeder")

	// Show usage (this will display help with custom app name)
	cli.Usage()

	// Run CLI (this will parse command line arguments)
	// err := cli.Run()
}

// ExampleLibraryUsage demonstrates using seeder as a library (no CLI)
func ExampleLibraryUsage() {
	// Create seeder manager
	manager := NewSeederManager()

	// Register seeders
	manager.RegisterSeeder("users", func() error {
		log.Println("Seeding users...")
		return nil
	})

	manager.RegisterSeeder("departments", func() error {
		log.Println("Seeding departments...")
		return nil
	})

	// Run specific seeder programmatically
	if err := manager.RunSeederByName("users"); err != nil {
		log.Printf("Error seeding users: %v", err)
	}

	// Run multiple seeders in specific order
	order := []string{"departments", "users"}
	if err := manager.RunSeedersInOrder(order); err != nil {
		log.Printf("Error running seeders in order: %v", err)
	}

	// Check if seeder is registered
	if manager.IsSeederRegistered("users") {
		log.Println("Users seeder is registered")
	}

	// Get list of registered seeders
	seeders := manager.GetRegisteredSeeders()
	log.Printf("Registered seeders: %v", seeders)
}
