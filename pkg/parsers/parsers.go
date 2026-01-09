package parsers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// DetectVersion detects the version for a given language
func DetectVersion(path, language, buildTool string) string {
	switch language {
	case "python":
		return detectPythonVersion(path, buildTool)
	case "java":
		return detectJavaVersion(path, buildTool)
	case "kotlin":
		return detectKotlinVersion(path, buildTool)
	case "node", "javascript", "typescript":
		return detectNodeVersion(path)
	case "go":
		return detectGoVersion(path)
	case "rust":
		return detectRustVersion(path)
	case "ruby":
		return detectRubyVersion(path)
	case "php":
		return detectPHPVersion(path)
	case "dotnet", "csharp":
		return detectDotNetVersion(path)
	case "swift":
		return detectSwiftVersion(path)
	case "scala":
		return detectScalaVersion(path, buildTool)
	default:
		return ""
	}
}

// DetectBuildTool detects the build tool for a given language
func DetectBuildTool(path, language string) string {
	switch language {
	case "python":
		return detectPythonBuildTool(path)
	case "java":
		return detectJavaBuildTool(path)
	case "kotlin":
		return detectKotlinBuildTool(path)
	case "node", "javascript", "typescript":
		return detectNodeBuildTool(path)
	case "go":
		return "go"
	case "rust":
		return "cargo"
	case "ruby":
		return "bundle"
	case "php":
		return "composer"
	case "dotnet", "csharp":
		return "dotnet"
	case "swift":
		return "swift"
	case "scala":
		return detectScalaBuildTool(path)
	default:
		return ""
	}
}

// Python detection
func detectPythonVersion(path, buildTool string) string {
	// Try .python-version
	if content := readFile(filepath.Join(path, ".python-version")); content != "" {
		return strings.TrimSpace(content)
	}

	// Try pyproject.toml
	if buildTool == "poetry" || buildTool == "pdm" || buildTool == "hatch" {
		if content := readFile(filepath.Join(path, "pyproject.toml")); content != "" {
			re := regexp.MustCompile(`python\s*=\s*["']([^"']+)["']`)
			if match := re.FindStringSubmatch(content); match != nil {
				version := strings.TrimPrefix(match[1], "^")
				version = strings.TrimPrefix(version, "~")
				version = strings.TrimPrefix(version, ">=")
				if parts := strings.Split(version, "."); len(parts) >= 2 {
					return parts[0] + "." + parts[1]
				}
			}
		}
	}

	// Try runtime.txt (Heroku)
	if content := readFile(filepath.Join(path, "runtime.txt")); content != "" {
		re := regexp.MustCompile(`python-(\d+\.\d+)`)
		if match := re.FindStringSubmatch(content); match != nil {
			return match[1]
		}
	}

	return ""
}

func detectPythonBuildTool(path string) string {
	if fileExists(path, "pyproject.toml") {
		content := readFile(filepath.Join(path, "pyproject.toml"))
		if strings.Contains(content, "[tool.poetry]") {
			return "poetry"
		}
		if strings.Contains(content, "[tool.pdm]") {
			return "pdm"
		}
		if strings.Contains(content, "[tool.hatch]") {
			return "hatch"
		}
		if strings.Contains(content, "[tool.uv]") || strings.Contains(content, "uv.lock") {
			return "uv"
		}
		return "pip"
	}
	if fileExists(path, "Pipfile") {
		return "pipenv"
	}
	if fileExists(path, "requirements.txt") || fileExists(path, "setup.py") {
		return "pip"
	}
	return "pip"
}

// Java detection
func detectJavaVersion(path, buildTool string) string {
	if buildTool == "maven" {
		if content := readFile(filepath.Join(path, "pom.xml")); content != "" {
			re := regexp.MustCompile(`<maven\.compiler\.source>(\d+)</maven\.compiler\.source>`)
			if match := re.FindStringSubmatch(content); match != nil {
				return match[1]
			}
			re = regexp.MustCompile(`<java\.version>(\d+)</java\.version>`)
			if match := re.FindStringSubmatch(content); match != nil {
				return match[1]
			}
		}
	}

	if buildTool == "gradle" {
		files := []string{"build.gradle", "build.gradle.kts"}
		for _, file := range files {
			if content := readFile(filepath.Join(path, file)); content != "" {
				re := regexp.MustCompile(`sourceCompatibility\s*=\s*['"]*(\d+)`)
				if match := re.FindStringSubmatch(content); match != nil {
					return match[1]
				}
				re = regexp.MustCompile(`JavaVersion\.VERSION_(\d+)`)
				if match := re.FindStringSubmatch(content); match != nil {
					return match[1]
				}
			}
		}
	}

	return ""
}

