package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"stainless_cli/internal/stainless"
	"stainless_cli/internal/utils"
	"strconv"
	"sync"
	"time"
)

type GenerateCommand struct {
	OrgName     string
	ProjectName string
	OpenAPI     string
	Config      string
	OutDir      string
	Language    string
}

func NewGenerateCommand(stl *stainless.Stainless) Command {
	return &GenerateCommand{
		OrgName:     stl.OrgName,
		ProjectName: stl.ProjectName,
	}
}

func (cmd *GenerateCommand) Parse(fs *flag.FlagSet) error {
	fs.StringVar(&cmd.OrgName, "org", cmd.OrgName, "Name of org")
	fs.StringVar(&cmd.ProjectName, "project", cmd.ProjectName, "Name of project")

	fs.StringVar(&cmd.OpenAPI, "openapi", "", "Path to OpenAPI spec")
	fs.StringVar(&cmd.Config, "config", "", "Path to Stainless API config")
	fs.StringVar(&cmd.OutDir, "out-dir", "./sdks", "Output directory")
	fs.StringVar(&cmd.Language, "language", "", "Language")

	return fs.Parse(os.Args[2:])
}

func (cmd *GenerateCommand) Exec(stl *stainless.Stainless) error {
	_, err := stl.Generate(cmd.OrgName, cmd.ProjectName, cmd.OpenAPI, cmd.Config, cmd.OutDir, cmd.Language)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	errorChan := make(chan error)
	var wg sync.WaitGroup

	languages := []string{}
	sdks, err := stl.ListSdks(cmd.OrgName, cmd.ProjectName)
	if err != nil {
		return err
	}
	if cmd.Language != "" {
		languages = append(languages, cmd.Language)
		validSdkLanguages := []string{}
		for _, sdk := range sdks.Items {
			validSdkLanguages = append(validSdkLanguages, sdk.Language)
		}

		for _, lang := range languages {
			if !slices.Contains(validSdkLanguages, lang) {
				return fmt.Errorf("cannot build sdk for language %s. does not exist in project", lang)
			}
		}
	} else {
		for _, sdk := range sdks.Items {
			languages = append(languages, sdk.Language)
		}
	}

	var spin *utils.Spinner
	if len(languages) == 1 {
		spin = utils.NewSpinner("generating " + languages[0] + " sdk...")
	} else {
		spin = utils.NewSpinner("generating " + strconv.Itoa(len(languages)) + " sdks...")
	}

	for _, lang := range languages {
		wg.Add(1)
		go func(language string) {
			defer wg.Done()
			err := checkBuildStatus(stl, cmd.OrgName, cmd.ProjectName, language)
			if err != nil {
				errorChan <- err
			}
		}(lang)
	}
	wg.Wait()
	fmt.Println("✅ done")
	close(errorChan)
	spin.Stop()

	for err := range errorChan {
		return err
	}

	if cmd.OutDir != "" {
		fmt.Println("❗️ ensure you have accepted the invitation to be a contributor to your sdk repos")
		for _, sdk := range sdks.Items {
			if slices.Contains(languages, sdk.Language) {
				repoDir := filepath.Join(cmd.OutDir, sdk.Language, sdk.ReleaseInformation.Repo)
				err := cloneOrPullFromGitHubRepo(repoDir, sdk.InternalRepositoryURL)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func checkBuildStatus(stl *stainless.Stainless, orgName, projectName, language string) error {
	// TODO (willmeyers): add timeout here to avoid infinite looping
	for {
		build, err := stl.RetrieveSdkStatus(orgName, projectName, language, "main")
		if err != nil {
			return err
		}
		if build.Status == "success" {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
}

func cloneOrPullFromGitHubRepo(repoDir, repoName string) error {
	if _, err := os.Stat(repoDir); os.IsNotExist(err) {
		err := os.MkdirAll(repoDir, 0755)
		if err != nil {
			return err
		}
	}

	cmd := exec.Command("git", "clone", repoName, repoDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%s\n", output)
		// If clone fails, try to pull
		cmd = exec.Command("git", "-C", repoDir, "pull")
		output, err = cmd.CombinedOutput()
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", output)
	}

	return nil
}
