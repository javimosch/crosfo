# AGENTS.md - Agent-First Go CLI Boilerplate

This document guides AI agents in understanding, extending, and maintaining this agent-first Go CLI boilerplate.

## Project Philosophy

This boilerplate implements **agent-first CLI design** following agent-friendly tools principles:

- **JSON-by-default**: All commands output JSON by default, even on TTY
- **`--human` opt-in**: Human-readable output only when explicitly requested
- **Semantic exit codes**: 0 (success), 80-89 (user errors), 90-99 (resource errors), 100-109 (integration errors), 110-119 (software errors)
- **Structured errors**: Error objects with code, type, recoverable field, and suggestions
- **Output separation**: stdout for data, stderr for logs/progress
- **No interactivity**: No prompts by default, `--no-interactive` is default behavior
- **Schema discovery**: `--schema` flag for JSON schema of each command output

## Project Structure

```
boilerplate-cli-ui-go/
├── main.go                  # CLI entry point and command routing (max 500 LOC)
├── server.go                # HTTP server and web UI (max 500 LOC)
├── daemon.go                # Process management (max 500 LOC)
├── go.mod                   # Go module dependencies
├── build.sh                 # Binary compilation script
├── README.md                # User documentation
└── AGENTS.md                # This file - project guide for agents
```

**Future Expansion Structure:**
```
├── cmd/                     # Command-specific packages
│   ├── greet.go             # Greet command implementation
│   ├── version.go           # Version command implementation
│   └── daemon/              # Daemon management commands
│       ├── start.go
│       ├── stop.go
│       └── status.go
├── pkg/                     # Reusable packages
│   ├── output/              # Output formatting (JSON/human)
│   │   └── formatter.go     # Output formatter (max 500 LOC)
│   ├── config/              # Configuration management
│   │   └── config.go        # Config loading (max 500 LOC)
│   ├── errors/              # Error definitions and handling
│   │   └── errors.go        # Semantic exit codes (max 500 LOC)
│   └── server/              # HTTP server
│       └── server.go        # Server implementation (max 500 LOC)
├── templates/               # Web UI frontend (if added)
│   └── index.html
├── schemas/                 # JSON schemas for command outputs
│   ├── greet.schema.json
│   ├── version.schema.json
│   └── status.schema.json
└── internal/                # Internal application code
    ├── daemon/              # Daemon process management
    │   └── daemon.go        (max 500 LOC)
    └── utils/               # Internal utilities
        └── utils.go         (max 500 LOC)
```

## Coding Rules

### File Size Limits

- **Max 500 LOC per Go file** - Split files that exceed this limit
- **Max 300 LOC per documentation file** - Keep documentation concise
- **Max 200 LOC per test file** - Split complex tests

### Module Organization

Each package has a single, well-defined responsibility:

