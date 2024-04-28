package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/willmeyers/stainless-cli/pkg/stainless"
	"github.com/willmeyers/stainless-cli/pkg/utils"
)

func runMembersCmd(cmd *cobra.Command, args []string) {
	credentials := cmd.Context().Value(utils.CredentialsCtxKey).(string)
	orgName := cmd.Context().Value(utils.OrgCtxKey).(string)
	client, err := stainless.New(stainless.WithCredentials(credentials))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error trying initialize stainless cli: %s", err)
		os.Exit(1)
	}
	members, err := client.ListMembers(orgName)
	if err != err {
		fmt.Fprintf(os.Stderr, "error listing members %s", err)
		return
	}

	out, err := utils.Prettify(members)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(out)
}

var membersCmd = &cobra.Command{
	Use:   "members",
	Short: "List members in org",
	Run:   runMembersCmd,
}
