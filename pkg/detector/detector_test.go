package detector

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewDetector(t *testing.T) {
	detector := NewDetector()
	if detector == nil {
		t.Fatal("NewDetector() returned nil")
	}
	if detector.config == nil {
		t.Error("Detector config should not be nil")
	}
}

func TestDetectLanguageFallback(t *testing.T) {
	detector := NewDetector()

	tests := []struct {
		name        string
		files       map[string]string
		expected    string
		shouldError bool
	}{
		{
			name: "Python with requirements.txt",
			files: map[string]string{
				"requirements.txt": "flask==2.0.0",
			},
			expected:    "python",
			shouldError: false,
		},
		{
			name: "Go with go.mod",
			files: map[string]string{
				"go.mod": "module example.com/test",
			},
			expected:    "go",
			shouldError: false,
		},
		{
			name: "Node.js with package.json",
			files: map[string]string{
				"package.json": `{"name": "test", "version": "1.0.0"}`,
			},
			expected:    "node",
			shouldError: false,
		},
		{
			name: "Java with pom.xml",
			files: map[string]string{
				"pom.xml": "<project></project>",
			},
			expected:    "java",
			shouldError: false,
		},
		{
			name: "Rust with Cargo.toml",
			files: map[string]string{
				"Cargo.toml": "[package]",
			},
			expected:    "rust",
			shouldError: false,
		},
		{
			name:        "No language files",
			files:       map[string]string{},
			expected:    "",
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory
			tmpDir, err := os.MkdirTemp("", "techstack-test-*")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			// Create test files
			for filename, content := range tt.files {
				filePath := filepath.Join(tmpDir, filename)
				if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
			}

			// Test detection
			lang, err := detector.detectLanguageFallback(tmpDir)

			if tt.shouldError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if lang != tt.expected {
				t.Errorf("Expected language %q, got %q", tt.expected, lang)
			}
		})
	}
}

func TestGenerateImageTag(t *testing.T) {
	detector := NewDetector()

	tests := []struct {
		language string
		version  string
		expected string
	}{
		{"python", "3.12", "python:3.12-slim"},
		{"go", "1.24", "golang:1.24-alpine"},
		{"node", "20", "node:20-alpine"},
		{"java", "17", "eclipse-temurin:17-jdk-alpine"},
		{"rust", "1.75", "rust:1.75-alpine"},
		{"dotnet", "8", "mcr.microsoft.com/dotnet/sdk:8-alpine"},
		{"unknown", "1.0", "unknown:1.0-alpine"}, // fallback
	}

	for _, tt := range tests {
		t.Run(tt.language, func(t *testing.T) {
			result := detector.generateImageTag(tt.language, tt.version)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestFileExists(t *testing.T) {
	// Create temporary directory with test files
	tmpDir, err := os.MkdirTemp("", "techstack-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name     string
		pattern  string
		expected bool
	}{
		{"Existing file", "test.txt", true},
		{"Non-existing file", "missing.txt", false},
		{"Wildcard match", "*.txt", true},
		{"Wildcard no match", "*.go", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fileExists(tmpDir, tt.pattern)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v for pattern %q", tt.expected, result, tt.pattern)
			}
		})
	}
}

func TestHasWildcard(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"test.txt", false},
		{"*.txt", true},
		{"test?.go", true},
		{"test[0-9].log", true},
		{"no-wildcard", false},
		{"path/to/file", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := hasWildcard(tt.input)
			if result != tt.expected {
				t.Errorf("hasWildcard(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestDetect(t *testing.T) {
	detector := NewDetector()

	tests := []struct {
		name        string
		files       map[string]string
		expected    string
		shouldError bool
	}{
		{
			name: "Python project",
			files: map[string]string{
				"requirements.txt": "flask==2.0.0",
			},
			expected:    "python",
			shouldError: false,
		},
		{
			name: "Go project",
			files: map[string]string{
				"go.mod": "module example.com/test\n\ngo 1.22",
			},
			expected:    "go",
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory
			tmpDir, err := os.MkdirTemp("", "techstack-test-*")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			// Create test files
			for filename, content := range tt.files {
				filePath := filepath.Join(tmpDir, filename)
				if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
			}

			// Test detection
			result, err := detector.Detect(tmpDir)

			if tt.shouldError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if result == nil {
				t.Fatal("Expected result, got nil")
			}

			if result.Language.Name != tt.expected {
				t.Errorf("Expected language %q, got %q", tt.expected, result.Language.Name)
			}

			// Verify other fields are populated
			if result.Language.CIImageTag == "" {
				t.Error("CIImageTag should not be empty")
			}
		})
	}
}

func TestDetectNonExistentPath(t *testing.T) {
	detector := NewDetector()

	_, err := detector.Detect("/path/that/does/not/exist")
	if err == nil {
		t.Error("Expected error for non-existent path")
	}
}
