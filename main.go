package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type LocalTTS struct {
	app            fyne.App
	window         fyne.Window
	fileEntry      *widget.Entry
	modelSelect    *widget.Select
	timestampSelect *widget.Select
	processBtn     *widget.Button
	saveBtn        *widget.Button
	resultsText    *widget.Entry
	progressBar    *widget.ProgressBar
	statusLabel    *widget.Label
	lastResults    string
	currentFile    string
}

func NewLocalTTS() *LocalTTS {
	myApp := app.New()
	myApp.SetIcon(nil) // You can add an icon resource here
	
	myWindow := myApp.NewWindow("LocalTTS - Offline Speech to Text")
	myWindow.Resize(fyne.NewSize(800, 600))
	
	return &LocalTTS{
		app:    myApp,
		window: myWindow,
	}
}

func (lt *LocalTTS) setupUI() {
	// Title
	title := widget.NewCard("LocalTTS", "Offline Speech-to-Text Transcription", nil)
	
	// File selection
	lt.fileEntry = widget.NewEntry()
	lt.fileEntry.SetPlaceHolder("Select an audio file...")
	lt.fileEntry.Disable()
	
	browseBtn := widget.NewButton("Browse", lt.browseFile)
	fileContainer := container.NewBorder(nil, nil, nil, browseBtn, lt.fileEntry)
	
	// Model selection
	lt.modelSelect = widget.NewSelect([]string{"tiny", "base", "small", "medium"}, nil)
	lt.modelSelect.SetSelected("base")
	modelContainer := container.NewBorder(nil, nil, widget.NewLabel("Model Size:"), nil, lt.modelSelect)
	
	// Timestamp granularity
	lt.timestampSelect = widget.NewSelect([]string{"word", "sentence"}, nil)
	lt.timestampSelect.SetSelected("word")
	timestampContainer := container.NewBorder(nil, nil, widget.NewLabel("Timestamps:"), nil, lt.timestampSelect)
	
	// Process button
	lt.processBtn = widget.NewButton("Process Audio", lt.processAudio)
	lt.processBtn.Importance = widget.HighImportance
	
	// Progress and status
	lt.progressBar = widget.NewProgressBar()
	lt.progressBar.Hide()
	lt.statusLabel = widget.NewLabel("Ready")
	
	// Results
	lt.resultsText = widget.NewMultiLineEntry()
	lt.resultsText.SetPlaceHolder("Transcription results will appear here...")
	lt.resultsText.Resize(fyne.NewSize(780, 300))
	
	// Save button
	lt.saveBtn = widget.NewButton("Save Results", lt.saveResults)
	lt.saveBtn.Disable()
	
	// Layout
	content := container.NewVBox(
		title,
		widget.NewSeparator(),
		widget.NewCard("", "Audio File", fileContainer),
		container.NewGridWithColumns(2, modelContainer, timestampContainer),
		lt.processBtn,
		widget.NewSeparator(),
		lt.statusLabel,
		lt.progressBar,
		widget.NewCard("", "Transcription Results", lt.resultsText),
		lt.saveBtn,
	)
	
	lt.window.SetContent(container.NewScroll(content))
}

func (lt *LocalTTS) browseFile() {
	dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, lt.window)
			return
		}
		if reader == nil {
			return
		}
		defer reader.Close()
		
		lt.currentFile = reader.URI().Path()
		lt.fileEntry.SetText(lt.currentFile)
		
		// Enable process button if file is selected
		if lt.currentFile != "" {
			lt.processBtn.Enable()
		}
	}, lt.window)
}

