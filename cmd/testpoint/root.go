package main

import (
	"github.com/spf13/cobra"
	"runtime/debug"
)

var version = ""

// getVersion returns the current version of the application,
// if it hasn't been overridden through the VERSION variable.
func getVersion() string {
	if version != "" {
		return version
	}
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	return info.Main.Version
}

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "testpoint",
		Short:   "Testpoint is a simple CLI tool for testing REST endpoints",
		Version: getVersion(),
	}

	cmd.Root().CompletionOptions.DisableDefaultCmd = true

	cmd.AddCommand(
		newSendCmd(),
		newCompareCmd(),
	)

	return cmd
}
