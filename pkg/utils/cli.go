package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func AskYesNo(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/N]: ", prompt)
		yN, err := reader.ReadString('\n')
		if err != nil {
			panic(1)
		}

		yN = strings.ToLower(strings.TrimSpace(yN))

		if yN == "y" || yN == "yes" {
			return true
		} else if yN == "n" || yN == "no" {
			return false
		} else {
			return false
		}
	}
}

func Prettify(data interface{}) (string, error) {
	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(val), nil
}
