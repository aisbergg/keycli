package cmd

import (
	"github.com/spf13/cobra"
)

var listUsersCmd = &cobra.Command{
	Use:           "users",
	Short:         "List all users",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		panic("not implemented")

		return nil
	},
}

func init() {
	listCmd.AddCommand(listUsersCmd)
	listUsersCmd.Flags().StringP("filter", "f", "", "Filter the results (e.g.: 'foo' in user.groups)")
	listUsersCmd.Flags().StringP("format", "m", "", "Output format for the results (e.g.: {{ user | json }})")
	listUsersCmd.Flags().String("sort", "", "Sort criteria (e.g.: user.created_at ~ user.name)")
	listUsersCmd.Flags().Bool("reverse", false, "Reverse the result output order")
}
