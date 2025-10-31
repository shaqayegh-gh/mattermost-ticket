## Contributing

Thank you for your interest in contributing! We welcome issues, feature requests, and pull requests.

### Development Setup
- Install Go 1.21+
- Run `make deps`
- Build with `make build` or `./build-plugin.sh`

### Pull Requests
- Fork the repository and create a feature branch
- Ensure builds pass: `make test` and `make check-style`
- Keep changes focused and add clear descriptions
- Link related issues in the PR description

### Commit Messages
- Use imperative mood: "Add X", "Fix Y"
- Reference issues with `Fixes #123` when applicable

### Code Style
- Follow Go best practices and `golangci-lint` defaults
- Prefer clarity over cleverness; add concise comments for non-obvious logic

### Release Process
- Bump `version` in `plugin.json` and `PLUGIN_VERSION` in `Makefile`
- Update `CHANGELOG.md` (or GitHub release notes)
- Create a git tag `vX.Y.Z`

### Reporting Issues
- Include Mattermost server version and plugin version
- Provide logs and reproduction steps when possible


