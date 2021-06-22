package cmd

import (
	"github.com/spf13/cobra"
)

// addCmd represents the base command for updating various resources.
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update information of a resource",
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.PersistentFlags().StringP("session", "s", "keycloak", "Name of the session to use")
}
