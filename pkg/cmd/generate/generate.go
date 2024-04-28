package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	GenerateCmd.PersistentFlags().StringP("openapi", "i", "", "Path to OpenAPI spec")
	GenerateCmd.PersistentFlags().StringP("config", "c", "", "Path to Stainless API config")
	GenerateCmd.PersistentFlags().StringP("out-dir", "o", "", "Output directory")
	GenerateCmd.PersistentFlags().IntP("timeout", "t", 30, "Timeout")
	GenerateCmd.AddCommand(generateAllCmd)
}

func runGenerateCmd(cmd *cobra.Command, args []string) {}

var GenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a new SDK from existing OpenAPI and Stainless specs",
	Long:  "Generate executes a new build given the parameters specified.",
	Run:   runGenerateCmd,
}
