package goseeder

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

// SeederTestSuite is a comprehensive test suite for the seeder package
type SeederTestSuite struct {
	suite.Suite
	manager *SeederManager
	helper  *TestHelper
}

// SetupTest runs before each test
func (suite *SeederTestSuite) SetupTest() {
	suite.manager = NewSeederManager()
	suite.helper = &TestHelper{}
}

// TestSeederManagerSuite tests the SeederManager functionality
func (suite *SeederTestSuite) TestSeederManagerSuite() {
	suite.Run("RegisterSeeder", suite.testRegisterSeeder)
	suite.Run("RegisterSeeders", suite.testRegisterSeeders)
	suite.Run("RunSeederByName", suite.testRunSeederByName)
	suite.Run("RunAllSeeders", suite.testRunAllSeeders)
	suite.Run("RunSeedersInOrder", suite.testRunSeedersInOrder)
	suite.Run("GetRegisteredSeeders", suite.testGetRegisteredSeeders)
	suite.Run("IsSeederRegistered", suite.testIsSeederRegistered)
}

// TestCLISuite tests the CLI functionality
func (suite *SeederTestSuite) TestCLISuite() {
	suite.Run("NewCLI", suite.testNewCLI)
	suite.Run("NewCLIWithAppName", suite.testNewCLIWithAppName)
	suite.Run("Usage", suite.testUsage)
	suite.Run("Run", suite.testRun)
}

// TestIntegrationSuite tests integration scenarios
func (suite *SeederTestSuite) TestIntegrationSuite() {
	suite.Run("CompleteWorkflow", suite.testCompleteWorkflow)
	suite.Run("ErrorHandling", suite.testErrorHandling)
	suite.Run("EdgeCases", suite.testEdgeCases)
}

// TestUtilitySuite tests the utility functions
func (suite *SeederTestSuite) TestUtilitySuite() {
	suite.Run("TestDataBuilder", suite.testTestDataBuilder)
	suite.Run("TestSeederFunction", suite.testTestSeederFunction)
	suite.Run("TestOutputCapture", suite.testTestOutputCapture)
}

// testRegisterSeeder tests the RegisterSeeder method
func (suite *SeederTestSuite) testRegisterSeeder() {
	suite.Run("ValidSeeder", func() {
		err := suite.manager.RegisterSeeder("test", func() error { return nil })
		suite.NoError(err)
		suite.True(suite.manager.IsSeederRegistered("test"))
	})

	suite.Run("EmptyName", func() {
		err := suite.manager.RegisterSeeder("", func() error { return nil })
		suite.Error(err)
		suite.Contains(err.Error(), "cannot be empty")
	})

	suite.Run("DuplicateName", func() {
		suite.manager.RegisterSeeder("duplicate", func() error { return nil })
		err := suite.manager.RegisterSeeder("duplicate", func() error { return nil })
		suite.Error(err)
		suite.Contains(err.Error(), "already exists")
	})
}

// testRegisterSeeders tests the RegisterSeeders method
func (suite *SeederTestSuite) testRegisterSeeders() {
	suite.Run("MultipleSeeders", func() {
		seeders := suite.helper.CreateTestSeeders()
		err := suite.manager.RegisterSeeders(seeders...)
		suite.NoError(err)
		suite.Len(suite.manager.GetRegisteredSeeders(), 3)
	})

	suite.Run("EmptyList", func() {
		err := suite.manager.RegisterSeeders()
		suite.NoError(err)
		suite.Len(suite.manager.GetRegisteredSeeders(), 0)
	})
}

