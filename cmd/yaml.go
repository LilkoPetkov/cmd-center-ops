package cmd

import (
	"github.com/spf13/cobra"
)

var yamlCmd = &cobra.Command{
	Use:   "yaml",
	Short: "⚙️ Manage YAML configuration files.",
}

func init() {
	rootCmd.AddCommand(yamlCmd)
}
