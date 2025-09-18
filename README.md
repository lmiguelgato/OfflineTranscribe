# OfflineTranscribe - Offline Speech-to-Text Tool

A standalone, offline speech-to-text transcription tool that works without an internet connection. Convert audio files to text with precise timestamps for easy navigation and reference.

## Features

- **Completely Offline**: No internet connection required after setup
- **Self-Contained Bundle**: Single executable with everything included (recommended)
- **Multiple Interfaces**: Command-line and web browser interfaces
- **Precise Timestamps**: Sentence-level timing information for easy navigation
- **Multiple Model Sizes**: Choose between speed and accuracy
- **Cross-Platform**: Works on Windows, macOS, and Linux
- **Multiple Audio Formats**: Supports WAV, MP3, MP4, FLAC, M4A, OGG

## Quick Start (Recommended)

### Option 1: Self-Contained Bundle 🚀

**For immediate use with zero setup:**

1. **Download the bundle** (available from releases or build locally):
   - `OfflineTranscribe-Bundle-CLI.exe` (232 MB) - Command-line interface
   - `OfflineTranscribe-Bundle-Web.exe` (236 MB) - Web interface

2. **Run immediately** - no installation, downloads, or setup required:
   ```bash
   # CLI Usage
   OfflineTranscribe-Bundle-CLI.exe recording.wav
   OfflineTranscribe-Bundle-CLI.exe recording.wav -model tiny -output transcript.txt
   
   # Web Interface
   OfflineTranscribe-Bundle-Web.exe
   # Then open http://localhost:8080 in your browser
   ```

**What's included in the bundle:**
- ✅ Whisper.cpp executable and all DLLs
- ✅ AI models (tiny and base) ready to use
- ✅ Web interface files
- ✅ All dependencies embedded

**Perfect for:**
- Sharing with colleagues who need immediate use
- Running on systems without admin rights
- Air-gapped or restricted environments
- Quick testing and evaluation

### Option 2: Manual Setup (Advanced Users)

### Option 2: Manual Setup (Advanced Users)

**For developers or users who want to customize models:**

#### 1. Download Required Files

