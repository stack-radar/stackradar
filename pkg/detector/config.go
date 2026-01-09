package detector

// LanguageConfig defines the configuration for a language
type LanguageConfig struct {
	Name           string
	FileIndicators []string
	ImageTemplate  string
	DefaultVersion string
}

// Config holds all language configurations
var Config = map[string]LanguageConfig{
	"python": {
		Name:           "python",
		FileIndicators: []string{"requirements.txt", "setup.py", "pyproject.toml", "Pipfile"},
		ImageTemplate:  "python:%s-slim",
		DefaultVersion: "3.12",
	},
	"java": {
		Name:           "java",
		FileIndicators: []string{"pom.xml", "build.gradle", "build.gradle.kts"},
		ImageTemplate:  "eclipse-temurin:%s-jdk-alpine",
		DefaultVersion: "17",
	},
	"kotlin": {
		Name:           "kotlin",
		FileIndicators: []string{"build.gradle.kts"},
		ImageTemplate:  "eclipse-temurin:%s-jdk-alpine",
		DefaultVersion: "17",
	},
	"node": {
		Name:           "node",
		FileIndicators: []string{"package.json"},
		ImageTemplate:  "node:%s-alpine",
		DefaultVersion: "20",
	},
	"javascript": {
		Name:           "javascript",
		FileIndicators: []string{"package.json"},
		ImageTemplate:  "node:%s-alpine",
		DefaultVersion: "20",
	},
	"typescript": {
		Name:           "typescript",
		FileIndicators: []string{"tsconfig.json", "package.json"},
		ImageTemplate:  "node:%s-alpine",
		DefaultVersion: "20",
	},
	"go": {
		Name:           "go",
		FileIndicators: []string{"go.mod"},
		ImageTemplate:  "golang:%s-alpine",
		DefaultVersion: "1.22",
	},
	"rust": {
		Name:           "rust",
		FileIndicators: []string{"Cargo.toml"},
		ImageTemplate:  "rust:%s-alpine",
		DefaultVersion: "1.75",
	},
	"ruby": {
		Name:           "ruby",
		FileIndicators: []string{"Gemfile"},
		ImageTemplate:  "ruby:%s-alpine",
		DefaultVersion: "3.3",
	},
	"php": {
		Name:           "php",
		FileIndicators: []string{"composer.json"},
		ImageTemplate:  "php:%s-cli-alpine",
		DefaultVersion: "8.3",
	},
	"dotnet": {
		Name:           "dotnet",
		FileIndicators: []string{"*.csproj", "*.sln", "*.slnx", "*/*.csproj"},
		ImageTemplate:  "mcr.microsoft.com/dotnet/sdk:%s-alpine",
		DefaultVersion: "8.0",
	},
	"csharp": {
		Name:           "csharp",
		FileIndicators: []string{"*.csproj", "*.sln", "*.slnx", "*/*.csproj"},
		ImageTemplate:  "mcr.microsoft.com/dotnet/sdk:%s-alpine",
		DefaultVersion: "8.0",
	},
	"swift": {
		Name:           "swift",
		FileIndicators: []string{"Package.swift"},
		ImageTemplate:  "swift:%s",
		DefaultVersion: "5.9",
	},
	"scala": {
		Name:           "scala",
		FileIndicators: []string{"build.sbt"},
		ImageTemplate:  "eclipse-temurin:%s-jdk-alpine",
		DefaultVersion: "17",
	},
}
