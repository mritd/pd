package helper

import (
	"errors"
	"os"
	"os/exec"
	"sort"
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
	return strings.TrimSuffix(strings.TrimPrefix(s, "{"), "}")
}

func convert2Snapshots(infoMap map[string]SnapshotInfo) Snapshots {
	var ss Snapshots
	for id, info := range infoMap {
		ss = append(ss, Snapshot{
			ID: deleteBrackets(id),
			SnapshotInfo: SnapshotInfo{
				Name:    info.Name,
				Date:    info.Date,
				State:   info.State,
				Current: info.Current,
				Parent:  deleteBrackets(info.Parent),
			},
		})
	}
	sort.Sort(ss)
	return ss
}
