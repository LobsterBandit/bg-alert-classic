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
		Use:   "bg-alert-classic",
		Short: "Classic WoW BG queue alerts",
		Long:  `Classic WoW battleground queue alerts via screen capture and OCR`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			fmt.Printf("\n%s %v\n\t%s\n\n", cmd.Name(), cmd.Long, cmd.Version)
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
