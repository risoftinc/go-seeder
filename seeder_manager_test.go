package seeder

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSeederManagerForManager is a mock implementation of SeederManager for testing
type MockSeederManagerForManager struct {
	mock.Mock
}

// TestSeederItem tests the SeederItem struct
func TestSeederItem(t *testing.T) {
	t.Run("SeederItem creation", func(t *testing.T) {
		name := "test_seeder"
		function := func() error { return nil }

		item := SeederItem{
			Name:     name,
			Function: function,
		}

		assert.Equal(t, name, item.Name)
		assert.NotNil(t, item.Function)
	})
}

// TestNewSeederManager tests the NewSeederManager function
func TestNewSeederManager(t *testing.T) {
	t.Run("Create new seeder manager", func(t *testing.T) {
		manager := NewSeederManager()

		assert.NotNil(t, manager)
		assert.NotNil(t, manager.seeders)
		assert.NotNil(t, manager.seederMap)
		assert.Len(t, manager.seeders, 0)
		assert.Len(t, manager.seederMap, 0)
	})
}

// TestRegisterSeeder tests the RegisterSeeder method
func TestRegisterSeeder(t *testing.T) {
	t.Run("Register valid seeder", func(t *testing.T) {
		manager := NewSeederManager()
		name := "test_seeder"
		function := func() error { return nil }

		err := manager.RegisterSeeder(name, function)

		assert.NoError(t, err)
		assert.Len(t, manager.seeders, 1)
		assert.Len(t, manager.seederMap, 1)
		assert.Equal(t, name, manager.seeders[0].Name)
		assert.True(t, manager.IsSeederRegistered(name))
	})

	t.Run("Register seeder with empty name", func(t *testing.T) {
		manager := NewSeederManager()
		function := func() error { return nil }

		err := manager.RegisterSeeder("", function)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "seeder name cannot be empty")
		assert.Len(t, manager.seeders, 0)
		assert.Len(t, manager.seederMap, 0)
	})

	t.Run("Register duplicate seeder", func(t *testing.T) {
		manager := NewSeederManager()
		name := "test_seeder"
		function1 := func() error { return nil }
		function2 := func() error { return nil }

		err1 := manager.RegisterSeeder(name, function1)
		assert.NoError(t, err1)

		err2 := manager.RegisterSeeder(name, function2)
		assert.Error(t, err2)
		assert.Contains(t, err2.Error(), "already exists")
		assert.Len(t, manager.seeders, 1)
		assert.Len(t, manager.seederMap, 1)
	})

	t.Run("Register multiple unique seeders", func(t *testing.T) {
		manager := NewSeederManager()

		err1 := manager.RegisterSeeder("seeder1", func() error { return nil })
		err2 := manager.RegisterSeeder("seeder2", func() error { return nil })
		err3 := manager.RegisterSeeder("seeder3", func() error { return nil })

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NoError(t, err3)
		assert.Len(t, manager.seeders, 3)
		assert.Len(t, manager.seederMap, 3)
	})
}

// TestRegisterSeeders tests the RegisterSeeders method
func TestRegisterSeeders(t *testing.T) {
	t.Run("Register multiple seeders successfully", func(t *testing.T) {
		manager := NewSeederManager()

		seeders := []SeederItem{
			{Name: "seeder1", Function: func() error { return nil }},
			{Name: "seeder2", Function: func() error { return nil }},
			{Name: "seeder3", Function: func() error { return nil }},
		}

		err := manager.RegisterSeeders(seeders...)

		assert.NoError(t, err)
		assert.Len(t, manager.seeders, 3)
		assert.Len(t, manager.seederMap, 3)
	})

	t.Run("Register seeders with duplicate name", func(t *testing.T) {
		manager := NewSeederManager()

		seeders := []SeederItem{
			{Name: "seeder1", Function: func() error { return nil }},
			{Name: "seeder1", Function: func() error { return nil }}, // Duplicate
		}

		err := manager.RegisterSeeders(seeders...)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to register seeder 'seeder1'")
		assert.Len(t, manager.seeders, 1) // Only first one should be registered
	})

	t.Run("Register empty seeders list", func(t *testing.T) {
		manager := NewSeederManager()

		err := manager.RegisterSeeders()

		assert.NoError(t, err)
		assert.Len(t, manager.seeders, 0)
		assert.Len(t, manager.seederMap, 0)
	})
}

