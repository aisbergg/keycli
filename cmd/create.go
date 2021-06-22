package cmd

import (
	"github.com/spf13/cobra"
)

// createCmd represents the base command for creating various resources.
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a resource",
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.PersistentFlags().StringP("session", "s", "keycloak", "Name of the session to use")
}
