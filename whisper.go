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
	resourceManager *ResourceManager
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

func NewWhisperTranscriber(resourceManager *ResourceManager) *WhisperTranscriber {
	return &WhisperTranscriber{
		executablePath: resourceManager.GetWhisperExecutable(),
		resourceManager: resourceManager,
	}
}

func (wt *WhisperTranscriber) LoadModel(modelSize string) error {
	// Use embedded models from resource manager
	modelPath := wt.resourceManager.GetModelPath(modelSize)
	
	// Check if model exists
	if _, err := os.Stat(modelPath); err != nil {
		// List available models for error message
		availableModels, _ := wt.resourceManager.ListAvailableModels()
		return fmt.Errorf("model '%s' not found. Available models: %v", modelSize, availableModels)
	}
	
	return nil
}

func (wt *WhisperTranscriber) TranscribeFile(inputFile string, modelSize string) (*TranscriptionResult, error) {
	// Check if file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("audio file not found: %s", inputFile)
	}
	
	// Prepare output file
	outputDir := filepath.Dir(inputFile)
	baseName := strings.TrimSuffix(filepath.Base(inputFile), filepath.Ext(inputFile))
	outputFile := filepath.Join(outputDir, baseName+"_whisper_output")
	
	// Get model path from resource manager
	modelPath := wt.resourceManager.GetModelPath(modelSize)
	
	// Build whisper command
	var args []string
	args = append(args, "-m", modelPath)
	args = append(args, "-f", inputFile)
	args = append(args, "-of", outputFile)
	args = append(args, "-osrt")  // Always use SRT format for sentence-level timestamps
	args = append(args, "-np")    // No print special tokens
	
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
	
	// Read the generated transcription file (always SRT format)
	transcriptionFile := outputFile + ".srt"
	content, readErr := os.ReadFile(transcriptionFile)
	if readErr != nil {
		// Fallback to other possible files
		possibleFiles := []string{
			outputFile + ".srt",
			outputFile + ".txt",
			filepath.Join(outputDir, baseName + ".srt"),
		}
		
		for _, possibleFile := range possibleFiles {
			if _, statErr := os.Stat(possibleFile); statErr == nil {
				content, readErr = os.ReadFile(possibleFile)
				transcriptionFile = possibleFile
				break
			}
		}
	}
	
	if readErr != nil {
		return nil, fmt.Errorf("whisper did not create expected output file.\nWhisper output: %s", string(output))
	}
	
	// Clean up temporary files
	os.Remove(transcriptionFile)
	
	// Parse SRT format for timestamps
	segments := wt.parseSRTFormat(string(content))
	
	return &TranscriptionResult{
		Text:     string(content),
		Segments: segments,
	}, nil
}

// parseSRTFormat parses Whisper's SRT subtitle format
func (wt *WhisperTranscriber) parseSRTFormat(content string) []Segment {
	var segments []Segment
	lines := strings.Split(content, "\n")
	
	i := 0
	for i < len(lines) {
		// Skip segment number
		for i < len(lines) && strings.TrimSpace(lines[i]) != "" && !strings.Contains(lines[i], "-->") {
			i++
		}
		
		if i >= len(lines) {
			break
		}
		
		// Parse timestamp line: 00:00:01,000 --> 00:00:03,500
		timestampLine := strings.TrimSpace(lines[i])
		if strings.Contains(timestampLine, "-->") {
			parts := strings.Split(timestampLine, "-->")
			if len(parts) == 2 {
				startTime := wt.parseSRTTimestamp(strings.TrimSpace(parts[0]))
				endTime := wt.parseSRTTimestamp(strings.TrimSpace(parts[1]))
				
				i++
				
				// Collect text lines until empty line
				var textLines []string
				for i < len(lines) && strings.TrimSpace(lines[i]) != "" {
					textLines = append(textLines, strings.TrimSpace(lines[i]))
					i++
				}
				
				if len(textLines) > 0 {
					text := strings.Join(textLines, " ")
					segment := Segment{
						Start: startTime,
						End:   endTime,
						Text:  text,
					}
					segments = append(segments, segment)
				}
			}
		}
		i++
	}
	
	return segments
}

// parseSRTTimestamp converts SRT timestamp format to seconds
func (wt *WhisperTranscriber) parseSRTTimestamp(timestamp string) float64 {
	// Format: 00:00:01,000
	timestamp = strings.ReplaceAll(timestamp, ",", ".")
	parts := strings.Split(timestamp, ":")
	
	if len(parts) != 3 {
		return 0
	}
	
	hours, _ := strconv.Atoi(parts[0])
	minutes, _ := strconv.Atoi(parts[1])
	seconds, _ := strconv.ParseFloat(parts[2], 64)
	
	return float64(hours*3600 + minutes*60) + seconds
}

func (wt *WhisperTranscriber) FormatResults(result *TranscriptionResult) string {
	var output strings.Builder
	
	// Sentence-level output using whisper's native segment timestamps
	for _, segment := range result.Segments {
		startTime := formatTimestamp(segment.Start)
		endTime := formatTimestamp(segment.End)
		output.WriteString(fmt.Sprintf("[%s - %s] %s\n\n", startTime, endTime, segment.Text))
	}
	
	return output.String()
}

func (wt *WhisperTranscriber) Close() {
	// Nothing to close for executable-based approach
}

func formatTimestamp(seconds float64) string {
	hours := int(seconds / 3600)
	minutes := int((seconds - float64(hours*3600)) / 60)
	secs := int(seconds - float64(hours*3600) - float64(minutes*60))
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secs)
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