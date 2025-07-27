package validators

import (
	"commandCenter/styles"
	"fmt"
	"github.com/spf13/cobra"
)

// VerifyStringInputs verifies and returns a string flag from the cobra command.
//
// Args:
//   - cmd: The cobra command.
//   - flag: The name of the string flag to verify.
//
// Returns:
//   - string: The value of the string flag.
//   - error: An error if the flag is not found or cannot be parsed.
func VerifyStringInputs(cmd *cobra.Command, flag string) (string, error) {
	passedFlag, err := cmd.Flags().GetString(flag)
	if err != nil {
		message := fmt.Errorf(styles.NewStyles().Error.Render("An error occurred while parsing flag '%s'.\nError: %s"), flag, err)

		return passedFlag, message
	}

	return passedFlag, nil
}

// VerifyBoolInputs verifies and returns a boolean flag from the cobra command.
//
// Args:
//   - cmd: The cobra command.
//   - flag: The name of the boolean flag to verify.
//
// Returns:
//   - bool: The value of the boolean flag.
//   - error: An error if the flag is not found or cannot be parsed.
func VerifyBoolInputs(cmd *cobra.Command, flag string) (bool, error) {
	passedFlag, err := cmd.Flags().GetBool(flag)
	if err != nil {
		message := fmt.Errorf(styles.NewStyles().Error.Render("An error occurred while parsing flag '%s'.\nError: %s"), flag, err)

		return passedFlag, message
	}

	return passedFlag, nil
}
