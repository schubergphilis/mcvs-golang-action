# Build Tags

Build tags (also known as build constraints) allow you to include or exclude Go files from compilation based on conditions. This action supports the following common build tag patterns:

- **`integration`**: For integration tests that require external services or databases
- **`component`**: For component tests that test multiple units working together
- **`e2e`**: For end-to-end tests that test the entire application flow
- **`lambda.norpc`**: For building AWS Lambda functions without RPC support

## Writing Tagged Tests

Put the tests in a file with a matching postfix, such as
`some_integration_test.go`, `some_component_test.go` or `some_e2e_test.go`, and
start the file with the corresponding header, e.g.:

```go
//go:build integration
```

## Using Build Tags

When running tests with specific build tags:

```zsh
# Run integration tests (unit tests are run as well)
task remote:test-integration --yes

# Run component tests
task remote:test-component --yes

# Run end-to-end tests
task remote:test-e2e --yes
```

`task remote:test --yes` only runs the unit tests.

When linting code with specific build tags, you may need to run the linter multiple times to cover all code paths:

```yml
- testing-type: "lint"  # Lint main code
- testing-type: "lint", build-tags: "integration"  # Lint integration test code
- testing-type: "lint", build-tags: "component"  # Lint component test code
```

This ensures that code in test files with different build tags is properly linted.
