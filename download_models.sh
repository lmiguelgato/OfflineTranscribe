#!/bin/bash

# Download script for Whisper models
# This script downloads the required Whisper models for offline speech-to-text

echo "============================================"
echo "OfflineTranscribe - Model Download Script"
echo "============================================"
echo

# Create models directory if it doesn't exist
if [ ! -d "models" ]; then
    echo "Creating models directory..."
    mkdir -p models
fi

echo "This script will download Whisper models for offline speech recognition."
echo
echo "Available models:"
echo "1. tiny   (~39 MB)  - Fastest, least accurate"
echo "2. base   (~142 MB) - Good balance (recommended)"
echo "3. small  (~466 MB) - Better accuracy, slower"
echo "4. medium (~1.5 GB) - Best accuracy, slowest"
echo

read -p "Which model would you like to download? (1-4) [2]: " choice
choice=${choice:-2}

case $choice in
    1)
        model="tiny"
        url="https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-tiny.bin"
        ;;
    2)
        model="base"
        url="https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-base.bin"
        ;;
    3)
        model="small"
        url="https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-small.bin"
        ;;
    4)
        model="medium"
        url="https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-medium.bin"
        ;;
    *)
        echo "Invalid choice. Using base model."
        model="base"
        url="https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-base.bin"
        ;;
esac

filename="models/ggml-${model}.bin"

echo
echo "Downloading ${model} model..."
echo "URL: ${url}"
echo "File: ${filename}"
echo

# Check if curl is available
if ! command -v curl &> /dev/null; then
    echo "Error: curl is not installed."
    echo "Please install curl or download the model manually:"
    echo "${url}"
    echo "Save it as: ${filename}"
    exit 1
fi

# Download the model
curl -L -o "${filename}" "${url}"

if [ $? -eq 0 ]; then
    echo
    echo "Download completed successfully!"
    echo "Model saved to: ${filename}"
    echo
    echo "You can now use OfflineTranscribe with the ${model} model."
else
    echo
    echo "Download failed! Please check your internet connection and try again."
    echo "You can also download manually from: ${url}"
fi

echo