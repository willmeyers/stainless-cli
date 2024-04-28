package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/willmeyers/stainless-cli/pkg/stainless"
	"github.com/willmeyers/stainless-cli/pkg/utils"
)

func runSdksCmd(cmd *cobra.Command, args []string) {
	credentials := cmd.Context().Value(utils.CredentialsCtxKey).(string)
	orgName := cmd.Context().Value(utils.OrgCtxKey).(string)
	projectName := cmd.Context().Value(utils.ProjectCtxKey).(string)
	client, err := stainless.New(stainless.WithCredentials(credentials))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error trying initialize stainless cli: %s", err)
		os.Exit(1)
	}

	sdks, err := client.ListSdks(orgName, projectName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error listing sdks: %s", err)
		return
	}

	out, err := utils.Prettify(sdks)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(out)
}

var sdksCmd = &cobra.Command{
	Use:   "sdks",
	Short: "List SDKs in project",
	Run:   runSdksCmd,
}