func detectJavaBuildTool(path string) string {
	if fileExists(path, "pom.xml") {
		return "maven"
	}
	if fileExists(path, "build.gradle") || fileExists(path, "build.gradle.kts") {
		return "gradle"
	}
	return "maven"
}

// Kotlin detection
func detectKotlinVersion(path, buildTool string) string {
	if buildTool == "gradle" {
		files := []string{"build.gradle.kts", "build.gradle"}
		for _, file := range files {
			if content := readFile(filepath.Join(path, file)); content != "" {
				re := regexp.MustCompile(`kotlin\("jvm"\)\s+version\s+"([\d.]+)"`)
				if match := re.FindStringSubmatch(content); match != nil {
					return match[1]
				}
			}
		}
	}
	return ""
}

func detectKotlinBuildTool(path string) string {
	if fileExists(path, "build.gradle.kts") || fileExists(path, "build.gradle") {
		return "gradle"
	}
	if fileExists(path, "pom.xml") {
		return "maven"
	}
	return "gradle"
}

// Node.js detection
func detectNodeVersion(path string) string {
	// Try .nvmrc
	if content := readFile(filepath.Join(path, ".nvmrc")); content != "" {
		version := strings.TrimSpace(content)
		version = strings.TrimPrefix(version, "v")
		return version
	}

	// Try package.json engines field
	if content := readFile(filepath.Join(path, "package.json")); content != "" {
		var pkg struct {
			Engines struct {
				Node string `json:"node"`
			} `json:"engines"`
		}
		if err := json.Unmarshal([]byte(content), &pkg); err == nil {
			if pkg.Engines.Node != "" {
				version := strings.TrimPrefix(pkg.Engines.Node, "^")
				version = strings.TrimPrefix(version, "~")
				version = strings.TrimPrefix(version, ">=")
				version = strings.TrimPrefix(version, "v")
				if parts := strings.Split(version, "."); len(parts) >= 1 {
					return parts[0]
				}
			}
		}
	}

	// Try .node-version
	if content := readFile(filepath.Join(path, ".node-version")); content != "" {
		version := strings.TrimSpace(content)
		version = strings.TrimPrefix(version, "v")
		return version
	}

	return ""
}

func detectNodeBuildTool(path string) string {
	if fileExists(path, "pnpm-lock.yaml") {
		return "pnpm"
	}
	if fileExists(path, "yarn.lock") {
		return "yarn"
	}
	if fileExists(path, "package-lock.json") {
		return "npm"
	}
	if fileExists(path, "bun.lockb") {
		return "bun"
	}
	return "npm"
}

// Go detection
func detectGoVersion(path string) string {
	if content := readFile(filepath.Join(path, "go.mod")); content != "" {
		re := regexp.MustCompile(`go\s+(\d+\.\d+)`)
		if match := re.FindStringSubmatch(content); match != nil {
			return match[1]
		}
	}
	return ""
}

// Rust detection
func detectRustVersion(path string) string {
	if content := readFile(filepath.Join(path, "rust-toolchain")); content != "" {
		return strings.TrimSpace(content)
	}
	if content := readFile(filepath.Join(path, "rust-toolchain.toml")); content != "" {
		re := regexp.MustCompile(`channel\s*=\s*"([^"]+)"`)
		if match := re.FindStringSubmatch(content); match != nil {
			return match[1]
		}
	}
	return ""
}

// Ruby detection
func detectRubyVersion(path string) string {
	if content := readFile(filepath.Join(path, ".ruby-version")); content != "" {
		return strings.TrimSpace(content)
	}
	if content := readFile(filepath.Join(path, "Gemfile")); content != "" {
		re := regexp.MustCompile(`ruby ['"](\d+\.\d+)`)
		if match := re.FindStringSubmatch(content); match != nil {
			return match[1]
		}
	}
	return ""
}

// PHP detection
func detectPHPVersion(path string) string {
	if content := readFile(filepath.Join(path, "composer.json")); content != "" {
		var composer struct {
			Require struct {
				PHP string `json:"php"`
			} `json:"require"`
		}
		if err := json.Unmarshal([]byte(content), &composer); err == nil {
			version := composer.Require.PHP
			version = strings.TrimPrefix(version, "^")
			version = strings.TrimPrefix(version, "~")
			version = strings.TrimPrefix(version, ">=")
			if parts := strings.Split(version, "."); len(parts) >= 2 {
				return parts[0] + "." + parts[1]
			}
		}
	}
	return ""
}

