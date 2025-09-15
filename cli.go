package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type OfflineTranscribe struct {
	currentFile string
}

func NewOfflineTranscribe() *OfflineTranscribe {
	return &OfflineTranscribe{}
}

func (ot *OfflineTranscribe) processAudio(inputFile, modelSize, timestampType string) (string, error) {
	fmt.Printf("Processing audio file: %s\n", inputFile)
	fmt.Printf("Model size: %s\n", modelSize)
	fmt.Printf("Timestamp type: %s\n", timestampType)
	
	// Check if file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return "", fmt.Errorf("file does not exist: %s", inputFile)
	}
	
	fmt.Printf("Loading Whisper model: %s...\n", modelSize)
	
	transcriber := NewWhisperTranscriber()
	defer transcriber.Close()
	
	// Load the model
	if err := transcriber.LoadModel(modelSize); err != nil {
		return "", fmt.Errorf("failed to load model: %v", err)
	}
	
	fmt.Println("Transcribing audio...")
	
	// Transcribe the audio
	result, err := transcriber.TranscribeFile(inputFile, timestampType)
	if err != nil {
		return "", fmt.Errorf("transcription failed: %v", err)
	}
	
	// Format the results
	formattedOutput := transcriber.FormatResults(result, timestampType)
	
	fmt.Println("Transcription complete!")
	return formattedOutput, nil
}

func (ot *OfflineTranscribe) saveResults(results, outputFile string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer file.Close()
	
	_, err = file.WriteString(results)
	if err != nil {
		return fmt.Errorf("failed to write results: %v", err)
	}
	
	return nil
}

func (ot *OfflineTranscribe) interactive() {
	fmt.Println("===========================================")
	fmt.Println("OfflineTranscribe - Offline Speech-to-Text Tool")
	fmt.Println("===========================================")
	fmt.Println()
	
	scanner := bufio.NewScanner(os.Stdin)
	
	// Get input file
	fmt.Print("Enter path to audio file (WAV, MP3, MP4): ")
	scanner.Scan()
	inputFile := strings.TrimSpace(scanner.Text())
	
	// Get model size
	fmt.Println("\nModel sizes:")
	fmt.Println("1. tiny   - Fastest, least accurate")
	fmt.Println("2. base   - Good balance (recommended)")
	fmt.Println("3. small  - Better accuracy, slower")
	fmt.Println("4. medium - Best accuracy, slowest")
	fmt.Print("Choose model (1-4) [2]: ")
	scanner.Scan()
	choice := strings.TrimSpace(scanner.Text())
	if choice == "" {
		choice = "2"
	}
	
	modelSizes := map[string]string{
		"1": "tiny",
		"2": "base",
		"3": "small",
		"4": "medium",
	}
	modelSize := modelSizes[choice]
	if modelSize == "" {
		modelSize = "base"
	}
	
	// Get timestamp type
	fmt.Println("\nTimestamp granularity:")
	fmt.Println("1. word     - Individual word timestamps")
	fmt.Println("2. sentence - Sentence-level timestamps")
	fmt.Print("Choose granularity (1-2) [1]: ")
	scanner.Scan()
	choice = strings.TrimSpace(scanner.Text())
	if choice == "" {
		choice = "1"
	}
	
	timestampTypes := map[string]string{
		"1": "word",
		"2": "sentence",
	}
	timestampType := timestampTypes[choice]
	if timestampType == "" {
		timestampType = "word"
	}
	
	fmt.Println()
	
	// Process audio
	results, err := ot.processAudio(inputFile, modelSize, timestampType)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	// Display results
	fmt.Println("\n=== TRANSCRIPTION RESULTS ===")
	fmt.Println(results)
	fmt.Println("=============================")
	
	// Save results
	fmt.Print("\nSave results to file? (y/N): ")
	scanner.Scan()
	save := strings.ToLower(strings.TrimSpace(scanner.Text()))
	
	if save == "y" || save == "yes" {
		baseName := strings.TrimSuffix(filepath.Base(inputFile), filepath.Ext(inputFile))
		outputFile := fmt.Sprintf("%s_transcription.txt", baseName)
		
		fmt.Printf("Enter output filename [%s]: ", outputFile)
		scanner.Scan()
		userFile := strings.TrimSpace(scanner.Text())
		if userFile != "" {
			outputFile = userFile
		}
		
		err = ot.saveResults(results, outputFile)
		if err != nil {
			fmt.Printf("Error saving file: %v\n", err)
		} else {
			fmt.Printf("Results saved to: %s\n", outputFile)
		}
	}
}

func printUsage() {
	fmt.Println("OfflineTranscribe - Offline Speech-to-Text Tool")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  OfflineTranscribe                                    - Interactive mode")
	fmt.Println("  OfflineTranscribe <input> [options]                  - CLI mode")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -model <size>    Model size: tiny, base, small, medium (default: base)")
	fmt.Println("  -type <type>     Timestamp type: word, sentence (default: word)")
	fmt.Println("  -output <file>   Output file (default: <input>_transcription.txt)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  OfflineTranscribe recording.wav")
	fmt.Println("  OfflineTranscribe recording.wav -model small -type sentence")
	fmt.Println("  OfflineTranscribe recording.wav -output transcript.txt")
}

func main() {
	ot := NewOfflineTranscribe()
	
	if len(os.Args) == 1 {
		// Interactive mode
		ot.interactive()
		return
	}
	
	if len(os.Args) == 2 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		printUsage()
		return
	}
	
	// CLI mode
	inputFile := os.Args[1]
	modelSize := "base"
	timestampType := "word"
	outputFile := ""
	
	// Parse command line arguments
	for i := 2; i < len(os.Args); i += 2 {
		if i+1 >= len(os.Args) {
			fmt.Printf("Error: option %s requires a value\n", os.Args[i])
			os.Exit(1)
		}
		
		switch os.Args[i] {
		case "-model":
			modelSize = os.Args[i+1]
		case "-type":
			timestampType = os.Args[i+1]
		case "-output":
			outputFile = os.Args[i+1]
		default:
			fmt.Printf("Error: unknown option %s\n", os.Args[i])
			printUsage()
			os.Exit(1)
		}
	}
	
	// Set default output file if not specified
	if outputFile == "" {
		baseName := strings.TrimSuffix(filepath.Base(inputFile), filepath.Ext(inputFile))
		outputFile = fmt.Sprintf("%s_transcription.txt", baseName)
	}
	
	// Process audio
	results, err := ot.processAudio(inputFile, modelSize, timestampType)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	
	// Save results
	err = ot.saveResults(results, outputFile)
	if err != nil {
		fmt.Printf("Error saving file: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("Transcription saved to: %s\n", outputFile)
}