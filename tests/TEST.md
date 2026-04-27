# BAA Testing Guide

This document explains how to run the test suite for BAA. 

## Running All Tests

To run the entire test suite, including internal tests and the consolidated tests in the `tests/` directory:

```bash
go test ./...
```

## Running Specific Tests

If you only want to run the tests in the consolidated `tests/` folder:

```bash
go test ./tests/...
```

## Test Structure

- `tests/sanitize_test.go`: Verifies error message cleaning and truncation logic.
- `tests/providers_test.go`: Ensures each package manager provider (apt, dnf, pacman, etc.) returns the correct commands and environment variables.
- `tests/views_test.go`: Tests UI-related helper functions like `LastErrorLine`.
- `tests/detector_test.go`: Uses a Mock Runner to verify that package manager detection logic works correctly without needing real binaries installed.

## Writing New Tests

New tests should be added to the `tests/` directory. Since these tests are in a separate package, they can only test **exported** members (functions/structs starting with a capital letter) of the internal packages.

If you need to test private/unexported members, you should create a `_test.go` file inside the package itself (standard Go practice), but for BAA's architectural goals, high-level API testing in the `tests/` folder is preferred.
