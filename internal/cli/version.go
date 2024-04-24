package cli

import (
	"flag"
	"fmt"
	"stainless_cli/internal/stainless"
)

func NewVersionCommand(stl *stainless.Stainless) Command {
	return &VersionCommand{}
}

type VersionCommand struct{}

func (c *VersionCommand) Parse(fs *flag.FlagSet) error {
	return nil
}

func (c *VersionCommand) Exec(stl *stainless.Stainless) error {
	fmt.Println("Stainless CLI version 0.0.1")
	return nil
}
