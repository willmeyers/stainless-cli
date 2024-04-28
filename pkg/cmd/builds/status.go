package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/willmeyers/stainless-cli/pkg/stainless"
	"github.com/willmeyers/stainless-cli/pkg/utils"
)

func init() {
	buildStatusCmd.Flags().StringP("language", "l", "", "Language name")
	buildStatusCmd.Flags().StringP("branch", "b", "main", "Branch name")
}

func runBuildStatusCmd(cmd *cobra.Command, args []string) {
	credentials := cmd.Context().Value(utils.CredentialsCtxKey).(string)
	orgName := cmd.Context().Value(utils.OrgCtxKey).(string)
	projectName := cmd.Context().Value(utils.ProjectCtxKey).(string)
	client, err := stainless.New(stainless.WithCredentials(credentials))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error trying initialize stainless cli: %s", err)
		os.Exit(1)
	}

	language, _ := cmd.Flags().GetString("language")
	branch, _ := cmd.Flags().GetString("branch")

	builds, err := client.GetBuildStatus(orgName, projectName, language, branch)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error listing builds: %s", err)
		return
	}

	out, err := utils.Prettify(builds)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(out)
}

var buildStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get build status of SDK",
	Run:   runBuildStatusCmd,
}
