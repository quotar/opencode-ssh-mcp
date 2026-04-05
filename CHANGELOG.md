# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial release of opencode-ssh-mcp
- MCP (Model Context Protocol) implementation
- SSH configuration support (`~/.ssh/config`, `~/.ssh/config.d/*.conf`)
- Connection pooling and session management
- Support for ProxyJump configurations
- Auto-retry mechanism (3 attempts with exponential backoff)

### Changed
- 

### Deprecated
- 

### Removed
- 

### Fixed
- 

### Security
- 

## [v1.0.0] - 2025-01-01

### Added
- Initial release
- SSH Config parser supporting OpenSSH format
- MCP server implementation compatible with opencode
- Connection management with multiple host support
- Command execution via SSH
- Session state management
- Status and disconnect functionality