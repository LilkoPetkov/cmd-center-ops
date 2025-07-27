package cmd

import (
	"commandCenter/styles"
	"commandCenter/validators"
	"fmt"
	"github.com/miekg/dns"
	"github.com/spf13/cobra"
	"log"
	"strings"
)

type DomainInterface interface {
	Resolve()
	ResolveAll()
	PrepareDnsCall(qtype string) *dns.Msg
}

type Domain struct {
	domainName  string
	qtype       string
	recordTypes []string
}

var resolve = &cobra.Command{
	Use:        "resolve",
	Short:      "Resolve a domain name.",
	Long:       "Get information about resolution of a domain name.",
	Aliases:    []string{"res", "resolv"},
	SuggestFor: []string{"resolve", "res", "resolv"},
	Example: `
      # Resolve the AAAA record for example.com (default)
      ops resolve -d example.com

      # Resolve the CNAME record for example.com
      ops resolve -d example.com -q cname

      # Resolve all main record types (A, AAAA, CNAME, NS, TXT) for example.com
      ops resolve -d example.com -a

      # Get help for the resolve command
      ops resolve --help
    `,

	Run: resolveDomain,
}

// init initializes the resolve command and its flags.
//
// Args:
//   - None
//
// Returns:
//   - None
func init() {
	dnsCmd.AddCommand(resolve)

	resolve.Flags().StringP("domain", "d", "example.com", "domain name to query for")
	resolve.Flags().StringP("qtype", "q", "AAAA", "record type to search for A/AAAA/cname/txt")
	resolve.Flags().BoolP("all", "a", false, "get information for all main records")

	resolve.MarkFlagsMutuallyExclusive("qtype", "all")
}

// resolveDomain is the main function for the resolve command.
//
// Args:
//   - cmd: The cobra command.
//   - args: The command arguments.
//
// Returns:
//   - None
func resolveDomain(cmd *cobra.Command, args []string) {
	domainName, err := validators.VerifyStringInputs(cmd, "domain")
	if err != nil {
		log.Fatalln(err)
	}

	qtype, err := validators.VerifyStringInputs(cmd, "qtype")
	if err != nil {
		log.Fatalln(err)
	}

	domain := Domain{
		domainName:  domainName,
		qtype:       qtype,
		recordTypes: []string{"ns", "a", "txt", "cname", "aaaa"},
	}

	all, err := validators.VerifyBoolInputs(cmd, "all")
	if err != nil {
		log.Fatalln(err)
	}

	if !all {
		ResolveDomain(domain)
	} else {
		ResolveAllRecords(domain)
	}

}

// ResolveDomain resolves a domain name.
//
// Args:
//   - D: The DomainInterface.
//
// Returns:
//   - None
func ResolveDomain(D DomainInterface) {
	D.Resolve()
}

// ResolveAllRecords resolves all records for a domain name.
//
// Args:
//   - D: The DomainInterface.
//
// Returns:
//   - None
func ResolveAllRecords(D DomainInterface) {
	D.ResolveAll()
}

// Resolve resolves a domain name.
//
// Args:
//   - None
//
// Returns:
//   - None
func (D Domain) Resolve() {
	in := D.PrepareDnsCall(D.qtype)

	for _, ans := range in.Answer {
		fmt.Printf(styles.NewStyles().Title.Render(
			"🛠️ %s Records 🛠️"), strings.ToUpper(D.qtype),
		)
		fmt.Println()
		fmt.Println(styles.NewStyles().Highlight.Render(ans.String()))
	}
}

// ResolveAll resolves all records for a domain name.
//
// Args:
//   - None
//
// Returns:
//   - None
func (D Domain) ResolveAll() {
	for _, dnsRecord := range D.recordTypes {
		in := D.PrepareDnsCall(dnsRecord)

		fmt.Printf(styles.NewStyles().Title.Render(
			"🛠️ %s Records 🛠️"), strings.ToUpper(dnsRecord),
		)
		fmt.Println()
		for _, ans := range in.Answer {
			fmt.Println(styles.NewStyles().Highlight.Render(ans.String()))

		}
	}
}

// PrepareDnsCall prepares a DNS call.
//
// Args:
//   - qtype: The query type.
//
// Returns:
//   - *dns.Msg: The DNS message.
func (D Domain) PrepareDnsCall(qtype string) *dns.Msg {
	qtypeToUpper := strings.ToUpper(qtype)
	record, ok := dns.StringToType[qtypeToUpper]
	if !ok {
		log.Fatalf(styles.NewStyles().Error.Render("Unknown record type: %d"), record)
	}

	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(D.domainName), record)
	m.RecursionDesired = true

	c := new(dns.Client)
	in, _, err := c.Exchange(m, "8.8.8.8:53")
	if err != nil {
		log.Fatalf(styles.NewStyles().Error.Render("An error occurred while resolving domain: %s"), err)
	}

	return in
}
