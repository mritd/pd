package helper

import (
	"context"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

func StartDockerVM(vm string, bindingHome bool) {
	logrus.Info("Starting VM...")

	info, err := GetVMInfo(vm)
	if err != nil {
		logrus.Fatalf("Failed to check vm status: %v", err)
	}
	if strings.ToLower(info.Status) != "running" {
		if err := stdPrlctl("start", vm); err != nil {
			logrus.Errorf("Failed to start VM [%s]: %v", vm, err)
		}
	} else {
		logrus.Warnf("VM [%s] already running...", vm)
	}

	logrus.Info("Waiting VM to startup...")
	tick := time.Tick(2 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			logrus.Fatalf("VM [%s] start up timeout!", vm)
		case <-tick:
			if _, err := prlctl("exec", vm, "uptime"); err == nil {
				goto MOUNT
			}
		}
	}

MOUNT:
	logrus.Info("Mount shared dir to VM...")
	if err := mount(vm, bindingHome); err != nil {
		logrus.Fatalf("Shared file system mount failed: %v", err)
	}
}

func mount(vm string, bindingHome bool) error {
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

	logrus.Info("Binding macOS user home dir to the same dir of VM to fix symlink...")
	err = stdPrlctl("exec", vm, "mount Home "+home+" -t prl_fs -o rw,sync,nosuid,nodev,noatime,ttl=250,share,x-mount.mkdir")
	if err != nil {
		return err
	}

	return nil
}

