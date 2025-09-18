@echo off
echo ============================================
echo OfflineTranscribe - Bundle Build Script
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
echo Step 1: Preparing bundle resources...
call prepare_bundle.bat
if %ERRORLEVEL% NEQ 0 (
    echo Error: Failed to prepare bundle resources
    pause
    exit /b 1
)

echo.
echo Step 2: Downloading dependencies...
go mod tidy

echo.
echo Step 3: Building self-contained CLI version...
go build -ldflags "-s -w" -o OfflineTranscribe-Bundle-CLI.exe cli.go whisper.go resources.go
if %ERRORLEVEL% NEQ 0 (
    echo Error: Failed to build CLI version
    pause
    exit /b 1
)

echo Building self-contained web version...
go build -ldflags "-s -w" -o OfflineTranscribe-Bundle-Web.exe web.go whisper.go resources.go
if %ERRORLEVEL% NEQ 0 (
    echo Error: Failed to build web version
    pause
    exit /b 1
)

echo.
echo Step 4: Creating distribution package...
if not exist "dist" mkdir dist
copy OfflineTranscribe-Bundle-CLI.exe dist\ >nul
copy OfflineTranscribe-Bundle-Web.exe dist\ >nul
copy README.md dist\ >nul 2>&1

echo.
echo Creating distribution README...
echo OfflineTranscribe - Self-Contained Distribution > dist\DISTRIBUTION_README.txt
echo ================================================= >> dist\DISTRIBUTION_README.txt
echo. >> dist\DISTRIBUTION_README.txt
echo This package contains fully self-contained executables that include: >> dist\DISTRIBUTION_README.txt
echo - Whisper AI speech recognition models >> dist\DISTRIBUTION_README.txt
echo - Whisper.cpp executable and dependencies >> dist\DISTRIBUTION_README.txt
echo - Web interface files >> dist\DISTRIBUTION_README.txt
echo. >> dist\DISTRIBUTION_README.txt
echo No additional downloads or installations required! >> dist\DISTRIBUTION_README.txt
echo. >> dist\DISTRIBUTION_README.txt
echo Usage: >> dist\DISTRIBUTION_README.txt
echo   OfflineTranscribe-Bundle-CLI.exe recording.wav >> dist\DISTRIBUTION_README.txt
echo   OfflineTranscribe-Bundle-Web.exe  (then open http://localhost:8080) >> dist\DISTRIBUTION_README.txt
echo. >> dist\DISTRIBUTION_README.txt
echo Supported audio formats: WAV, MP3, MP4, FLAC, M4A, OGG >> dist\DISTRIBUTION_README.txt
echo Available models: tiny (fast), base (balanced), small/medium (accurate) >> dist\DISTRIBUTION_README.txt

echo.
echo ============================================
echo Bundle build completed successfully!
echo ============================================
echo.
echo Self-contained files created:
echo - OfflineTranscribe-Bundle-CLI.exe  (Fully portable CLI)
echo - OfflineTranscribe-Bundle-Web.exe  (Fully portable Web interface)
echo.
echo Distribution package: dist\
echo.
echo These executables can be shared and run on any Windows machine
echo without requiring any additional downloads or installations!
echo.
echo File sizes:
for %%f in (OfflineTranscribe-Bundle-*.exe) do (
    echo - %%f: 
    dir %%f | findstr %%f
)
echo.
pause