package cmd

import (
	"github.com/spf13/cobra"
)

var getUserCmd = &cobra.Command{
	Use:           "user USER...",
	Short:         "Get information for one or more users",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		panic("not implemented")
		// sessionName, _ := cmd.Flags().GetString("session")
		// sessionName = strings.TrimSpace(sessionName)
		// if sessionName == "" {
		// 	sessionName = "keycloak"
		// }

		// format, _ := cmd.Flags().GetString("format")

		return nil
	},
}

func init() {
	getCmd.AddCommand(getUserCmd)
	getUserCmd.Flags().StringP("format", "m", "", "Output format for the results (e.g.: {{ user | json }})")
}
