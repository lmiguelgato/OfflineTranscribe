# OfflineTranscribe - Offline Speech-to-Text Tool

A standalone, offline speech-to-text transcription tool built with Go that processes audio files and generates timestamped transcriptions without requiring internet connectivity.

## Features

- **Fully Offline**: Works without internet connection
- **Multiple Interfaces**: Command-line and web browser GUI
- **High Performance**: Built with Go for speed and efficiency
- **Timestamped Output**: Generates precise timestamps for easy navigation
- **Flexible Granularity**: Choose between word-level or sentence-level timestamps
- **User-Friendly**: Simple interfaces for both technical and non-technical users
- **Multiple Formats**: Supports WAV, MP3, MP4, and other audio formats
- **Standalone**: Single executable files with no dependencies
- **Cross-Platform**: Works on Windows, macOS, and Linux

## Quick Start

### Download and Run (Easiest)
1. Download the appropriate executable for your system from the releases page
2. **For Web GUI**: Double-click `OfflineTranscribe-Web.exe` and open http://localhost:8080 in your browser
3. **For CLI**: Run `OfflineTranscribe-CLI.exe` in terminal/command prompt

### Available Versions
- **OfflineTranscribe-Web.exe** - Modern web interface (recommended for most users)
- **OfflineTranscribe-CLI.exe** - Command-line interface for advanced users
- **OfflineTranscribe-Standalone.exe** - Optimized single executable

## Usage

### Web Interface (Recommended)
```bash
# Start the web server
./OfflineTranscribe-Web.exe

# Open your browser to http://localhost:8080
# Drag and drop your audio file or click to browse
# Select model size and timestamp type
# Click "Process Audio" and wait for results
```

### Command Line Interface
```bash
# Interactive mode
./OfflineTranscribe-CLI.exe

# Direct processing
./OfflineTranscribe-CLI.exe recording.wav
./OfflineTranscribe-CLI.exe recording.wav -model small -type sentence
./OfflineTranscribe-CLI.exe recording.wav -output transcript.txt
```

### CLI Options
- `-model <size>`: Model size (`tiny`, `base`, `small`, `medium`)
- `-type <type>`: Timestamp type (`word`, `sentence`) 
- `-output <file>`: Output file path

## Output Format

### Word-level timestamps:
```
[00:00:12.340] Hello
[00:00:12.580] there,
[00:00:13.120] how
[00:00:13.340] are
[00:00:13.560] you
```

### Sentence-level timestamps:
```
[00:00:12.340 - 00:00:15.120] Hello there, how are you doing today?

[00:00:15.340 - 00:00:18.560] I hope everything is going well for you.
```

## Model Sizes

| Model  | Size  | Speed | Accuracy | Recommended For |
|--------|-------|-------|----------|-----------------|
| tiny   | 39MB  | Fastest | Basic | Quick testing |
| base   | 74MB  | Fast | Good | Most users |
| small  | 244MB | Medium | Better | Quality focused |
| medium | 769MB | Slow | Best | Maximum accuracy |

## Building from Source

### Prerequisites
- Go 1.21 or higher

### Build All Versions
```bash
# Windows
build.bat

# Linux/macOS
chmod +x build.sh
./build.sh
```

### Manual Build
```bash
# CLI version
go build -o OfflineTranscribe-CLI cli.go

# Web version  
go build -o OfflineTranscribe-Web web.go

# Optimized standalone
go build -ldflags "-s -w" -o OfflineTranscribe-Standalone cli.go
```

## Development

### Project Structure
```
OfflineTranscribe/
├── cli.go              # Command-line interface
├── web.go              # Web server and API
├── index.html          # Web interface
├── build.bat           # Windows build script
├── build.sh            # Linux/macOS build script
├── go.mod              # Go dependencies
└── README.md           # Documentation
```

### Key Components
- **CLI Interface**: Interactive and direct command-line processing
- **Web Interface**: Modern drag-and-drop browser interface
- **Audio Processing**: Placeholder for Whisper.cpp integration
- **Timestamp Generation**: Word and sentence-level timing
- **File Handling**: Temporary file management and cleanup

## Roadmap

### Current (Demo Version)
- ✅ CLI and Web interfaces
- ✅ File upload and processing simulation
- ✅ Timestamp formatting
- ✅ Cross-platform builds

### Next Steps (Production Version)
- [ ] Integrate actual Whisper.cpp for real transcription
- [ ] Add model downloading and caching
- [ ] Implement progress tracking for long files
- [ ] Add batch processing capabilities
- [ ] Include audio format conversion

## Technical Details

- **Language**: Go 1.21+
- **Web Framework**: Standard library HTTP server
- **Frontend**: Vanilla JavaScript with modern CSS
- **Audio Processing**: Ready for Whisper.cpp integration
- **Build**: Single static binaries with no external dependencies

## Deployment

The built executables are completely self-contained:
- No Go installation required on target machines
- No external dependencies or libraries needed
- Can be distributed as single files
- Work on fresh operating system installations

## License

[Add your license information here]

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## Support

For issues, questions, or contributions, please [create an issue](link-to-issues) on the repository.

---

**Note**: This is currently a demonstration version with simulated transcription. Integration with actual Whisper.cpp speech recognition will provide real audio processing capabilities.