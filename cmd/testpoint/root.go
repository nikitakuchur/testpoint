package main

import (
	"github.com/spf13/cobra"
)

var version = "undefined"

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "testpoint",
		Short:   "Testpoint is a simple CLI tool for testing REST endpoints",
		Version: version,
	}

	cmd.Root().CompletionOptions.DisableDefaultCmd = true

	cmd.AddCommand(
		newSendCmd(),
		newCompareCmd(),
	)

	return cmd
}
