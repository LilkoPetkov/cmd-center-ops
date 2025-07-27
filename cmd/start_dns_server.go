package cmd

import (
	"commandCenter/styles"
	"commandCenter/validators"
	"fmt"
	"log"

	"github.com/miekg/dns"
	"github.com/spf13/cobra"
)

var startServerCmd = &cobra.Command{
	Use:        "dns",
	Short:      "Start DNS server on specified port.",
	Long:       "Start a DNS server on specified port that can be used for DNS testing.",
	Aliases:    []string{"server", "init", "initialize"},
	SuggestFor: []string{"serve", "ser", "serve"},
	Example: `
      # Start a DNS server on the default port 8888
      ops server dns

      # Start a DNS server on the standard DNS port 53
      ops server dns -p 53

      # Start a DNS server on a custom port, e.g., 5353 for mDNS testing
      ops server dns -p 5353

      # Get help for the DNS server command
      ops server dns --help
    `,

	Run: startDNSServer,
}

// init initializes the startServerCmd and its flags.
//
// Args:
//   - None
//
// Returns:
//   - None
func init() {
	startServerCmd.Flags().StringP("port", "p", "8888", "port for the DNS server")

	connectCmd.AddCommand(startServerCmd)
}

// startDNSServer starts a DNS server on the specified port.
//
// Args:
//   - cmd: The cobra command.
//   - args: The command arguments.
//
// Returns:
//   - None
func startDNSServer(cmd *cobra.Command, args []string) {
	port, err := validators.VerifyStringInputs(cmd, "port")
	if err != nil {
		log.Fatalln(err)
	}

	server := &dns.Server{
		Addr:      fmt.Sprintf(":%s", port),
		Net:       "udp",
		UDPSize:   65535,
		ReusePort: true,
	}

	fmt.Printf(styles.NewStyles().Highlight.Render("Startiung DNS server on port %s"), port)
	err = server.ListenAndServe()
	if err != nil {
		fmt.Printf(styles.NewStyles().Error.Render("Failed to start server: %s\n"), err.Error())
	}
}
