package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type WebServer struct {
	port            string
	resourceManager *ResourceManager
	transcriber     *WhisperTranscriber
}

type TranscriptionRequest struct {
	ModelSize string `json:"modelSize"`
}

type TranscriptionResponse struct {
	Success bool   `json:"success"`
	Results string `json:"results"`
	Error   string `json:"error,omitempty"`
}

func NewWebServer(port string, resourceManager *ResourceManager) *WebServer {
	transcriber := NewWhisperTranscriber(resourceManager)
	return &WebServer{
		port:            port,
		resourceManager: resourceManager,
		transcriber:     transcriber,
	}
}

func (ws *WebServer) handleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, ws.resourceManager.GetIndexHTML())
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

	if modelSize == "" {
		modelSize = "base"
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
	results, err := ws.processAudio(tempFile, modelSize)
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

func (ws *WebServer) processAudio(inputFile, modelSize string) (string, error) {
	log.Printf("Processing audio file: %s with model: %s", inputFile, modelSize)

	// Load the model
	if err := ws.transcriber.LoadModel(modelSize); err != nil {
		return "", fmt.Errorf("failed to load model: %v", err)
	}

	// Transcribe the audio
	result, err := ws.transcriber.TranscribeFile(inputFile, modelSize)
	if err != nil {
		return "", fmt.Errorf("transcription failed: %v", err)
	}

	// Format the results
	formattedOutput := ws.transcriber.FormatResults(result)

	return formattedOutput, nil
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

	// Initialize resource manager
	resourceManager, err := NewResourceManager()
	if err != nil {
		log.Fatalf("Failed to initialize resources: %v", err)
	}
	defer resourceManager.Cleanup()

	// Verify resources
	if err := resourceManager.VerifyResources(); err != nil {
		log.Fatalf("Resource verification failed: %v", err)
	}

	server := NewWebServer(port, resourceManager)
	server.Start()
}
