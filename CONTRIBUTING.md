# Contributing to Palantir Proxy

Thank you for your interest in contributing to Palantir Proxy! We welcome contributions of all forms, including bug reports, feature requests, documentation improvements, and pull requests.

## How to Contribute

### 1. Reporting Bugs and Feature Requests
- Please open a new Issue describing the bug or feature request in detail.
- Provide a clear description, reproduction steps, and expected behavior.

### 2. Pull Request Guidelines
- **Create a Branch**: Create a feature branch (e.g., `feat/grpc-streaming`, `fix/circuit-race`).
- **Follow Code Guidelines**:
  - Keep functions focused and between 4 to 20 lines of code.
  - Structure Go structs with memory alignment ordered from largest byte size to smallest.
  - Do not use emojis or decorative icons in documentation or commit messages.
- **Write Tests**: Every bug fix and new feature must be backed by fast, isolated unit tests.
- **Run Formatter**: Ensure code is formatted with `gofmt ./...` before submitting.

## Development Environment

To run the project locally, ensure you have Go 1.14+ installed.
