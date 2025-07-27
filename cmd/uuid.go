package cmd

import (
	"commandCenter/styles"
	"commandCenter/validators"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var uuidCmd = &cobra.Command{
	Use:   "uuid",
	Short: "🎲 Generate new UUID with supported types: uuid4, uuid6, uuid7, clock",
	Example: `
      # Generate a default UUID (uuid4)
      ops uuid

      # Generate a UUID version 6
      ops uuid -t uuid6

      # Generate a UUID version 7
      ops uuid --type=uuid7

      # Generate a clock-based UUID
      ops uuid -t clock

      # Show help for the uuid command
      ops uuid --help
    `,

	Run: generateUUID,
}

// init initializes the uuid command and its flags.
//
// Args:
//   - None
//
// Returns:
//   - None
func init() {
	rootCmd.AddCommand(uuidCmd)

	uuidCmd.Flags().StringP("type", "t", "uuid4", "Generate new UUID")
}

// createUUID allows the user to create a UUID4 identifier
//
// Args:
//   - None
//
// Returns:
//   - None
func createUUID() {
	log.Println(styles.StyliseMessage("UUID❹ : "+uuid.NewString(), styles.FormatStyle.Highlight))
}

// createClockUUID allows the user to create a clock sequence
//
// Params:
// - None
//
// Returns:
// - int: the clock sequence UUID
func createClockUUID() {
	log.Println(styles.StyliseMessage("UUID🕕: "+string(uuid.ClockSequence()), styles.FormatStyle.Highlight))
}

// validateUUID validates a UUID string.
//
// Args:
//   - uuidString: The UUID string to validate.
//
// Returns:
//   - error: An error if the UUID is invalid, otherwise nil.
func validateUUID(uuidString string) error {

	return uuid.Validate(uuidString)
}

// createv6UUID allows the user to create a UUID version 6
//
// Params:
// - None
//
// Returns:
// - string: the version 6 UUID
func createv6UUID() {
	if newUUID, err := uuid.NewV6(); err != nil {
		log.Fatalf(styles.StyliseMessage(fmt.Sprintf("Could not generate UUID: %s", err), styles.FormatStyle.Error))
	} else {
		log.Println(styles.StyliseMessage("UUID❻:: "+newUUID.String(), styles.FormatStyle.Highlight))
	}
}

// createv7UUID allows the user to create a UUID version 7
//
// Params:
// - None
//
// Returns:
// - string: the version 7 UUID
func createv7UUID() {
	if newUUID, err := uuid.NewV6(); err != nil {
		log.Fatalf(styles.StyliseMessage(fmt.Sprintf("Could not generate UUID: %s", err), styles.FormatStyle.Error))
	} else {
		log.Println(styles.StyliseMessage("UUID❼:: "+newUUID.String(), styles.FormatStyle.Highlight))
	}
}

func generateUUID(cmd *cobra.Command, args []string) {
	uuid, err := validators.VerifyStringInputs(cmd, "type")
	if err != nil {
		log.Fatalln(err)
	}

	switch uuid {
	case "uuid4":
		createUUID()
	case "uuid6":
		createv6UUID()
	case "uuid7":
		createv7UUID()
	case "clock":
		createClockUUID()
	default:
		log.Println(styles.StyliseMessage("🛑 UUID types can be uuid4, uuid6, uuid7, clock", styles.FormatStyle.Error))
	}
}
