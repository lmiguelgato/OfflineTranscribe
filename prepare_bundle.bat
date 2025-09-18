@echo off
echo ============================================
echo OfflineTranscribe - Bundle Preparation
echo ============================================
echo.

REM Create bundle directory structure
echo Creating bundle directory structure...
if not exist "bundle" mkdir bundle
if not exist "bundle\resources" mkdir bundle\resources
if not exist "bundle\resources\models" mkdir bundle\resources\models
if not exist "bundle\resources\whisper" mkdir bundle\resources\whisper

echo.
echo Copying Whisper executable and dependencies...
copy "whisper-bin-x64\Release\whisper-cli.exe" "bundle\resources\whisper\" >nul 2>&1
copy "whisper-bin-x64\Release\*.dll" "bundle\resources\whisper\" >nul 2>&1

if not exist "bundle\resources\whisper\whisper-cli.exe" (
    echo Error: Whisper executable not found!
    echo Please ensure whisper-cli.exe is in whisper-bin-x64\Release\
    pause
    exit /b 1
)

echo.
echo Copying AI models...
copy "models\*.bin" "bundle\resources\models\" >nul 2>&1

REM Check if we have at least one model
if not exist "bundle\resources\models\*.bin" (
    echo Warning: No AI models found in models\ directory
    echo Downloading base model...
    if not exist "models" mkdir models
    powershell -Command "& { Invoke-WebRequest -Uri 'https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-base.bin' -OutFile 'models\ggml-base.bin' }"
    copy "models\ggml-base.bin" "bundle\resources\models\" >nul 2>&1
)

echo.
echo Copying web interface files...
copy "index.html" "bundle\resources\" >nul 2>&1

echo.
echo Creating embedded resource directory listings...
echo // Bundle resource manifest > bundle\resources\manifest.txt
echo Models: >> bundle\resources\manifest.txt
dir /b "bundle\resources\models\*.bin" >> bundle\resources\manifest.txt
echo. >> bundle\resources\manifest.txt
echo Whisper files: >> bundle\resources\manifest.txt
dir /b "bundle\resources\whisper\*.*" >> bundle\resources\manifest.txt

echo.
echo ============================================
echo Bundle preparation completed!
echo ============================================
echo.
echo Bundle contents:
echo - bundle\resources\models\     (AI models)
echo - bundle\resources\whisper\   (Whisper executable + DLLs)
echo - bundle\resources\index.html (Web interface)
echo.
echo Ready for embedding in Go application.
echo.
pause