package models

import (
	"strings"
	"testing"
)

func TestLanguageStruct(t *testing.T) {
	lang := Language{
		Name:       "python",
		Version:    "3.12",
		BuildTool:  "poetry",
		CIImageTag: "python:3.12-slim",
	}

	if lang.Name != "python" {
		t.Errorf("Expected Name to be 'python', got %q", lang.Name)
	}
	if lang.Version != "3.12" {
		t.Errorf("Expected Version to be '3.12', got %q", lang.Version)
	}
	if lang.BuildTool != "poetry" {
		t.Errorf("Expected BuildTool to be 'poetry', got %q", lang.BuildTool)
	}
	if lang.CIImageTag != "python:3.12-slim" {
		t.Errorf("Expected CIImageTag to be 'python:3.12-slim', got %q", lang.CIImageTag)
	}
}

func TestToEnv(t *testing.T) {
	ts := TechStack{
		Language: Language{
			Name:       "go",
			Version:    "1.24",
			BuildTool:  "go",
			CIImageTag: "golang:1.24-alpine",
		},
	}

	env := ts.ToEnv()

	// Check format
	lines := strings.Split(strings.TrimSpace(env), "\n")
	if len(lines) != 4 {
		t.Errorf("Expected 4 lines, got %d", len(lines))
	}

	// Check each line
	expectedLines := map[string]bool{
		"LANGUAGE_NAME=go":                true,
		"LANGUAGE_VERSION=1.24":           true,
		"BUILD_TOOL=go":                   true,
		"CI_IMAGE_TAG=golang:1.24-alpine": true,
	}

	for _, line := range lines {
		if !expectedLines[line] {
			t.Errorf("Unexpected line: %q", line)
		}
	}
}

func TestToEnvWithEmptyValues(t *testing.T) {
	ts := TechStack{
		Language: Language{
			Name:       "unknown",
			Version:    "",
			BuildTool:  "",
			CIImageTag: "",
		},
	}

	env := ts.ToEnv()

	// Should still generate all lines, just with empty values
	lines := strings.Split(strings.TrimSpace(env), "\n")
	if len(lines) != 4 {
		t.Errorf("Expected 4 lines even with empty values, got %d", len(lines))
	}

	// Verify keys exist
	for _, line := range lines {
		if !strings.Contains(line, "=") {
			t.Errorf("Line should contain '=': %q", line)
		}
	}
}
