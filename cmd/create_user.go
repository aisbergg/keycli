package cmd

import (
	"github.com/spf13/cobra"
)

var createUserCmd = &cobra.Command{
	Use:           "user",
	Short:         "Create a user",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		panic("not implemented")

		return nil
	},
}

func init() {
	createCmd.AddCommand(createUserCmd)
	addUserInfoFlags(createUserCmd)
}

func addUserInfoFlags(command *cobra.Command) {
	command.Flags().StringP("firstname", "f", "", "First name")
	command.Flags().StringP("lastname", "l", "", "Last name")
	command.Flags().StringP("email", "e", "", "Email address")
	command.Flags().Bool("enabled", true, "Enable/Disable user")
	command.Flags().Bool("email-verified", false, "Email verification status")
	command.Flags().Bool("totp", false, "Enable/Disable TOTP ")

	command.Flags().StringSliceP("groups", "g", []string{}, "Associated groups as comma separated")
	command.Flags().StringSliceP("required-actions", "a", []string{}, "Required user actions as comma separated list")
	command.Flags().StringSliceP("realm-roles", "r", []string{}, "Realm roles as comma separated list")

	command.Flags().String("federation-link", "", "")
	command.Flags().String("service-account-client-id", "", "")
	command.Flags().StringArray("attribute", []string{}, "Attribute, can be specified multiple times (e.g.: key=value)")
	command.Flags().StringArray("access", []string{}, "???, can be specified multiple times (e.g.: key=value)")
	command.Flags().StringSlice("disableable-credential-types", []string{}, "Disableable credential types comma separated")
}

func parseUserInfoFlags() {

}
