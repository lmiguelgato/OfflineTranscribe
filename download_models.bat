@echo off
REM Download script for Whisper models
REM This script downloads the required Whisper models for offline speech-to-text

echo ============================================
echo OfflineTranscribe - Model Download Script
echo ============================================
echo.

REM Create models directory if it doesn't exist
if not exist "models" (
    echo Creating models directory...
    mkdir models
)

echo This script will download Whisper models for offline speech recognition.
echo.
echo Available models:
echo 1. tiny   (~39 MB)  - Fastest, least accurate
echo 2. base   (~142 MB) - Good balance (recommended)
echo 3. small  (~466 MB) - Better accuracy, slower  
echo 4. medium (~1.5 GB) - Best accuracy, slowest
echo.

set /p choice="Which model would you like to download? (1-4) [2]: "
if "%choice%"=="" set choice=2

if "%choice%"=="1" (
    set model=tiny
    set url=https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-tiny.bin
) else if "%choice%"=="2" (
    set model=base
    set url=https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-base.bin
) else if "%choice%"=="3" (
    set model=small
    set url=https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-small.bin
) else if "%choice%"=="4" (
    set model=medium
    set url=https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-medium.bin
) else (
    echo Invalid choice. Using base model.
    set model=base
    set url=https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-base.bin
)

set filename=models\ggml-%model%.bin

echo.
echo Downloading %model% model...
echo URL: %url%
echo File: %filename%
echo.

REM Check if curl is available
curl --version >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo Error: curl is not installed or not in PATH.
    echo Please install curl or download the model manually:
    echo %url%
    echo Save it as: %filename%
    pause
    exit /b 1
)

REM Download the model
curl -L -o "%filename%" "%url%"

if %ERRORLEVEL% EQU 0 (
    echo.
    echo Download completed successfully!
    echo Model saved to: %filename%
    echo.
    echo You can now use OfflineTranscribe with the %model% model.
) else (
    echo.
    echo Download failed! Please check your internet connection and try again.
    echo You can also download manually from: %url%
)

echo.
pause