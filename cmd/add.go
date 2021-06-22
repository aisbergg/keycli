package cmd

import (
	"github.com/spf13/cobra"
)

// addCmd represents the base command for adding various resources to other ones.
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a resource to another one",
	// TODO: long description
	Long: ``,
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.PersistentFlags().StringP("session", "s", "keycloak", "Name of the session to use")
}
