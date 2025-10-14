# GitHub Actions CI/CD Workflow

## Overview

This repository now has an automated build and test pipeline for the CLI flavor of OfflineTranscribe using GitHub Actions.

## Workflow File

`.github/workflows/cli-build-test.yml`

## Features

### Build Matrix
The workflow tests builds across multiple environments:
- **Operating Systems**: Ubuntu (Linux), Windows, macOS
- **Go Versions**: 1.21, 1.22, 1.23

This ensures compatibility across different platforms and Go versions.

### Workflow Triggers
The workflow runs on:
- Push to `main` or `develop` branches
- Pull requests targeting `main` or `develop` branches
- Manual trigger via `workflow_dispatch`

### Pipeline Steps

#### 1. Code Quality Checks
- **Go Vet**: Static analysis to find potential bugs (continues on error due to GUI dependencies)
- **Go Fmt**: Enforces consistent code formatting
- **golangci-lint**: Additional linting checks (optional, in separate job)

#### 2. Dependency Management
- **Go Module Download**: Downloads all dependencies
- **Go Module Verify**: Verifies dependencies integrity
- **Go Module Tidy**: Ensures go.mod is up to date
- **Caching**: Caches Go modules to speed up builds

#### 3. Build Process
- Creates placeholder resources needed for embedded files
- Builds CLI binary for each platform:
  - Linux: `OfflineTranscribe-cli`
  - Windows: `OfflineTranscribe-cli.exe`
  - macOS: `OfflineTranscribe-cli`
- Uses optimized build flags: `-ldflags "-s -w"` for smaller binaries

#### 4. Testing
- Runs unit tests with race detection
- Generates code coverage reports
- Uploads coverage to Codecov (optional)

#### 5. Artifacts
- Uploads built binaries for each OS/Go version combination
- Artifacts are retained for 7 days for download and testing

## Local Testing

To test the build process locally:

```bash
# Install dependencies
go mod download
go mod tidy

# Create placeholder resources (needed for embedded files)
mkdir -p bundle/resources/models bundle/resources/whisper
echo "placeholder" > bundle/resources/models/placeholder.txt
echo "placeholder" > bundle/resources/whisper/placeholder.txt

# Format code
gofmt -s -w .

# Build CLI
go build -ldflags "-s -w" -o OfflineTranscribe-cli cli.go whisper.go resources.go

# Run tests (if available)
go test -v ./...
```

## Build Requirements

### Embedded Resources
The build requires placeholder files in:
- `bundle/resources/models/placeholder.txt`
- `bundle/resources/whisper/placeholder.txt`

These placeholders allow the Go embed directive to work without requiring actual model and whisper binary files during the build process.

### Dependencies
All dependencies are managed via Go modules:
- See `go.mod` for full dependency list
- Main dependencies include Fyne for GUI components (used by main.go, not CLI)

## Workflow Outputs

### Artifacts
For each successful build, the workflow uploads:
- Binary artifact named: `OfflineTranscribe-cli-{os}-go{version}`
- Available for download from the Actions tab

### Code Coverage
- Coverage reports are generated during testing
- Uploaded to Codecov for visualization (if configured)

## Notes

### GUI Dependencies
The repository includes GUI components (main.go with Fyne) that require system libraries not available in CI environments. The workflow handles this gracefully:
- `go vet` continues on error
- Tests continue on error if they fail due to missing GUI libraries

### Platform-Specific Considerations
- **Linux**: No special requirements
- **Windows**: Uses PowerShell for some checks
- **macOS**: Standard build, no special requirements

## Troubleshooting

### Build Failures
If builds fail, check:
1. Go version compatibility (ensure code works with Go 1.21+)
2. Placeholder resource files exist
3. go.mod is properly formatted

### Test Failures
The workflow continues on test errors because:
- No test files currently exist in the repository
- GUI dependencies may cause test compilation to fail

## Future Enhancements

Potential improvements:
1. Add actual unit tests for CLI functionality
2. Add integration tests with sample audio files
3. Implement deployment to GitHub Releases
4. Add Docker container builds
5. Set up automated model downloads for testing
6. Add benchmark tests for performance tracking