// TestGetRegisteredSeeders tests the GetRegisteredSeeders method
func TestGetRegisteredSeeders(t *testing.T) {
	t.Run("Get empty seeders list", func(t *testing.T) {
		manager := NewSeederManager()

		seeders := manager.GetRegisteredSeeders()

		assert.NotNil(t, seeders)
		assert.Len(t, seeders, 0)
	})

	t.Run("Get registered seeders", func(t *testing.T) {
		manager := NewSeederManager()

		manager.RegisterSeeder("seeder1", func() error { return nil })
		manager.RegisterSeeder("seeder2", func() error { return nil })
		manager.RegisterSeeder("seeder3", func() error { return nil })

		seeders := manager.GetRegisteredSeeders()

		assert.Len(t, seeders, 3)
		assert.Contains(t, seeders, "seeder1")
		assert.Contains(t, seeders, "seeder2")
		assert.Contains(t, seeders, "seeder3")
	})

	t.Run("Get seeders in registration order", func(t *testing.T) {
		manager := NewSeederManager()

		manager.RegisterSeeder("first", func() error { return nil })
		manager.RegisterSeeder("second", func() error { return nil })
		manager.RegisterSeeder("third", func() error { return nil })

		seeders := manager.GetRegisteredSeeders()

		assert.Equal(t, "first", seeders[0])
		assert.Equal(t, "second", seeders[1])
		assert.Equal(t, "third", seeders[2])
	})
}

