package cli

import (
	"errors"
	"flag"
	"os"
	"stainless_cli/internal/stainless"
)

type Command interface {
	Parse(fs *flag.FlagSet) error
	Exec(stl *stainless.Stainless) error
}

func ExecCommand(args []string) error {
	if len(args) < 2 {
		return errors.New("please provide a command. run help for a list of commands")
	}

	cmdName := args[1]

	options := [](func(*stainless.Stainless) error){}
	if cmdName != "login" {
		options = append(options, stainless.WithAuthCookies(os.Getenv("STAINLESS_NEXTAUTH_SESSION_COOKIES")))
		options = append(options, stainless.WithDefaultOrgName())
		options = append(options, stainless.WithDefaultProjectName())
	}

	stl, err := stainless.New(options...)
	if err != nil {
		return err
	}

	commands := map[string]func(stl *stainless.Stainless) Command{
		"version":  NewVersionCommand,
		"help":     NewHelpCommand,
		"login":    NewLoginCommand,
		"orgs":     NewOrgsCommand,
		"projects": NewProjectsCommand,
		"generate": NewGenerateCommand,
		"builds":   NewListBuildsCommand,
		"sdks":     NewSdkCommand,
	}

	if newCommandFunc, ok := commands[cmdName]; ok {
		cmd := newCommandFunc(stl)
		fs := flag.NewFlagSet(cmdName, flag.ExitOnError)
		err := cmd.Parse(fs)
		if err != nil {
			return err
		}

		err = cmd.Exec(stl)
		if err != nil {
			return err
		}
	}

	return nil
}
