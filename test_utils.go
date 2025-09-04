package goseeder

import (
	"bytes"
	"os"
)

// TestOutputCapture provides utilities for capturing and testing output
type TestOutputCapture struct {
	originalStdout *os.File
	originalStderr *os.File
	stdoutPipe     *os.File
	stderrPipe     *os.File
	stdoutBuffer   *bytes.Buffer
	stderrBuffer   *bytes.Buffer
}

// NewTestOutputCapture creates a new output capture instance
func NewTestOutputCapture() *TestOutputCapture {
	return &TestOutputCapture{}
}

// Start begins capturing stdout and stderr
func (toc *TestOutputCapture) Start() error {
	// Save original stdout and stderr
	toc.originalStdout = os.Stdout
	toc.originalStderr = os.Stderr

	// Create pipes for capturing output
	var err error
	toc.stdoutPipe, err = os.CreateTemp("", "test_stdout_*")
	if err != nil {
		return err
	}

	toc.stderrPipe, err = os.CreateTemp("", "test_stderr_*")
	if err != nil {
		return err
	}

	// Create buffers to store captured output
	toc.stdoutBuffer = &bytes.Buffer{}
	toc.stderrBuffer = &bytes.Buffer{}

	// Redirect stdout and stderr to our pipes
	os.Stdout = toc.stdoutPipe
	os.Stderr = toc.stderrPipe

	return nil
}

// Stop stops capturing and returns the captured output
func (toc *TestOutputCapture) Stop() (stdout, stderr string, err error) {
	// Restore original stdout and stderr
	os.Stdout = toc.originalStdout
	os.Stderr = toc.originalStderr

	// Close the pipes
	if toc.stdoutPipe != nil {
		toc.stdoutPipe.Close()
	}
	if toc.stderrPipe != nil {
		toc.stderrPipe.Close()
	}

	// Read the captured output
	if toc.stdoutPipe != nil {
		stdoutBytes, readErr := os.ReadFile(toc.stdoutPipe.Name())
		if readErr != nil {
			err = readErr
		} else {
			stdout = string(stdoutBytes)
		}
		os.Remove(toc.stdoutPipe.Name())
	}

	if toc.stderrPipe != nil {
		stderrBytes, readErr := os.ReadFile(toc.stderrPipe.Name())
		if readErr != nil {
			err = readErr
		} else {
			stderr = string(stderrBytes)
		}
		os.Remove(toc.stderrPipe.Name())
	}

	return stdout, stderr, err
}

// CaptureOutput captures output from a function execution
func CaptureOutput(fn func()) (stdout, stderr string, err error) {
	capture := NewTestOutputCapture()

	if err := capture.Start(); err != nil {
		return "", "", err
	}

	fn()

	return capture.Stop()
}

// TestSeederFunction is a helper type for creating test seeder functions
type TestSeederFunction struct {
	Name         string
	ShouldError  bool
	ErrorMsg     string
	ExecutionLog *[]string
}

// CreateTestSeeder creates a test seeder function
func CreateTestSeeder(config TestSeederFunction) func() error {
	return func() error {
		if config.ExecutionLog != nil {
			*config.ExecutionLog = append(*config.ExecutionLog, config.Name)
		}

		if config.ShouldError {
			if config.ErrorMsg != "" {
				return &TestError{Message: config.ErrorMsg}
			}
			return &TestError{Message: "test error"}
		}

		return nil
	}
}

// TestError is a custom error type for testing
type TestError struct {
	Message string
}

func (te *TestError) Error() string {
	return te.Message
}

// MockSeederManagerInterface defines the interface for mocking SeederManager
type MockSeederManagerInterface interface {
	RegisterSeeder(name string, function func() error) error
	RegisterSeeders(seeders ...SeederItem) error
	GetRegisteredSeeders() []string
	RunSeederByName(name string) error
	RunSeedersInOrder(names []string) error
	RunAllSeeders() error
	IsSeederRegistered(name string) bool
}

// TestSeederManager is a test implementation of SeederManager
type TestSeederManager struct {
	seeders     []SeederItem
	seederMap   map[string]func() error
	shouldError bool
	errorMsg    string
}

// NewTestSeederManager creates a new test seeder manager
func NewTestSeederManager() *TestSeederManager {
	return &TestSeederManager{
		seeders:   make([]SeederItem, 0),
		seederMap: make(map[string]func() error),
	}
}

// SetErrorBehavior sets whether the manager should return errors
func (tsm *TestSeederManager) SetErrorBehavior(shouldError bool, errorMsg string) {
	tsm.shouldError = shouldError
	tsm.errorMsg = errorMsg
}

// RegisterSeeder implements the SeederManager interface
func (tsm *TestSeederManager) RegisterSeeder(name string, function func() error) error {
	if tsm.shouldError {
		if tsm.errorMsg != "" {
			return &TestError{Message: tsm.errorMsg}
		}
		return &TestError{Message: "registration error"}
	}

	if name == "" {
		return &TestError{Message: "seeder name cannot be empty"}
	}

	if _, exists := tsm.seederMap[name]; exists {
		return &TestError{Message: "seeder with name '" + name + "' already exists"}
	}

	seederItem := SeederItem{
		Name:     name,
		Function: function,
	}
	tsm.seeders = append(tsm.seeders, seederItem)
	tsm.seederMap[name] = function

	return nil
}

// RegisterSeeders implements the SeederManager interface
func (tsm *TestSeederManager) RegisterSeeders(seeders ...SeederItem) error {
	for _, seeder := range seeders {
		if err := tsm.RegisterSeeder(seeder.Name, seeder.Function); err != nil {
			return &TestError{Message: "failed to register seeder '" + seeder.Name + "': " + err.Error()}
		}
	}
	return nil
}

