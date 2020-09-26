package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "v1.0.0-alpha1"
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "classic-bg-ocr version number",
	Long:  `Print the version number of classic-bg-ocr`,
	Run:   runVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func runVersion(cmd *cobra.Command, args []string) {
	fmt.Println(version)
}
