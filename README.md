# Go Seeder

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.19-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

A flexible and extensible seeder package for Go applications. Perfect for seeding your application with initial data, test data, or sample data.

## ✨ Features

- 🚀 **Flexible Registration**: Register individual seeders or multiple seeders at once
- 🖥️ **CLI Support**: Built-in command line interface with customizable app names
- 📚 **Library Mode**: Use as a library without CLI for programmatic seeding
- ✅ **Validation**: Automatic validation of seeder names and uniqueness
- 🛡️ **Error Handling**: Comprehensive error handling and logging
- 📋 **Execution Order**: Support for running seeders in specific order
- 📦 **Zero Dependencies**: No external dependencies

## 📦 Installation

```bash
go get github.com/risoftinc/goseeder
```

## 🚀 Quick Start

### 1. Basic Usage (Library Mode)

```go
package main

import (
    "log"
    "github.com/risoftinc/goseeder"
)

func main() {
    // Create seeder manager
    manager := goseeder.NewSeederManager()
    
    // Register a seeder
    manager.RegisterSeeder("users", func() error {
        log.Println("Seeding users...")
        // Your seeding logic here
        return nil
    })
    
    // Run the seeder
    if err := manager.RunSeederByName("users"); err != nil {
        log.Fatal(err)
    }
}
```

### 2. CLI Mode

```go
package main

import (
    "log"
    "github.com/risoftinc/goseeder"
)

func main() {
    // Create seeder manager
    manager := goseeder.NewSeederManager()
    
    // Register seeders
    manager.RegisterSeeder("users", func() error {
        log.Println("Seeding users...")
        return nil
    })
    
    // Create CLI with custom app name
    cli := goseeder.NewCLIWithAppName(manager, "my-app seeder")
    
    // Run CLI (parses command line arguments)
    if err := cli.Run(); err != nil {
        log.Fatal(err)
    }
}
```

## 🖥️ Command Line Usage

When using CLI mode, you can run seeders with the following commands:

```bash
# Show help and available seeders
./your-app

# Run all seeders
./your-app -type=all

# Run specific seeder
./your-app -type=users
```

### Example Output

```
============================================================
DATABASE SEEDER - MY-APP SEEDER
============================================================

Usage:
  my-app seeder -type=all     # Run all seeders
  my-app seeder -type=<name>  # Run specific seeder
  my-app seeder               # Show this help

Available seeders (in execution order):
----------------------------------------
  1. users
     Command: my-app seeder -type=users

  2. departments
     Command: my-app seeder -type=departments

Quick commands:
  my-app seeder -type=all     # Run all seeders
============================================================
```

## 📚 API Reference

### SeederManager

#### `NewSeederManager() *SeederManager`
Creates a new seeder manager instance.

#### `RegisterSeeder(name string, function func() error) error`
Registers a single seeder with validation for unique names.

**Parameters:**
- `name`: Unique name for the seeder
- `function`: Function that performs the seeding operation

**Returns:**
- `error`: Returns error if name is empty or already exists

#### `RegisterSeeders(seeders ...SeederItem) error`
Registers multiple seeders at once using variadic function.

**Parameters:**
- `seeders`: Variadic list of SeederItem structs

**Returns:**
- `error`: Returns error if any seeder registration fails

#### `RunSeederByName(name string) error`
Runs a specific seeder by name.

**Parameters:**
- `name`: Name of the seeder to run

**Returns:**
- `error`: Returns error if seeder not found or execution fails

#### `RunAllSeeders() error`
Runs all registered seeders in registration order.

**Returns:**
- `error`: Returns error if any seeder execution fails

#### `RunSeedersInOrder(names []string) error`
Runs multiple seeders in the specified order.

**Parameters:**
- `names`: Slice of seeder names in execution order

**Returns:**
- `error`: Returns error if any seeder execution fails

#### `GetRegisteredSeeders() []string`
Returns a list of all registered seeder names.

**Returns:**
- `[]string`: Slice of registered seeder names

#### `IsSeederRegistered(name string) bool`
Checks if a seeder with the given name is registered.

**Parameters:**
- `name`: Name of the seeder to check

**Returns:**
- `bool`: True if seeder is registered, false otherwise

### CLI

#### `NewCLI(manager *SeederManager) *CLI`
Creates a new CLI instance with default app name "seeder".

#### `NewCLIWithAppName(manager *SeederManager, appName string) *CLI`
Creates a new CLI instance with custom app name.

**Parameters:**
- `manager`: SeederManager instance
- `appName`: Custom name for the application (used in help text)

