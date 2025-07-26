package validators

import (
	"commandCenter/styles"
	"fmt"
	"github.com/spf13/cobra"
)

func VerifyStringInputs(cmd *cobra.Command, flag string) (string, error) {
	passedFlag, err := cmd.Flags().GetString(flag)
	if err != nil {
		message := fmt.Errorf(styles.NewStyles().Error.Render("An error occurred while parsing flag '%s'.\nError: %s"), flag, err)

		return passedFlag, message
	}

	return passedFlag, nil
}

func VerifyBoolInputs(cmd *cobra.Command, flag string) (bool, error) {
	passedFlag, err := cmd.Flags().GetBool(flag)
	if err != nil {
		message := fmt.Errorf(styles.NewStyles().Error.Render("An error occurred while parsing flag '%s'.\nError: %s"), flag, err)

		return passedFlag, message
	}

	return passedFlag, nil
}
