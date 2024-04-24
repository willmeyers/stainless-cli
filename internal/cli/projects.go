package cli

import (
	"flag"
	"fmt"
	"os"
	"stainless_cli/internal/stainless"
)

type ProjectsCommand struct {
	OrgName string
}

func NewProjectsCommand(stl *stainless.Stainless) Command {
	return &ProjectsCommand{
		OrgName: stl.OrgName,
	}
}

func (cmd *ProjectsCommand) Parse(fs *flag.FlagSet) error {
	fs.StringVar(&cmd.OrgName, "org", cmd.OrgName, "Name of org")

	return fs.Parse(os.Args[2:])
}

func (cmd *ProjectsCommand) Exec(stl *stainless.Stainless) error {
	projects, err := stl.ListProjects(cmd.OrgName)
	if err != nil {
		return err
	}

	fmt.Println("Org\tName")
	for _, project := range projects.Items {
		fmt.Printf("%s\t%s\n", project.Org, project.Name)
	}

	return nil
}