// TestRunSeederByName tests the RunSeederByName method
func TestRunSeederByName(t *testing.T) {
	t.Run("Run existing seeder successfully", func(t *testing.T) {
		manager := NewSeederManager()
		executed := false

		manager.RegisterSeeder("test_seeder", func() error {
			executed = true
			return nil
		})

		err := manager.RunSeederByName("test_seeder")

		assert.NoError(t, err)
		assert.True(t, executed)
	})

	t.Run("Run non-existing seeder", func(t *testing.T) {
		manager := NewSeederManager()

		err := manager.RunSeederByName("non_existing")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("Run seeder that returns error", func(t *testing.T) {
		manager := NewSeederManager()
		expectedError := errors.New("seeder error")

		manager.RegisterSeeder("error_seeder", func() error {
			return expectedError
		})

		err := manager.RunSeederByName("error_seeder")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed")
		assert.Contains(t, err.Error(), "error_seeder")
	})
}

// TestRunSeedersInOrder tests the RunSeedersInOrder method
func TestRunSeedersInOrder(t *testing.T) {
	t.Run("Run seeders in order successfully", func(t *testing.T) {
		manager := NewSeederManager()
		executionOrder := []string{}

		manager.RegisterSeeder("first", func() error {
			executionOrder = append(executionOrder, "first")
			return nil
		})
		manager.RegisterSeeder("second", func() error {
			executionOrder = append(executionOrder, "second")
			return nil
		})
		manager.RegisterSeeder("third", func() error {
			executionOrder = append(executionOrder, "third")
			return nil
		})

		order := []string{"first", "second", "third"}
		err := manager.RunSeedersInOrder(order)

		assert.NoError(t, err)
		assert.Equal(t, []string{"first", "second", "third"}, executionOrder)
	})

	t.Run("Run seeders with one failing", func(t *testing.T) {
		manager := NewSeederManager()
		executionOrder := []string{}
		expectedError := errors.New("second seeder failed")

		manager.RegisterSeeder("first", func() error {
			executionOrder = append(executionOrder, "first")
			return nil
		})
		manager.RegisterSeeder("second", func() error {
			executionOrder = append(executionOrder, "second")
			return expectedError
		})
		manager.RegisterSeeder("third", func() error {
			executionOrder = append(executionOrder, "third")
			return nil
		})

		order := []string{"first", "second", "third"}
		err := manager.RunSeedersInOrder(order)

		assert.Error(t, err)
		assert.Equal(t, []string{"first", "second"}, executionOrder) // third should not execute
	})

	t.Run("Run empty order list", func(t *testing.T) {
		manager := NewSeederManager()

		err := manager.RunSeedersInOrder([]string{})

		assert.NoError(t, err)
	})

	t.Run("Run order with non-existing seeder", func(t *testing.T) {
		manager := NewSeederManager()

		manager.RegisterSeeder("existing", func() error { return nil })

		order := []string{"existing", "non_existing"}
		err := manager.RunSeedersInOrder(order)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

// TestRunAllSeeders tests the RunAllSeeders method
func TestRunAllSeeders(t *testing.T) {
	t.Run("Run all seeders successfully", func(t *testing.T) {
		manager := NewSeederManager()
		executionOrder := []string{}

		manager.RegisterSeeder("seeder1", func() error {
			executionOrder = append(executionOrder, "seeder1")
			return nil
		})
		manager.RegisterSeeder("seeder2", func() error {
			executionOrder = append(executionOrder, "seeder2")
			return nil
		})
		manager.RegisterSeeder("seeder3", func() error {
			executionOrder = append(executionOrder, "seeder3")
			return nil
		})

		err := manager.RunAllSeeders()

		assert.NoError(t, err)
		assert.Equal(t, []string{"seeder1", "seeder2", "seeder3"}, executionOrder)
	})

	t.Run("Run all seeders with one failing", func(t *testing.T) {
		manager := NewSeederManager()
		executionOrder := []string{}
		expectedError := errors.New("seeder2 failed")

		manager.RegisterSeeder("seeder1", func() error {
			executionOrder = append(executionOrder, "seeder1")
			return nil
		})
		manager.RegisterSeeder("seeder2", func() error {
			executionOrder = append(executionOrder, "seeder2")
			return expectedError
		})
		manager.RegisterSeeder("seeder3", func() error {
			executionOrder = append(executionOrder, "seeder3")
			return nil
		})

		err := manager.RunAllSeeders()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "seeder2")
		assert.Equal(t, []string{"seeder1", "seeder2"}, executionOrder) // seeder3 should not execute
	})

	t.Run("Run all seeders when none registered", func(t *testing.T) {
		manager := NewSeederManager()

		err := manager.RunAllSeeders()

		assert.NoError(t, err)
	})
}

// TestIsSeederRegistered tests the IsSeederRegistered method
func TestIsSeederRegistered(t *testing.T) {
	t.Run("Check registered seeder", func(t *testing.T) {
		manager := NewSeederManager()

		manager.RegisterSeeder("test_seeder", func() error { return nil })

		assert.True(t, manager.IsSeederRegistered("test_seeder"))
	})

	t.Run("Check non-registered seeder", func(t *testing.T) {
		manager := NewSeederManager()

		assert.False(t, manager.IsSeederRegistered("non_existing"))
	})

	t.Run("Check empty name", func(t *testing.T) {
		manager := NewSeederManager()

		assert.False(t, manager.IsSeederRegistered(""))
	})
}

// TestSeederManagerIntegration tests integration scenarios
func TestSeederManagerIntegration(t *testing.T) {
	t.Run("Complete workflow", func(t *testing.T) {
		manager := NewSeederManager()
		executionLog := []string{}

		// Register multiple seeders
		seeders := []SeederItem{
			{
				Name: "users",
				Function: func() error {
					executionLog = append(executionLog, "users")
					return nil
				},
			},
			{
				Name: "departments",
				Function: func() error {
					executionLog = append(executionLog, "departments")
					return nil
				},
			},
			{
				Name: "roles",
				Function: func() error {
					executionLog = append(executionLog, "roles")
					return nil
				},
			},
		}

		// Register all seeders
		err := manager.RegisterSeeders(seeders...)
		assert.NoError(t, err)

		// Verify registration
		registered := manager.GetRegisteredSeeders()
		assert.Len(t, registered, 3)
		assert.True(t, manager.IsSeederRegistered("users"))
		assert.True(t, manager.IsSeederRegistered("departments"))
		assert.True(t, manager.IsSeederRegistered("roles"))

		// Run all seeders
		err = manager.RunAllSeeders()
		assert.NoError(t, err)

		// Verify execution order
		assert.Equal(t, []string{"users", "departments", "roles"}, executionLog)

		// Run specific seeder
		executionLog = []string{}
		err = manager.RunSeederByName("departments")
		assert.NoError(t, err)
		assert.Equal(t, []string{"departments"}, executionLog)

		// Run in custom order
		executionLog = []string{}
		customOrder := []string{"roles", "users", "departments"}
		err = manager.RunSeedersInOrder(customOrder)
		assert.NoError(t, err)
		assert.Equal(t, []string{"roles", "users", "departments"}, executionLog)
	})
}
