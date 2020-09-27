package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	debug   bool
	version = "v1.0.0-alpha1"

	rootCmd = &cobra.Command{
		Use:   "classic-bg-ocr",
		Short: "Classic WoW BG timer capture and analysis",
		Long:  `Classic WoW battleground timer alerts via screen capture and OCR`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			fmt.Printf("\nclassic-bg-ocr %v\n\tWoW Classic BG timer screen capture and analysis\n\n", cmd.Version)
		},
	}
)

func init() {
	rootCmd.Version = version
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "print debug messages")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Fatal error: %v\n", err)
		os.Exit(1)
	}
}
