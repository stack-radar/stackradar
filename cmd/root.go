package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "1.0.0"

var rootCmd = &cobra.Command{
	Use:   "stackradar",
	Short: "StackRadar - Enterprise tech stack detection for CI/CD automation",
	Long: `StackRadar analyzes repositories to detect programming languages,
versions, build tools, and generates appropriate CI/CD Docker image tags.

Supports: Python, Java, Kotlin, Node.js, Go, Rust, Ruby, PHP, .NET, Swift, Scala`,
	Version: version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