- **cmd/**: Command implementations and CLI-specific logic
- **pkg/**: Reusable packages that can be imported by other projects
- **internal/**: Application-specific code not meant for external use
- **templates/**: Web UI frontend assets
- **schemas/**: JSON schemas for command output validation

### Naming Conventions

- **Go files**: `snake_case.go` for files
- **Functions**: `PascalCase` for exported, `camelCase` for unexported
- **Constants**: `PascalCase` for exported, `camelCase` for unexported
- **Interfaces**: `PascalCase` (often ending in -er)
- **Packages**: `lowercase` single word when possible
- **Errors**: `ErrPascalCase` for exported error variables

### Go-Specific Patterns

**Error Handling:**
```go
// Always handle errors explicitly
result, err := performOperation()
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}
```

**Interface Design:**
```go
// Keep interfaces small and focused
type Formatter interface {
    Output(data interface{}, exitCode int) error
    OutputError(errorData map[string]interface{}, exitCode int) error
}
```

**Package Structure:**
```go
// pkg/output/formatter.go
package output

type Formatter struct {
    humanMode bool
}

func NewFormatter(humanMode bool) *Formatter {
    return &Formatter{humanMode: humanMode}
}
```

### Agent-First Output Patterns

**Default JSON Output:**
```go
// Always output JSON by default
data := map[string]interface{}{
    "result": "success",
    "timestamp": time.Now().UTC().Format(time.RFC3339),
}
formatter.Output(data, EXIT_SUCCESS)
```

**Structured Errors:**
```go
// Use semantic exit codes and structured errors
if port < 1 || port > 65535 {
    return NewInvalidArgumentError(
        "Invalid port number",
        map[string]interface{}{
            "port": port,
            "valid_range": "1-65535",
        },
    )
}
```

**Output Separation:**
```go
// Data goes to stdout, logs to stderr
formatter.Output(data, exitCode)  // stdout
fmt.Fprintln(os.Stderr, "Processing...")  // stderr
```

### Semantic Exit Codes

Define semantic exit codes in a dedicated package:

```go
// pkg/errors/errors.go
package errors

const (
    EXIT_SUCCESS           = 0
    EXIT_INVALID_ARGUMENT  = 85
    EXIT_BAD_PERMISSIONS   = 86
    EXIT_RESOURCE_NOT_FOUND = 92
    EXIT_CONNECTION_TIMEOUT = 105
    EXIT_INTERNAL_ERROR    = 110
)

type CLIError struct {
    Code       int
    Type       string
    Message    string
    Details    map[string]interface{}
    Recoverable bool
    Suggestions []string
}

func (e *CLIError) Error() string {
    return e.Message
}

func (e *CLIError) ToDict() map[string]interface{} {
    return map[string]interface{}{
        "error": map[string]interface{}{
            "code":        e.Code,
            "type":        e.Type,
            "message":     e.Message,
            "details":     e.Details,
            "recoverable": e.Recoverable,
            "suggestions": e.Suggestions,
        },
    }
}
```

### Error Handling Pattern

```go
func handleCommand(args []string) error {
    result, err := performOperation()
    if err != nil {
        // Check if it's a CLI error
        if cliErr, ok := err.(*CLIError); ok {
            return formatter.OutputError(cliErr.ToDict(), cliErr.Code)
        }
        // Unexpected errors become internal errors
        internalErr := NewInternalError(fmt.Sprintf("Unexpected error: %v", err))
        return formatter.OutputError(internalErr.ToDict(), internalErr.Code)
    }
    
    return formatter.Output(result, EXIT_SUCCESS)
}
```

## Adding New Commands

### 1. Create Command Package

```go
// cmd/mycommand/mycommand.go
package mycommand

import (
    "fmt"
    "time"
)

type Handler struct {
    formatter Formatter
}

func NewHandler(formatter Formatter) *Handler {
    return &Handler{formatter: formatter}
}

func (h *Handler) Execute(param string) error {
    data := map[string]interface{}{
        "result": "success",
        "param":  param,
        "timestamp": time.Now().UTC().Format(time.RFC3339),
    }
    return h.formatter.Output(data, EXIT_SUCCESS)
}
```

### 2. Add Command Flags in main.go

```go
// In main.go
mycommandCmd := flag.NewFlagSet("mycommand", flag.ExitOnError)
param := mycommandCmd.String("param", "", "Parameter description")
human := mycommandCmd.Bool("human", false, "Human-readable output")
```

### 3. Add Command Routing

```go
// In main() function
case "mycommand":
    mycommandCmd.Parse(args[1:])
    handler := mycommand.NewHandler(formatter)
    if err := handler.Execute(*param); err != nil {
        return err
    }
```

### 4. Create JSON Schema

Create `schemas/mycommand.schema.json`:
```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "properties": {
    "result": {"type": "string"},
    "param": {"type": "string"},
    "timestamp": {"type": "string", "format": "date-time"}
  },
  "required": ["result", "param", "timestamp"]
}
```

## Configuration Management

### Environment Variables

Prefix all environment variables with `BOILERPLATE_`:

```bash
BOILERPLATE_PORT=8080
BOILERPLATE_HOST=127.0.0.1
BOILERPLATE_LOG_LEVEL=INFO
BOILERPLATE_PID_FILE=/tmp/boilerplate-cli-ui-go.pid
BOILERPLATE_LOG_FILE=/tmp/boilerplate-cli-ui-go.log
BOILERPLATE_NO_INTERACTIVE=1
```

### Configuration Package

```go
// pkg/config/config.go
package config

import (
    "os"
    "strconv"
)

type Config struct {
    Port     int
    Host     string
    LogLevel string
    PIDFile  string
    LogFile  string
}

func New() *Config {
    return &Config{
        Port:     getEnvWithDefault("BOILERPLATE_PORT", 8080),
        Host:     getEnvWithDefault("BOILERPLATE_HOST", "127.0.0.1"),
        LogLevel: getEnvWithDefault("BOILERPLATE_LOG_LEVEL", "INFO"),
        PIDFile:  getEnvWithDefault("BOILERPLATE_PID_FILE", "/tmp/boilerplate-cli-ui-go.pid"),
        LogFile:  getEnvWithDefault("BOILERPLATE_LOG_FILE", "/tmp/boilerplate-cli-ui-go.log"),
    }
}

func getEnvWithDefault(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intValue, err := strconv.Atoi(value); err == nil {
            return intValue
        }
    }
    return defaultValue
}

func getEnvStringWithDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func (c *Config) Override(overrides map[string]interface{}) {
    if port, ok := overrides["port"].(int); ok {
        c.Port = port
    }
    if host, ok := overrides["host"].(string); ok {
        c.Host = host
    }
}
```

## Testing Guidelines

### Test Structure

```go
// cmd/mycommand/mycommand_test.go
package mycommand_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "yourproject/pkg/output"
    "yourproject/cmd/mycommand"
)

func TestMyCommand(t *testing.T) {
    formatter := output.NewFormatter(false)
    handler := mycommand.NewHandler(formatter)
    
    err := handler.Execute("test_param")
    assert.NoError(t, err)
}
```

### Agent-Friendly Test Patterns

- Test JSON schema validation
- Test semantic exit codes
- Test error format structure
- Test `--human` mode output
- Test stderr/stdout separation
- Test environment variable handling

### Table-Driven Tests

```go
func TestValidatePort(t *testing.T) {
    tests := []struct {
        name    string
        port    int
        wantErr bool
    }{
        {"valid port", 8080, false},
        {"invalid port", -1, true},
        {"out of range", 70000, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidatePort(tt.port)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

## Development Workflow

### Local Development

```bash
# Run CLI
go run main.go greet --name Alice

# Run with human output
go run main.go greet --name Alice --human

# Start server (foreground)
go run main.go start --port 8080

# Start server (daemon)
go run main.go start --port 8080 --daemon
```

### Building Binary

```bash
# Make executable
chmod +x build.sh

# Build optimized binary
./build.sh

# Or manually
go build -ldflags "-s -w" -o boilerplate-cli-ui-go-optimized main.go server.go daemon.go
```

### Cross-Platform Builds

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o boilerplate-linux-amd64

# macOS
GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o boilerplate-darwin-amd64

# Windows
GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o boilerplate-windows-amd64.exe
```

## Go-Specific Development Guidelines

### Daemon Process Management

**Critical Pattern**: Always wait for process termination before cleanup

```go
// internal/daemon/daemon.go
package daemon

import (
    "os"
    "os/signal"
    "syscall"
    "time"
)

func (d *Daemon) Stop() error {
    pid, err := d.readPID()
    if err != nil {
        return err
    }
    
    process, err := os.FindProcess(pid)
    if err != nil {
        return err
    }
    
    // Send SIGTERM for graceful shutdown
    if err := process.Signal(syscall.SIGTERM); err != nil {
        return err
    }
    
    // Wait for process to terminate (max 5 seconds)
    timeout := time.After(5 * time.Second)
    ticker := time.NewTicker(100 * time.Millisecond)
    defer ticker.Stop()
    
    for {
        select {
        case <-timeout:
            // Force kill if graceful shutdown fails
            process.Signal(syscall.SIGKILL)
            time.Sleep(100 * time.Millisecond)
        case <-ticker.C:
            // Check if process still exists
            if err := process.Signal(syscall.Signal(0)); err != nil {
                // Process has terminated
                return d.cleanup()
            }
        }
    }
}
```

### HTTP Server Response Encoding

**Pattern**: Use proper JSON encoding with headers

```go
// pkg/server/server.go
package server

import (
    "encoding/json"
    "net/http"
)

func (s *Server) sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) error {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    
    encoder := json.NewEncoder(w)
    encoder.SetIndent("", "  ") // Pretty print for human readability
    return encoder.Encode(data)
}

func (s *Server) sendErrorResponse(w http.ResponseWriter, err error, statusCode int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    
    errorResponse := map[string]interface{}{
        "error": map[string]interface{}{
            "message": err.Error(),
            "code":    statusCode,
        },
    }
    
    json.NewEncoder(w).Encode(errorResponse)
}
```

### Configuration Management

**Pattern**: Environment variables with struct tags

```go
// pkg/config/config.go
package config

import (
    "os"
    "strconv"
)

type Config struct {
    Port     int    `env:"BOILERPLATE_PORT" default:"8080"`
    Host     string `env:"BOILERPLATE_HOST" default:"127.0.0.1"`
    LogLevel string `env:"BOILERPLATE_LOG_LEVEL" default:"INFO"`
}

func Load() *Config {
    cfg := &Config{}
    
    if port := os.Getenv("BOILERPLATE_PORT"); port != "" {
        if p, err := strconv.Atoi(port); err == nil {
            cfg.Port = p
        }
    }
    
    if host := os.Getenv("BOILERPLATE_HOST"); host != "" {
        cfg.Host = host
    }
    
    if logLevel := os.Getenv("BOILERPLATE_LOG_LEVEL"); logLevel != "" {
        cfg.LogLevel = logLevel
    }
    
    return cfg
}
```

### Error Handling Strategy

**Structured Error Pattern**:
```go
func performOperation() (map[string]interface{}, error) {
    result, err := doSomething()
    if err != nil {
        return nil, &CLIError{
            Code:    EXIT_CONNECTION_TIMEOUT,
            Type:    "connection_timeout",
            Message: "Failed to connect to service",
            Details: map[string]interface{}{
                "endpoint": "https://api.example.com",
            },
            Recoverable: true,
            Suggestions: []string{
                "Check network connectivity",
                "Verify endpoint URL",
                "Try again later",
            },
        }
    }
    
    return result, nil
}
```

### Signal Handling

**Pattern**: Graceful shutdown on SIGTERM/SIGINT

```go
func main() {
    // Setup signal handling
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    
    // Start server in goroutine
    go func() {
        if err := server.Start(); err != nil {
            log.Fatalf("Server error: %v", err)
        }
    }()
    
    // Wait for signal
    <-sigChan
    log.Println("Shutting down gracefully...")
    
    // Cleanup
    server.Stop()
}
```

## Agent-Friendly Design Checklist

When extending this boilerplate, ensure:

- [ ] All commands default to JSON output
- [ ] `--human` flag provides human-readable output
- [ ] Semantic exit codes for all error paths
- [ ] Structured error output with recovery hints
- [ ] Output separation (stdout data, stderr logs)
- [ ] No interactive prompts by default
- [ ] JSON schemas for all command outputs
- [ ] `--help-json` provides machine-readable help
- [ ] `--schema` provides schema discovery
- [ ] Environment variables for configuration
- [ ] Max 500 LOC per file
- [ ] Clear package responsibilities
- [ ] Comprehensive error handling
- [ ] Proper Go idioms and patterns

## Common Patterns

### Reading Configuration

```go
import "yourproject/pkg/config"

cfg := config.Load()
port := cfg.Port  // From env or default
cfg.Override(config.Override{Port: 3000})  // CLI override
```

### Formatting Output

```go
import "yourproject/pkg/output"

formatter := output.NewFormatter(false)
formatter.Output(map[string]interface{}{"result": "success"}, EXIT_SUCCESS)
```

### Handling Errors

```go
import "yourproject/pkg/errors"

result, err := performOperation()
if err != nil {
    if cliErr, ok := err.(*errors.CLIError); ok {
        return formatter.OutputError(cliErr.ToDict(), cliErr.Code)
    }
    return fmt.Errorf("unexpected error: %w", err)
}
```

## Performance Optimization

### Memory Management

```go
// Use sync.Pool for frequently allocated objects
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func processRequest(data []byte) {
    buf := bufferPool.Get().(*bytes.Buffer)
    defer func() {
        buf.Reset()
        bufferPool.Put(buf)
    }()
    
    buf.Write(data)
    // Process buffer
}
```

### Concurrent Processing

```go
// Use worker pools for concurrent tasks
func processItems(items []Item) {
    workerCount := 4
    jobs := make(chan Item, len(items))
    results := make(chan Result, len(items))
    
    // Start workers
    for i := 0; i < workerCount; i++ {
        go worker(jobs, results)
    }
    
    // Send jobs
    for _, item := range items {
        jobs <- item
    }
    close(jobs)
    
    // Collect results
    for i := 0; i < len(items); i++ {
        <-results
    }
}
```

## Binary Size Optimization

### Build Flags

```bash
# Strip debug information
go build -ldflags "-s -w" -o output

# Further optimization with UPX (optional)
upx --best --lzma output
```

### Dependency Management

```go
// Use specific imports to reduce binary size
import (
    "net/http"  // Instead of importing entire packages
    // Avoid unused dependencies
)
```

## Agent-First CLI Principles Reference

This boilerplate implements these core principles from agent-friendly CLI design:

1. **Machine-Friendly Escape Hatches**: All commands support `--no-interactive` and environment variables
2. **Output as API Contracts**: JSON output has stable schemas with versioning
3. **Semantic Exit Codes**: Uses 80-119 range for structured error communication
4. **Structured Output**: Default JSON with `--human` opt-in for readability
5. **Real-Time Feedback**: Progress on stderr, data on stdout

For detailed principles, see the Python boilerplate's `docs/AGENTS_FRIENDLY_TOOLS.md`.

## Future Enhancements

- [ ] Add comprehensive JSON schema validation
- [ ] Add `--help-json` command for machine-readable help
- [ ] Add `--schema` flag for schema discovery
- [ ] Add web UI frontend (React CDN or similar)
- [ ] Add configuration file support (YAML/TOML)
- [ ] Add authentication for web UI
- [ ] Add HTTPS support
- [ ] Add systemd service file generation
- [ ] Add metrics/monitoring endpoints
- [ ] Add integration tests

## Related Projects

- [SuperCLI](https://github.com/javimosch/supercli) - Universal CLI framework
- [supercli-clis](https://github.com/jarancibia/supercli-clis) - Collection of SuperCLI plugins
- [boilerplate-cli-ui-python](https://github.com/javimosch/boilerplate-cli-ui-python) - Python version of this boilerplate

## References

1. **InfoQ Article**: "Keep the Terminal Relevant: Patterns for AI Agent Driven CLIs" (August 2025)
2. **Square Engineering**: "Command Line Observability with Semantic Exit Codes" (January 2023)
3. **Command Line Interface Guidelines** (clig.dev)
4. **Effective Go** - https://golang.org/doc/effective_go
5. **Go Code Review Comments** - https://github.com/golang/go/wiki/CodeReviewComments