**Download Whisper Executable:**
- Go to [Whisper.cpp Releases](https://github.com/ggerganov/whisper.cpp/releases)
- Download the appropriate version for your system
- Extract the executable to the OfflineTranscribe directory or subdirectory

**Supported locations for whisper executable:**
- `whisper.exe` or `whisper` (in main directory)
- `whisper-bin-x64/Release/whisper-cli.exe` (Windows build directory)
- `whisper-bin-x64/Release/main.exe` (Alternative Windows build)
- `bin/whisper.exe` or `bin/whisper` (bin subdirectory)

**Download Models:**
Run the included download script:
```bash
# Windows
download_models.bat

# Linux/Mac
chmod +x download_models.sh
./download_models.sh
```

Or download manually:
- [Tiny Model (~39 MB)](https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-tiny.bin) - Fastest, least accurate
- [Base Model (~142 MB)](https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-base.bin) - Good balance (recommended)
- [Small Model (~466 MB)](https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-small.bin) - Better accuracy, slower
- [Medium Model (~1.5 GB)](https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-medium.bin) - Best accuracy, slowest

Save models in the `models/` directory.

#### 2. Usage

**Interactive Mode (Easiest):**
```bash
./OfflineTranscribe-cli.exe    # Windows
./OfflineTranscribe-cli        # Linux/Mac
```

**Command Line Mode:**
```bash
# Basic usage
./OfflineTranscribe-cli.exe recording.wav

# With options
./OfflineTranscribe-cli.exe recording.wav -model tiny -output transcript.txt
```

**Web Interface:**
```bash
./OfflineTranscribe-web.exe    # Windows
./OfflineTranscribe-web        # Linux/Mac
```
Then open http://localhost:8080 in your browser.

## Command Line Options

- `-model <size>`: Model size (tiny, base, small, medium) - default: base
- `-output <file>`: Output file path - default: `<input>_transcription.txt`

## Output Format

**Sentence-level timestamps:**
```
[00:00:01.240 - 00:00:03.680] Hello there, this is a sample transcription.

[00:00:04.120 - 00:00:06.200] Each sentence has its own time range.
```

## Building Your Own Bundle

**To create self-contained executables with embedded dependencies:**

### Prerequisites
- Go 1.21 or later
- Git

### Steps
```bash
# Clone the repository
git clone <repository-url>
cd OfflineTranscribe

# Download dependencies
go mod tidy

# Build the self-contained bundle
./build_bundle.bat          # Windows
./build_bundle.sh           # Linux/Mac (coming soon)
```

This creates:
- `OfflineTranscribe-Bundle-CLI.exe` - Self-contained CLI
- `OfflineTranscribe-Bundle-Web.exe` - Self-contained web interface
- `dist/` folder with distribution files

### Manual Build (Advanced)
```bash
# Prepare resources
./prepare_bundle.bat

# Build regular versions
go build -o OfflineTranscribe-cli.exe cli.go whisper.go
go build -o OfflineTranscribe-web.exe web.go whisper.go

# Build self-contained versions
go build -o OfflineTranscribe-Bundle-CLI.exe cli.go whisper.go resources.go
go build -o OfflineTranscribe-Bundle-Web.exe web.go whisper.go resources.go

# Cross-platform builds
GOOS=linux GOARCH=amd64 go build -o OfflineTranscribe-Bundle-CLI-linux cli.go whisper.go resources.go
GOOS=darwin GOARCH=amd64 go build -o OfflineTranscribe-Bundle-CLI-mac cli.go whisper.go resources.go
```

## Project Structure

```
OfflineTranscribe/
├── cli.go                 # Command-line interface
├── web.go                 # Web server interface  
├── whisper.go             # Whisper integration
├── resources.go           # Embedded resource management
├── index.html             # Web interface frontend
├── build_bundle.bat       # Bundle build script (Windows)
├── prepare_bundle.bat     # Resource preparation script
├── download_models.bat    # Windows model download script
├── download_models.sh     # Unix model download script
├── go.mod                 # Go module definition
├── README.md              # This file
├── models/                # Downloaded Whisper models (manual setup)
│   ├── ggml-tiny.bin
│   ├── ggml-base.bin
│   ├── ggml-small.bin
│   └── ggml-medium.bin
├── bundle/                # Bundle resources (auto-created)
│   └── resources/
│       ├── models/        # Embedded models
│       ├── whisper/       # Embedded Whisper executable
│       └── index.html     # Embedded web interface
├── dist/                  # Distribution package (auto-created)
│   ├── OfflineTranscribe-Bundle-CLI.exe
│   ├── OfflineTranscribe-Bundle-Web.exe
│   └── DISTRIBUTION_README.txt
└── whisper-bin-x64/       # Whisper executable (manual setup)
    └── Release/
        └── whisper-cli.exe
```

## Model Comparison

| Model  | Size   | Speed | Accuracy | RAM Usage | Best For | Bundle |
|--------|--------|-------|----------|-----------|----------|---------|
| Tiny   | ~39 MB | Fast  | Basic    | ~390 MB   | Quick drafts, real-time | ✅ Included |
| Base   | ~142 MB| Good  | Good     | ~500 MB   | General use (recommended) | ✅ Included |
| Small  | ~466 MB| Slow  | Better   | ~750 MB   | Important recordings | Manual setup |
| Medium | ~1.5 GB| Slower| Best     | ~1.5 GB   | Maximum accuracy needed | Manual setup |

*Bundle = Available in self-contained executables*

## Supported Audio Formats

- WAV (recommended for best quality)
- MP3
- MP4
- M4A
- FLAC
- OGG

## Use Cases

- **Meeting Transcription**: Convert recorded meetings to searchable text
- **Interview Analysis**: Transcribe interviews with precise timestamps
- **Lecture Notes**: Convert recorded lectures to text with time references
- **Podcast Transcription**: Create text versions of audio content
- **Accessibility**: Generate captions and transcripts for audio content
- **Content Creation**: Extract quotes and segments from longer recordings

## Privacy & Security

- **Completely Offline**: Your audio never leaves your computer
- **No Data Collection**: No telemetry or usage tracking
- **Local Processing**: All transcription happens on your machine
- **No Internet Required**: Works in air-gapped environments

## Troubleshooting

### Bundle Version Issues

**"Failed to initialize resources"**
- The bundle executables are large (~230MB) and may take a moment to start
- Ensure you have sufficient disk space for temporary file extraction
- Check antivirus software isn't blocking execution

### Manual Setup Issues

**"whisper executable not found"**
- Download whisper.cpp from the releases page
- Place the executable in the OfflineTranscribe directory
- Ensure it's named `whisper.exe` (Windows) or `whisper` (Unix)

**"model file not found"**
- Run the download script: `download_models.bat` or `./download_models.sh`
- Or manually download models to the `models/` directory

### General Issues

**"transcription failed"**
- Check that your audio file is in a supported format
- Ensure the audio file isn't corrupted
- Try with a different model size

**Web interface not loading**
- Ensure port 8080 isn't blocked by firewall
- Try a different port: `./OfflineTranscribe-web.exe 3000`
- For bundle version: The web interface is embedded and extracted automatically

**Performance issues**
- Use smaller models for faster processing
- Convert audio to WAV format for best compatibility
- Use shorter audio segments (under 30 minutes) for better performance
- Close other applications to free up RAM when using larger models

## Distribution & Sharing

### Sharing the Bundle
The self-contained bundles are perfect for distribution:
- **Email**: Bundle files can be shared directly (note: large file size)
- **USB Drive**: Copy to portable drives for offline use
- **Network Share**: Place on shared folders for team access
- **Air-Gapped Systems**: Works without internet connectivity

### File Sizes
- `OfflineTranscribe-Bundle-CLI.exe`: ~232 MB
- `OfflineTranscribe-Bundle-Web.exe`: ~236 MB

Large sizes include:
- Complete Whisper AI models (142MB for base model)
- Whisper.cpp executable and all required libraries
- Web interface and runtime dependencies

## License

This project is open source. See LICENSE file for details.

Whisper.cpp is developed by Georgi Gerganov and contributors under MIT license.

## Contributing

Contributions are welcome! Please feel free to submit pull requests or open issues for bugs and feature requests.