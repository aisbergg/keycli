package cmd

import (
	"github.com/spf13/cobra"
)

// removeCmd represents the base command for removing various resources from
// other ones.
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a resource from another one",
}

func init() {
	rootCmd.AddCommand(removeCmd)
	removeCmd.PersistentFlags().StringP("session", "s", "keycloak", "Name of the session to use")
}
