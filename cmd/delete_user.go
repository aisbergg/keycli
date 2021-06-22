package cmd

import (
	"github.com/spf13/cobra"
)

var deleteUsersCmd = &cobra.Command{
	Use:   "users",
	Short: "Delete one or more users",
	// TODO: long description
	Long:          ``,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		panic("not implemented")

		return nil
	},
}

func init() {
	deleteCmd.AddCommand(deleteUsersCmd)
	deleteUsersCmd.Flags().BoolP("ignore-error", "i", false, "Don't exit with an error, when a user doesn't exist")
}
