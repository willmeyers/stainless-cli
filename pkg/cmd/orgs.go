package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/willmeyers/stainless-cli/pkg/stainless"
	"github.com/willmeyers/stainless-cli/pkg/utils"
)

func runOrgsCmd(cmd *cobra.Command, args []string) {
	credentials := cmd.Context().Value(utils.CredentialsCtxKey)
	client, err := stainless.New(stainless.WithCredentials(credentials.(string)))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error trying initialize stainless cli: %s", err)
		os.Exit(1)
	}

	orgs, err := client.ListOrgs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error listing orgs: %s", err)
		os.Exit(1)
	}

	out, err := utils.Prettify(orgs)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(out)
}

var orgsCmd = &cobra.Command{
	Use:   "orgs",
	Short: "Orgs",
	Long:  "TODO (willmeyers): implement long desc",
	Run:   runOrgsCmd,
}
