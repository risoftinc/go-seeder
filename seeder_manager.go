package seeder

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

// SeederItem represents a single seeder with its name and function
type SeederItem struct {
	Name     string
	Function func() error
}

// SeederManager manages all registered seeders
type SeederManager struct {
	db        *gorm.DB
	seeders   []SeederItem
	seederMap map[string]func() error
}

// NewSeederManager creates a new seeder manager instance
func NewSeederManager(db *gorm.DB) *SeederManager {
	return &SeederManager{
		db:        db,
		seeders:   make([]SeederItem, 0),
		seederMap: make(map[string]func() error),
	}
}

// RegisterSeeder registers a new seeder with validation for unique names
func (sm *SeederManager) RegisterSeeder(name string, function func() error) error {
	// Validate name is not empty
	if name == "" {
		return fmt.Errorf("seeder name cannot be empty")
	}

	// Check if name already exists
	if _, exists := sm.seederMap[name]; exists {
		return fmt.Errorf("seeder with name '%s' already exists", name)
	}

	// Add to slice and map
	seederItem := SeederItem{
		Name:     name,
		Function: function,
	}
	sm.seeders = append(sm.seeders, seederItem)
	sm.seederMap[name] = function

	log.Printf("Registered seeder: %s", name)
	return nil
}

// RegisterSeeders registers multiple seeders at once using variadic function
func (sm *SeederManager) RegisterSeeders(seeders ...SeederItem) error {
	for _, seeder := range seeders {
		if err := sm.RegisterSeeder(seeder.Name, seeder.Function); err != nil {
			return fmt.Errorf("failed to register seeder '%s': %w", seeder.Name, err)
		}
	}
	return nil
}

// GetRegisteredSeeders returns a list of all registered seeder names
func (sm *SeederManager) GetRegisteredSeeders() []string {
	names := make([]string, len(sm.seeders))
	for i, seeder := range sm.seeders {
		names[i] = seeder.Name
	}
	return names
}

// RunSeederByName runs a specific seeder by name
func (sm *SeederManager) RunSeederByName(name string) error {
	if function, exists := sm.seederMap[name]; exists {
		log.Printf("Running seeder: %s", name)
		if err := function(); err != nil {
			return fmt.Errorf("seeder '%s' failed: %w", name, err)
		}
		log.Printf("Seeder '%s' completed successfully", name)
		return nil
	}
	return fmt.Errorf("seeder with name '%s' not found", name)
}

// RunSeedersInOrder runs multiple seeders in the specified order
func (sm *SeederManager) RunSeedersInOrder(names []string) error {
	for _, name := range names {
		if err := sm.RunSeederByName(name); err != nil {
			return err
		}
	}
	return nil
}

// RunAllSeeders runs all registered seeders in order
func (sm *SeederManager) RunAllSeeders() error {
	log.Println("Running all seeders...")

	// Run all registered seeders in order
	for _, seeder := range sm.seeders {
		log.Printf("Running seeder: %s", seeder.Name)
		if err := seeder.Function(); err != nil {
			return fmt.Errorf("seeder '%s' failed: %w", seeder.Name, err)
		}
		log.Printf("Seeder '%s' completed successfully", seeder.Name)
	}

	log.Println("All seeders completed successfully!")
	return nil
}

// IsSeederRegistered checks if a seeder with the given name is registered
func (sm *SeederManager) IsSeederRegistered(name string) bool {
	_, exists := sm.seederMap[name]
	return exists
}
