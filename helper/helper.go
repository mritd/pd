package helper

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

func prlctl(cmds ...string) (string, error) {
	cmd := exec.Command("prlctl", cmds...)

	bs, err := cmd.CombinedOutput()
	if err != nil {
		if bs != nil {
			return "", errors.New(strings.TrimSpace(string(bs)))
		}
		return "", err
	}

	return strings.TrimSpace(string(bs)), nil
}

func stdPrlctl(cmds ...string) error {
	cmd := exec.Command("prlctl", cmds...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func deleteBrackets(s string) string {
	return strings.TrimFunc(s, func(r rune) bool {
		return r == '{' || r == '}'
	})
}
