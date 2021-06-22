package cmd

import (
	"github.com/spf13/cobra"
)

// deleteCmd represents the base command for deleting various resources.
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete one or more resources",
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.PersistentFlags().StringP("session", "s", "keycloak", "Name of the session to use")
}
