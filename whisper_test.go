package main

import (
	"testing"
)

// Test SearchWord functionality
func TestSearchWord(t *testing.T) {
	wt := &WhisperTranscriber{}
	
	// Create test transcription result
	result := &TranscriptionResult{
		Segments: []Segment{
			{Start: 1.0, End: 3.0, Text: "Hello there, how are you?"},
			{Start: 3.5, End: 5.0, Text: "I am doing great today."},
			{Start: 5.5, End: 7.0, Text: "Hello world!"},
			{Start: 7.5, End: 9.0, Text: "This is a test."},
		},
	}
	
	// Test 1: Search for a word that appears multiple times
	matches := wt.SearchWord(result, "hello")
	if len(matches) != 2 {
		t.Errorf("Expected 2 matches for 'hello', got %d", len(matches))
	}
	if matches[0].Text != "Hello there, how are you?" {
		t.Errorf("Expected first match to be 'Hello there, how are you?', got '%s'", matches[0].Text)
	}
	if matches[1].Text != "Hello world!" {
		t.Errorf("Expected second match to be 'Hello world!', got '%s'", matches[1].Text)
	}
	
	// Test 2: Search for a word that appears once
	matches = wt.SearchWord(result, "test")
	if len(matches) != 1 {
		t.Errorf("Expected 1 match for 'test', got %d", len(matches))
	}
	
	// Test 3: Search for a word that doesn't exist
	matches = wt.SearchWord(result, "nonexistent")
	if len(matches) != 0 {
		t.Errorf("Expected 0 matches for 'nonexistent', got %d", len(matches))
	}
	
	// Test 4: Test case insensitivity
	matches = wt.SearchWord(result, "HELLO")
	if len(matches) != 2 {
		t.Errorf("Expected 2 matches for 'HELLO' (case insensitive), got %d", len(matches))
	}
	
	// Test 5: Test punctuation handling
	matches = wt.SearchWord(result, "world")
	if len(matches) != 1 {
		t.Errorf("Expected 1 match for 'world' (with punctuation), got %d", len(matches))
	}
}

// Test FormatSearchResults functionality
func TestFormatSearchResults(t *testing.T) {
	wt := &WhisperTranscriber{}
	
	// Test 1: No matches
	matches := []Segment{}
	result := wt.FormatSearchResults(matches, "test")
	if result != "No matches found for 'test'\n" {
		t.Errorf("Expected 'No matches found' message, got: %s", result)
	}
	
	// Test 2: One match
	matches = []Segment{
		{Start: 1.0, End: 3.0, Text: "This is a test."},
	}
	result = wt.FormatSearchResults(matches, "test")
	if result == "" {
		t.Error("Expected non-empty result for one match")
	}
	
	// Test 3: Multiple matches
	matches = []Segment{
		{Start: 1.0, End: 3.0, Text: "First test."},
		{Start: 5.0, End: 7.0, Text: "Second test."},
	}
	result = wt.FormatSearchResults(matches, "test")
	if result == "" {
		t.Error("Expected non-empty result for multiple matches")
	}
}
