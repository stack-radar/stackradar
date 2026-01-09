package parsers

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectBuildTool(t *testing.T) {
	tests := []struct {
		name     string
		language string
		files    map[string]string
		expected string
	}{
		{"Python poetry", "python", map[string]string{"pyproject.toml": "[tool.poetry]\\nname = 'test'"}, "poetry"},
		{"Python pip", "python", map[string]string{"requirements.txt": ""}, "pip"},
		{"Java Maven", "java", map[string]string{"pom.xml": ""}, "maven"},
		{"Java Gradle", "java", map[string]string{"build.gradle": ""}, "gradle"},
		{"Node npm", "node", map[string]string{"package.json": "{}", "package-lock.json": ""}, "npm"},
		{"Node yarn", "node", map[string]string{"package.json": "{}", "yarn.lock": ""}, "yarn"},
		{"Node pnpm", "node", map[string]string{"package.json": "{}", "pnpm-lock.yaml": ""}, "pnpm"},
		{"Rust cargo", "rust", map[string]string{"Cargo.toml": ""}, "cargo"},
		{"Go", "go", map[string]string{"go.mod": ""}, "go"},
		{"Ruby bundle", "ruby", map[string]string{"Gemfile": ""}, "bundle"},
		{"PHP composer", "php", map[string]string{"composer.json": ""}, "composer"},
		{".NET dotnet", "dotnet", map[string]string{"project.csproj": ""}, "dotnet"},
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
			result := DetectBuildTool(tmpDir, tt.language)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestDetectGoVersion(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "Go 1.22",
			content:  "module example.com/test\n\ngo 1.22",
			expected: "1.22",
		},
		{
			name:     "Go 1.24",
			content:  "module example.com/test\n\ngo 1.24",
			expected: "1.24",
		},
		{
			name:     "Go with toolchain",
			content:  "module example.com/test\n\ngo 1.22\n\ntoolchain go1.22.0",
			expected: "1.22",
		},
		{
			name:     "No version",
			content:  "module example.com/test",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "techstack-test-*")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			goModPath := filepath.Join(tmpDir, "go.mod")
			if err := os.WriteFile(goModPath, []byte(tt.content), 0644); err != nil {
				t.Fatalf("Failed to create go.mod: %v", err)
			}

			result := detectGoVersion(tmpDir)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestDetectPythonVersion(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		content  string
		expected string
	}{
		{
			name:     ".python-version",
			filename: ".python-version",
			content:  "3.11.5",
			expected: "3.11.5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "techstack-test-*")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			filePath := filepath.Join(tmpDir, tt.filename)
			if err := os.WriteFile(filePath, []byte(tt.content), 0644); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			result := detectPythonVersion(tmpDir, "")
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestDetectNodeVersion(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		content  string
		expected string
	}{
		{
			name:     ".nvmrc",
			filename: ".nvmrc",
			content:  "20.10.0",
			expected: "20.10.0",
		},
		{
			name:     "package.json with engines",
			filename: "package.json",
			content:  `{"engines": {"node": ">=18.0.0"}}`,
			expected: "18",
		},
		{
			name:     ".node-version",
			filename: ".node-version",
			content:  "v20.11.0",
			expected: "20.11.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "techstack-test-*")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			filePath := filepath.Join(tmpDir, tt.filename)
			if err := os.WriteFile(filePath, []byte(tt.content), 0644); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			result := detectNodeVersion(tmpDir)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestDetectVersion(t *testing.T) {
	tests := []struct {
		name      string
		language  string
		buildTool string
		files     map[string]string
		expected  string
	}{
		{
			name:      "Python",
			language:  "python",
			buildTool: "poetry",
			files: map[string]string{
				".python-version": "3.12.0",
			},
			expected: "3.12.0",
		},
		{
			name:      "Go",
			language:  "go",
			buildTool: "go",
			files: map[string]string{
				"go.mod": "module test\n\ngo 1.24",
			},
			expected: "1.24",
		},
		{
			name:      "Node.js",
			language:  "node",
			buildTool: "npm",
			files: map[string]string{
				".nvmrc": "20",
			},
			expected: "20",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "techstack-test-*")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			for filename, content := range tt.files {
				filePath := filepath.Join(tmpDir, filename)
				if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
			}

			result := DetectVersion(tmpDir, tt.language, tt.buildTool)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}
