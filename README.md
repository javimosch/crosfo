# Crosfo - Cross Follows

Crosfo helps solo entrepreneurs cross-follow and cross-like each other's content to grow together. Join communities and amplify your reach.

## Features

- **Community Management**: Create and manage communities for cross-promotion
- **Entry Tracking**: Track your cross-follows and cross-likes
- **Admin Functionality**: Community admins can edit community details and manage members
- **Premium UI**: Modern, responsive design with smooth animations
- **Mobile-First**: Optimized for mobile devices with desktop polish
- **Real-time Updates**: Live status updates and notifications

## Architecture

Crosfo is built as a CLI-first application with a web interface:

```
Crosfo
├── main.go - CLI entry point and command routing
├── server.go - HTTP server and web UI
├── daemon.go - Process management
├── pkg/
│   ├── database/ - SQLite database operations
│   └── handlers/ - API endpoint handlers
└── templates/ - HTML templates for web UI
```

## Tech Stack

- **Backend**: Go 1.21+
- **Database**: SQLite
- **Frontend**: Vanilla HTML/CSS/JavaScript
- **Design**: Premium design system with CSS variables
- **Deployment**: systemd service with Traefik reverse proxy

## Design System

Crosfo uses a premium design system with:

- **Color Palette**: Off-white backgrounds, charcoal text, teal accent
- **Typography**: SF Pro Display/Geist Sans with multiple weights
- **Visual Depth**: 3-level shadow system with tinted shadows
- **Animations**: Smooth transitions, staggered entry, spotlight borders
- **Accessibility**: Skip links, focus rings, semantic HTML
- **Responsive**: Mobile-first with desktop enhancements

See [DESIGN.md](DESIGN.md) for complete design guidelines.

## Installation

### Build from Source

```bash
# Clone the repository
git clone https://github.com/javimosch/crosfo.git
cd crosfo

# Build the binary
go build -o crosfo main.go server.go daemon.go

# Or use the build script
chmod +x build.sh
./build.sh
```

### Systemd Service

Create a systemd service file at `/etc/systemd/system/crosfo.service`:

```ini
[Unit]
Description=Crosfo (Cross Follows) App
After=network.target

[Service]
Type=simple
User=dk1
WorkingDirectory=/home/dk1/ffaf
ExecStart=/home/dk1/crosfo-bin start -port 8081
Restart=always

[Install]
WantedBy=multi-user.target
```

Enable and start the service:

```bash
sudo systemctl enable crosfo.service
sudo systemctl start crosfo.service
```

## Usage

### CLI Commands

**Start HTTP server (foreground):**
```bash
./crosfo start -port 8081
```

**Start as daemon (background):**
```bash
./crosfo start -port 8081 -daemon
```

**Stop daemon:**
```bash
./crosfo stop
```

**Check daemon status:**
```bash
./crosfo status
```

**Show version:**
```bash
./crosfo version
```

### Web Interface

When the server is running, access the UI at:
- `http://localhost:8081` (default port)

### API Endpoints

- `GET /` - Web UI (community list)
- `GET /c/{community}` - Community details page
- `GET /api/communities` - List all communities
- `GET /api/community/{name}` - Get community details
- `GET /api/entries/{community}` - Get community entries
- `POST /api/thumb-up` - Thumb up an entry
- `GET /api/thumbs-up` - Get thumbs up count
- `GET /api/user-entries` - Get user's entries
- `POST /api/auth` - Authenticate user
- `POST /api/entry/update` - Update an entry
- `POST /api/entry/delete` - Delete an entry
- `GET /api/community-admins/` - Get community admins
- `POST /api/community-admins/` - Add community admin
- `DELETE /api/community-admins/` - Remove community admin
- `POST /api/community/update` - Update community details

## Database

Crosfo uses SQLite for data storage. The database file is located at `./ffaf.db` by default.

### Schema

- `users` - User accounts
- `communities` - Community information
- `entries` - Cross-follow/like entries
- `community_admins` - Community admin relationships

## Development

### Running in Development

```bash
# Start server in foreground
go run main.go server.go daemon.go start -port 8081

# Or build and run
go build -o crosfo main.go server.go daemon.go
./crosfo start -port 8081
```

### Adding New Features

1. **Database changes**: Add functions in `pkg/database/database.go`
2. **API endpoints**: Add handlers in `pkg/handlers/handlers.go`
3. **Routes**: Register routes in `server.go`
4. **Frontend**: Update templates in `templates/`

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

3. **Setup systemd service** (see Installation section)

4. **Configure reverse proxy** (Traefik/Nginx)

### Traefik Configuration

Example Traefik dynamic configuration:

```yaml
http:
  routers:
    crosfo:
      rule: "Host(`crosfo.intrane.fr`)"
      service: crosfo
      entryPoints:
        - websecure
      tls:
        certResolver: letsencrypt
        domains:
          - main: crosfo.intrane.fr

  services:
    crosfo:
      loadBalancer:
        servers:
          - url: "http://127.0.0.1:8081"
```

## Contributing

Contributions are welcome! Please follow these guidelines:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is provided as-is for educational and development purposes.

## Live Demo

Crosfo is live at [https://crosfo.intrane.fr](https://crosfo.intrane.fr)

## Roadmap

### Planned Features

**Screenshot Validation System**
- Upload screenshots as proof of follows/likes
- Mutual validation between community members
- Screenshot review and approval workflow
- Automated image verification (basic validation)
- Trust score system based on validation history

**Enhanced Authentication**
- OAuth integration (Google, GitHub, Twitter)
- Two-factor authentication (2FA)
- Session management with refresh tokens
- Password reset functionality

**Advanced Community Features**
- Community categories and tags
- Community discovery and search
- Community rules and guidelines
- Member reputation system
- Activity feed and notifications

**Analytics Dashboard**
- Follow/like statistics and trends
- Member engagement metrics
- Community growth analytics
- Export data to CSV/JSON
- Visual charts and graphs

**Mobile App**
- Native iOS and Android applications
- Push notifications for new entries
- Offline mode with sync
- Quick screenshot capture and upload

**Integration Features**
- API for third-party integrations
- Webhook notifications for events
- Slack/Discord bot integration
- Email notifications and digests

**Gamification**
- Points and badges system
- Leaderboards within communities
- Achievement unlocks
- Streak tracking for consistent participation

**Advanced Moderation**
- Report system for violations
- Moderation queue for admins
- Automatic spam detection
- Temporary bans and warnings

**Multi-Platform Support**
- Support for more social platforms (Instagram, TikTok, LinkedIn)
- Platform-specific validation rules
- Cross-platform activity tracking
- Unified dashboard for all platforms

**Collaboration Features**
- Direct messaging between members
- Group challenges and campaigns
- Collaborative entry creation
- Comment system on entries

## Support

For issues and questions, please open an issue on GitHub.