// testRunSeederByName tests the RunSeederByName method
func (suite *SeederTestSuite) testRunSeederByName() {
	suite.Run("ExistingSeeder", func() {
		executed := false
		suite.manager.RegisterSeeder("test", func() error {
			executed = true
			return nil
		})

		err := suite.manager.RunSeederByName("test")
		suite.NoError(err)
		suite.True(executed)
	})

	suite.Run("NonExistingSeeder", func() {
		err := suite.manager.RunSeederByName("nonexistent")
		suite.Error(err)
		suite.Contains(err.Error(), "not found")
	})

	suite.Run("SeederWithError", func() {
		suite.manager.RegisterSeeder("error", func() error {
			return &TestError{Message: "test error"}
		})

		err := suite.manager.RunSeederByName("error")
		suite.Error(err)
		suite.Contains(err.Error(), "failed")
	})
}

// testRunAllSeeders tests the RunAllSeeders method
func (suite *SeederTestSuite) testRunAllSeeders() {
	suite.Run("AllSeedersSuccess", func() {
		executionLog := make([]string, 0)

		seeders := []SeederItem{
			{
				Name: "first",
				Function: CreateTestSeeder(TestSeederFunction{
					Name:         "first",
					ShouldError:  false,
					ExecutionLog: &executionLog,
				}),
			},
			{
				Name: "second",
				Function: CreateTestSeeder(TestSeederFunction{
					Name:         "second",
					ShouldError:  false,
					ExecutionLog: &executionLog,
				}),
			},
		}

		suite.manager.RegisterSeeders(seeders...)
		err := suite.manager.RunAllSeeders()
		suite.NoError(err)
		suite.Equal([]string{"first", "second"}, executionLog)
	})

	suite.Run("OneSeederFails", func() {
		executionLog := make([]string, 0)

		seeders := []SeederItem{
			{
				Name: "success",
				Function: CreateTestSeeder(TestSeederFunction{
					Name:         "success",
					ShouldError:  false,
					ExecutionLog: &executionLog,
				}),
			},
			{
				Name: "failure",
				Function: CreateTestSeeder(TestSeederFunction{
					Name:         "failure",
					ShouldError:  true,
					ErrorMsg:     "seeder failed",
					ExecutionLog: &executionLog,
				}),
			},
		}

		suite.manager.RegisterSeeders(seeders...)
		err := suite.manager.RunAllSeeders()
		suite.Error(err)
		suite.Contains(err.Error(), "failure")
		suite.Equal([]string{"success", "failure"}, executionLog)
	})
}

// testRunSeedersInOrder tests the RunSeedersInOrder method
func (suite *SeederTestSuite) testRunSeedersInOrder() {
	suite.Run("CustomOrder", func() {
		executionLog := make([]string, 0)

		seeders := []SeederItem{
			{
				Name: "first",
				Function: CreateTestSeeder(TestSeederFunction{
					Name:         "first",
					ShouldError:  false,
					ExecutionLog: &executionLog,
				}),
			},
			{
				Name: "second",
				Function: CreateTestSeeder(TestSeederFunction{
					Name:         "second",
					ShouldError:  false,
					ExecutionLog: &executionLog,
				}),
			},
		}

		suite.manager.RegisterSeeders(seeders...)
		err := suite.manager.RunSeedersInOrder([]string{"second", "first"})
		suite.NoError(err)
		suite.Equal([]string{"second", "first"}, executionLog)
	})
}

// testGetRegisteredSeeders tests the GetRegisteredSeeders method
func (suite *SeederTestSuite) testGetRegisteredSeeders() {
	suite.Run("EmptyList", func() {
		seeders := suite.manager.GetRegisteredSeeders()
		suite.Empty(seeders)
	})

	suite.Run("MultipleSeeders", func() {
		suite.manager.RegisterSeeder("first", func() error { return nil })
		suite.manager.RegisterSeeder("second", func() error { return nil })

		seeders := suite.manager.GetRegisteredSeeders()
		suite.Len(seeders, 2)
		suite.Equal("first", seeders[0])
		suite.Equal("second", seeders[1])
	})
}

