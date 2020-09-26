package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	debug bool

	rootCmd = &cobra.Command{
		Use:   "classic-bg-ocr",
		Short: "Classic WoW BG timer capture and analysis",
		Long:  `Classic WoW battleground timer alerts`,
	}
)

func init() {
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "print debug messages")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Fatal error: %v\n", err)
		os.Exit(1)
	}
}
