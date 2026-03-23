package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestMaxInt(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"a greater", 5, 3, 5},
		{"b greater", 2, 7, 7},
		{"equal", 4, 4, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maxInt(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("maxInt(%d, %d) = %d, want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLen   int
		expected string
	}{
		{"shorter", "hello", 10, "hello"},
		{"exact", "hello", 5, "hello"},
		{"longer", "hello world", 5, "hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncateString(tt.input, tt.maxLen)
			if result != tt.expected {
				t.Errorf("truncateString(%q, %d) = %q, want %q", tt.input, tt.maxLen, result, tt.expected)
			}
		})
	}
}

func TestBinarySearchFilePosition(t *testing.T) {
	// Create a temporary test file
	tempFile, err := os.CreateTemp("", "test*.txt")
	if err != nil {
		t.Fatal("Failed to create temp file:", err)
	}
	defer os.Remove(tempFile.Name())

	// Write test data
	testData := "line001\nline002\nline003\nline004\nline005\nline006\nline007\nline008\nline009\nline010\n"
	if _, err := tempFile.Write([]byte(testData)); err != nil {
		t.Fatal("Failed to write to temp file:", err)
	}
	tempFile.Close()

	// Open file for reading
	file, err := os.Open(tempFile.Name())
	if err != nil {
		t.Fatal("Failed to open temp file:", err)
	}
	defer file.Close()

	result := binarySearchFilePosition(file, 58, "line004", "line008")
	if result <= 0 {
		t.Errorf("binarySearchFilePosition() = %d, want positive position", result)
	}
}

func TestScanLinesInRange(t *testing.T) {
	// Create a temporary test file
	tempFile, err := os.CreateTemp("", "test*.txt")
	if err != nil {
		t.Fatal("Failed to create temp file:", err)
	}
	defer os.Remove(tempFile.Name())

	// Write test data
	testData := "line001\nline002\nline003\nline004\nline005\nline006\nline007\nline008\nline009\nline010\n"
	if _, err := tempFile.Write([]byte(testData)); err != nil {
		t.Fatal("Failed to write to temp file:", err)
	}
	tempFile.Close()

	// Open file for reading
	file, err := os.Open(tempFile.Name())
	if err != nil {
		t.Fatal("Failed to open temp file:", err)
	}
	defer file.Close()

	// Capture output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	scanLinesInRange(file, 0, "line003", "line007", "line005")

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "line005") {
		t.Errorf("Expected to find 'line005' in output, got: %q", output)
	}
}

func TestSeekToLineStart(t *testing.T) {
	// Create a temporary test file
	tempFile, err := os.CreateTemp("", "test*.txt")
	if err != nil {
		t.Fatal("Failed to create temp file:", err)
	}
	defer os.Remove(tempFile.Name())

	// Write test data
	testData := "line001\nline002\nline003\nline004\nline005\n"
	if _, err := tempFile.Write([]byte(testData)); err != nil {
		t.Fatal("Failed to write to temp file:", err)
	}
	tempFile.Close()

	// Open file for reading
	file, err := os.Open(tempFile.Name())
	if err != nil {
		t.Fatal("Failed to open temp file:", err)
	}
	defer file.Close()

	// Test seeking to middle of file
	seekToLineStart(file, 10)
	
	// Read current position
	pos, _ := file.Seek(0, io.SeekCurrent)
	if pos != 8 { // Should be at start of "line002" content (after "line001\n")
		t.Errorf("seekToLineStart() positioned at %d, expected 8", pos)
	}
}

func TestReadLinePrefix(t *testing.T) {
	// Create a temporary test file
	tempFile, err := os.CreateTemp("", "test*.txt")
	if err != nil {
		t.Fatal("Failed to create temp file:", err)
	}
	defer os.Remove(tempFile.Name())

	// Write test data
	testData := "line001\nline002\nline003\n"
	if _, err := tempFile.Write([]byte(testData)); err != nil {
		t.Fatal("Failed to write to temp file:", err)
	}
	tempFile.Close()

	// Open file for reading
	file, err := os.Open(tempFile.Name())
	if err != nil {
		t.Fatal("Failed to open temp file:", err)
	}
	defer file.Close()

	// Seek to middle of first line
	file.Seek(3, io.SeekStart)
	
	result := readLinePrefix(file, 4)
	if result != "e001" {
		t.Errorf("readLinePrefix() = %q, want %q", result, "e001")
	}
}

func TestBinarySearchDateBoundary(t *testing.T) {
	// Create a temporary test file with date boundary crossing
	tempFile, err := os.CreateTemp("", "test_date_boundary*.txt")
	if err != nil {
		t.Fatal("Failed to create temp file:", err)
	}
	defer os.Remove(tempFile.Name())

	// Write test data that crosses a date boundary (like Mar 23 -> Mar 24)
	testData := "Mar 23 23:59:50 line1\nMar 23 23:59:55 line2\nMar 24 00:00:05 line3\nMar 24 00:00:10 line4\nMar 24 00:00:15 line5\n"
	if _, err := tempFile.Write([]byte(testData)); err != nil {
		t.Fatal("Failed to write to temp file:", err)
	}
	tempFile.Close()

	// Open file for reading
	file, err := os.Open(tempFile.Name())
	if err != nil {
		t.Fatal("Failed to open temp file:", err)
	}
	defer file.Close()

	// Test binary search across date boundary
	result := binarySearchFilePosition(file, int64(len(testData)), "Mar 23 23:59:00", "Mar 24 00:00:20")
	if result <= 0 {
		t.Errorf("binarySearchFilePosition() = %d, want positive position for date boundary search", result)
	}

	// Test that we can actually find content in the range
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	scanLinesInRange(file, result, "Mar 23 23:59:00", "Mar 24 00:00:20", "line3")

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "line3") {
		t.Errorf("Expected to find 'line3' in date boundary search output, got: %q", output)
	}
}