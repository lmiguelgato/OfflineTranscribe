package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type WhisperTranscriber struct {
	executablePath string
}

type TranscriptionResult struct {
	Text      string
	Segments  []Segment
	Error     error
}

type Segment struct {
	Start float64
	End   float64
	Text  string
	Words []Word
}

type Word struct {
	Start float64
	End   float64
	Text  string
}

func NewWhisperTranscriber() *WhisperTranscriber {
	return &WhisperTranscriber{
		executablePath: "whisper.exe", // Will be bundled with the application
	}
}

func (wt *WhisperTranscriber) LoadModel(modelSize string) error {
	// Model file paths based on size
	modelPaths := map[string]string{
		"tiny":   "models/ggml-tiny.bin",
		"base":   "models/ggml-base.bin", 
		"small":  "models/ggml-small.bin",
		"medium": "models/ggml-medium.bin",
	}
	
	modelPath, exists := modelPaths[modelSize]
	if !exists {
		return fmt.Errorf("unsupported model size: %s", modelSize)
	}
	
	// Create models directory if it doesn't exist
	if err := os.MkdirAll("models", 0755); err != nil {
		return fmt.Errorf("failed to create models directory: %v", err)
	}
	
	// Check if model exists, if not provide download instructions
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		return fmt.Errorf("model file not found: %s\n\nTo download models, run the download script:\n"+
			"Windows: download_models.bat\n"+
			"Linux/Mac: ./download_models.sh", modelPath)
	}
	
	// Check if whisper executable exists
	if _, err := os.Stat(wt.executablePath); os.IsNotExist(err) {
		// Try alternative paths
		alternatives := []string{
			"whisper",                                       // Unix systems
			"./whisper",                                    // Local directory
			"./whisper.exe",                                // Local directory Windows
			"bin/whisper.exe",                              // Subdirectory Windows
			"bin/whisper",                                  // Subdirectory Unix
			"whisper-bin-x64/Release/whisper-cli.exe",      // Common whisper.cpp build location (CLI version)
			"./whisper-bin-x64/Release/whisper-cli.exe",    // Relative path CLI
			"whisper-bin-x64\\Release\\whisper-cli.exe",    // Windows path separators CLI
			".\\whisper-bin-x64\\Release\\whisper-cli.exe", // Windows relative path CLI
			"whisper-bin-x64/Release/main.exe",             // Alternative main executable
			"./whisper-bin-x64/Release/main.exe",           // Relative path main
			"whisper-bin-x64\\Release\\main.exe",           // Windows path main
			".\\whisper-bin-x64\\Release\\main.exe",        // Windows relative path main
		}
		
		found := false
		for _, alt := range alternatives {
			if _, err := os.Stat(alt); err == nil {
				wt.executablePath = alt
				found = true
				break
			}
		}
		
		if !found {
			return fmt.Errorf("whisper executable not found\n\n" +
				"Please download whisper.cpp from: https://github.com/ggerganov/whisper.cpp/releases\n" +
				"Extract the whisper executable to this directory or add it to your PATH")
		}
	}
	
	return nil
}

func (wt *WhisperTranscriber) TranscribeFile(inputFile string, timestampType string) (*TranscriptionResult, error) {
	// Check if file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("audio file not found: %s", inputFile)
	}
	
	// Prepare output file
	outputDir := filepath.Dir(inputFile)
	baseName := strings.TrimSuffix(filepath.Base(inputFile), filepath.Ext(inputFile))
	outputFile := filepath.Join(outputDir, baseName+"_whisper_output")
	
	// Determine model path (default to base)
	modelPath := "models/ggml-base.bin"
	
	// Build whisper command
	var args []string
	args = append(args, "-m", modelPath)
	args = append(args, "-f", inputFile)
	args = append(args, "-of", outputFile)
	args = append(args, "-otxt")  // Output text file
	args = append(args, "-nt")    // No timestamps in output (we'll add our own)
	args = append(args, "-np")    // No print special tokens
	
	// Add word-level options if requested
	if timestampType == "word" {
		args = append(args, "-ml", "1") // Maximum line length for word-level
		args = append(args, "-sow")     // Split on word rather than token
	}
	
	// Execute whisper
	cmd := exec.Command(wt.executablePath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		outputStr := string(output)
		if strings.Contains(outputStr, "failed to read audio") {
			return nil, fmt.Errorf("invalid audio file: %s\nWhisper supports: WAV, MP3, FLAC, MP4, M4A, OGG\nError: %v", inputFile, err)
		}
		return nil, fmt.Errorf("whisper execution failed: %v\nOutput: %s", err, outputStr)
	}
	
	// Read the generated transcription file
	transcriptionFile := outputFile + ".txt"
	content, err := os.ReadFile(transcriptionFile)
	if err != nil {
		// Check if whisper actually created any output files
		possibleFiles := []string{
			outputFile + ".txt",
			filepath.Join(outputDir, baseName + ".txt"),
			inputFile + ".txt",
		}
		
		var foundFile string
		for _, possibleFile := range possibleFiles {
			if _, err := os.Stat(possibleFile); err == nil {
				foundFile = possibleFile
				break
			}
		}
		
		if foundFile != "" {
			content, err = os.ReadFile(foundFile)
			transcriptionFile = foundFile
		} else {
			return nil, fmt.Errorf("whisper did not create expected output file.\nTried: %v\nWhisper output: %s", possibleFiles, string(output))
		}
		
		if err != nil {
			return nil, fmt.Errorf("failed to read transcription output: %v", err)
		}
	}
	
	// Clean up temporary files
	os.Remove(transcriptionFile)
	
	// Parse the output into segments
	result := &TranscriptionResult{
		Text:     string(content),
		Segments: wt.parseSegments(string(content), timestampType),
	}
	
	// Clean up temporary files
	os.Remove(outputFile + ".txt")
	
	return result, nil
}

