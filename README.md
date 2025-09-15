# OfflineTranscribe - Offline Speech-to-Text Tool

A standalone, offline speech-to-text transcription tool that works without an internet connection. Convert audio files to text with precise timestamps for easy navigation and reference.

## Features

- **Completely Offline**: No internet connection required after setup
- **Standalone Executable**: No dependencies to install for end users
- **Multiple Interfaces**: Command-line and web browser interfaces
- **Precise Timestamps**: Sentence-level timing information for easy navigation
- **Multiple Model Sizes**: Choose between speed and accuracy
- **Cross-Platform**: Works on Windows, macOS, and Linux
- **Multiple Audio Formats**: Supports WAV, MP3, MP4, and more

## Quick Start

### 1. Download Required Files

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

### 2. Usage

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
./OfflineTranscribe-cli.exe recording.wav -model small -output transcript.txt
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

## Building from Source

**Prerequisites:**
- Go 1.21 or later
- Git

**Build Steps:**
```bash
# Clone the repository
git clone <repository-url>
cd OfflineTranscribe

# Download dependencies
go mod tidy

# Build CLI version
go build -o OfflineTranscribe-cli.exe cli.go whisper.go

# Build web version  
go build -o OfflineTranscribe-web.exe web.go whisper.go

# For other platforms
GOOS=linux GOARCH=amd64 go build -o OfflineTranscribe-cli-linux cli.go whisper.go
GOOS=darwin GOARCH=amd64 go build -o OfflineTranscribe-cli-mac cli.go whisper.go
```

## Project Structure

```
OfflineTranscribe/
├── cli.go                 # Command-line interface
├── web.go                 # Web server interface  
├── whisper.go             # Whisper integration
├── index.html             # Web interface frontend
├── download_models.bat    # Windows model download script
├── download_models.sh     # Unix model download script
├── go.mod                 # Go module definition
├── README.md              # This file
├── models/                # Downloaded Whisper models
│   ├── ggml-tiny.bin
│   ├── ggml-base.bin
│   ├── ggml-small.bin
│   └── ggml-medium.bin
└── whisper.exe            # Whisper executable (download separately)
```

## Model Comparison

| Model  | Size   | Speed | Accuracy | RAM Usage | Best For |
|--------|--------|-------|----------|-----------|----------|
| Tiny   | ~39 MB | Fast  | Basic    | ~390 MB   | Quick drafts, real-time |
| Base   | ~142 MB| Good  | Good     | ~500 MB   | General use (recommended) |
| Small  | ~466 MB| Slow  | Better   | ~750 MB   | Important recordings |
| Medium | ~1.5 GB| Slower| Best     | ~1.5 GB   | Maximum accuracy needed |

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

**"whisper executable not found"**
- Download whisper.cpp from the releases page
- Place the executable in the OfflineTranscribe directory
- Ensure it's named `whisper.exe` (Windows) or `whisper` (Unix)

**"model file not found"**
- Run the download script: `download_models.bat` or `./download_models.sh`
- Or manually download models to the `models/` directory

**"transcription failed"**
- Check that your audio file is in a supported format
- Ensure the audio file isn't corrupted
- Try with a different model size

**Web interface not loading**
- Check that `index.html` is in the same directory as the executable
- Ensure port 8080 isn't blocked by firewall
- Try a different port: `./OfflineTranscribe-web.exe 3000`

## Performance Tips

- Use smaller models for faster processing
- Convert audio to WAV format for best compatibility
- Use shorter audio segments (under 30 minutes) for better performance
- Close other applications to free up RAM when using larger models

## License

This project is open source. See LICENSE file for details.

Whisper.cpp is developed by Georgi Gerganov and contributors under MIT license.

## Contributing

Contributions are welcome! Please feel free to submit pull requests or open issues for bugs and feature requests.