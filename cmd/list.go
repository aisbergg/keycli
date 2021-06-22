package cmd

import (
	"github.com/spf13/cobra"
)

// listCmd represents the base command for listing various resources.
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List resource items",
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.PersistentFlags().StringP("session", "s", "keycloak", "Name of the session to use")
}