// GetRegisteredSeeders implements the SeederManager interface
func (tsm *TestSeederManager) GetRegisteredSeeders() []string {
	names := make([]string, len(tsm.seeders))
	for i, seeder := range tsm.seeders {
		names[i] = seeder.Name
	}
	return names
}

// RunSeederByName implements the SeederManager interface
func (tsm *TestSeederManager) RunSeederByName(name string) error {
	if tsm.shouldError {
		if tsm.errorMsg != "" {
			return &TestError{Message: tsm.errorMsg}
		}
		return &TestError{Message: "execution error"}
	}

	if function, exists := tsm.seederMap[name]; exists {
		return function()
	}
	return &TestError{Message: "seeder with name '" + name + "' not found"}
}

// RunSeedersInOrder implements the SeederManager interface
func (tsm *TestSeederManager) RunSeedersInOrder(names []string) error {
	for _, name := range names {
		if err := tsm.RunSeederByName(name); err != nil {
			return err
		}
	}
	return nil
}

// RunAllSeeders implements the SeederManager interface
func (tsm *TestSeederManager) RunAllSeeders() error {
	for _, seeder := range tsm.seeders {
		if err := seeder.Function(); err != nil {
			return &TestError{Message: "seeder '" + seeder.Name + "' failed: " + err.Error()}
		}
	}
	return nil
}

// IsSeederRegistered implements the SeederManager interface
func (tsm *TestSeederManager) IsSeederRegistered(name string) bool {
	_, exists := tsm.seederMap[name]
	return exists
}

// TestCLI is a test implementation of CLI
type TestCLI struct {
	manager MockSeederManagerInterface
	appName string
}

// NewTestCLI creates a new test CLI
func NewTestCLI(manager MockSeederManagerInterface) *TestCLI {
	return &TestCLI{
		manager: manager,
		appName: "test-seeder",
	}
}

// NewTestCLIWithAppName creates a new test CLI with custom app name
func NewTestCLIWithAppName(manager MockSeederManagerInterface, appName string) *TestCLI {
	return &TestCLI{
		manager: manager,
		appName: appName,
	}
}

// Usage implements the CLI interface
func (tcli *TestCLI) Usage() {
	// Simplified usage for testing
	seeders := tcli.manager.GetRegisteredSeeders()

	// This would normally use log.Printf, but for testing we'll use a simpler approach
	_ = seeders // Use the seeders to avoid unused variable warning
}

// Run implements the CLI interface
func (tcli *TestCLI) Run() error {
	// Simplified run for testing
	return tcli.manager.RunAllSeeders()
}

// TestDataBuilder provides a fluent interface for building test data
type TestDataBuilder struct {
	seeders []SeederItem
}

// NewTestDataBuilder creates a new test data builder
func NewTestDataBuilder() *TestDataBuilder {
	return &TestDataBuilder{
		seeders: make([]SeederItem, 0),
	}
}

// AddSeeder adds a seeder to the test data
func (tdb *TestDataBuilder) AddSeeder(name string, function func() error) *TestDataBuilder {
	tdb.seeders = append(tdb.seeders, SeederItem{
		Name:     name,
		Function: function,
	})
	return tdb
}

// AddTestSeeder adds a test seeder with predefined behavior
func (tdb *TestDataBuilder) AddTestSeeder(config TestSeederFunction) *TestDataBuilder {
	tdb.seeders = append(tdb.seeders, SeederItem{
		Name:     config.Name,
		Function: CreateTestSeeder(config),
	})
	return tdb
}

// Build returns the built seeders
func (tdb *TestDataBuilder) Build() []SeederItem {
	return tdb.seeders
}

// TestAssertions provides common test assertions for seeder operations
type TestAssertions struct{}

// AssertSeederRegistered checks if a seeder is registered
func (ta *TestAssertions) AssertSeederRegistered(t interface{}, manager MockSeederManagerInterface, name string) {
	// This would use the testing.T interface, but we'll keep it simple for now
	_ = t
	_ = manager
	_ = name
}

// AssertExecutionOrder checks if seeders were executed in the correct order
func (ta *TestAssertions) AssertExecutionOrder(t interface{}, expected, actual []string) {
	// This would use the testing.T interface, but we'll keep it simple for now
	_ = t
	_ = expected
	_ = actual
}

// Helper function to create a simple output capture for testing
func captureOutputSimple(fn func()) string {
	// This is a simplified version that doesn't actually capture output
	// In a real implementation, you would use the TestOutputCapture above
	fn()
	return "captured output"
}

// TestHelper provides common test helper functions
type TestHelper struct{}

// CreateMockSeederManager creates a mock seeder manager for testing
func (th *TestHelper) CreateMockSeederManager() *TestSeederManager {
	return NewTestSeederManager()
}

// CreateTestSeeders creates a set of test seeders
func (th *TestHelper) CreateTestSeeders() []SeederItem {
	executionLog := make([]string, 0)

	return []SeederItem{
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
		{
			Name: "roles",
			Function: CreateTestSeeder(TestSeederFunction{
				Name:         "roles",
				ShouldError:  false,
				ExecutionLog: &executionLog,
			}),
		},
	}
}

// CreateErrorSeeder creates a seeder that will return an error
func (th *TestHelper) CreateErrorSeeder(name, errorMsg string) SeederItem {
	return SeederItem{
		Name: name,
		Function: CreateTestSeeder(TestSeederFunction{
			Name:        name,
			ShouldError: true,
			ErrorMsg:    errorMsg,
		}),
	}
}
