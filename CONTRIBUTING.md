# Contributing to RoxKV

Thank you for your interest in contributing to **RoxKV** and **RoxAI**.

RoxKV is a lightweight key-value database written in Go, inspired by systems such as Redis. RoxAI extends the project with an AI-powered layer for interacting with, monitoring, and understanding the database.

Contributions of all sizes are welcome — bug fixes, documentation improvements, new commands, tests, performance improvements, developer tooling, and larger features.

## Getting Started

### 1. Fork the Repository

Fork the repository and clone your fork:

```bash
git clone <your-fork-url>
cd roxkv
```

Add the original repository as an upstream remote:

```bash
git remote add upstream <original-repository-url>
```

Verify your remotes:

```bash
git remote -v
```

### 2. Create a Branch

Do not make changes directly on `main`.

Create a branch describing what you are working on:

```bash
git checkout -b feature/my-feature
```

Examples:

```text
feature/new-command
feature/roxai-tool
fix/ttl-expiration
fix/resp-parser
docs/docker-setup
refactor/storage-engine
test/resp-parser
```

Keep each branch focused on one logical change.

## Running the Project

Check the README for the latest setup instructions.

If you are using the containerized setup:

```bash
cp .env.example .env
docker compose up --build
```

For local Go development, ensure the Go version specified by the project's `go.mod` is installed.

Install dependencies:

```bash
go mod download
```

Run the appropriate RoxKV entrypoint according to the current repository structure.

## Project Areas

Before modifying the project, inspect the relevant package and understand how it interacts with the rest of the system.

Major areas may include:

* RESP parsing and serialization
* Command handling
* In-memory storage
* TTL and expiration
* Persistence and snapshots
* TCP networking
* Pub/Sub
* Monitoring and metrics
* RoxAI
* AI tools and routing
* Ollama integration
* Web/API layer
* Dashboard
* Docker and deployment
* Documentation

Avoid introducing dependencies between components unless they are necessary.

## Making Changes

### Keep Changes Focused

A pull request should solve one clear problem.

Avoid combining unrelated changes such as:

* fixing a parser bug
* redesigning the dashboard
* adding a command
* refactoring storage

into the same pull request.

Smaller pull requests are easier to review, test, and maintain.

### Follow Existing Code Style

For Go code, format your changes before committing:

```bash
gofmt -w .
```

Run:

```bash
go vet ./...
```

and:

```bash
go test ./...
```

when applicable.

Follow the existing naming, package, error-handling, and project conventions unless the purpose of the contribution is specifically to improve them.

### Error Handling

Do not silently ignore errors.

Prefer explicit error handling and meaningful error messages.

Avoid unnecessary panics in code paths where an error can be handled normally.

### Dependencies

Avoid adding dependencies when the same functionality can reasonably be implemented using the Go standard library or existing project dependencies.

If a new dependency is necessary, explain why in the pull request.

## Working on RoxKV Commands

When adding or modifying a command:

1. Understand how the command reaches the command handler.
2. Validate arguments correctly.
3. Return the correct response format.
4. Handle invalid input.
5. Consider concurrency implications.
6. Consider persistence implications.
7. Consider TTL behavior where applicable.
8. Add or update tests.
9. Update documentation.

Do not implement only the successful path.

## Working on the RESP Parser

The RESP parser is a critical part of RoxKV.

Changes to parsing or serialization should account for:

* malformed input
* incomplete input
* empty values
* multiple arguments
* binary-safe behavior where supported
* connection behavior
* protocol compatibility
* edge cases

Parser changes should include tests whenever possible.

## Working on Storage

The storage engine can be accessed concurrently.

Changes involving keys, values, TTLs, persistence, or internal data structures must consider concurrency and synchronization.

Be especially careful when modifying locking behavior.

Avoid introducing race conditions or holding locks longer than necessary.

Run relevant tests with Go's race detector when possible:

```bash
go test -race ./...
```

## Working on RoxAI

RoxAI should remain separated from the fundamental correctness of RoxKV.

The database should not depend on an LLM to perform core database operations.

When adding RoxAI tools or capabilities:

* Keep tool responsibilities narrow.
* Validate tool arguments.
* Return structured and useful results.
* Avoid sending unnecessary context to the LLM.
* Keep model-specific assumptions configurable.
* Avoid hardcoding Ollama addresses or model names.
* Handle unavailable AI services gracefully.

AI-generated responses should never silently replace deterministic database behavior.

## Docker Contributions

Containerization changes should preserve portability.

Avoid:

```text
hardcoded host paths
host-specific IP addresses
unnecessary host networking
machine-specific configuration
committing model files
committing secrets
```

Where possible, the complete application should remain runnable using:

```bash
docker compose up --build
```

If you modify Docker configuration, verify the Compose configuration:

```bash
docker compose config
```

## Tests

Add tests for new functionality when practical.

Before submitting a pull request, run:

```bash
go test ./...
```

For concurrency-sensitive changes:

```bash
go test -race ./...
```

A bug fix should ideally include a test demonstrating the bug and verifying the fix.

## Commit Messages

Use concise and descriptive commit messages.

Good examples:

```text
feat: add unsubscribe command
fix: prevent expired keys from being returned
fix: handle malformed RESP arrays
docs: add Docker setup instructions
refactor: simplify command routing
test: add RESP parser edge cases
```

Avoid messages such as:

```text
update
changes
fixed stuff
final
working now
```

## Pull Requests

Before opening a pull request:

* Make sure the project builds.
* Format your code.
* Run relevant tests.
* Remove debugging code.
* Remove unnecessary commented-out code.
* Update documentation when behavior changes.
* Check that no secrets or local configuration files are included.
* Keep the PR focused.

Your pull request description should explain:

### What changed?

Describe the implementation.

### Why?

Explain the problem being solved.

### How was it tested?

Include the commands, tests, or manual verification performed.

### Additional Notes

Mention limitations, compatibility concerns, follow-up work, or design decisions when relevant.

## Reporting Bugs

When reporting a bug, include enough information to reproduce it.

Useful information includes:

* RoxKV version/commit
* Operating system
* Architecture
* Go version
* Docker version, if relevant
* Command/input that caused the issue
* Expected behavior
* Actual behavior
* Relevant logs
* Minimal reproduction steps

Do not include passwords, API keys, tokens, private data, or other secrets in issues or logs.

## Feature Requests

Feature requests are welcome.

Explain:

* The problem you are trying to solve.
* The proposed behavior.
* Why it belongs in RoxKV/RoxAI.
* Possible alternatives, if applicable.

For large architectural changes, consider opening an issue for discussion before implementing the complete feature.

## Security Issues

Do not publicly disclose vulnerabilities that could put users or deployments at risk.

If a private security reporting mechanism is listed in the repository's security policy or project documentation, use that mechanism.

Do not include credentials, private keys, tokens, or sensitive deployment information in issues.

## Documentation

Documentation contributions are welcome.

This includes:

* command documentation
* examples
* architecture explanations
* setup instructions
* Docker documentation
* typo fixes
* troubleshooting information

Documentation should reflect actual project behavior.

## Code of Conduct

All contributors are expected to follow the project's `CODE_OF_CONDUCT.md`.

Be respectful when discussing implementations, reviewing code, reporting bugs, and disagreeing on technical decisions.

## License

By contributing to RoxKV, you agree that your contributions will be distributed under the same license as the project.
