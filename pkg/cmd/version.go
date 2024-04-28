package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func runVersionCmd(cmd *cobra.Command, args []string) {
	fmt.Println(cmd.Root().Version)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Output the version number of this Stainless CLI",
	Run:   runVersionCmd,
}
