package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/willmeyers/stainless-cli/pkg/utils"
)

func runLogoutCmd(cmd *cobra.Command, args []string) {
	y := utils.AskYesNo("Are you sure you want to delete all cached credentials? You will have to login again.")
	if y {
		viper.Set("credentials", "")
		err := viper.WriteConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error trying to write config: %s", err)
			os.Exit(1)
		}

	}
	os.Exit(0)
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Delete any cached authentication credentials.",
	Run:   runLogoutCmd,
}
