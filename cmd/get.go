package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/stack-radar/stackradar/pkg/detector"
	"gopkg.in/yaml.v3"
)

var (
	path      string
	directory string
	format    string
	output    string
	quiet     bool
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Detect tech stack from local repository",
	Long:  `Analyzes a local repository to detect programming language, version, build tool, and CI image tag.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Use directory if set, otherwise use path
		repoPath := path
		if cmd.Flags().Changed("directory") {
			repoPath = directory
		}

		if !quiet {
			fmt.Println("🔍 TechStack Detector")
			fmt.Println()
		}

		// Detect tech stack
		d := detector.NewDetector()
		result, err := d.Detect(repoPath)
		if err != nil {
			return fmt.Errorf("detection failed: %w", err)
		}

		// Format output
		var outputText string
		switch format {
		case "json":
			data, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %w", err)
			}
			outputText = string(data)
		case "env":
			outputText = result.ToEnv()
		default: // yaml
			data, err := yaml.Marshal(result)
			if err != nil {
				return fmt.Errorf("failed to marshal YAML: %w", err)
			}
			outputText = string(data)
		}

		// Write output
		if output != "" {
			if err := os.WriteFile(output, []byte(outputText), 0644); err != nil {
				return fmt.Errorf("failed to write output file: %w", err)
			}
			if !quiet {
				fmt.Printf("✅ Output written to: %s\n", output)
			}
		} else {
			fmt.Print(outputText)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().StringVarP(&path, "path", "p", ".", "Local repository path")
	getCmd.Flags().StringVarP(&directory, "directory", "d", ".", "Local repository directory")
	getCmd.Flags().StringVarP(&format, "format", "f", "yaml", "Output format (yaml, json, env)")
	getCmd.Flags().StringVarP(&output, "output", "o", "", "Output file")
	getCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Suppress messages")
}