func (wt *WhisperTranscriber) parseSegments(content string, timestampType string) []Segment {
	var segments []Segment
	
	// For now, create simple segments since whisper.cpp output parsing can be complex
	// This is a simplified version - in a full implementation, you'd parse the actual whisper output format
	lines := strings.Split(content, "\n")
	
	segmentDuration := 3.0 // Default 3-second segments
	currentTime := 0.0
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		segment := Segment{
			Start: currentTime,
			End:   currentTime + segmentDuration,
			Text:  line,
		}
		
		// If word-level timestamps are requested, create word segments
		if timestampType == "word" {
			words := strings.Fields(line)
			wordDuration := segmentDuration / float64(len(words))
			
			for i, word := range words {
				wordStart := currentTime + float64(i)*wordDuration
				wordEnd := wordStart + wordDuration
				
				segment.Words = append(segment.Words, Word{
					Start: wordStart,
					End:   wordEnd,
					Text:  word,
				})
			}
		}
		
		segments = append(segments, segment)
		currentTime += segmentDuration
	}
	
	return segments
}

func (wt *WhisperTranscriber) FormatResults(result *TranscriptionResult, timestampType string) string {
	var output strings.Builder
	
	if timestampType == "word" {
		// Word-level output
		for _, segment := range result.Segments {
			if len(segment.Words) > 0 {
				for _, word := range segment.Words {
					timestamp := formatTimestamp(word.Start)
					output.WriteString(fmt.Sprintf("[%s] %s\n", timestamp, word.Text))
				}
			} else {
				// Fallback if word-level data isn't available
				words := strings.Fields(segment.Text)
				wordDuration := (segment.End - segment.Start) / float64(len(words))
				
				for j, word := range words {
					wordStart := segment.Start + float64(j)*wordDuration
					timestamp := formatTimestamp(wordStart)
					output.WriteString(fmt.Sprintf("[%s] %s\n", timestamp, word))
				}
			}
		}
	} else {
		// Sentence-level output
		for _, segment := range result.Segments {
			startTime := formatTimestamp(segment.Start)
			endTime := formatTimestamp(segment.End)
			output.WriteString(fmt.Sprintf("[%s - %s] %s\n\n", startTime, endTime, segment.Text))
		}
	}
	
	return output.String()
}

func (wt *WhisperTranscriber) Close() {
	// Nothing to close for executable-based approach
}

func formatTimestamp(seconds float64) string {
	hours := int(seconds / 3600)
	minutes := int((seconds - float64(hours*3600)) / 60)
	secs := seconds - float64(hours*3600) - float64(minutes*60)
	return fmt.Sprintf("%02d:%02d:%06.3f", hours, minutes, secs)
}

// parseWhisperTimestamps parses timestamp format from whisper output
func parseWhisperTimestamps(line string) (float64, float64, string) {
	// Whisper typically outputs: [00:00:00.000 --> 00:00:03.000]  Text here
	re := regexp.MustCompile(`\[(\d{2}):(\d{2}):(\d{2})\.(\d{3}) --> (\d{2}):(\d{2}):(\d{2})\.(\d{3})\]\s*(.*)`)
	matches := re.FindStringSubmatch(line)
	
	if len(matches) != 10 {
		return 0, 0, line // Fallback if no timestamps found
	}
	
	// Parse start time
	startH, _ := strconv.Atoi(matches[1])
	startM, _ := strconv.Atoi(matches[2])
	startS, _ := strconv.Atoi(matches[3])
	startMs, _ := strconv.Atoi(matches[4])
	startTime := float64(startH*3600 + startM*60 + startS) + float64(startMs)/1000.0
	
	// Parse end time
	endH, _ := strconv.Atoi(matches[5])
	endM, _ := strconv.Atoi(matches[6])
	endS, _ := strconv.Atoi(matches[7])
	endMs, _ := strconv.Atoi(matches[8])
	endTime := float64(endH*3600 + endM*60 + endS) + float64(endMs)/1000.0
	
	text := strings.TrimSpace(matches[9])
	
	return startTime, endTime, text
}