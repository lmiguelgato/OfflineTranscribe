@echo off
echo Building OfflineTranscribe Go Applications...
echo.

REM Check if Go is installed
go version >nul 2>&1
if errorlevel 1 (
    echo Error: Go is not installed or not in PATH
    echo Please install Go from https://golang.org
    pause
    exit /b 1
)

echo Building CLI version...
go build -o OfflineTranscribe-CLI.exe cli.go
if exist "OfflineTranscribe-CLI.exe" (
    echo ✓ CLI version built successfully: OfflineTranscribe-CLI.exe
) else (
    echo ✗ Failed to build CLI version
)

echo.
echo Building Web GUI version...
go build -o OfflineTranscribe-Web.exe web.go
if exist "OfflineTranscribe-Web.exe" (
    echo ✓ Web GUI version built successfully: OfflineTranscribe-Web.exe
) else (
    echo ✗ Failed to build Web GUI version
)

echo.
echo Building static binary (all dependencies included)...
go build -ldflags "-s -w" -o OfflineTranscribe-Standalone.exe cli.go
if exist "OfflineTranscribe-Standalone.exe" (
    echo ✓ Standalone version built successfully: OfflineTranscribe-Standalone.exe
) else (
    echo ✗ Failed to build standalone version
)

echo.
echo ============================================
echo Build completed!
echo.
echo Available executables:
echo   OfflineTranscribe-CLI.exe        - Command line interface
echo   OfflineTranscribe-Web.exe        - Web browser interface
echo   OfflineTranscribe-Standalone.exe - Optimized standalone version
echo.
echo These are completely self-contained and can be distributed
echo to users without requiring Go installation.
echo ============================================

echo.
pause