# AGENTS.md - Crosfo Agent Guide

This document guides AI agents in understanding, extending, and maintaining the Crosfo application.

## Project Overview

Crosfo is a CLI-first web application for solo entrepreneurs to cross-follow and cross-like each other's content. It combines a Go backend with a premium web interface.

## Project Structure

```
crosfo/
├── main.go                  # CLI entry point and command routing
├── server.go                # HTTP server and web UI
├── daemon.go                # Process management
├── go.mod                   # Go module dependencies
├── go.sum                   # Go dependency checksums
├── build.sh                 # Binary compilation script
├── README.md                # User documentation
├── DESIGN.md                # Design system and guidelines
├── AGENTS.md                # This file - agent guide
├── pkg/
│   ├── database/            # SQLite database operations
│   │   └── database.go     # Database functions (users, communities, entries, admins)
│   └── handlers/            # API endpoint handlers
│       └── handlers.go      # HTTP request handlers
└── templates/               # HTML templates for web UI
    ├── index.html           # Community list page
    ├── community.html       # Community details page
    └── 404.html             # Custom 404 page
```

## Coding Rules

### File Size Limits

- **Max 500 LOC per Go file** - Split files that exceed this limit
- **Max 300 LOC per documentation file** - Keep documentation concise
- **Max 200 LOC per test file** - Split complex tests

### Module Organization

Each package has a single, well-defined responsibility:

