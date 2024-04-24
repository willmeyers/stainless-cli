package cli

import (
	"flag"
	"fmt"
	"stainless_cli/internal/stainless"
)

func NewHelpCommand(stl *stainless.Stainless) Command {
	return &HelpCommand{}
}

type HelpCommand struct{}

func (c *HelpCommand) Parse(fs *flag.FlagSet) error {
	return nil
}

func (c *HelpCommand) Exec(stl *stainless.Stainless) error {
	fmt.Println(`Usage: stainless [command]

Commands:
  login     Log in to your account
  orgs      List organizations
  projects  List projects
  generate  Generate SDKs
  builds    List builds
  sdks      List SDKs and SDK build status
  version   Show the version of the CLI
  help      Show this help message`)
	return nil
}
