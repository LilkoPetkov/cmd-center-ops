package cmd

import (
	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{
	Use:   "server",
	Short: "Start and manage network services like Telnet, DNS, or TCP.",
	Long: `A multifunction network utility that can start and manage Telnet sessions, run a custom DNS server, 
or launch a lightweight TCP server. Useful for testing, debugging, and developing network protocols or services.`,
}

// init initializes the server command.
//
// Args:
//   - None
//
// Returns:
//   - None
func init() {
	rootCmd.AddCommand(connectCmd)
}
