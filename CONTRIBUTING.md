# Contributing to Kubiya Terraform Provider

Thank you for your interest in contributing to the Kubiya Terraform Provider! This document provides guidelines and instructions for contributing to this project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [How to Contribute](#how-to-contribute)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Pull Request Process](#pull-request-process)
- [Reporting Bugs](#reporting-bugs)
- [Suggesting Enhancements](#suggesting-enhancements)

## Code of Conduct

This project adheres to the Contributor Covenant Code of Conduct. By participating, you are expected to uphold this code. Please report unacceptable behavior to the project maintainers.

## Getting Started

1. Fork the repository on GitHub
2. Clone your fork locally
3. Create a new branch for your changes
4. Make your changes
5. Push your changes to your fork
6. Submit a pull request

## Development Setup

### Prerequisites

- Go 1.21 or later
- Terraform 1.0 or later
- Git

### Local Development

1. Clone the repository:
   ```bash
   git clone https://github.com/kubiya-terraform/terraform-provider-kubiya.git
   cd terraform-provider-kubiya
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Build the provider:
   ```bash
   go build -o terraform-provider-kubiya
   ```

4. Set up local testing:
   ```bash
   # Create a local Terraform plugins directory
   mkdir -p ~/.terraform.d/plugins/registry.terraform.io/kubiya/kubiya/0.0.1/darwin_arm64

   # Copy the built provider
   cp terraform-provider-kubiya ~/.terraform.d/plugins/registry.terraform.io/kubiya/kubiya/0.0.1/darwin_arm64/
   ```

5. Test with Terraform:
   ```bash
   cd examples
   terraform init
   terraform plan
   ```

## How to Contribute

### Types of Contributions

We welcome various types of contributions:

- Bug fixes
- New features
- Documentation improvements
- Test coverage improvements
- Performance optimizations
- Code refactoring

### Before You Start

1. Check existing issues and pull requests to avoid duplicates
2. For major changes, open an issue first to discuss your proposal
3. Make sure you understand the project's architecture and conventions

## Coding Standards

### Go Code Style

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` to format your code
- Run `go vet` to check for common mistakes
- Use meaningful variable and function names
- Add comments for exported functions and complex logic

### Terraform Provider Conventions

- Follow [HashiCorp's Terraform Provider Development Guidelines](https://www.terraform.io/docs/extend/best-practices/index.html)
- Use the Terraform Plugin Framework patterns
- Implement proper error handling and validation
- Add appropriate logging for debugging

### Code Organization

- Keep files focused on a single resource or data source
- Place shared utilities in appropriate packages
- Follow the existing directory structure:
  - `internal/provider/` - Provider implementation and resources
  - `internal/clients/` - API client implementations
  - `internal/entities/` - Data models and entities
  - `examples/` - Example Terraform configurations
  - `docs/` - Provider documentation

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific tests
go test -v ./internal/provider -run TestAccAgent

# Run integration tests
TF_ACC=1 go test -v ./internal/provider -run TestAccAgent
```

### Writing Tests

1. **Unit Tests**: Write unit tests for utility functions and business logic
   ```go
   func TestHelperFunction(t *testing.T) {
       result := helperFunction(input)
       if result != expected {
           t.Errorf("Expected %v, got %v", expected, result)
       }
   }
   ```

2. **Acceptance Tests**: Write acceptance tests for resources
   ```go
   func TestAccKubiyaAgent_basic(t *testing.T) {
       resource.Test(t, resource.TestCase{
           PreCheck:     func() { testAccPreCheck(t) },
           Providers:    testAccProviders,
           CheckDestroy: testAccCheckKubiyaAgentDestroy,
           Steps: []resource.TestStep{
               {
                   Config: testAccKubiyaAgentConfig_basic(),
                   Check: resource.ComposeTestCheckFunc(
                       testAccCheckKubiyaAgentExists("kubiya_agent.test"),
                   ),
               },
           },
       })
   }
   ```

### Test Coverage

- Aim for at least 70% code coverage
- Focus on critical paths and edge cases
- Test error conditions and validation logic

## Pull Request Process

1. **Update Documentation**: Update relevant documentation including:
   - README.md if adding new features
   - CHANGELOG.md following [Keep a Changelog](https://keepachangelog.com/) format
   - Code comments and inline documentation

2. **Ensure Tests Pass**: Make sure all tests pass before submitting
   ```bash
   go test ./...
   ```

3. **Code Quality**: Run linters and formatters
   ```bash
   gofmt -s -w .
   go vet ./...
   ```

4. **Commit Messages**: Write clear, descriptive commit messages
   - Use present tense ("Add feature" not "Added feature")
   - Use imperative mood ("Move cursor to..." not "Moves cursor to...")
   - Reference issues and pull requests when relevant
   - Example: "Fix agent creation validation (#123)"

5. **Pull Request Description**: Provide a clear description including:
   - What changes were made
   - Why the changes were made
   - How to test the changes
   - Related issues or pull requests

6. **Review Process**:
   - Maintainers will review your pull request
   - Address any feedback or requested changes
   - Once approved, a maintainer will merge your PR

## Reporting Bugs

When reporting bugs, please include:

1. **Description**: Clear description of the bug
2. **Steps to Reproduce**: Detailed steps to reproduce the issue
3. **Expected Behavior**: What you expected to happen
4. **Actual Behavior**: What actually happened
5. **Environment**:
   - Terraform version
   - Provider version
   - Operating system
   - Go version (if building from source)
6. **Configuration**: Relevant Terraform configuration (sanitized)
7. **Logs**: Any relevant logs or error messages

Use the bug report template in `.github/ISSUE_TEMPLATE/bug_report.md`

## Suggesting Enhancements

When suggesting enhancements:

1. **Use Case**: Describe the use case for the enhancement
2. **Proposed Solution**: Explain your proposed solution
3. **Alternatives**: Describe alternatives you've considered
4. **Additional Context**: Add any other context or screenshots

Use the feature request template in `.github/ISSUE_TEMPLATE/feature_request.md`

## Development Workflow

1. Create a feature branch from `main`:
   ```bash
   git checkout -b feature/my-new-feature
   ```

2. Make your changes and commit them:
   ```bash
   git add .
   git commit -m "Add new feature"
   ```

3. Keep your branch up to date:
   ```bash
   git fetch origin
   git rebase origin/main
   ```

4. Push to your fork:
   ```bash
   git push origin feature/my-new-feature
   ```

5. Open a pull request against the `main` branch

## Questions?

If you have questions about contributing:

- Check existing issues and discussions
- Review the documentation
- Reach out to the maintainers

## License

By contributing to this project, you agree that your contributions will be licensed under the Apache License 2.0.

## Thank You!

Your contributions help make this project better for everyone. Thank you for taking the time to contribute!
