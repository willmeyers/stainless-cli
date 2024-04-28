package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/willmeyers/stainless-cli/pkg/stainless"
)

type Spinner struct {
	s *spinner.Spinner
}

func NewSpinner(suffix string) *Spinner {
	s := spinner.New(spinner.CharSets[0], 100*time.Millisecond)
	s.Suffix = " " + suffix
	s.Start()
	return &Spinner{s: s}
}

func (sp *Spinner) Stop() {
	sp.s.Stop()
}

func pollSDKBuildStatusUtilSuccess(client *stainless.Client, orgName, projectName, language, branch string, timeout int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Generation has exceeded set timeout. SDKs are still queued to be generated.")
			return fmt.Errorf("polling timed out")
		default:
			build, err := client.GetBuildStatus(orgName, projectName, language, branch)
			if err != nil {
				return err
			}
			if build.Status == "success" {
				return nil
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func showGenerateBuildStatus(
	client *stainless.Client,
	orgName string,
	projectName string,
	outDir string,
	language string,
	branch string,
	timeout int,
) error {
	errorChan := make(chan error)
	var wg sync.WaitGroup

	languages := []string{}
	sdks, err := client.ListSdks(orgName, projectName)
	if err != nil {
		return err
	}
	if language != "" {
		languages = append(languages, language)
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

	spin := NewSpinner("generating sdks...")

	for _, lang := range languages {
		wg.Add(1)
		go func(language string) {
			defer wg.Done()
			err := pollSDKBuildStatusUtilSuccess(client, orgName, projectName, language, branch, timeout)
			if err != nil {
				errorChan <- err
			}
		}(lang)
	}
	wg.Wait()
	fmt.Println("done")
	close(errorChan)
	spin.Stop()

	for err := range errorChan {
		return err
	}

	if outDir != "" {
		for _, sdk := range sdks.Items {
			if slices.Contains(languages, sdk.Language) {
				repoDir := filepath.Join(outDir, sdk.Language, sdk.ReleaseInformation.Repo)
				err := cloneOrPullFromGitRepo(repoDir, sdk.InternalRepositoryURL)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func cloneOrPullFromGitRepo(repoDir, repoName string) error {
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
