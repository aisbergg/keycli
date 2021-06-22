package cmd

import (
	"strings"

	"github.com/aisbergg/keycli/pkg/interface/cli"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout [SESSION]",
	Short: "Logout of Keycloak and thereby ending the session",
	Example: `  # End the default unnamed session
  logout

  # End the named session baz
  logout baz

  # Delete session keys, even if the session couldn't be successfully ended
  logout -f`,
	Args:          cobra.MaximumNArgs(1),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		//
		// parse flags and args
		//
		name := "keycloak"
		if len(args) > 0 {
			name = args[0]
		}
		name = strings.TrimSpace(name)

		force, _ := cmd.Flags().GetBool("force")

		//
		// perform logout
		//
		return cli.Logout(name, force)
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
	logoutCmd.Flags().BoolP("force", "f", false, "Force deletion of tokens, even if the session couldn't be terminated")
}
