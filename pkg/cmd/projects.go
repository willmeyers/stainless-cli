package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/willmeyers/stainless-cli/pkg/stainless"
	"github.com/willmeyers/stainless-cli/pkg/utils"
)

func runProjectsCmd(cmd *cobra.Command, args []string) {
	credentials := cmd.Context().Value(utils.CredentialsCtxKey).(string)
	orgName := cmd.Context().Value(utils.OrgCtxKey).(string)
	client, err := stainless.New(stainless.WithCredentials(credentials))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error trying initialize stainless cli: %s", err)
		os.Exit(1)
	}

	projects, err := client.ListProjects(orgName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error listing projects: %s", err)
		return
	}

	out, err := utils.Prettify(projects)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(out)
}

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "List projects in org",
	Run:   runProjectsCmd,
}
