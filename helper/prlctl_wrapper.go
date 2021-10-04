package helper

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
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

func ListVMInfo(all bool) ([]VMInfo, error) {
	var data string
	var err error
	if all {
		data, err = prlctl("list", "-f", "-j", "-a")
	} else {
		data, err = prlctl("list", "-f", "-j")
	}
	if err != nil {
		return nil, err
	}

	var vms []VMInfo
	err = jsoniter.Unmarshal([]byte(data), &vms)
	if err != nil {
		return nil, err
	}
	return vms, nil
}

func ListVM(all bool) {
	vms, err := ListVMInfo(all)
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
	var wg sync.WaitGroup
	wg.Add(len(vms))

	for _, vm := range vms {
		go func(vm string) {
			defer wg.Done()
			logrus.Infof("Starting vm %s...", vm)
			if vm == "docker" {
				StartDockerVM(vm, false)
			} else {
				if err := stdPrlctl("start", vm); err != nil {
					logrus.Errorf("Failed to start vm [%s]: %v", vm, err)
				}
			}
		}(vm)
	}

	wg.Wait()
}

func StopVM(vms []string, force bool) {
	var wg sync.WaitGroup
	wg.Add(len(vms))

	for _, vm := range vms {
		go func(vm string, force bool) {
			defer wg.Done()
			logrus.Infof("Stopping VM %s...", vm)

			var err error
			if force {
				err = stdPrlctl("stop", vm, "--kill")
			} else {
				err = stdPrlctl("stop", vm)
			}
			if err != nil {
				logrus.Errorf("Failed to stop VM [%s]: %v", vm, err)
			}
		}(vm, force)
	}

	wg.Wait()
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

func CreateSnapshot(vms []string, name string) {
	var wg sync.WaitGroup
	wg.Add(len(vms))
	for _, vm := range vms {
		go func(vm, name string) {
			defer wg.Done()
			logrus.Infof("Create [%s] snapshot(%s)...", name, vm)

			if err := stdPrlctl("snapshot", vm, "-n", name); err != nil {
				logrus.Errorf("VM [%s] create snapshot failed: %v", vm, err)
			}
		}(vm, name)
	}
	wg.Wait()
}

func DeleteSnapshot(vms []string, name string) {
	var wg sync.WaitGroup
	wg.Add(len(vms))

	for _, vm := range vms {
		go func(vm, name string) {
			defer wg.Done()
			logrus.Infof("Delete %s snapshot --> %s...", name, vm)

			data, err := prlctl("snapshot-list", vm, "-j")
			if err != nil {
				logrus.Errorf("Failed to get VM %s snapshots: %v", vm, err)
				return
			}

			var sps map[string]SnapshotInfo
			err = jsoniter.Unmarshal([]byte(data), &sps)
			if err != nil {
				logrus.Errorf("Get VM %s snapshot info failed: %v", vm, err)
				return
			}
			var spID string
			for id, sp := range sps {
				if sp.Name == name {
					spID = id
					break
				}
			}
			if spID == "" {
				logrus.Warnf("VM %s snapshot(%s) not found.", vm, name)
				return
			}

			err = stdPrlctl("snapshot-delete", vm, "-i", spID)
			if err != nil {
				logrus.Errorf("Failed to delete VM %s snapshot: %v", vm, err)
			}
		}(vm, name)
	}

	wg.Wait()
}

func SwitchSnapshot(vms []string, name string) {
	var wg sync.WaitGroup
	wg.Add(len(vms))

	for _, vm := range vms {
		go func(vm, name string) {
			defer wg.Done()
			logrus.Infof("Swicth VM %s snapshot to %s...", name, vm)

			data, err := prlctl("snapshot-list", vm, "-j")
			if err != nil {
				logrus.Errorf("Failed to get VM %s snapshots: %v", vm, err)
				return
			}

			var sps map[string]SnapshotInfo
			err = jsoniter.Unmarshal([]byte(data), &sps)
			if err != nil {
				logrus.Errorf("Get VM %s snapshot info failed: %v", vm, err)
				return
			}
			var spID string
			for id, sp := range sps {
				if sp.Name == name {
					spID = id
					break
				}
			}
			if spID == "" {
				logrus.Errorf("VM %s snapshot(%s) not found.", vm, name)
				return
			}

			err = stdPrlctl("snapshot-switch", vm, "-i", spID)
			if err != nil {
				logrus.Errorf("Failed to switch VM %s snapshot: %v", vm, err)
			}
		}(vm, name)
	}

	wg.Wait()
}

func GetVMInfo(vm string) (VMInfo, error) {
	data, err := prlctl("list", "-f", "-j", vm)
	if err != nil {
		return VMInfo{}, err
	}
	if data == "" {
		return VMInfo{}, fmt.Errorf("VM %s not found", vm)
	}

	var vms []VMInfo
	err = jsoniter.Unmarshal([]byte(data), &vms)
	if err != nil {
		return VMInfo{}, err
	}
	if len(vms) == 0 {
		return VMInfo{}, fmt.Errorf("VM %s not found", vm)
	}

	return vms[0], nil
}

func StartDockerVM(vm string, bindingHome bool) {
	info, err := GetVMInfo(vm)
	if err != nil {
		logrus.Fatalf("Failed to check vm status: %v", err)
	}
	if strings.ToLower(info.Status) != "running" {
		if _, err := prlctl("start", vm); err != nil {
			logrus.Errorf("Failed to start vm [%s]: %v", vm, err)
		}
	} else {
		logrus.Warnf("VM %s already running...", vm)
	}

	logrus.Info("Waiting vm to startup...")
	tick := time.Tick(2 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			logrus.Errorf("VM [%s] start up timeout!", vm)
			return
		case <-tick:
			if _, err := prlctl("exec", vm, "uptime"); err == nil {
				goto MOUNT
			}
		}
	}

MOUNT:
	if err := DockerMount(vm, bindingHome); err != nil {
		logrus.Errorf("Shared file system mount failed: %v", err)
	}
}

func DockerMount(vm string, bindingHome bool) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	logrus.Info("Mount macOS user home dir to vm shared dir...")
	_, err = prlctl("set", vm, "--shf-host-defined", "home")
	if err != nil {
		return err
	}

	if bindingHome {
		logrus.Info("Binding macOS user home dir to vm root home dir...")
		err = stdPrlctl("exec", vm, "mount Home /root -t prl_fs -o rw,sync,nosuid,nodev,noatime,ttl=250,share")
		if err != nil {
			return err
		}
	}

	logrus.Info("Binding macOS user home dir to the same dir of VM to fix docker volume...")
	err = stdPrlctl("exec", vm, "mount Home "+home+" -t prl_fs -o rw,sync,nosuid,nodev,noatime,ttl=250,share,x-mount.mkdir")
	if err != nil {
		return err
	}

	return nil
}