#### `Run() error`
Executes the seeder based on command line arguments.

**Returns:**
- `error`: Returns error if execution fails

#### `Usage()`
Prints usage information and available seeders.

### SeederItem

```go
type SeederItem struct {
    Name     string
    Function func() error
}
```

Represents a single seeder with its name and function.

## 🔧 Advanced Examples

### Variadic Registration

```go
seeders := []goseeder.SeederItem{
    {Name: "users", Function: func() error {
        log.Println("Seeding users...")
        return nil
    }},
    {Name: "departments", Function: func() error {
        log.Println("Seeding departments...")
        return nil
    }},
}

if err := manager.RegisterSeeders(seeders...); err != nil {
    log.Fatal(err)
}
```

### Custom App Name for CLI

```go
// For a web application
cli := seeder.NewCLIWithAppName(manager, "go run cmd/seeder/main.go")

// For a compiled binary
cli := seeder.NewCLIWithAppName(manager, "./my-app seeder")

// For a Docker container
cli := seeder.NewCLIWithAppName(manager, "docker run my-app seeder")
```

### Real-world Example

```go
package main

import (
    "log"
    "github.com/risoftinc/goseeder"
)

type User struct {
    ID    uint
    Name  string
    Email string
}

func main() {
    // Create seeder manager
    manager := goseeder.NewSeederManager()
    
    // Register user seeder
    manager.RegisterSeeder("users", func() error {
        users := []User{
            {Name: "John Doe", Email: "john@example.com"},
            {Name: "Jane Smith", Email: "jane@example.com"},
        }
        
        // Your seeding logic here
        // For example, save to database, create files, etc.
        for _, user := range users {
            log.Printf("Creating user: %s (%s)", user.Name, user.Email)
            // Add your logic to save the user
        }
        
        log.Printf("Created %d users", len(users))
        return nil
    })
    
    // Run the seeder
    if err := manager.RunSeederByName("users"); err != nil {
        log.Fatal(err)
    }
}
```

## 🏗️ Project Structure

```
your-project/
├── cmd/
│   └── seeder/
│       └── main.go          # CLI entry point
├── database/
│   └── seeders/
│       ├── main_seeder.go   # Main seeder registration
│       ├── user_seeder.go   # User seeder implementation
│       └── department_seeder.go # Department seeder implementation
└── go.mod
```

## 🎯 Best Practices

### 1. **Naming Convention**
Use descriptive names for seeders:
```go
// Good
manager.RegisterSeeder("users", userSeeder)
manager.RegisterSeeder("departments", departmentSeeder)

// Avoid
manager.RegisterSeeder("seed1", seeder1)
manager.RegisterSeeder("data", dataSeeder)
```

### 2. **Error Handling**
Always handle errors returned by seeder functions:
```go
if err := manager.RunSeederByName("users"); err != nil {
    log.Printf("Failed to seed users: %v", err)
    return err
}
```

### 3. **Dependencies**
Consider execution order when seeders have dependencies:
```go
// Run departments first, then users (users might reference departments)
order := []string{"departments", "users"}
if err := manager.RunSeedersInOrder(order); err != nil {
    log.Fatal(err)
}
```

### 4. **Idempotency**
Make seeders idempotent (safe to run multiple times):
```go
manager.RegisterSeeder("users", func() error {
    // Check if users already exist
    // Add your logic to check if data already exists
    if usersAlreadyExist() {
        log.Println("Users already exist, skipping...")
        return nil
    }
    
    // Create users only if they don't exist
    return createUsers()
})
```

### 5. **Logging**
Use consistent logging for better debugging:
```go
manager.RegisterSeeder("users", func() error {
    log.Println("Starting user seeding...")
    
    // Your seeding logic
    
    log.Println("User seeding completed successfully")
    return nil
})
```

## 🔍 Testing

```go
func TestSeeder(t *testing.T) {
    // Create seeder manager
    manager := goseeder.NewSeederManager()
    
    // Register test seeder
    manager.RegisterSeeder("test", func() error {
        // Test seeding logic
        return nil
    })
    
    // Run seeder
    if err := manager.RunSeederByName("test"); err != nil {
        t.Fatalf("Seeder failed: %v", err)
    }
    
    // Verify results
    // Your assertions here
}
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Go](https://golang.org/) - The amazing programming language

## 📞 Support

If you have any questions or need help, please:

1. Check the [Issues](https://github.com/risoftinc/seeder/issues) page
2. Create a new issue if your question isn't already answered
3. Join our community discussions

---

**Made with ❤️ by [Risoft Inc](https://github.com/risoftinc)**