func (lt *LocalTTS) loadModel(modelSize string) error {
	lt.statusLabel.SetText(fmt.Sprintf("Loading %s model...", modelSize))
	lt.progressBar.Show()
	
	// In a real implementation, you would download/load the model here
	// For now, we'll simulate the model loading
	modelPath := fmt.Sprintf("models/ggml-%s.bin", modelSize)
	
	// Check if model exists, if not, we would download it
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		// Create models directory
		os.MkdirAll("models", 0755)
		
		// In a real implementation, download the model here
		lt.statusLabel.SetText(fmt.Sprintf("Downloading %s model... (this may take a while)", modelSize))
		
		// Simulate download time
		time.Sleep(2 * time.Second)
		
		// Create a placeholder file (in real implementation, this would be the actual model)
		file, err := os.Create(modelPath)
		if err != nil {
			return fmt.Errorf("failed to create model file: %v", err)
		}
		file.Close()
	}
	
	// Load the model (placeholder - in real implementation use whisper.New)
	var err error
	// lt.model, err = whisper.New(modelPath)
	if err != nil {
		return fmt.Errorf("failed to load model: %v", err)
	}
	
	lt.progressBar.Hide()
	lt.statusLabel.SetText("Model loaded successfully")
	return nil
}

func (lt *LocalTTS) processAudio() {
	if lt.currentFile == "" {
		dialog.ShowError(fmt.Errorf("please select an audio file"), lt.window)
		return
	}
	
	modelSize := lt.modelSelect.Selected
	if modelSize == "" {
		dialog.ShowError(fmt.Errorf("please select a model size"), lt.window)
		return
	}
	
	// Disable UI during processing
	lt.processBtn.Disable()
	lt.saveBtn.Disable()
	
	go func() {
		// Load model if needed
		if err := lt.loadModel(modelSize); err != nil {
			lt.statusLabel.SetText("Error loading model")
			dialog.ShowError(err, lt.window)
			lt.processBtn.Enable()
			return
		}
		
		// Process audio
		lt.statusLabel.SetText("Transcribing audio...")
		lt.progressBar.Show()
		
		// In a real implementation, you would use:
		// context := lt.model.NewContext()
		// err := context.Process(lt.currentFile, nil, nil)
		
		// For demonstration, simulate processing
		time.Sleep(3 * time.Second)
		
		// Generate sample results
		timestampType := lt.timestampSelect.Selected
		lt.generateSampleResults(timestampType)
		
		lt.progressBar.Hide()
		lt.statusLabel.SetText("Transcription complete")
		lt.processBtn.Enable()
		lt.saveBtn.Enable()
	}()
}

func (lt *LocalTTS) generateSampleResults(timestampType string) {
	// This is placeholder content - in a real implementation, this would come from Whisper
	var results string
	
	if timestampType == "word" {
		results = `[00:00:01.240] Hello
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
		results = `[00:00:01.240 - 00:00:03.680] Hello there, this is a sample transcription with word-level timestamps.

[00:00:04.120 - 00:00:06.200] Each word has its own precise timestamp for easy navigation.

[00:00:06.500 - 00:00:09.800] This sentence-level format groups words together for better readability.

[00:00:10.100 - 00:00:13.500] You can quickly find specific sections using the time references provided.`
	}
	
	lt.lastResults = results
	lt.resultsText.SetText(results)
}

func (lt *LocalTTS) saveResults() {
	if lt.lastResults == "" {
		dialog.ShowError(fmt.Errorf("no results to save"), lt.window)
		return
	}
	
	// Generate default filename
	baseName := "transcription"
	if lt.currentFile != "" {
		baseName = strings.TrimSuffix(filepath.Base(lt.currentFile), filepath.Ext(lt.currentFile)) + "_transcription"
	}
	
	dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, lt.window)
			return
		}
		if writer == nil {
			return
		}
		defer writer.Close()
		
		_, err = writer.Write([]byte(lt.lastResults))
		if err != nil {
			dialog.ShowError(err, lt.window)
			return
		}
		
		dialog.ShowInformation("Success", fmt.Sprintf("Results saved to %s", writer.URI().Path()), lt.window)
	}, lt.window)
}

func (lt *LocalTTS) Run() {
	lt.setupUI()
	lt.window.ShowAndRun()
}

func main() {
	app := NewLocalTTS()
	app.Run()
}