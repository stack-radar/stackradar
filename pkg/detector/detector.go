package detector

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/stack-radar/stackradar/pkg/models"
	"github.com/stack-radar/stackradar/pkg/parsers"
)

// Detector handles tech stack detection
type Detector struct {
	linguistAvailable bool
	config            map[string]LanguageConfig
}

// NewDetector creates a new Detector instance
func NewDetector() *Detector {
	return &Detector{
		linguistAvailable: checkLinguist(),
		config:            Config,
	}
}

// LinguistAvailable returns whether GitHub Linguist is available
func (d *Detector) LinguistAvailable() bool {
	return d.linguistAvailable
}

func checkLinguist() bool {
	cmd := exec.Command("github-linguist", "--version")
	err := cmd.Run()
	return err == nil
}

// Detect analyzes a repository and returns tech stack information
func (d *Detector) Detect(path string) (*models.TechStack, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve path: %w", err)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("path does not exist: %s", path)
	}

	// 1. Detect language
	language, err := d.detectLanguage(absPath)
	if err != nil {
		return nil, err
	}

	// 2. Detect build tool
	buildTool := parsers.DetectBuildTool(absPath, language)

	// 3. Detect version
	version := parsers.DetectVersion(absPath, language, buildTool)
	if version == "" {
		// Use default version from config or fallback to "latest"
		if cfg, ok := d.config[language]; ok && cfg.DefaultVersion != "" {
			version = cfg.DefaultVersion
		} else {
			version = "latest"
		}
	}

	// 4. Generate CI image tag
	ciImageTag := d.generateImageTag(language, version)

	return &models.TechStack{
		Language: models.Language{
			Name:       language,
			Version:    version,
			BuildTool:  buildTool,
			CIImageTag: ciImageTag,
		},
	}, nil
}

// detectLanguage tries multiple detection methods
func (d *Detector) detectLanguage(path string) (string, error) {
	// Try Linguist first if available
	if d.linguistAvailable {
		lang, err := d.detectLanguageWithLinguist(path)
		if err == nil && lang != "" {
			return lang, nil
		}
	}

	// Fallback to file-based detection
	return d.detectLanguageFallback(path)
}

// detectLanguageWithLinguist uses GitHub Linguist for detection
func (d *Detector) detectLanguageWithLinguist(path string) (string, error) {
	cmd := exec.Command("github-linguist", "--json")
	cmd.Dir = path
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	var languages map[string]int
	if err := json.Unmarshal(output, &languages); err != nil {
		return "", err
	}

	if len(languages) == 0 {
		return "", fmt.Errorf("no languages detected")
	}

	// Find language with most bytes
	var maxLang string
	var maxBytes int
	for lang, bytes := range languages {
		if bytes > maxBytes {
			maxLang = lang
			maxBytes = bytes
		}
	}

	return strings.ToLower(maxLang), nil
}

// detectLanguageFallback uses configuration-based file pattern matching
func (d *Detector) detectLanguageFallback(path string) (string, error) {
	// Iterate through all configured languages
	for langName, langConfig := range d.config {
		for _, filePattern := range langConfig.FileIndicators {
			if fileExists(path, filePattern) {
				return langName, nil
			}
		}
	}

	return "", fmt.Errorf("unable to detect programming language")
}

// generateImageTag creates Docker image tag from configuration
func (d *Detector) generateImageTag(language, version string) string {
	// Look up the image template from configuration
	if cfg, ok := d.config[language]; ok && cfg.ImageTemplate != "" {
		return fmt.Sprintf(cfg.ImageTemplate, version)
	}

	// Default fallback for unconfigured languages
	return fmt.Sprintf("%s:%s-alpine", language, version)
}

// fileExists checks if a file or pattern exists in the given path
func fileExists(basePath, pattern string) bool {
	// Check for simple filename (no wildcards)
	if !hasWildcard(pattern) {
		_, err := os.Stat(filepath.Join(basePath, pattern))
		return err == nil
	}

	// Handle glob patterns
	matches, err := filepath.Glob(filepath.Join(basePath, pattern))
	return err == nil && len(matches) > 0
}

// hasWildcard checks if a string contains glob wildcards
func hasWildcard(s string) bool {
	return strings.ContainsAny(s, "*?[]")
}
