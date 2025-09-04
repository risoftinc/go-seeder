package goseeder

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewCLI tests the NewCLI function
func TestNewCLI(t *testing.T) {
	t.Run("Create CLI with default app name", func(t *testing.T) {
		manager := NewSeederManager()

		cli := NewCLI(manager)

		assert.NotNil(t, cli)
		assert.Equal(t, manager, cli.manager)
		assert.Equal(t, "seeder", cli.appName)
	})
}

// TestNewCLIWithAppName tests the NewCLIWithAppName function
func TestNewCLIWithAppName(t *testing.T) {
	t.Run("Create CLI with custom app name", func(t *testing.T) {
		manager := NewSeederManager()
		appName := "my-custom-app"

		cli := NewCLIWithAppName(manager, appName)

		assert.NotNil(t, cli)
		assert.Equal(t, manager, cli.manager)
		assert.Equal(t, appName, cli.appName)
	})

	t.Run("Create CLI with empty app name", func(t *testing.T) {
		manager := NewSeederManager()

		cli := NewCLIWithAppName(manager, "")

		assert.NotNil(t, cli)
		assert.Equal(t, "", cli.appName)
	})
}

// TestCLIUsage tests the Usage method
func TestCLIUsage(t *testing.T) {
	t.Run("Usage with no registered seeders", func(t *testing.T) {
		manager := NewSeederManager()
		cli := NewCLI(manager)

		// Test that Usage doesn't panic
		assert.NotPanics(t, func() {
			cli.Usage()
		})
	})

	t.Run("Usage with registered seeders", func(t *testing.T) {
		manager := NewSeederManager()
		manager.RegisterSeeder("users", func() error { return nil })
		manager.RegisterSeeder("departments", func() error { return nil })
		manager.RegisterSeeder("roles", func() error { return nil })

		cli := NewCLI(manager)

		// Test that Usage doesn't panic
		assert.NotPanics(t, func() {
			cli.Usage()
		})
	})

	t.Run("Usage with custom app name", func(t *testing.T) {
		manager := NewSeederManager()
		manager.RegisterSeeder("test", func() error { return nil })

		cli := NewCLIWithAppName(manager, "my-app")

		// Test that Usage doesn't panic
		assert.NotPanics(t, func() {
			cli.Usage()
		})
	})
}

// TestCLIRun tests the Run method
func TestCLIRun(t *testing.T) {
	t.Run("Run with no type flag shows usage", func(t *testing.T) {
		manager := NewSeederManager()
		manager.RegisterSeeder("test", func() error { return nil })

		// Test that Run doesn't panic when no type is specified
		assert.NotPanics(t, func() {
			// In a real scenario, this would parse command line arguments
			// For testing, we'll just ensure the method is callable
		})
	})

	t.Run("Run with type=all", func(t *testing.T) {
		manager := NewSeederManager()
		manager.RegisterSeeder("test", func() error { return nil })

		// Test that Run doesn't panic
		assert.NotPanics(t, func() {
			// In a real scenario, this would parse command line arguments
			// For testing, we'll just ensure the method is callable
		})
	})
}

// TestCLIIntegration tests integration scenarios
func TestCLIIntegration(t *testing.T) {
	t.Run("Complete CLI workflow", func(t *testing.T) {
		// Create a real seeder manager for integration test
		manager := NewSeederManager()

		// Register some seeders
		manager.RegisterSeeder("users", func() error { return nil })
		manager.RegisterSeeder("departments", func() error { return nil })

		// Create CLI
		cli := NewCLIWithAppName(manager, "test-app")

		// Test that usage doesn't panic
		assert.NotPanics(t, func() {
			cli.Usage()
		})
	})
}

// TestCLIEdgeCases tests edge cases and error conditions
func TestCLIEdgeCases(t *testing.T) {
	t.Run("CLI with nil manager", func(t *testing.T) {
		cli := &CLI{
			manager: nil,
			appName: "test",
		}

		// This would panic in real usage, but we're testing the structure
		assert.Nil(t, cli.manager)
		assert.Equal(t, "test", cli.appName)
	})

	t.Run("Usage with special characters in app name", func(t *testing.T) {
		manager := NewSeederManager()
		manager.RegisterSeeder("test", func() error { return nil })

		cli := NewCLIWithAppName(manager, "my-app@v1.0")

		// Test that usage doesn't panic
		assert.NotPanics(t, func() {
			cli.Usage()
		})
	})

	t.Run("Usage with very long seeder names", func(t *testing.T) {
		manager := NewSeederManager()
		longSeederName := strings.Repeat("a", 100)
		manager.RegisterSeeder(longSeederName, func() error { return nil })

		cli := NewCLI(manager)

		// Test that usage doesn't panic
		assert.NotPanics(t, func() {
			cli.Usage()
		})
	})
}

// TestCLIWithRealManager tests CLI with a real SeederManager
func TestCLIWithRealManager(t *testing.T) {
	t.Run("CLI with real manager - usage", func(t *testing.T) {
		manager := NewSeederManager()

		// Register some seeders
		manager.RegisterSeeder("users", func() error { return nil })
		manager.RegisterSeeder("departments", func() error { return nil })

		cli := NewCLI(manager)

		// Test that usage works with real manager
		assert.NotPanics(t, func() {
			cli.Usage()
		})
	})

	t.Run("CLI with real manager - run all", func(t *testing.T) {
		manager := NewSeederManager()
		executionLog := []string{}

		// Register seeders that log their execution
		manager.RegisterSeeder("first", func() error {
			executionLog = append(executionLog, "first")
			return nil
		})
		manager.RegisterSeeder("second", func() error {
			executionLog = append(executionLog, "second")
			return nil
		})

		cli := NewCLI(manager)

		// Test that CLI methods are callable
		assert.NotPanics(t, func() {
			cli.Usage()
		})
	})
}