// .NET detection
func detectDotNetVersion(path string) string {
	// Check for global.json first
	if content := readFile(filepath.Join(path, "global.json")); content != "" {
		var globalJSON struct {
			SDK struct {
				Version string `json:"version"`
			} `json:"sdk"`
		}
		if err := json.Unmarshal([]byte(content), &globalJSON); err == nil {
			if globalJSON.SDK.Version != "" {
				parts := strings.Split(globalJSON.SDK.Version, ".")
				if len(parts) >= 2 {
					if parts[1] == "0" {
						return parts[0]
					}
					return parts[0] + "." + parts[1]
				}
			}
		}
	}

	// Search for .csproj files (recursively up to 3 levels)
	patterns := []string{"*.csproj", "*/*.csproj", "*/*/*.csproj"}
	versions := make(map[int]int) // major -> minor

	for _, pattern := range patterns {
		matches, _ := filepath.Glob(filepath.Join(path, pattern))
		for _, file := range matches {
			if content := readFile(file); content != "" {
				// Handle TargetFrameworks (plural)
				re := regexp.MustCompile(`<TargetFrameworks>([^<]+)</TargetFrameworks>`)
				if match := re.FindStringSubmatch(content); match != nil {
					frameworks := strings.Split(match[1], ";")
					for _, fw := range frameworks {
						if major, minor := parseNetVersion(fw); major >= 5 {
							if existingMinor, ok := versions[major]; !ok || minor > existingMinor {
								versions[major] = minor
							}
						}
					}
				}

				// Handle TargetFramework (singular)
				re = regexp.MustCompile(`<TargetFramework>([^<]+)</TargetFramework>`)
				if match := re.FindStringSubmatch(content); match != nil {
					if major, minor := parseNetVersion(match[1]); major >= 5 {
						if existingMinor, ok := versions[major]; !ok || minor > existingMinor {
							versions[major] = minor
						}
					}
				}
			}
		}
	}

	// Return the latest version
	if len(versions) > 0 {
		maxMajor := 0
		for major := range versions {
			if major > maxMajor {
				maxMajor = major
			}
		}
		minor := versions[maxMajor]
		if minor == 0 {
			return fmt.Sprintf("%d", maxMajor)
		}
		return fmt.Sprintf("%d.%d", maxMajor, minor)
	}

	return ""
}

func parseNetVersion(framework string) (int, int) {
	// Only match modern .NET: net5.0, net6.0, net8.0, net10.0, etc.
	// Explicitly avoid netstandard, netcoreapp, net4xx (Framework)
	re := regexp.MustCompile(`^net(\d+)\.(\d+)$`)
	if match := re.FindStringSubmatch(framework); match != nil {
		major, _ := strconv.Atoi(match[1])
		minor, _ := strconv.Atoi(match[2])
		// Only modern .NET (5+)
		if major >= 5 {
			return major, minor
		}
	}
	return 0, 0
}

// Swift detection
func detectSwiftVersion(path string) string {
	if content := readFile(filepath.Join(path, "Package.swift")); content != "" {
		re := regexp.MustCompile(`swift-tools-version:\s*(\d+\.\d+)`)
		if match := re.FindStringSubmatch(content); match != nil {
			return match[1]
		}
	}
	return ""
}

// Scala detection
func detectScalaVersion(path, buildTool string) string {
	if buildTool == "sbt" {
		if content := readFile(filepath.Join(path, "build.sbt")); content != "" {
			re := regexp.MustCompile(`scalaVersion\s*:=\s*"(\d+\.\d+)`)
			if match := re.FindStringSubmatch(content); match != nil {
				return match[1]
			}
		}
	}
	if buildTool == "gradle" {
		files := []string{"build.gradle", "build.gradle.kts"}
		for _, file := range files {
			if content := readFile(filepath.Join(path, file)); content != "" {
				re := regexp.MustCompile(`scala\s+plugin.*version\s+"(\d+\.\d+)"`)
				if match := re.FindStringSubmatch(content); match != nil {
					return match[1]
				}
			}
		}
	}
	return ""
}

func detectScalaBuildTool(path string) string {
	if fileExists(path, "build.sbt") {
		return "sbt"
	}
	if fileExists(path, "build.gradle") || fileExists(path, "build.gradle.kts") {
		return "gradle"
	}
	return "sbt"
}

// Helper functions
func readFile(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(content)
}

func fileExists(basePath, filename string) bool {
	_, err := os.Stat(filepath.Join(basePath, filename))
	return err == nil
}

// POM represents a Maven POM file structure
type POM struct {
	XMLName    xml.Name `xml:"project"`
	Properties struct {
		JavaVersion         string `xml:"java.version"`
		MavenCompilerSource string `xml:"maven.compiler.source"`
		MavenCompilerTarget string `xml:"maven.compiler.target"`
	} `xml:"properties"`
}
