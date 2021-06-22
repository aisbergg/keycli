package cmd

import (
	"github.com/spf13/cobra"
)

var updateUsersCmd = &cobra.Command{
	Use:           "user USER",
	Short:         "Update information of a user",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		panic("not implemented")

		return nil
	},
}

func init() {
	updateCmd.AddCommand(updateUsersCmd)
	addUserInfoFlags(updateUsersCmd)
}