- **pkg/database/**: SQLite database operations and queries
- **pkg/handlers/**: HTTP request handlers and API endpoints
- **templates/**: HTML templates for the web interface

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

**Database Operations:**
```go
// Use prepared statements to prevent SQL injection
stmt, err := db.Prepare("INSERT INTO users (username, password) VALUES (?, ?)")
if err != nil {
    return err
}
defer stmt.Close()

_, err = stmt.Exec(username, hashedPassword)
if err != nil {
    return err
}
```

**HTTP Handlers:**
```go
// Set proper content-type headers
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusOK)
json.NewEncoder(w).Encode(response)
```

## Database Schema

### Tables

**users**
- `id` (INTEGER PRIMARY KEY)
- `username` (TEXT UNIQUE)
- `password` (TEXT)
- `created_at` (TEXT)

**communities**
- `id` (INTEGER PRIMARY KEY)
- `name` (TEXT UNIQUE)
- `description` (TEXT)
- `created_at` (TEXT)

**entries**
- `id` (INTEGER PRIMARY KEY)
- `community_name` (TEXT)
- `username` (TEXT)
- `content` (TEXT)
- `content_type` (TEXT)
- `url` (TEXT)
- `created_at` (TEXT)

**community_admins**
- `id` (INTEGER PRIMARY KEY)
- `community_name` (TEXT)
- `username` (TEXT)
- `created_at` (TEXT)

## API Endpoints

### Community Management
- `GET /api/communities` - List all communities
- `GET /api/community/{name}` - Get community details
- `POST /api/community/update` - Update community details (admin only)

### Entry Management
- `GET /api/entries/{community}` - Get community entries
- `POST /api/thumb-up` - Thumb up an entry
- `GET /api/thumbs-up` - Get thumbs up count
- `GET /api/user-entries` - Get user's entries
- `POST /api/entry/update` - Update an entry
- `POST /api/entry/delete` - Delete an entry

### Admin Management
- `GET /api/community-admins/` - Get community admins
- `POST /api/community-admins/` - Add community admin
- `DELETE /api/community-admins/` - Remove community admin

### Authentication
- `POST /api/auth` - Authenticate user

## Adding New Features

### 1. Database Changes

Add functions in `pkg/database/database.go`:

```go
func NewFeature(db *sql.DB, param string) error {
    stmt, err := db.Prepare("INSERT INTO features (param) VALUES (?)")
    if err != nil {
        return err
    }
    defer stmt.Close()

    _, err = stmt.Exec(param)
    return err
}
```

### 2. API Handler

Add handler in `pkg/handlers/handlers.go`:

```go
func HandleNewFeature(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Parse request
    param := r.FormValue("param")

    // Call database function
    db := database.GetDB()
    if err := database.NewFeature(db, param); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Return success
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
```

### 3. Register Route

Add route in `server.go`:

```go
mux.HandleFunc("/api/new-feature", handlers.HandleNewFeature)
```

### 4. Frontend Updates

Update templates in `templates/` as needed.

## Design System

Crosfo uses a premium design system documented in `DESIGN.md`. Key principles:

- **Color Palette**: Off-white (#FAFAFA), charcoal (#1A1A1A), teal (#00B4D8)
- **Typography**: SF Pro Display/Geist Sans with multiple weights
- **Visual Depth**: 3-level shadow system with tinted shadows
- **Animations**: Smooth transitions, staggered entry, spotlight borders
- **Accessibility**: Skip links, focus rings, semantic HTML
- **Responsive**: Mobile-first with desktop enhancements

When updating the UI, follow the design tokens in `DESIGN.md` and maintain consistency.

## Development Workflow

### Local Development

```bash
# Start server in foreground
go run main.go server.go daemon.go start -port 8081

# Or build and run
go build -o crosfo main.go server.go daemon.go
./crosfo start -port 8081
```

### Building Binary

```bash
# Build optimized binary
go build -ldflags "-s -w" -o crosfo main.go server.go daemon.go

# Or use build script
chmod +x build.sh
./build.sh
```

### Testing

```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...
```

## Deployment

### Production Deployment

1. **Build optimized binary:**
```bash
go build -ldflags "-s -w" -o crosfo main.go server.go daemon.go
```

2. **Copy to server:**
```bash
scp crosfo user@server:/path/to/crosfo-bin
```

3. **Restart systemd service:**
```bash
ssh user@server "sudo systemctl restart crosfo.service"
```

### Environment Variables

Crosfo uses environment variables for configuration (see `.env.example`):

- `CROSFO_PORT` - Server port (default: 8081)
- `CROSFO_HOST` - Server host (default: 127.0.0.1)
- `CROSFO_DB_PATH` - Database file path (default: ./ffaf.db)
- `CROSFO_LOG_LEVEL` - Log level (default: INFO)
- `CROSFO_LOG_FILE` - Log file path (default: /tmp/crosfo.log)
- `CROSFO_PID_FILE` - PID file path (default: /tmp/crosfo.pid)

## Security Considerations

- **Password Hashing**: Use bcrypt for password hashing
- **SQL Injection**: Always use prepared statements
- **XSS Prevention**: Sanitize user input in templates
- **CSRF Protection**: Consider adding CSRF tokens for forms
- **Authentication**: Implement proper session management

## Common Patterns

### Database Query

```go
func GetCommunity(db *sql.DB, name string) (*Community, error) {
    var community Community
    err := db.QueryRow("SELECT id, name, description, created_at FROM communities WHERE name = ?", name).Scan(
        &community.ID, &community.Name, &community.Description, &community.CreatedAt,
    )
    if err != nil {
        return nil, err
    }
    return &community, nil
}
```

### Error Response

```go
func sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "error": message,
    })
}
```

### Success Response

```go
func sendSuccessResponse(w http.ResponseWriter, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(data)
}
```

## Agent-Friendly Design Checklist

When extending Crosfo, ensure:

- [ ] Follow file size limits (max 500 LOC per file)
- [ ] Use prepared statements for all database operations
- [ ] Handle errors explicitly and return appropriate HTTP status codes
- [ ] Follow naming conventions for Go code
- [ ] Update documentation when adding new features
- [ ] Follow the design system in DESIGN.md for UI changes
- [ ] Test changes locally before deploying
- [ ] Update .env.example if adding new environment variables
- [ ] Keep database migrations in mind when changing schema

## Troubleshooting

### Database Locked Error

If you encounter "database is locked" errors:

1. Check if multiple processes are accessing the database
2. Ensure proper connection closing with `defer db.Close()`
3. Consider using WAL mode for better concurrency

### Port Already in Use

If the port is already in use:

```bash
# Find process using the port
lsof -i :8081

# Kill the process
kill <PID>
```

### Template Not Found

If templates are not loading:

1. Check that the working directory is correct
2. Verify template files exist in the `templates/` directory
3. Check file permissions

## Future Enhancements

- [ ] Add comprehensive test coverage
- [ ] Implement proper session management
- [ ] Add CSRF protection
- [ ] Add rate limiting for API endpoints
- [ ] Implement database migrations
- [ ] Add logging framework
- [ ] Add metrics/monitoring
- [ ] Add webhook support for notifications

## Related Projects

- [SuperCLI](https://github.com/javimosch/supercli) - Universal CLI framework
- [supercli-clis](https://github.com/jarancibia/supercli-clis) - Collection of SuperCLI plugins

## License

This project is provided as-is for educational and development purposes.
