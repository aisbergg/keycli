package cmd

import (
	"fmt"
	"os"

	keycli "github.com/aisbergg/keycli/pkg"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "keycli",
	Short: "A command line interface for Keycloak",
	Long:  `Keycli is a CLI program for Keycloak. It makes use of Keycloaks built-in REST API to query information and execute commands over HTTPS. This allows the tool to list/add/update/delete users/groups/clients/roles conveniently from the command line. Thus it can be used as a replacement for the clumsy web-ui for the most common management operations.`,
	Run: func(cmd *cobra.Command, args []string) {
		//
		// parse flags and args
		//
		printVersion, _ := cmd.Flags().GetBool("version")
		if printVersion {
			fmt.Println(keycli.Version)
			return
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	cmd, err := rootCmd.ExecuteC()
	if err != nil {
		debug, _ := cmd.Flags().GetBool("debug")
		errMsg := formatError(err, debug)
		fmt.Fprintln(os.Stderr, errMsg)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().Bool("version", false, "Print program version and quit")
	rootCmd.PersistentFlags().Bool("debug", false, "Turn on debug mode (verbose output and stack traces)")
}

func formatError(err error, debug bool) string {
	if debug {
		// return full stack trace
		err = errors.WithStack(err)
		return fmt.Sprintf("%+v", err)
	}

	return err.Error()
}
