package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/willmeyers/stainless-cli/pkg/stainless"
	"github.com/willmeyers/stainless-cli/pkg/utils"
)

func runApiKeysCmd(cmd *cobra.Command, args []string) {
	credentials := cmd.Context().Value(utils.CredentialsCtxKey).(string)
	orgName := cmd.Context().Value(utils.OrgCtxKey).(string)
	client, err := stainless.New(stainless.WithCredentials(credentials))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error trying initialize stainless cli: %s", err)
		os.Exit(1)
	}
	apiKeys, err := client.ListApiKeys(orgName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error listing api keys %s", err)
		return
	}

	out, err := utils.Prettify(apiKeys)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(out)
}

var apiKeysCmd = &cobra.Command{
	Use:   "apikeys",
	Short: "Quickly get API keys",
	Run:   runApiKeysCmd,
}
