# .github Directory

This directory contains GitHub-specific configurations and documentation for the OfflineTranscribe project.

## Contents

### 📁 workflows/
Contains GitHub Actions workflow definitions for CI/CD automation.

- **cli-build-test.yml** - Main CI/CD pipeline for building and testing the CLI flavor of OfflineTranscribe
  - Runs on push/PR to main and develop branches
  - Tests across 3 operating systems and 3 Go versions
  - Produces build artifacts for each platform
  - See [WORKFLOW.md](WORKFLOW.md) for details

### 📄 Documentation Files

- **WORKFLOW.md** - Comprehensive documentation for the CI/CD pipeline
  - Overview of workflow features
  - Trigger events and pipeline steps
  - Local testing instructions
  - Troubleshooting guide
  - Future enhancement ideas

- **WORKFLOW_DIAGRAM.md** - Visual representation of the CI/CD pipeline
  - Trigger event diagrams
  - Build matrix visualization
  - Step-by-step pipeline flow
  - Success criteria and failure handling
  - Artifact generation details
  - Performance metrics

## Quick Links

- [View Workflow Runs](../../actions)
- [Download Build Artifacts](../../actions)
- [Workflow Documentation](WORKFLOW.md)
- [Visual Pipeline Diagrams](WORKFLOW_DIAGRAM.md)

## CI/CD Pipeline at a Glance

```
Trigger: Push/PR to main or develop
    ↓
Build Matrix: 3 OSes × 3 Go versions = 9 jobs
    ↓
Steps: Setup → Dependencies → Quality Checks → Build → Test → Artifacts
    ↓
Output: 9 platform-specific CLI binaries (7-day retention)
```

### Supported Platforms

- **Ubuntu Latest** (Linux)
- **Windows Latest**
- **macOS Latest**

### Supported Go Versions

- Go 1.21
- Go 1.22
- Go 1.23

### Build Outputs

Each successful workflow run produces 9 downloadable artifacts:
- `OfflineTranscribe-cli-ubuntu-latest-go1.21`
- `OfflineTranscribe-cli-ubuntu-latest-go1.22`
- `OfflineTranscribe-cli-ubuntu-latest-go1.23`
- `OfflineTranscribe-cli-windows-latest-go1.21`
- `OfflineTranscribe-cli-windows-latest-go1.22`
- `OfflineTranscribe-cli-windows-latest-go1.23`
- `OfflineTranscribe-cli-macos-latest-go1.21`
- `OfflineTranscribe-cli-macos-latest-go1.22`
- `OfflineTranscribe-cli-macos-latest-go1.23`

## For Contributors

When contributing to OfflineTranscribe:

1. **Format your code** before committing:
   ```bash
   gofmt -s -w .
   ```

2. **Test locally** before pushing:
   ```bash
   go build -ldflags "-s -w" -o OfflineTranscribe-cli cli.go whisper.go resources.go
   ```

3. **Check the workflow status** after pushing - the CI/CD pipeline will automatically run

4. **Download artifacts** from the Actions tab to test your changes on different platforms

## Workflow Badge

Add this badge to your README.md to show build status:

```markdown
![CLI Build and Test](https://github.com/lmiguelgato/OfflineTranscribe/workflows/CLI%20Build%20and%20Test/badge.svg)
```

## Additional Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Go GitHub Actions](https://github.com/actions/setup-go)
- [golangci-lint GitHub Action](https://github.com/golangci/golangci-lint-action)
