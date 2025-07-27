package cmd

import (
	"github.com/spf13/cobra"
)

var yamlCmd = &cobra.Command{
	Use:   "yaml",
	Short: "⚙️ Manage YAML configuration files.",
}

// init initializes the yaml command.
//
// Args:
//   - None
//
// Returns:
//   - None
func init() {
	rootCmd.AddCommand(yamlCmd)
}
