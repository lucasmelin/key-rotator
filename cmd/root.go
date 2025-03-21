package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "key-rotator",
	Short:        "Rotate your secrets from the command line.",
	SilenceUsage: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddGroup(&cobra.Group{
		ID:    "core-commands",
		Title: "Core Commands:",
	}, &cobra.Group{
		ID:    "additional-commands",
		Title: "Additional Commands:",
	})
	rootCmd.SetHelpCommandGroupID("additional-commands")
	rootCmd.SetCompletionCommandGroupID("additional-commands")
	rootCmd.AddCommand(rotateCmd)
	rootCmd.AddCommand(versionCmd)
}
