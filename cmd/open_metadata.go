package cmd

import (
	"github.com/spf13/cobra"
)

var metadataCmd = &cobra.Command{
	Use:   "ometadata",
	Short: "Allows execution of Open MetaData commands.",
	Long: `Please do note that the commands require a python virtual environemnt and UV as
package manager for the setup and installation of Open Metadata and connector dependencies.
You would also need to prepare your Connector and the ingestion configuration yaml file
prior to the ingestion.`,
	Aliases: []string{"metadata", "open_meta_data"},
}

// init initializes the metadata command.
//
// Args:
//   - None
//
// Returns:
//   - None
func init() {
	rootCmd.AddCommand(metadataCmd)
}
