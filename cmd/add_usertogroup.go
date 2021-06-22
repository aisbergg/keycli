package cmd

import (
	"github.com/spf13/cobra"
)

var addUserToGroupCmd = &cobra.Command{
	Use:           "usertogroup",
	Short:         "Add one or more users to a group",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		panic("not implemented")

		return nil
	},
}

func init() {
	addCmd.AddCommand(addUserToGroupCmd)
}
