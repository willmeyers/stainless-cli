package cli

import (
	"flag"
	"fmt"
	"os"
	"stainless_cli/internal/stainless"
)

type OrgsCommand struct{}

func NewOrgsCommand(stl *stainless.Stainless) Command {
	return &OrgsCommand{}
}

func (cmd *OrgsCommand) Parse(fs *flag.FlagSet) error {
	return fs.Parse(os.Args[2:])
}

func (cmd *OrgsCommand) Exec(stl *stainless.Stainless) error {
	orgs, err := stl.ListOrgs()
	if err != nil {
		return err
	}

	fmt.Println("Display Name\tName")
	for _, org := range orgs.Items {
		fmt.Printf("%s\t%s\n", org.DisplayName, org.Name)
	}

	return nil
}
