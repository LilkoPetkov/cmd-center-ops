package cmd

import (
	"commandCenter/styles"
	"commandCenter/validators"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/spf13/cobra"
)

type DestinationInterface interface {
	telnet()
}

type Destination struct {
	host string
	port string
}

var testServerConnectionCmd = &cobra.Command{
	Use:        "telnet",
	Short:      "Test connection to a server on specific port, similar to telnet.",
	Aliases:    []string{"test_connection", "conn"},
	SuggestFor: []string{"test_connection", "conn", "telnet"},
	Example: `
      # Test connection to localhost on default port 443
      ops telnet

      # Test HTTPS port on a public server
      ops telnet -n google.com -p 443

      # Check if an internal service is reachable on a custom port
      ops telnet -n 10.0.0.15 -p 8080

      # Test if SSH is open on a remote machine
      ops telnet -n example.org -p 22

      # Troubleshoot DNS server connectivity
      ops telnet -n 8.8.8.8 -p 53

      # Get help for this command
      ops telnet --help
    `,

	Run: testConnection,
}

// init initializes the testServerConnectionCmd and its flags.
//
// Args:
//   - None
//
// Returns:
//   - None
func init() {
	connectCmd.AddCommand(testServerConnectionCmd)

	testServerConnectionCmd.Flags().StringP("hostname", "n", "localhost", "host for which to test the connection")
	testServerConnectionCmd.Flags().StringP("port", "p", "443", "port which would be used for the connection")
}

// testConnection is the main function for the telnet command.
//
// Args:
//   - cmd: The cobra command.
//   - args: The command arguments.
//
// Returns:
//   - None
func testConnection(cmd *cobra.Command, args []string) {
	host, err := validators.VerifyStringInputs(cmd, "hostname")
	if err != nil {
		log.Fatalln(err)
	}

	port, err := validators.VerifyStringInputs(cmd, "port")
	if err != nil {
		log.Fatalln(err)
	}

	destination := Destination{
		host: host,
		port: port,
	}

	TestServerConnection(destination)
}

// telnet connects to a host and port.
//
// Args:
//   - None
//
// Returns:
//   - None
func (D Destination) telnet() {
	destination := fmt.Sprintf("%s:%s", D.host, D.port)

	conn, err := net.DialTimeout("tcp", destination, 5*time.Second)
	if err != nil {
		log.Fatalf(styles.NewStyles().Error.Render("Could not connect to host '%s' on port '%s'"), D.host, D.port)
	}

	defer conn.Close()

	fmt.Printf(styles.NewStyles().Highlight.Render("Successfully connected to %s"), destination)
}

// TestServerConnection tests the connection to a server.
//
// Args:
//   - D: The DestinationInterface.
//
// Returns:
//   - None
func TestServerConnection(D DestinationInterface) {
	D.telnet()
}
