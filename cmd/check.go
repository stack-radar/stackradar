package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stack-radar/stackradar/pkg/detector"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check dependencies",
	Long:  `Checks if required dependencies like GitHub Linguist are available.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Checking dependencies...")
		fmt.Println()

		d := detector.NewDetector()
		if d.LinguistAvailable() {
			fmt.Println("✓ GitHub Linguist: Available")
		} else {
			fmt.Println("⚠ GitHub Linguist: Not available")
		}
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
