package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/willmeyers/stainless-cli/pkg/stainless"
	"github.com/willmeyers/stainless-cli/pkg/utils"
)

func runGenerateAllCmd(cmd *cobra.Command, args []string) {
	credentials := cmd.Context().Value(utils.CredentialsCtxKey).(string)
	orgName := cmd.Context().Value(utils.OrgCtxKey).(string)
	projectName := cmd.Context().Value(utils.ProjectCtxKey).(string)
	client, err := stainless.New(stainless.WithCredentials(credentials))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error trying initialize stainless cli: %s", err)
		os.Exit(1)
	}

	timeout, _ := cmd.Flags().GetInt("timeout")
	openAPISpec, _ := cmd.Flags().GetString("openapi")
	stainlessConfig, _ := cmd.Flags().GetString("config")
	outDir, _ := cmd.Flags().GetString("out-dir")

	generate, err := client.GenerateSdk(
		orgName,
		projectName,
		openAPISpec,
		stainlessConfig,
		outDir,
		"",
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error generating sdk: %s", err)
		return
	}

	out, err := utils.Prettify(generate)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(out)

	err = showGenerateBuildStatus(
		client,
		orgName,
		projectName,
		outDir,
		"",
		"main",
		timeout,
	)
	if err != nil {
		fmt.Println(err)
		return
	}
}

var generateAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Generate all SDKs",
	Run:   runGenerateAllCmd,
}
