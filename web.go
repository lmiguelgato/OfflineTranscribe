package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type WebServer struct {
	port string
}

type TranscriptionRequest struct {
	ModelSize     string `json:"modelSize"`
	TimestampType string `json:"timestampType"`
}

type TranscriptionResponse struct {
	Success bool   `json:"success"`
	Results string `json:"results"`
	Error   string `json:"error,omitempty"`
}

const indexHTML = "index.html"

func NewWebServer(port string) *WebServer {
	return &WebServer{port: port}
}

func (ws *WebServer) handleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func (ws *WebServer) handleTranscribe(w http.ResponseWriter, r *http.Request) {
	// Add CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	
	// Handle preflight request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(100 << 20) // 100 MB max
	if err != nil {
		ws.sendJSONResponse(w, TranscriptionResponse{
			Success: false,
			Error:   "Failed to parse form data",
		})
		return
	}

	// Get file
	file, header, err := r.FormFile("audioFile")
	if err != nil {
		ws.sendJSONResponse(w, TranscriptionResponse{
			Success: false,
			Error:   "No audio file provided",
		})
		return
	}
	defer file.Close()

	// Get options
	modelSize := r.FormValue("modelSize")
	timestampType := r.FormValue("timestampType")

	if modelSize == "" {
		modelSize = "base"
	}
	if timestampType == "" {
		timestampType = "word"
	}

	// Save uploaded file temporarily
	tempDir := os.TempDir()
	tempFile := filepath.Join(tempDir, header.Filename)
	
	outFile, err := os.Create(tempFile)
	if err != nil {
		ws.sendJSONResponse(w, TranscriptionResponse{
			Success: false,
			Error:   "Failed to save uploaded file",
		})
		return
	}
	defer func() {
		outFile.Close()
		os.Remove(tempFile) // Clean up
	}()

	_, err = io.Copy(outFile, file)
	if err != nil {
		ws.sendJSONResponse(w, TranscriptionResponse{
			Success: false,
			Error:   "Failed to save uploaded file",
		})
		return
	}
	outFile.Close()

	// Process the audio file
	results, err := ws.processAudio(tempFile, modelSize, timestampType)
	if err != nil {
		ws.sendJSONResponse(w, TranscriptionResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	ws.sendJSONResponse(w, TranscriptionResponse{
		Success: true,
		Results: results,
	})
}

func (ws *WebServer) processAudio(inputFile, modelSize, timestampType string) (string, error) {
	log.Printf("Processing audio file: %s with model: %s, timestamps: %s", inputFile, modelSize, timestampType)
	
	// Simulate processing time based on model size
	var processingTime time.Duration
	switch modelSize {
	case "tiny":
		processingTime = 2 * time.Second
	case "base":
		processingTime = 3 * time.Second
	case "small":
		processingTime = 5 * time.Second
	case "medium":
		processingTime = 8 * time.Second
	default:
		processingTime = 3 * time.Second
	}
	
	time.Sleep(processingTime)
	
	// Generate sample results (in real implementation, this would be actual Whisper processing)
	return ws.generateSampleResults(timestampType), nil
}

func (ws *WebServer) generateSampleResults(timestampType string) string {
	if timestampType == "word" {
		return `[00:00:01.240] Hello
[00:00:01.480] there,
[00:00:01.720] this
[00:00:01.960] is
[00:00:02.120] a
[00:00:02.280] sample
[00:00:02.560] transcription
[00:00:03.000] with
[00:00:03.200] word-level
[00:00:03.680] timestamps.
[00:00:04.120] Each
[00:00:04.280] word
[00:00:04.440] has
[00:00:04.600] its
[00:00:04.760] own
[00:00:04.920] precise
[00:00:05.240] timestamp
[00:00:05.600] for
[00:00:05.760] easy
[00:00:05.960] navigation.`
	} else {
		return `[00:00:01.240 - 00:00:03.680] Hello there, this is a sample transcription with word-level timestamps.

[00:00:04.120 - 00:00:06.200] Each word has its own precise timestamp for easy navigation.

[00:00:06.500 - 00:00:09.800] This sentence-level format groups words together for better readability.

[00:00:10.100 - 00:00:13.500] You can quickly find specific sections using the time references provided.`
	}
}

func (ws *WebServer) sendJSONResponse(w http.ResponseWriter, response TranscriptionResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
}

func (ws *WebServer) Start() {
	http.HandleFunc("/", ws.handleIndex)
	http.HandleFunc("/transcribe", ws.handleTranscribe)
	
	fmt.Printf("OfflineTranscribe Web Interface starting on http://localhost:%s\n", ws.port)
	fmt.Println("Open your web browser and navigate to the URL above")
	
	log.Fatal(http.ListenAndServe(":"+ws.port, nil))
}

func main() {
	port := "8080"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	
	server := NewWebServer(port)
	server.Start()
}