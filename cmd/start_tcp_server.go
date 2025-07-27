package cmd

import (
	"commandCenter/validators"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/spf13/cobra"
)

var startTCPServerCmd = &cobra.Command{
	Use:   "tcp",
	Short: "Start TCP server on specified port.",
	Long:  "Start a TCP server on specified port that can be used for network testing.",
	Example: `
      # Start a TCP server on default port 8888
      ops server tcp

      # Start a TCP server on port 9000
      ops server tcp -p 9000

      # Start a TCP server on a common HTTP port for testing
      ops server tcp -p 80

      # Start a TCP server on a high port for local development
      ops server tcp -p 3000

      # Get help for this command
      ops server tcp --help
    `,

	Run: startTCPServer,
}

// init initializes the startTCPServerCmd and its flags.
//
// Args:
//   - None
//
// Returns:
//   - None
func init() {
	startTCPServerCmd.Flags().StringP("port", "p", "8888", "port for the TCP server")

	connectCmd.AddCommand(startTCPServerCmd)
}

// handleConnection handles a single TCP connection.
//
// Args:
//   - c: The TCP connection.
//
// Returns:
//   - None
func handleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	packet := make([]byte, 4096)
	tmp := make([]byte, 4096)
	defer c.Close()
	for {
		_, err := c.Read(tmp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			break
		}
		packet = append(packet, tmp...)
	}
	c.Write(packet)
}

// startTCPServer starts a TCP server on the specified port.
//
// Args:
//   - cmd: The cobra command.
//   - args: The command arguments.
//
// Returns:
//   - None
func startTCPServer(cmd *cobra.Command, args []string) {
	port, err := validators.VerifyStringInputs(cmd, "port")
	if err != nil {
		log.Fatalln(err)
	}

	listener, err := net.Listen("tcp6", ":"+port)
	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	fmt.Printf("TCP6 server started on port: %s", port)
	for {
		c, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}

		go handleConnection(c)
	}

}
