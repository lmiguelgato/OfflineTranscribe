#!/bin/bash#!/bin/bash



echo "Building LocalTTS Go Applications..."echo "Building LocalTTS Standalone Executable..."

echoecho



# Check if Go is installed# Check if Python is installed

if ! command -v go &> /dev/null; thenif ! command -v python3 &> /dev/null; then

    echo "Error: Go is not installed or not in PATH"    echo "Error: Python 3 is not installed or not in PATH"

    echo "Please install Go from https://golang.org"    echo "Please install Python 3.8 or higher"

    exit 1    exit 1

fifi



echo "Building CLI version..."# Install required packages

go build -o LocalTTS-CLI cli.goecho "Installing dependencies..."

if [ -f "LocalTTS-CLI" ]; thenpip3 install -r requirements.txt

    echo "✓ CLI version built successfully: LocalTTS-CLI"pip3 install pyinstaller

else

    echo "✗ Failed to build CLI version"# Create standalone executable

fiecho

echo "Creating standalone executable..."

echopyinstaller --onefile --windowed --name "LocalTTS" main.py

echo "Building Web GUI version..."

go build -o LocalTTS-Web web.go# Check if build was successful

if [ -f "LocalTTS-Web" ]; thenif [ -f "dist/LocalTTS" ]; then

    echo "✓ Web GUI version built successfully: LocalTTS-Web"    echo

else    echo "============================================"

    echo "✗ Failed to build Web GUI version"    echo "Build completed successfully!"

fi    echo

    echo "The standalone executable is located at:"

echo    echo "  dist/LocalTTS"

echo "Building static binary (all dependencies included)..."    echo

go build -ldflags "-s -w" -o LocalTTS-Standalone cli.go    echo "You can now distribute this single file to users."

if [ -f "LocalTTS-Standalone" ]; then    echo "No Python installation required on target machines."

    echo "✓ Standalone version built successfully: LocalTTS-Standalone"    echo "============================================"

elseelse

    echo "✗ Failed to build standalone version"    echo

fi    echo "Build failed. Check the output above for errors."

fi

echo

echo "============================================"echo
echo "Build completed!"
echo
echo "Available executables:"
echo "  LocalTTS-CLI        - Command line interface"
echo "  LocalTTS-Web        - Web browser interface"  
echo "  LocalTTS-Standalone - Optimized standalone version"
echo
echo "These are completely self-contained and can be distributed"
echo "to users without requiring Go installation."
echo "============================================"
echo