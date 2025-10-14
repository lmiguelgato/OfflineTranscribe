# CI/CD Workflow Visualization

## Workflow Trigger Events
```
┌─────────────────────────────────────────┐
│  Trigger Events                         │
│  • Push to main/develop                 │
│  • Pull Request to main/develop         │
│  • Manual workflow_dispatch             │
└─────────────────┬───────────────────────┘
                  │
                  ▼
```

## Build Matrix (9 Parallel Jobs)
```
┌──────────────────────────────────────────────────────────────┐
│                      Build Matrix                             │
├──────────────────────────────────────────────────────────────┤
│                                                               │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │  Ubuntu      │  │  Windows     │  │  macOS       │      │
│  │  + Go 1.21   │  │  + Go 1.21   │  │  + Go 1.21   │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│                                                               │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │  Ubuntu      │  │  Windows     │  │  macOS       │      │
│  │  + Go 1.22   │  │  + Go 1.22   │  │  + Go 1.22   │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│                                                               │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │  Ubuntu      │  │  Windows     │  │  macOS       │      │
│  │  + Go 1.23   │  │  + Go 1.23   │  │  + Go 1.23   │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│                                                               │
└──────────────────────────────────────────────────────────────┘
```

## Pipeline Steps (Each Job)
```
┌────────────────────────────────────────────────────────┐
│  Step 1: Setup                                         │
│  ├─ Checkout code                                      │
│  ├─ Set up Go ($version)                               │
│  ├─ Display Go version                                 │
│  └─ Cache Go modules                                   │
└────────────┬───────────────────────────────────────────┘
             │
             ▼
┌────────────────────────────────────────────────────────┐
│  Step 2: Dependencies                                  │
│  ├─ Download dependencies (go mod download)           │
│  ├─ Verify dependencies (go mod verify)               │
│  └─ Tidy dependencies (go mod tidy)                   │
└────────────┬───────────────────────────────────────────┘
             │
             ▼
┌────────────────────────────────────────────────────────┐
│  Step 3: Code Quality                                  │
│  ├─ Run go vet (continue-on-error)                    │
│  └─ Run go fmt check                                   │
└────────────┬───────────────────────────────────────────┘
             │
             ▼
┌────────────────────────────────────────────────────────┐
│  Step 4: Build Preparation                             │
│  ├─ Create bundle/resources/models/                    │
│  ├─ Create bundle/resources/whisper/                   │
│  ├─ Add placeholder.txt to models/                     │
│  └─ Add placeholder.txt to whisper/                    │
└────────────┬───────────────────────────────────────────┘
             │
             ▼
┌────────────────────────────────────────────────────────┐
│  Step 5: Build CLI                                     │
│  └─ go build -ldflags "-s -w" -o                      │
│     OfflineTranscribe-cli[.exe]                       │
│     cli.go whisper.go resources.go                    │
└────────────┬───────────────────────────────────────────┘
             │
             ▼
┌────────────────────────────────────────────────────────┐
│  Step 6: Verification                                  │
│  ├─ Verify CLI binary exists                          │
│  └─ Display binary info                               │
└────────────┬───────────────────────────────────────────┘
             │
             ▼
┌────────────────────────────────────────────────────────┐
│  Step 7: Testing                                       │
│  └─ Run tests with race detection                     │
│     (continue-on-error if no tests)                   │
└────────────┬───────────────────────────────────────────┘
             │
             ▼
┌────────────────────────────────────────────────────────┐
│  Step 8: Artifacts                                     │
│  ├─ Upload CLI binary                                  │
│  │  Name: OfflineTranscribe-cli-{os}-go{version}     │
│  │  Retention: 7 days                                 │
│  └─ Upload coverage to Codecov (Ubuntu+1.21 only)    │
└────────────────────────────────────────────────────────┘
```

## Separate Lint Job
```
┌────────────────────────────────────────────────────────┐
│  Lint Job (Runs in parallel)                          │
│  ├─ Checkout code                                      │
│  ├─ Set up Go 1.21                                     │
│  └─ Run golangci-lint                                  │
│     (continue-on-error)                                │
└────────────────────────────────────────────────────────┘
```

## Success Criteria

Each job succeeds if:
- ✅ Code checks out successfully
- ✅ Dependencies download and verify
- ✅ Code formatting passes (or no formatting issues found)
- ✅ Build completes without errors
- ✅ Binary file is created and verified
- ⚠️  go vet passes (or allowed to fail)
- ⚠️  Tests pass (or allowed to fail if no tests exist)

## Failure Handling

Jobs continue on error for:
- `go vet` - May fail due to GUI dependencies
- `go test` - May fail if no test files exist
- `golangci-lint` - May fail on strict checks

Jobs fail immediately on:
- Checkout failure
- Go setup failure
- Dependency download failure
- Build failure
- Binary verification failure

## Artifacts Generated

For each successful build:
```
OfflineTranscribe-cli-ubuntu-latest-go1.21
OfflineTranscribe-cli-ubuntu-latest-go1.22
OfflineTranscribe-cli-ubuntu-latest-go1.23
OfflineTranscribe-cli-windows-latest-go1.21
OfflineTranscribe-cli-windows-latest-go1.22
OfflineTranscribe-cli-windows-latest-go1.23
OfflineTranscribe-cli-macos-latest-go1.21
OfflineTranscribe-cli-macos-latest-go1.22
OfflineTranscribe-cli-macos-latest-go1.23
```

Each artifact contains:
- Compiled binary (OfflineTranscribe-cli or OfflineTranscribe-cli.exe)
- Available for download for 7 days
- Can be used for testing or distribution

## Performance

Typical workflow run time:
- **Cold run** (no cache): ~5-8 minutes per job
- **Warm run** (with cache): ~2-4 minutes per job
- **Total (9 jobs parallel)**: ~5-8 minutes

Cache effectiveness:
- Go modules cached by: `${{ hashFiles('**/go.sum') }}`
- Build cache location: `~/.cache/go-build` and `~/go/pkg/mod`
- Cache hit significantly reduces dependency download time
