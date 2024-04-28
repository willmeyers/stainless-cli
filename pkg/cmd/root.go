package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	builds "github.com/willmeyers/stainless-cli/pkg/cmd/builds"
	generate "github.com/willmeyers/stainless-cli/pkg/cmd/generate"
	stainless "github.com/willmeyers/stainless-cli/pkg/stainless"
	"github.com/willmeyers/stainless-cli/pkg/utils"
)

var ConfigFile string
var Credentials string

var rootCmd = &cobra.Command{
	Use:           "stainless",
	Short:         "CLI interface for Stainless API",
	Long:          "An unofficial command-line tool to interact with Stainless API.",
	SilenceUsage:  true,
	SilenceErrors: true,
	Annotations: map[string]string{
		"help":     "help",
		"login":    "setup",
		"logout":   "setup",
		"orgs":     "web",
		"projects": "web",
		"apikeys":  "web",
		"members":  "web",
		"sdks":     "web",
		"builds":   "web",
		"generate": "web",
	},
	Version: "1.0.0",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		populateCmdContextPreRun(cmd, args)
	},
}

func Execute() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initializeConfig)

	rootCmd.PersistentFlags().String("org", "", "Manually set org name")
	rootCmd.PersistentFlags().String("project", "", "Manually set project name")
	rootCmd.PersistentFlags().StringVarP(&Credentials, "auth-credentials", "a", "", "Manaully set authentication cookies used to authenticate requests to Stainless API.")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(generate.GenerateCmd)
	rootCmd.AddCommand(orgsCmd)
	rootCmd.AddCommand(projectsCmd)
	rootCmd.AddCommand(apiKeysCmd)
	rootCmd.AddCommand(membersCmd)
	rootCmd.AddCommand(builds.BuildsCmd)
	rootCmd.AddCommand(sdksCmd)

	viper.SetDefault("credentials", Credentials)
}

func initializeConfig() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}

	viper.SetConfigName(".stainless-cli")
	viper.SetConfigType("json")
	viper.AddConfigPath(home)

	if ConfigFile != "" {
		viper.SetConfigFile(ConfigFile)
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.WriteConfig()
			viper.SafeWriteConfig()
		} else {
			fmt.Fprintf(os.Stderr, "error trying to read in config: %s\n", err)
			os.Exit(1)
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "error trying to read in config: %s\n", err)
		os.Exit(1)
	}
}

func populateCmdContextPreRun(cmd *cobra.Command, args []string) {
	annotation := cmd.Root().Annotations[cmd.Name()]
	if annotation == "help" || annotation == "setup" {
		return
	}

	credentials, err := cmd.Flags().GetString("a")
	if err != nil {
		cachedCookies := viper.Get("credentials")
		if cachedCookies != nil {
			credentials = cachedCookies.(string)
		}
	}

	ctx := context.WithValue(cmd.Context(), utils.CredentialsCtxKey, credentials)
	cmd.SetContext(ctx)

	client, err := stainless.New(stainless.WithCredentials(credentials))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error trying to init stainless cli: %s\n", err)
		os.Exit(1)
	}

	var defaultOrgName string
	orgName, _ := cmd.Flags().GetString("org")
	if orgName == "" {
		retrievedOrgName, err := client.GetDefaultOrg()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error trying to get default org: %s\n", err)
			os.Exit(1)
		}
		defaultOrgName = retrievedOrgName
	} else {
		defaultOrgName = orgName
	}
	ctx = context.WithValue(cmd.Context(), utils.OrgCtxKey, defaultOrgName)
	cmd.SetContext(ctx)

	var defaultProjectName string
	projectName, _ := cmd.Flags().GetString("project")
	if projectName == "" {
		retrievedProjectName, err := client.GetDefaultProject(defaultOrgName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error trying to get default org: %s\n", err)
			os.Exit(1)
		}
		defaultProjectName = retrievedProjectName
	} else {
		defaultProjectName = projectName
	}
	ctx = context.WithValue(cmd.Context(), utils.ProjectCtxKey, defaultProjectName)
	cmd.SetContext(ctx)
}