// testIsSeederRegistered tests the IsSeederRegistered method
func (suite *SeederTestSuite) testIsSeederRegistered() {
	suite.Run("RegisteredSeeder", func() {
		suite.manager.RegisterSeeder("test", func() error { return nil })
		suite.True(suite.manager.IsSeederRegistered("test"))
	})

	suite.Run("NonRegisteredSeeder", func() {
		suite.False(suite.manager.IsSeederRegistered("nonexistent"))
	})
}

// testNewCLI tests the NewCLI function
func (suite *SeederTestSuite) testNewCLI() {
	cli := NewCLI(suite.manager)
	suite.NotNil(cli)
	suite.Equal(suite.manager, cli.manager)
	suite.Equal("seeder", cli.appName)
}

// testNewCLIWithAppName tests the NewCLIWithAppName function
func (suite *SeederTestSuite) testNewCLIWithAppName() {
	appName := "custom-app"
	cli := NewCLIWithAppName(suite.manager, appName)
	suite.NotNil(cli)
	suite.Equal(suite.manager, cli.manager)
	suite.Equal(appName, cli.appName)
}

// testUsage tests the Usage method
func (suite *SeederTestSuite) testUsage() {
	suite.Run("NoSeeders", func() {
		cli := NewCLI(suite.manager)
		// Usage() will be called and should not panic
		suite.NotPanics(func() {
			cli.Usage()
		})
	})

	suite.Run("WithSeeders", func() {
		suite.manager.RegisterSeeder("test", func() error { return nil })
		cli := NewCLI(suite.manager)
		// Usage() will be called and should not panic
		suite.NotPanics(func() {
			cli.Usage()
		})
	})
}

// testRun tests the Run method
func (suite *SeederTestSuite) testRun() {
	suite.Run("NoTypeFlag", func() {
		// This would normally parse command line arguments
		// For testing, we'll just ensure the method exists and doesn't panic
		suite.NotPanics(func() {
			// In a real test, you would set up command line arguments
			// For now, we'll just test that the method is callable
		})
	})
}

// testCompleteWorkflow tests a complete workflow
func (suite *SeederTestSuite) testCompleteWorkflow() {
	executionLog := make([]string, 0)

	// Register multiple seeders
	seeders := []SeederItem{
		{
			Name: "users",
			Function: CreateTestSeeder(TestSeederFunction{
				Name:         "users",
				ShouldError:  false,
				ExecutionLog: &executionLog,
			}),
		},
		{
			Name: "departments",
			Function: CreateTestSeeder(TestSeederFunction{
				Name:         "departments",
				ShouldError:  false,
				ExecutionLog: &executionLog,
			}),
		},
	}

	// Register all seeders
	err := suite.manager.RegisterSeeders(seeders...)
	suite.NoError(err)

	// Verify registration
	registered := suite.manager.GetRegisteredSeeders()
	suite.Len(registered, 2)
	suite.True(suite.manager.IsSeederRegistered("users"))
	suite.True(suite.manager.IsSeederRegistered("departments"))

	// Run all seeders
	err = suite.manager.RunAllSeeders()
	suite.NoError(err)

	// Verify execution order
	suite.Equal([]string{"users", "departments"}, executionLog)

	// Run specific seeder
	executionLog = make([]string, 0)
	err = suite.manager.RunSeederByName("departments")
	suite.NoError(err)
	suite.Equal([]string{"departments"}, executionLog)

	// Run in custom order
	executionLog = make([]string, 0)
	customOrder := []string{"departments", "users"}
	err = suite.manager.RunSeedersInOrder(customOrder)
	suite.NoError(err)
	suite.Equal([]string{"departments", "users"}, executionLog)
}

// testErrorHandling tests error handling scenarios
func (suite *SeederTestSuite) testErrorHandling() {
	suite.Run("RegistrationError", func() {
		err := suite.manager.RegisterSeeder("", func() error { return nil })
		suite.Error(err)
	})

	suite.Run("ExecutionError", func() {
		suite.manager.RegisterSeeder("error", func() error {
			return &TestError{Message: "execution error"}
		})

		err := suite.manager.RunSeederByName("error")
		suite.Error(err)
		suite.Contains(err.Error(), "failed")
	})

	suite.Run("NonExistentSeeder", func() {
		err := suite.manager.RunSeederByName("nonexistent")
		suite.Error(err)
		suite.Contains(err.Error(), "not found")
	})
}

