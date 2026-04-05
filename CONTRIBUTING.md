# Contributing to opencode-ssh-mcp

Thank you for considering contributing to opencode-ssh-mcp! This document outlines the process for contributing to the project.

## Development Environment Setup

### Prerequisites
- Go 1.21+
- git
- make (optional)

### Setting up the environment
```bash
# Clone the repository
git clone https://github.com/quotar/opencode-ssh-mcp.git
cd opencode-ssh-mcp

# Install dependencies
go mod download

# Build the project
go build -o opencode-ssh-mcp
```

## Code Style

- Follow Go language idioms
- Use `gofmt` to format code
- Use meaningful variable and function names
- Include appropriate comments, especially for complex logic
- Follow the existing code structure and patterns

## Testing

Before submitting changes, please ensure:
- All existing tests pass
- New functionality is tested
- Run the complete test suite

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

## Submitting Pull Requests

1. Fork the repository
2. Create a new branch from main
3. Make changes
4. Run tests
5. Submit a PR with description of changes

### Commit Message Format
- Use imperative mood ("Add feature" not "Added feature")
- Keep first line short (<72 characters)
- Add detailed description if needed after a blank line

## Issue Reporting

When reporting issues, please include:
- Operating system and Go version
- Steps to reproduce
- Expected behavior
- Actual behavior
- Relevant logs or error messages

## Community

Please follow our code of conduct to create a welcoming and inclusive environment.

All kinds of contributions are welcome, including:
- Bug fixes
- Feature implementations
- Documentation improvements
- Issue and PR reviews