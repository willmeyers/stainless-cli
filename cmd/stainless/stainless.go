package main

import (
	"log"
	"os"
	"stainless_cli/internal/cli"
)

func main() {
	err := cli.ExecCommand(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