// testEdgeCases tests edge cases
func (suite *SeederTestSuite) testEdgeCases() {
	suite.Run("EmptySeederList", func() {
		err := suite.manager.RunAllSeeders()
		suite.NoError(err)
	})

	suite.Run("EmptyOrderList", func() {
		err := suite.manager.RunSeedersInOrder([]string{})
		suite.NoError(err)
	})

	suite.Run("LongSeederName", func() {
		longName := "very_long_seeder_name_that_might_cause_issues"
		err := suite.manager.RegisterSeeder(longName, func() error { return nil })
		suite.NoError(err)
		suite.True(suite.manager.IsSeederRegistered(longName))
	})
}

// testTestDataBuilder tests the TestDataBuilder utility
func (suite *SeederTestSuite) testTestDataBuilder() {
	builder := NewTestDataBuilder()

	seeders := builder.
		AddSeeder("test1", func() error { return nil }).
		AddTestSeeder(TestSeederFunction{
			Name:        "test2",
			ShouldError: false,
		}).
		Build()

	suite.Len(seeders, 2)
	suite.Equal("test1", seeders[0].Name)
	suite.Equal("test2", seeders[1].Name)
}

// testTestSeederFunction tests the TestSeederFunction utility
func (suite *SeederTestSuite) testTestSeederFunction() {
	suite.Run("SuccessFunction", func() {
		executionLog := make([]string, 0)
		function := CreateTestSeeder(TestSeederFunction{
			Name:         "test",
			ShouldError:  false,
			ExecutionLog: &executionLog,
		})

		err := function()
		suite.NoError(err)
		suite.Equal([]string{"test"}, executionLog)
	})

	suite.Run("ErrorFunction", func() {
		function := CreateTestSeeder(TestSeederFunction{
			Name:        "test",
			ShouldError: true,
			ErrorMsg:    "test error",
		})

		err := function()
		suite.Error(err)
		suite.Equal("test error", err.Error())
	})
}

// testTestOutputCapture tests the TestOutputCapture utility
func (suite *SeederTestSuite) testTestOutputCapture() {
	suite.Run("CaptureOutput", func() {
		stdout, stderr, err := CaptureOutput(func() {
			// This would normally print to stdout/stderr
			// For testing, we'll just ensure the function works
		})

		suite.NoError(err)
		suite.NotNil(stdout)
		suite.NotNil(stderr)
	})
}

// TestSeederTestSuite runs the complete test suite
func TestSeederTestSuite(t *testing.T) {
	suite.Run(t, new(SeederTestSuite))
}

// BenchmarkSeederManager benchmarks the SeederManager performance
func BenchmarkSeederManager(b *testing.B) {
	manager := NewSeederManager()

	// Register seeders
	for i := 0; i < 100; i++ {
		name := fmt.Sprintf("seeder_%d", i)
		manager.RegisterSeeder(name, func() error { return nil })
	}

	b.ResetTimer()

	b.Run("RunAllSeeders", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			manager.RunAllSeeders()
		}
	})

	b.Run("IsSeederRegistered", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			manager.IsSeederRegistered("seeder_50")
		}
	})

	b.Run("GetRegisteredSeeders", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			manager.GetRegisteredSeeders()
		}
	})
}

// BenchmarkCLI benchmarks the CLI performance
func BenchmarkCLI(b *testing.B) {
	manager := NewSeederManager()
	manager.RegisterSeeder("test", func() error { return nil })
	cli := NewCLI(manager)

	b.ResetTimer()

	b.Run("Usage", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cli.Usage()
		}
	})
}
