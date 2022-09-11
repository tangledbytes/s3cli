package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:     "s3cli",
	Short:   "S3CLI is a stupid simple CLI for S3",
	Version: "v0.0.3",
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.AddCommand(
		ApiCmd,
		ListAPICmd,
	)
}
