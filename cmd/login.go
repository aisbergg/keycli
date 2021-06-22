package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/aisbergg/keycli/pkg/interface/cli"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login [SESSION]",
	Short: "Login to Keycloak and create a new session",
	Long: `Login to Keycloak and create a new session.

The session is tied to a specific server and realm. Multiple sessions for
different servers and realms can be opened and named using the SESSION argument.
Other commands accept the session name as an option, which will effectively
execute the commands in the context of the given session.`,
	Example: `  # Ask for url, user and password and then login
  login

  # Create a session using given url, realm (foo), user (bar) and password
  login -l 'https://sso.example.org' -r foo -u bar -p secret

  # Create a session and name it 'baz'
  login baz`,
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

		url, _ := cmd.Flags().GetString("url")
		if url == "" {
			fmt.Printf("Keycloak URL: ")
			fmt.Scanln(&url)
		}
		url = strings.TrimSpace(url)
		if !strings.HasSuffix(url, "/") {
			url = url + "/"
		}

		realm, _ := cmd.Flags().GetString("realm")
		realm = strings.TrimSpace(realm)
		if realm == "" {
			return errors.New("realm must not be empty")
		}

		secretKey, _ := cmd.Flags().GetString("secret-key")
		secretKey = strings.TrimSpace(secretKey)

		user, _ := cmd.Flags().GetString("user")
		user = strings.TrimSpace(user)
		password, _ := cmd.Flags().GetString("password")
		if secretKey == "" {
			if user == "" {
				fmt.Printf("Keycloak Admin User: ")
				fmt.Scanln(&user)
			}
			if password == "" {
				fmt.Printf("Keycloak Admin Password: ")
				p, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
				fmt.Println()
				password = string(p)
			}
		}

		skipVerify, _ := cmd.Flags().GetBool("skip-verify")

		clientID, _ := cmd.Flags().GetString("client-id")
		clientID = strings.TrimSpace(clientID)

		//
		// perform login
		//
		return cli.Login(name, url, realm, clientID, secretKey, user, password, skipVerify)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringP("url", "l", "", "Base URL of Keycloak server (e.g.: https://sso.example.org/)")
	loginCmd.Flags().StringP("realm", "r", "master", "Keycloak realm")
	loginCmd.Flags().StringP("user", "u", "", "Keycloak admin user")
	loginCmd.Flags().StringP("password", "p", "", "Keycloak admin user password")
	loginCmd.Flags().StringP("secret-key", "s", "", "Keycloak admin secret key")
	loginCmd.Flags().Bool("skip-verify", false, "Skip TLS certificate verification")
	loginCmd.Flags().String("client-id", "admin-cli", "Client ID to be used")
}
