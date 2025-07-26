package cmd

import (
	"github.com/spf13/cobra"
)

var dnsCmd = &cobra.Command{
	Use:   "dns",
	Short: "🛰️ Get information about resolution of a domain name.",
}

func init() {
	rootCmd.AddCommand(dnsCmd)
}
