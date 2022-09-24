package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/ViBiOh/lucca/pkg/lucca"
	"github.com/spf13/cobra"
)

var (
	subdomain   string
	username    string
	password    string
	dryRun      bool
	luccaClient lucca.App
)

var rootCmd = &cobra.Command{
	Use:   "lucca",
	Short: "Run Lucca action fro the CLI",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if parent := cmd.Parent(); parent != nil && parent.Name() == "completion" {
			return
		}

		luccaClient, err = lucca.NewFromValues(subdomain, username, password)

		return err
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		principal, err := luccaClient.Principal(context.Background())
		if err != nil {
			return fmt.Errorf("principal: %w", err)
		}

		fmt.Printf("Hello %s\n", principal.FirstName)

		return nil
	},
}

func init() {
	flags := rootCmd.PersistentFlags()

	flags.StringVarP(&subdomain, "subdomain", "", "", "Subdomain")
	flags.StringVarP(&username, "username", "", "", "Username")
	flags.StringVarP(&password, "password", "", "", "Password")
	flags.BoolVarP(&dryRun, "dry-run", "", false, "Dry run")

	rootCmd.AddCommand(birthdaysCmd)
	rootCmd.AddCommand(companyBirthdaysCmd)

	rootCmd.AddCommand(leaveCmd)
	initLeave()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
}
