# Quick Start

This guide will help you get started with axon quickly with a basic example.

## Basic Usage

Here's a simple example to get you started:

```go
package main

import (
    "fmt"
    "log"
    "github.com/kolosys/axon"
)

func main() {
    // Basic usage example
    fmt.Println("Welcome to axon!")
    
    // TODO: Add your code here
}
```

## Common Use Cases

### Using axon

**Import Path:** `github.com/kolosys/axon`



```go
package main

import (
    "fmt"
    "github.com/kolosys/axon"
)

func main() {
    // Example usage of axon
    fmt.Println("Using axon package")
}
```

#### Available Types
- **Conn** - Conn represents a WebSocket connection with type-safe message handling
- **Frame** - Frame represents a WebSocket frame
- **Metrics** - Metrics tracks WebSocket connection metrics
- **MetricsSnapshot** - MetricsSnapshot represents a snapshot of metrics at a point in time
- **UpgradeOptions** - UpgradeOptions configures WebSocket connection upgrade options
- **Upgrader** - Upgrader handles WebSocket connection upgrades

For detailed API documentation, see the [axon API Reference](../api-reference/axon.md).

## Step-by-Step Tutorial

### Step 1: Import the Package

First, import the necessary packages in your Go file:

```go
import (
    "fmt"
    "github.com/kolosys/axon"
)
```

### Step 2: Initialize

Set up the basic configuration:

```go
func main() {
    // Initialize your application
    fmt.Println("Initializing axon...")
}
```

### Step 3: Use the Library

Implement your specific use case:

```go
func main() {
    // Your implementation here
}
```

## Running Your Code

To run your Go program:

```bash
go run main.go
```

To build an executable:

```bash
go build -o myapp
./myapp
```

## Configuration Options

axon can be configured to suit your needs. Check the [Core Concepts](../core-concepts/) section for detailed information about configuration options.

## Error Handling

Always handle errors appropriately:

```go
result, err := someFunction()
if err != nil {
    log.Fatalf("Error: %v", err)
}
```

## Best Practices

- Always handle errors returned by library functions
- Check the API documentation for detailed parameter information
- Use meaningful variable and function names
- Add comments to document your code

## Complete Example

Here's a complete working example:

```go
package main

import (
    "fmt"
    "log"
    "github.com/kolosys/axon"
)

func main() {
    fmt.Println("Starting axon application...")
    
    // Add your implementation here
    
    fmt.Println("Application completed successfully!")
}
```

## Next Steps

Now that you've seen the basics, explore:

- **[Core Concepts](../core-concepts/)** - Understanding the library architecture
- **[API Reference](../api-reference/)** - Complete API documentation
- **[Examples](../examples/README.md)** - More detailed examples
- **[Advanced Topics](../advanced/)** - Performance tuning and advanced patterns

## Getting Help

If you run into issues:

1. Check the [API Reference](../api-reference/)
2. Browse the [Examples](../examples/README.md)
3. Visit the [GitHub Issues](https://github.com/kolosys/axon/issues) page

