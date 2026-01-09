# 🔍 TechStack Detector

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)](https://go.dev/)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

> **Automatically detect programming languages, versions, build tools, and generate standardized CI/CD Docker image tags from any repository.**

TechStack Detector is a fast, zero-configuration command-line tool that analyzes your codebase and provides complete tech stack information. Perfect for automating CI/CD pipelines, standardizing Docker builds, multi-repo management, and DevOps automation.

### 🎯 Why TechStack Detector?

- **Zero Configuration** - Works out of the box with sensible defaults
- **Multi-Language** - Supports 12+ languages and 20+ build tools
- **CI/CD Ready** - Generates official Docker Hub image tags automatically
- **Lightning Fast** - Written in Go, optimized for performance
- **DevOps Friendly** - Output in YAML, JSON, or environment variables
- **Open Source** - MIT licensed, community-driven



## ⚡ Quick Start

```bash
# Detect current directory
stackradar get

# Detect specific directory
stackradar get --path /path/to/repo

# Output formats: yaml (default), json, env
stackradar get --format json

# Save to file and use in CI/CD
stackradar get --format env --output .env
```

**Example Output:**
```yaml
language:
  name: python
  version: "3.12"
  build_tool: poetry
  ci_image_tag: python:3.12-slim
```

## 📦 Installation

### Option 1: Download Pre-built Binary (Recommended)

**Linux (Ubuntu/Debian):**
```bash
wget https://github.com/stack-radar/stackradar/releases/latest/download/stackradar-linux-amd64
chmod +x stackradar-linux-amd64
sudo mv stackradar-linux-amd64 /usr/local/bin/stackradar
```

**macOS:**
```bash
# Apple Silicon (M1/M2/M3)
curl -L https://github.com/stack-radar/stackradar/releases/latest/download/stackradar-darwin-arm64 -o stackradar

# Intel Macs:
curl -L https://github.com/stack-radar/stackradar/releases/latest/download/stackradar-darwin-amd64 -o stackradar

chmod +x stackradar
sudo mv stackradar /usr/local/bin/
```

**Windows:**
Download `stackradar-windows-amd64.exe` from [releases](https://github.com/stack-radar/stackradar/releases) and add to PATH.

### Option 2: Build from Source

Requires Go 1.24.12+ or Go 1.25.6+:

```bash
git clone https://github.com/stack-radar/stackradar.git
cd stackradar
make build
sudo cp build/stackradar /usr/local/bin/
```

### Option 3: Go Install

```bash
go install github.com/stack-radar/stackradar@latest
```

## 🗂️ Supported Tech Stack

| Language | Version Detection | Build Tools | CI Image Format |
|----------|------------------|-------------|-----------------|
| **Python** | `.python-version`, `pyproject.toml` | pip, poetry, pdm, pipenv, hatch, uv | `python:{version}-slim` |
| **Java** | `pom.xml`, `build.gradle` | maven, gradle | `eclipse-temurin:{version}-jdk-alpine` |
| **Kotlin** | `build.gradle.kts` | gradle, maven | `eclipse-temurin:{version}-jdk-alpine` |
| **Node.js** | `.nvmrc`, `package.json` | npm, yarn, pnpm, bun | `node:{version}-alpine` |
| **TypeScript** | `.nvmrc`, `package.json` | npm, yarn, pnpm, bun | `node:{version}-alpine` |
| **Go** | `go.mod` | go | `golang:{version}-alpine` |
| **Rust** | `rust-toolchain`, `Cargo.toml` | cargo | `rust:{version}-alpine` |
| **Ruby** | `.ruby-version`, `Gemfile` | bundle | `ruby:{version}-alpine` |
| **PHP** | `composer.json` | composer | `php:{version}-cli-alpine` |
| **.NET/C#** | `*.csproj`, `global.json` | dotnet | `mcr.microsoft.com/dotnet/sdk:{version}-alpine` |
| **Swift** | `Package.swift` | swift | `swift:{version}` |
| **Scala** | `build.sbt` | sbt, gradle | `hseeberger/scala-sbt:{version}` |

### 📝 Notes on Detection

- **Build Tool Versions**: Not detected. Build tools (Maven, Gradle, npm, etc.) are identified by name only, as their versions are typically managed by project config files (e.g., `gradlew`, `package-lock.json`)
- **Multi-Language Projects**: Currently detects primary language. Multi-language support is on the roadmap
- **Version Fallbacks**: If version cannot be detected, uses sensible defaults from config

## 💡 Usage Examples

### Basic Detection

```bash
$ stackradar get --path ./my-python-project
language:
  name: python
  version: "3.12"
  build_tool: poetry
  ci_image_tag: python:3.12-slim
```

### JSON Output for CI/CD Pipelines

```bash
$ stackradar get --format json
{
  "language": {
    "name": "go",
    "version": "1.24",
    "build_tool": "go",
    "ci_image_tag": "golang:1.24-alpine"
  }
}
```

### Environment Variables for Shell Scripts

```bash
$ stackradar get --format env
LANGUAGE_NAME=node
LANGUAGE_VERSION=20
BUILD_TOOL=pnpm
CI_IMAGE_TAG=node:20-alpine

# Use in shell scripts
eval $(stackradar get --format env)
docker build --build-arg BASE_IMAGE=$CI_IMAGE_TAG .
```

### 🚀 CI/CD Integration Examples

#### GitHub Actions

```yaml
name: Build
on: [push]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Detect Tech Stack
        id: stackradar
        run: |
          wget -q https://github.com/stack-radar/stackradar/releases/latest/download/stackradar-linux-amd64 -O stackradar
          chmod +x stackradar
          ./stackradar get --format env >> $GITHUB_ENV
      
      - name: Build with detected image
        run: |
          echo "Building with $CI_IMAGE_TAG"
          docker build --build-arg BASE_IMAGE=$CI_IMAGE_TAG -t myapp .
```

#### GitLab CI

```yaml
detect:
  stage: detect
  script:
    - wget -q https://github.com/stack-radar/stackradar/releases/latest/download/stackradar-linux-amd64 -O stackradar
    - chmod +x stackradar
    - ./stackradar get --format env > techstack.env
  artifacts:
    reports:
      dotenv: techstack.env

build:
  stage: build
  image: ${CI_IMAGE_TAG}
  script:
    - echo "Building with ${BUILD_TOOL}"
    - make build
```

## 👨‍💻 Developer Setup

### Prerequisites

- **Go 1.24.12** or higher  # go version go1.24.1 darwin/arm64 
- Make (optional but recommended) 
- Git

### Clone and Build

```bash
# Clone repository
git clone https://github.com/stack-radar/stackradar.git
cd stackradar

# Build for current platform
make build

# Run tests
make test

# Build for all platforms
make build-all
```

### Available Make Commands

```bash
make build          # Build for current platform
make build-linux    # Build for Linux (amd64 + arm64)
make build-mac      # Build for macOS (Intel + Apple Silicon)
make build-windows  # Build for Windows (amd64)
make build-all      # Build for all platforms
make install        # Install to $GOPATH/bin
make test           # Run tests
make test-coverage  # Run tests with coverage report
make fmt            # Format code
make vet            # Run go vet
make lint           # Run golangci-lint (if installed)
make clean          # Remove build artifacts
make security       # Run all security scans
```

### Project Structure

```
stackradar/
├── cmd/                    # CLI commands (cobra)
│   ├── root.go            # Root command
│   ├── get.go             # Detection command
│   └── check.go           # Validation command
├── pkg/
│   ├── detector/
│   │   ├── config.go      # Language configurations
│   │   └── detector.go    # Detection logic
│   ├── parsers/
│   │   └── parsers.go     # Version parsers
│   └── models/
│       └── models.go      # Data models
├── .github/workflows/
│   └── build.yml          # CI/CD pipeline with security scans
├── main.go                # Entry point
├── go.mod                 # Go module
├── Makefile              # Build automation
└── README.md
```

### Running Tests

```bash
# Run all unit tests
make test

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
make test-coverage
# Opens coverage.html in browser

# Run specific package tests
go test -v ./pkg/detector
go test -v ./pkg/parsers
go test -v ./pkg/models

# Run specific test
go test -v ./pkg/detector -run TestDetect

# Run tests with race detection
go test -race ./...

# Run integration tests with real repositories
./test.sh

# Test specific repository
./build/stackradar get --path /path/to/repo
```

**Test Coverage:**
- `pkg/detector/detector_test.go` - Core detection logic
- `pkg/models/models_test.go` - Data models
- `pkg/parsers/parsers_test.go` - Version parsers

**What's Tested:**
- Language detection for Python, Go, Node.js, Java, Rust
- Build tool identification (Maven, Gradle, npm, yarn, pnpm, poetry, etc.)
- Version parsing from various config files
- CI image tag generation
- File pattern matching and wildcards
- Error handling for invalid paths
- Environment variable output formatting

## 🔒 Security Scanning

### Local Security Testing

Run comprehensive security checks before committing:

```bash
# All-in-one security scan
make security

# Individual security tools:

# 1. Trivy - Comprehensive vulnerability scanner
trivy fs .

# 2. govulncheck - Go vulnerability database check
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...

# 3. GoSec - Go security checker (SAST)
go install github.com/securego/gosec/v2/cmd/gosec@latest
gosec ./...

# 4. Go vet - Official Go static analysis
go vet ./...

# 5. staticcheck - Advanced Go linter
go install honnef.co/go/tools/cmd/staticcheck@latest
staticcheck ./...
```

### Installing Security Tools

**Trivy (Vulnerability Scanner):**
```bash
# macOS
brew install trivy

# Linux (Ubuntu/Debian)
sudo apt-get install wget apt-transport-https gnupg lsb-release
wget -qO - https://aquasecurity.github.io/trivy-repo/deb/public.key | sudo apt-key add -
echo "deb https://aquasecurity.github.io/trivy-repo/deb $(lsb_release -sc) main" | sudo tee -a /etc/apt/sources.list.d/trivy.list
sudo apt-get update
sudo apt-get install trivy

# Windows (using Chocolatey)
choco install trivy
```

**Go Security Tools:**
```bash
# govulncheck - Official Go vulnerability checker
go install golang.org/x/vuln/cmd/govulncheck@latest

# GoSec - Security-focused SAST tool
go install github.com/securego/gosec/v2/cmd/gosec@latest

# staticcheck - Advanced linter
go install honnef.co/go/tools/cmd/staticcheck@latest

# golangci-lint - Meta-linter (includes multiple security checks)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### CI/CD Security Pipeline

Our GitHub Actions workflow includes automated security scanning:

**Security Checks in CI/CD:**
1. **Trivy** - Scans for vulnerabilities in dependencies and code
2. **govulncheck** - Checks against Go vulnerability database
3. **GoSec** - Static analysis security testing (SAST)
4. **Go vet** - Official Go static analyzer

Results are automatically uploaded to GitHub Security tab (SARIF format).

View security scan results:
- Go to your repository → **Security** tab → **Code scanning alerts**

## 🔄 CI/CD Information

### GitHub Actions Workflow

**Go Version Matrix:**
- **Go 1.25.6** - Latest stable (January 2026)
- **Go 1.24.12** - Maintenance version (EOL when Go 1.26 releases)

**Build Matrix:**
- **Linux**: amd64, arm64
- **macOS**: amd64 (Intel), arm64 (Apple Silicon)
- **Windows**: amd64

**Automated Workflow Steps:**
1. **Security Scan**: Trivy, govulncheck, GoSec
2. **Test**: Unit tests with race detection
3. **Build**: Cross-compile for all platforms
4. **Release**: Auto-create GitHub releases (on version tags)

### Triggering CI/CD

```bash
# Test on pull request
git push origin feature-branch

# Create release
git tag v1.0.0
git push origin v1.0.0
```

### Release Process

1. Tag version: `git tag v1.0.0`
2. Push tag: `git push origin v1.0.0`
3. GitHub Actions automatically:
   - Runs security scans
   - Runs tests
   - Builds binaries for all platforms
   - Creates GitHub release with:
     - Binaries for all platforms
     - SHA256 checksums
     - Auto-generated release notes

## 🤝 Contributing

Contributions welcome! To add a new language:

1. **Add to config** (`pkg/detector/config.go`):
   ```go
   "newlang": {
       Name:           "newlang",
       FileIndicators: []string{"newlang.json", "*.newlang"},
       ImageTemplate:  "newlang:%s-alpine",
       DefaultVersion: "1.0",
   }
   ```

2. **Add version parser** (`pkg/parsers/parsers.go`):
   ```go
   func detectNewlangVersion(basePath string) string {
       // Parse version from language files
       return "1.0"
   }
   ```

3. **Update `DetectVersion()`** to route to your parser
4. **Add tests** in respective `*_test.go` files
5. **Submit PR** with documentation

## 🎯 Use Cases

- **Monorepo Management**: Detect different tech stacks across multiple services
- **CI/CD Automation**: Automatically select appropriate build images
- **Migration Planning**: Inventory tech stacks across organization
- **Dependency Auditing**: Track language versions and build tools
- **Docker Image Selection**: Generate standardized base images
- **DevOps Tooling**: Build internal tools that need language detection
- **Project Onboarding**: Quickly understand new codebases

## 🔧 Architecture

**Key Design Principles:**
- **Configuration-driven**: Add languages by editing config, not code
- **Minimal dependencies**: Only Cobra and YAML
- **Fast & lightweight**: ~5MB binary, written in Go
- **Maintainable**: Clean separation of concerns

**Detection Flow:**
1. Check for GitHub Linguist (optional)
2. Fallback to configuration-based file detection
3. Parse language-specific version files
4. Detect build tool from project files
5. Generate CI image tag from template

## 📜 License

MIT License - Free for personal and commercial use.

See [LICENSE](LICENSE) for details.

## 🌟 Star Us!

If you find this project useful, please ⭐ star us on GitHub!

## 📊 Project Information

- **Language**: Go 1.24+
- **License**: MIT
- **Platforms**: Linux, macOS, Windows
- **Architectures**: amd64, arm64
- **Binary Size**: ~5MB
- **Dependencies**: Minimal (Cobra + YAML)

## 💬 Community & Support

- 💡 [Discussions](https://github.com/stack-radar/stackradar/discussions) - Questions and ideas
- 🐛 [Issues](https://github.com/stack-radar/stackradar/issues) - Bug reports
- ✨ [Pull Requests](https://github.com/stack-radar/stackradar/pulls) - Contributions
- 📦 [Releases](https://github.com/stack-radar/stackradar/releases) - Downloads

---

<div align="center">

**Made with ❤️ by the open source community**

⭐ Star us on GitHub — it helps!

[Report Bug](https://github.com/stack-radar/stackradar/issues) · [Request Feature](https://github.com/stack-radar/stackradar/discussions) · [View Releases](https://github.com/stack-radar/stackradar/releases)

</div>
