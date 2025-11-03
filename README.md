# Mattermost Ticket Plugin

[![CI](https://github.com/shaqayegh-gh/mattermost-ticket/actions/workflows/ci.yml/badge.svg)](https://github.com/shaqayegh-gh/mattermost-ticket/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Mattermost](https://img.shields.io/badge/Mattermost-Plugin-blue)](https://developers.mattermost.com/integrate/plugins/)

A Golang plugin for Mattermost that allows users to create tickets using the `/ticket` slash command with interactive dialogs.

## Features

- **Slash Command**: Use `/ticket` to create tickets
- **Interactive Dialog**: User-friendly form with dropdowns and text input
- **Message Priority**: Set native Mattermost priority (Standard, Important, Urgent)
- **Team Management**: Configurable team members via System Console
- **Auto Mentions**: Automatically mentions team members and "all" users
- **Multiple Teams**: Support for 14 different teams (Issuance, Design, DevOps, etc.)
- **Project Selection**: 30+ backend and frontend projects
- **Environment Support**: Development, Stage, Production environments
 - **Allowed Channels**: Restrict usage to a list of channels

## Quick Start

**Tested with:**
- Mattermost Team Edition: `10.12.0`
- PostgreSQL: `13-alpine`

### 1. Build the Plugin

```bash
# Clone the repository
git clone https://github.com/shaqayegh-gh/mattermost-ticket.git
cd mattermost-ticket-plugin

# Build the plugin
./build-plugin.sh
```

This creates `mattermost-ticket-plugin.tar.gz` in the project root.

Releases and changelogs are published on the [Releases page](https://github.com/shaqayegh-gh/mattermost-ticket/releases).


### 2. Install the Plugin

#### Method A: Upload via UI (Recommended)

1. Go to **System Console** → **Plugins** → **Plugin Management**
2. Click **Upload Plugin**
3. Select `mattermost-ticket-plugin.tar.gz`
4. Click **Enable** after upload

#### Method B: Copy to Container

```bash
# Copy plugin to running container
docker cp mattermost-ticket-plugin.tar.gz mattermost:/mattermost/plugins/

# Restart Mattermost to load plugin
docker restart mattermost
```

### 3. Configure Settings

1. Go to **System Console** → **Plugins** → **Ticket Plugin**
2. Set **Allowed Channels for Tickets** (`DefaultChannel`) with a comma or newline separated list of channel names that can use `/ticket`. Leave empty to allow any channel.

   Examples:

   - Single line (comma separated): `tickets, support, dev-help`
   - Multi-line:

     ```
     tickets
     support
     dev-help
     ```
3. In **"Team Members Configuration"** (`TeamMembersConfig`), enter JSON:

```json
{
  "issuance": ["john.doe", "jane.smith"],
  "design": ["designer1", "designer2"],
  "marketing": ["marketer1", "marketer2"],
  "devops": ["devops1", "devops2"],
  "all": ["admin", "manager"]
}
```

4. Optionally set **Team Options (Dropdown)** (`TeamOptionsConfig`) and **Project Options (Dropdown)** (`ProjectOptionsConfig`) as JSON arrays to override the built-in defaults. Then click **Save**.

**Important Notes:**
- Use actual Mattermost usernames (not display names)
- Users in `"all"` will be mentioned in every ticket
- Team names must match exactly: `issuance`, `design`, `marketing`, etc.

### 4. Test the Plugin

1. Go to any Mattermost channel
2. Type `/ticket`
3. Fill out the dialog:
   - **Team**: Select a team (e.g., "Issuance")
   - **Project**: Select a project (e.g., "Backend API")
   - **Environment**: Select environment (Development/Stage/Production)
   - **Description**: Enter issue description
4. Click **Create Ticket**

The ticket will be posted with mentions for team members and "all" users. If a priority was selected, the post will use Mattermost's native message priority and show the corresponding label in the UI.

## Configuration

### Team Options (Configurable)

You can override the team dropdown via System Console → Plugins → Ticket Plugin.

- Set **"Team Options (Dropdown)"** (`TeamOptionsConfig`) with a JSON array:

```json
[
  { "Text": "Issuance", "Value": "issuance" },
  { "Text": "QA", "Value": "qa" },
  { "Text": "DevOps", "Value": "devops" }
]
```

If empty, the built-in defaults below are used:

- Develop
- Design
- Scrum Master
- QA
- Product
- Marketing
- Support
- DevOps
- HR
- Others...

### Project Options (Configurable)

You can override the project dropdown via System Console → Plugins → Ticket Plugin.

- Set **"Project Options (Dropdown)"** (`ProjectOptionsConfig`) with a JSON array:

```json
[
  { "Text": "Sodoor PWA Frontend", "Value": "sodoor-pwa-frontend" },
  { "Text": "Estate API Backend", "Value": "estate-api-backend" },
  { "Text": "Others...", "Value": "others" }
]
```

If empty, the built-in defaults below are used:

- Backend
- Frontend
- Others...

### Environment Options (Fixed in Code)

- Development
- Stage
- Production

## Customization

### Adding New Teams/Projects

Edit `server/constants.go` and modify the variables:

```go
var (
    teamOptions = []*model.PostActionOptions{
        {Text: "Your New Team", Value: "your-new-team"},
        // ... existing teams
    }

    projectOptions = []*model.PostActionOptions{
        {Text: "Your New Project", Value: "your-new-project"},
        // ... existing projects
    }
)
```

Then rebuild: `./build-plugin.sh`

### Team Members Configuration
### Message Priority (Native)

- The dialog includes a "Message Priority" select with these options: `Standard`, `Important`, `Urgent`.
- The selected value is written to the post `props.priority.priority`, which uses Mattermost's native priority feature so the label shows in the UI (no custom rendering required).
- If no priority is chosen, `Standard` is used by default.

### Allowed Channels

- You can restrict where `/ticket` can be used by configuring **Allowed Channels for Tickets** (`DefaultChannel`).
- Provide channel names separated by commas or newlines. Matching is by channel name (case-insensitive).
- Leave empty to allow the command in any channel.


Team members are configured via Mattermost System Console, not in code. This allows admins to:

- Add/remove team members without code changes
- Use actual Mattermost usernames
- Configure "all" users who get mentioned in every ticket
- Update configuration in real-time

## Troubleshooting

### Plugin Not Loading

1. Check Mattermost logs:
   ```bash
   docker logs mattermost
   ```

2. Verify plugin is enabled:
   - System Console → Plugins → Plugin Management
   - Ensure "Ticket Plugin" is enabled

3. Check plugin file exists:
   ```bash
   docker exec mattermost ls -la /mattermost/plugins/
   ```

### Team Members Not Mentioned

1. **Check usernames**: Must match exactly (case-sensitive)
2. **Check JSON syntax**: Validate at https://jsonlint.com/
3. **Check team names**: Must match exactly (e.g., `issuance`, not `Issuance`)

### Build Issues

```bash
# Clean and rebuild
rm -rf server/dist/
./build-plugin.sh
```

### Docker Issues

```bash
# Restart containers
docker restart mattermost
docker restart mattermost-postgres

# Check container status
docker ps

# View logs
docker logs mattermost
```

## Development

### Project Structure

```
mattermost-ticket-plugin/
├── server/
│   └── main.go          # Main plugin code
├── plugin.json          # Plugin manifest
├── go.mod              # Go dependencies
├── go.sum              # Go checksums
├── build-plugin.sh     # Build script
└── README.md           # This file
```

### Building from Source

```bash
# Install dependencies
go mod tidy

# Build plugin
./build-plugin.sh

# Test locally (if Mattermost dependencies are resolved)
go vet ./server
```

## License

This project is licensed under the [MIT License](LICENSE).

## Support

For issues and questions:
1. Check the troubleshooting section above
2. Verify your Mattermost version compatibility
3. Check Mattermost logs for errors
4. Ensure proper Docker setup and networking

Commercial support and responsible disclosure: see [SECURITY.md](SECURITY.md). For contributing, see [CONTRIBUTING.md](CONTRIBUTING.md) and [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md).