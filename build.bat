@echo off
echo ============================================
echo OfflineTranscribe - Build Script
echo ============================================
echo.

echo Checking Go installation...
go version
if %ERRORLEVEL% NEQ 0 (
    echo Error: Go is not installed or not in PATH
    echo Please install Go from https://golang.org/dl/
    pause
    exit /b 1
)

echo.
echo Downloading dependencies...
go mod tidy

echo.
echo Building CLI version...
go build -ldflags "-s -w" -o OfflineTranscribe-cli.exe cli.go whisper.go
if %ERRORLEVEL% NEQ 0 (
    echo Error: Failed to build CLI version
    pause
    exit /b 1
)

echo Building web version...
go build -ldflags "-s -w" -o OfflineTranscribe-web.exe web.go whisper.go
if %ERRORLEVEL% NEQ 0 (
    echo Error: Failed to build web version
    pause
    exit /b 1
)

echo.
echo ============================================
echo Build completed successfully!
echo ============================================
echo.
echo Files created:
echo - OfflineTranscribe-cli.exe  (Command-line interface)
echo - OfflineTranscribe-web.exe  (Web interface)
echo.
echo Next steps:
echo 1. Whisper executable found in whisper-bin-x64\Release\
echo 2. Download speech recognition models: run download_models.bat
echo 3. Start transcribing real audio files!
echo.
echo Note: Use actual audio files (WAV, MP3, etc.) - not text files
echo.
pause