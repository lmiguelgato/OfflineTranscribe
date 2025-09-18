package main

import (
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

//go:embed index.html bundle/resources/models bundle/resources/whisper
var embeddedResources embed.FS

// ResourceManager handles extraction and management of embedded resources
type ResourceManager struct {
	tempDir       string
	whisperPath   string
	modelsDir     string
	indexHTMLPath string
}

// NewResourceManager creates a new resource manager and extracts embedded files
func NewResourceManager() (*ResourceManager, error) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "OfflineTranscribe-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %v", err)
	}

	rm := &ResourceManager{
		tempDir:       tempDir,
		whisperPath:   filepath.Join(tempDir, "whisper"),
		modelsDir:     filepath.Join(tempDir, "models"),
		indexHTMLPath: filepath.Join(tempDir, "index.html"),
	}

	// Extract all resources
	if err := rm.extractResources(); err != nil {
		rm.Cleanup() // Clean up on error
		return nil, err
	}

	return rm, nil
}

// extractResources extracts all embedded files to the temporary directory
func (rm *ResourceManager) extractResources() error {
	// Create subdirectories
	if err := os.MkdirAll(rm.whisperPath, 0755); err != nil {
		return fmt.Errorf("failed to create whisper directory: %v", err)
	}
	if err := os.MkdirAll(rm.modelsDir, 0755); err != nil {
		return fmt.Errorf("failed to create models directory: %v", err)
	}

	// Extract whisper executables and DLLs
	whisperFiles, err := embeddedResources.ReadDir("bundle/resources/whisper")
	if err != nil {
		return fmt.Errorf("failed to read whisper resources: %v", err)
	}

	for _, file := range whisperFiles {
		if file.IsDir() {
			continue
		}
		
		sourcePath := "bundle/resources/whisper/" + file.Name()
		destPath := filepath.Join(rm.whisperPath, file.Name())
		
		if err := rm.extractFile(sourcePath, destPath); err != nil {
			return fmt.Errorf("failed to extract %s: %v", file.Name(), err)
		}
		
		// Make whisper executable on Unix systems
		if file.Name() == "whisper-cli.exe" || file.Name() == "whisper-cli" {
			if runtime.GOOS != "windows" {
				if err := os.Chmod(destPath, 0755); err != nil {
					return fmt.Errorf("failed to make whisper executable: %v", err)
				}
			}
		}
	}

	// Extract AI models
	modelFiles, err := embeddedResources.ReadDir("bundle/resources/models")
	if err != nil {
		return fmt.Errorf("failed to read model resources: %v", err)
	}

	for _, file := range modelFiles {
		if file.IsDir() {
			continue
		}
		
		sourcePath := "bundle/resources/models/" + file.Name()
		destPath := filepath.Join(rm.modelsDir, file.Name())
		
		if err := rm.extractFile(sourcePath, destPath); err != nil {
			return fmt.Errorf("failed to extract model %s: %v", file.Name(), err)
		}
	}

	// Extract index.html (optional - only needed for web interface)
	err = rm.extractFile("index.html", rm.indexHTMLPath)
	if err != nil {
		// Log warning but don't fail for CLI usage
		fmt.Printf("Warning: Could not extract index.html (only needed for web interface): %v\n", err)
	}

	return nil
}

// extractFile extracts a single file from embedded resources
func (rm *ResourceManager) extractFile(sourcePath, destPath string) error {
	sourceFile, err := embeddedResources.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// GetWhisperExecutable returns the path to the extracted whisper executable
func (rm *ResourceManager) GetWhisperExecutable() string {
	whisperExe := "whisper-cli.exe"
	if runtime.GOOS != "windows" {
		whisperExe = "whisper-cli"
	}
	return filepath.Join(rm.whisperPath, whisperExe)
}

// GetModelsDir returns the path to the extracted models directory
func (rm *ResourceManager) GetModelsDir() string {
	return rm.modelsDir
}

// GetIndexHTML returns the path to the extracted index.html
func (rm *ResourceManager) GetIndexHTML() string {
	return rm.indexHTMLPath
}

// GetTempDir returns the base temporary directory
func (rm *ResourceManager) GetTempDir() string {
	return rm.tempDir
}

// Cleanup removes all extracted temporary files
func (rm *ResourceManager) Cleanup() error {
	if rm.tempDir != "" {
		return os.RemoveAll(rm.tempDir)
	}
	return nil
}

// ListAvailableModels returns a list of available AI models
func (rm *ResourceManager) ListAvailableModels() ([]string, error) {
	entries, err := os.ReadDir(rm.modelsDir)
	if err != nil {
		return nil, err
	}

	var models []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".bin" {
			// Extract model name (remove ggml- prefix and .bin suffix)
			name := entry.Name()
			name = name[5:]  // Remove "ggml-"
			name = name[:len(name)-4]  // Remove ".bin"
			models = append(models, name)
		}
	}
	
	return models, nil
}

// GetModelPath returns the full path to a specific model file
func (rm *ResourceManager) GetModelPath(modelName string) string {
	return filepath.Join(rm.modelsDir, fmt.Sprintf("ggml-%s.bin", modelName))
}

// VerifyResources checks that all essential resources were extracted correctly
func (rm *ResourceManager) VerifyResources() error {
	// Check whisper executable
	whisperPath := rm.GetWhisperExecutable()
	if _, err := os.Stat(whisperPath); err != nil {
		return fmt.Errorf("whisper executable not found: %v", err)
	}

	// Check at least one model exists
	models, err := rm.ListAvailableModels()
	if err != nil {
		return fmt.Errorf("failed to list models: %v", err)
	}
	if len(models) == 0 {
		return fmt.Errorf("no AI models found")
	}

	// Check index.html
	if _, err := os.Stat(rm.indexHTMLPath); err != nil {
		return fmt.Errorf("index.html not found: %v", err)
	}

	return nil
}