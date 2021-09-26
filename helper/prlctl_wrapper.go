package helper

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
)

type VMInfo struct {
	Name         string `json:"name"`
	Status       string `json:"status"`
	IPConfigured string `json:"ip_configured"`
	UUID         string `json:"uuid"`
}

type SnapshotInfo struct {
	Name    string `json:"name"`
	Date    string `json:"date"`
	State   string `json:"state"`
	Current bool   `json:"current"`
	Parent  string `json:"parent"`
}

func ListVMs(all bool) {
	var data string
	var err error
	if all {
		data, err = prlctl("list", "-f", "-j", "-a")
	} else {
		data, err = prlctl("list", "-f", "-j")
	}
	if err != nil {
		logrus.Fatal(err)
	}

	var vms []VMInfo
	err = jsoniter.Unmarshal([]byte(data), &vms)
	if err != nil {
		logrus.Fatal(err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")

	table.SetHeader([]string{"NAME", "STATUS", "IP", "UUID"})
	for _, vm := range vms {
		table.Append([]string{vm.Name, vm.Status, vm.IPConfigured, vm.UUID})
	}
	table.Render()
}

func StartVM(vms []string) {
	for _, vm := range vms {
		err := stdPrlctl("start", vm)
		if err != nil {
			logrus.Fatal(err)
		}
	}
}

func StopVM(vms []string, force bool) {
	var err error
	for _, vm := range vms {
		if force {
			err = stdPrlctl("stop", vm, "--kill")
		} else {
			err = stdPrlctl("stop", vm)
		}
		if err != nil {
			logrus.Fatal(err)
		}
	}
}

func CreateSnapshot(vms []string, name string) {
	var err error
	for _, vm := range vms {
		logrus.Infof("Create [%s] snapshot(%s)...", name,vm, )
		err = stdPrlctl("snapshot", vm, "-n", name)
		if err != nil {
			logrus.Fatal(err)
		}
	}
}

func ListSnapshot(vm string) {
	data, err := prlctl("snapshot-list", vm, "-j")
	if err != nil {
		logrus.Fatal(err)
	}

	var sps map[string]SnapshotInfo
	err = jsoniter.Unmarshal([]byte(data), &sps)
	if err != nil {
		logrus.Fatal(err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")

	table.SetHeader([]string{"NAME", "CURRENT", "UUID", "STATE", "DATE"})
	for id, sp := range sps {
		table.Append([]string{sp.Name, strconv.FormatBool(sp.Current), deleteBrackets(id), sp.State, sp.Date})
	}
	table.Render()
}

func DeleteSnapshot(vms []string, name string) {
	for _, vm := range vms {
		logrus.Infof("Delete [%s] snapshot(%s)...", name, vm)

		data, err := prlctl("snapshot-list", vm, "-j")
		if err != nil {
			logrus.Fatal(err)
		}

		var sps map[string]SnapshotInfo
		err = jsoniter.Unmarshal([]byte(data), &sps)
		if err != nil {
			logrus.Fatal(err)
		}
		var spID string
		for id, sp := range sps {
			if sp.Name == name {
				spID = id
				break
			}
		}
		if spID == "" {
			logrus.Warnf("VM [%s] snapshot(%s) not found.", vm,name)
			continue
		}

		err = stdPrlctl("snapshot-delete", vm, "-i", spID)
		if err != nil {
			logrus.Fatal(err)
		}
	}
}

func SwitchSnapshot(vms []string, name string) {
	for _, vm := range vms {
		logrus.Infof("Swicth [%s] to snapshot %s...", name, vm)

		data, err := prlctl("snapshot-list", vm, "-j")
		if err != nil {
			logrus.Fatal(err)
		}

		var sps map[string]SnapshotInfo
		err = jsoniter.Unmarshal([]byte(data), &sps)
		if err != nil {
			logrus.Fatal(err)
		}
		var spID string
		for id, sp := range sps {
			if sp.Name == name {
				spID = id
				break
			}
		}
		if spID == "" {
			logrus.Fatalf("VM [%s] snapshot not found.", vm)
		}

		err = stdPrlctl("snapshot-switch", vm, "-i", spID)
		if err != nil {
			logrus.Fatal(err)
		}
	}
}
