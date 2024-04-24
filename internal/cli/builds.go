package cli

import (
	"flag"
	"fmt"
	"os"
	"stainless_cli/internal/stainless"
	"strconv"
)

type ListBuildsCommand struct {
	OrgName     string
	ProjectName string
}

func NewListBuildsCommand(stl *stainless.Stainless) Command {
	return &ListBuildsCommand{
		OrgName:     stl.OrgName,
		ProjectName: stl.ProjectName,
	}
}

func (cmd *ListBuildsCommand) Parse(fs *flag.FlagSet) error {
	fs.StringVar(&cmd.OrgName, "org", cmd.OrgName, "Name of org")
	fs.StringVar(&cmd.ProjectName, "project", cmd.ProjectName, "Name of project")
	return fs.Parse(os.Args[2:])
}

func (cmd *ListBuildsCommand) Exec(stl *stainless.Stainless) error {
	fmt.Println(cmd)
	builds, err := stl.ListBuilds(cmd.OrgName, cmd.ProjectName)
	if err != nil {
		return err
	}

	for _, build := range builds.Builds {
		sdkString := ""
		for i, sdk := range build.Sdks {
			sdkString += fmt.Sprintf("%s(%s)", sdk.Language, sdk.Status)
			if i != len(build.Sdks)-1 {
				sdkString += " "
			}
		}
		fmt.Printf("%s\t%s\t%s\t%s\n", strconv.Itoa(build.ID), build.Org, build.Project, sdkString)
	}

	return nil
}
