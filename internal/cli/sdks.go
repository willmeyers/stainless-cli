package cli

import (
	"flag"
	"fmt"
	"os"
	"stainless_cli/internal/stainless"
)

type SdkCommand struct {
	OrgName     string
	ProjectName string
	Language    string
	Branch      string
}

func NewSdkCommand(stl *stainless.Stainless) Command {
	return &SdkCommand{
		OrgName:     stl.OrgName,
		ProjectName: stl.ProjectName,
	}
}

func (cmd *SdkCommand) Parse(fs *flag.FlagSet) error {
	fs.StringVar(&cmd.OrgName, "org", cmd.OrgName, "Name of org")
	fs.StringVar(&cmd.ProjectName, "project", cmd.ProjectName, "Name of project")

	fs.StringVar(&cmd.Language, "language", "", "Name of language ")
	fs.StringVar(&cmd.Branch, "branch", "main", "Name of branch")

	return fs.Parse(os.Args[2:])
}

func (cmd *SdkCommand) Exec(stl *stainless.Stainless) error {
	if cmd.Language != "" && cmd.Branch != "" {
		sdkBuild, err := stl.RetrieveSdkStatus(cmd.OrgName, cmd.ProjectName, cmd.Language, cmd.Branch)
		if err != nil {
			return err
		}
		fmt.Println("ID\tStatus\tHas Generated\tStarted Generating At\t\t\tEnded At\t\t\t\tDiagnostics File URL")
		fmt.Printf("%d\t%s\t%v\t\t%s\t%s%s\n", sdkBuild.ID, sdkBuild.Status, sdkBuild.HasGenerated, sdkBuild.StartedGeneratingAt, sdkBuild.EndedAt, sdkBuild.DiagnosticsFileURL)

	}

	if cmd.Language == "" {
		sdks, err := stl.ListSdks(cmd.OrgName, cmd.ProjectName)
		if err != nil {
			return err
		}
		fmt.Println("ID\tOrg\t\tProject\t\t\tLanguage\tGit Repo")
		for _, sdk := range sdks.Items {
			fmt.Printf("%d\t%s\t%v\t%s\t\t%s\n", sdk.ID, sdk.Org, sdk.Project, sdk.Language, sdk.InternalRepositoryURL)
		}
	}

	return nil
}
