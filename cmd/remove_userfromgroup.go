package cmd

import (
	"github.com/spf13/cobra"
)

var removeUserFromGroupCmd = &cobra.Command{
	Use:           "userfromgroup GROUP USER...",
	Short:         "Remove one or more users from a group",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		panic("not implemented")

		return nil
	},
}

func init() {
	removeCmd.AddCommand(removeUserFromGroupCmd)
	removeUserFromGroupCmd.Flags().BoolP("ignore-error", "i", false, "Don't exit with an error, when an user is not assigned to the group")
}
