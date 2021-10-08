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

type Snapshot struct {
	ID string `json:"id"`
	SnapshotInfo
}

type Snapshots []Snapshot

func (ss Snapshots) Len() int {
	return len(ss)
}

func (ss Snapshots) Less(i, j int) bool {
	t1, _ := time.Parse("2006-01-02 15:04:05", ss[i].Date)
	t2, _ := time.Parse("2006-01-02 15:04:05", ss[j].Date)
	return t1.Before(t2)
}

func (ss Snapshots) Swap(i, j int) {
	ss[i], ss[j] = ss[j], ss[i]
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
				_, err = prlctl("stop", vm, "--kill")
			} else {
				_, err = prlctl("stop", vm)
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
	ss := convert2Snapshots(sps)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")

	table.SetHeader([]string{"NAME", "CURRENT", "UUID", "STATE", "DATE"})
	for _, sp := range ss {
		table.Append([]string{sp.Name, strconv.FormatBool(sp.Current), sp.ID, sp.State, sp.Date})
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

			ss := convert2Snapshots(sps)
			var spID string
			for _, sp := range ss {
				if sp.Name == name || sp.ID == name {
					spID = sp.ID
					break
				}
			}
			if spID == "" {
				if name == "latest" && len(ss) > 0 {
					spID = ss[len(ss)-1].ID
				} else {
					logrus.Warnf("VM %s snapshot(%s) not found.", vm, name)
					return
				}
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

			ss := convert2Snapshots(sps)
			var spID string
			for _, sp := range ss {
				if sp.Name == name || sp.ID == name {
					spID = sp.ID
					break
				}
			}
			if spID == "" {
				if name == "latest" && len(ss) > 0 {
					spID = ss[len(ss)-1].ID
				} else {
					logrus.Errorf("VM %s snapshot(%s) not found.", vm, name)
					return
				}
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

func SetVMCPU(vms []string, count int) {
	var wg sync.WaitGroup
	wg.Add(len(vms))
	for _, vm := range vms {
		go func(vm string, count int) {
			defer wg.Done()
			logrus.Infof("Set the number of CPU cores of the VM [%s] to %d...", vm, count)

			if _, err := prlctl("set", vm, "--cpus", strconv.Itoa(count)); err != nil {
				logrus.Errorf("Failed to set VM %s CPU cores: %v", vm, err)
			}
		}(vm, count)
	}
	wg.Wait()
}

func SetVMRAM(vms []string, size int) {
	var wg sync.WaitGroup
	wg.Add(len(vms))
	for _, vm := range vms {
		go func(vm string, size int) {
			defer wg.Done()
			logrus.Infof("Set the memory size of the VM [%s] to %d mb...", vm, size)

			if _, err := prlctl("set", vm, "--memsize", strconv.Itoa(size)); err != nil {
				logrus.Errorf("Failed to set VM %s Memsize: %v", vm, err)
			}
		}(vm, size)
	}
	wg.Wait()
}