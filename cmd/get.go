package cmd

import (
	"github.com/spf13/cobra"
)

// getCmd represents the base command for getting various resources.
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get one or more resources",
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.PersistentFlags().StringP("session", "s", "keycloak", "Name of the session to use")
}
