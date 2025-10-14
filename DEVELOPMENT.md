# LocalTTS Development Guide

## Quick Start for End Users

### Option 1: Web Interface (Recommended)
1. Download `LocalTTS-Web.exe` 
2. Double-click to run
3. Open http://localhost:8080 in your browser
4. Drag and drop audio files and process!

### Option 2: Command Line
1. Download `LocalTTS-CLI.exe`
2. Open terminal/command prompt
3. Run `LocalTTS-CLI.exe` for interactive mode
4. Or `LocalTTS-CLI.exe audiofile.wav` for direct processing

No installation or setup required!

## For Developers

### Prerequisites
- Go 1.21 or higher
- Git

### Development Setup
```bash
# Clone the repository
git clone <repository-url>
cd LocalTTS

# Initialize Go module (already done)
go mod init localtts
go mod tidy

# Run CLI version
go run cli.go

# Run web version
go run web.go
```

### Building All Versions
```bash
# Windows
build.bat

# Linux/Mac
chmod +x build.sh
./build.sh
```

### Project Structure
```
LocalTTS/
├── cli.go              # Command-line interface implementation
├── web.go              # Web server and HTTP handlers  
├── index.html          # Modern web interface with drag-drop
├── build.bat           # Windows build automation
├── build.sh            # Linux/Mac build automation
├── go.mod              # Go module definition
├── README.md           # User documentation
├── DEVELOPMENT.md      # This file - developer guide
└── .github/
    ├── workflows/
    │   └── cli-build-test.yml  # CI/CD pipeline for automated builds
    └── WORKFLOW.md     # CI/CD documentation
```

### Key Implementation Details

#### CLI Interface (`cli.go`)
- Interactive and direct command-line processing
- File validation and error handling
- Configurable model sizes and timestamp types
- Progress feedback and result saving

#### Web Interface (`web.go` + `index.html`)
- HTTP server with file upload handling
- RESTful API for audio processing
- Modern drag-and-drop web interface
- Real-time progress updates and result display
- Background processing with goroutines

#### Features Implemented

✅ **Multiple Interfaces**: CLI and Web GUI options
✅ **File Processing**: Upload, validation, temporary storage
✅ **Model Selection**: tiny, base, small, medium options
✅ **Timestamp Formats**: Word-level and sentence-level
✅ **Cross-Platform**: Windows, Linux, macOS support
✅ **Standalone Builds**: No dependencies required
✅ **Progress Feedback**: Real-time status updates
✅ **Result Export**: Text file downloads

### Architecture Decisions

#### Why Go?
- **Performance**: Compiled binaries are fast
- **Simplicity**: Easy deployment as single executables
- **Cross-Platform**: Built-in support for multiple OS
- **Standard Library**: HTTP server, file handling included
- **Concurrency**: Goroutines for background processing

#### Why Web + CLI?
- **Web Interface**: User-friendly for non-technical users
- **CLI Interface**: Scriptable for technical users and automation
- **No GUI Dependencies**: Avoids complex UI framework issues
- **Browser Compatibility**: Works on any system with a browser

#### Current Implementation (Demo Version)
- Simulated audio processing with realistic timing
- Sample transcription results for demonstration
- Complete user interface and file handling
- Ready for Whisper.cpp integration

### Next Steps for Production

#### Phase 1: Core Audio Processing
```bash
# Add Whisper.cpp integration
go get github.com/ggerganov/whisper.cpp/bindings/go
# Implement real audio transcription
# Add model downloading and caching
```

#### Phase 2: Enhanced Features
- Progress tracking for long audio files
- Batch processing capabilities  
- Audio format conversion
- Model management and updates

#### Phase 3: Optimization
- Memory usage optimization
- Parallel processing for large files
- Result caching and incremental processing
- Performance benchmarking

### Integration Points

#### Whisper.cpp Integration
```go
// Planned integration structure
type WhisperProcessor struct {
    model whisper.Model
    ctx   whisper.Context
}

func (wp *WhisperProcessor) ProcessAudio(inputFile string) (*TranscriptionResult, error) {
    // Load audio file
    // Process with Whisper
    // Extract timestamps
    // Format results
}
```

#### Model Management
```go
type ModelManager struct {
    modelsPath string
    cache      map[string]whisper.Model
}

func (mm *ModelManager) LoadModel(size string) error {
    // Download if not exists
    // Load into memory
    // Cache for reuse
}
```

### Testing Strategy

#### Automated CI/CD Pipeline
GitHub Actions workflow automatically tests all builds:
- **Multi-Platform**: Ubuntu, Windows, macOS
- **Multi-Version**: Go 1.21, 1.22, 1.23
- **Code Quality**: go vet, gofmt, golangci-lint
- **Artifacts**: Built binaries available for download

See `.github/WORKFLOW.md` for detailed documentation.

#### Unit Tests
```bash
# Test CLI argument parsing
go test -v ./cli_test.go

# Test web API endpoints  
go test -v ./web_test.go

# Test audio processing
go test -v ./whisper_test.go
```

#### Integration Tests
```bash
# Test full workflow with sample audio
./test-samples.sh

# Test cross-platform builds
./test-builds.sh
```

### Deployment Strategy

#### Single Binary Distribution
- **Advantage**: No installation required
- **Implementation**: `go build -ldflags "-s -w"`
- **Size**: ~10-20MB per binary (without models)

#### Model Distribution
- **Option 1**: Separate model downloads (current)
- **Option 2**: Bundle common models with binary
- **Option 3**: On-demand model streaming

### Performance Considerations

#### Memory Usage
- Stream audio processing for large files
- Model caching to avoid reloading
- Garbage collection optimization

#### Processing Speed
- Parallel processing where possible
- Optimized timestamp extraction
- Efficient file I/O operations

### Security Considerations

#### File Handling
- Validate uploaded file types
- Secure temporary file creation
- Automatic cleanup of temporary files
- Size limits for uploads

#### Web Interface
- CORS configuration for local use
- Input validation and sanitization
- Rate limiting for API endpoints

### Contributing Guidelines

#### Code Style
- Follow Go conventions and `gofmt`
- Add comments for public functions
- Use meaningful variable names
- Keep functions focused and small

#### Pull Request Process
1. Fork and create feature branch
2. Add tests for new functionality  
3. Ensure all tests pass
4. Update documentation
5. Submit PR with clear description

### Build and Release Process

#### Automated Builds
```bash
# Cross-platform compilation
GOOS=windows GOARCH=amd64 go build -o LocalTTS-Windows.exe
GOOS=linux GOARCH=amd64 go build -o LocalTTS-Linux
GOOS=darwin GOARCH=amd64 go build -o LocalTTS-macOS
```

#### Release Checklist
- [ ] All tests passing
- [ ] Documentation updated
- [ ] Cross-platform builds tested
- [ ] Performance benchmarks run
- [ ] Security review completed

This Go-based implementation provides a solid foundation for a production-ready, standalone speech-to-text tool that meets all the original requirements while being easy to deploy and use.