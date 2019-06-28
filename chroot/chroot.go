package chroot

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/rhdedgar/pleg-watcher/containerinfo"
	"github.com/rhdedgar/pleg-watcher/models"
)

// SysCmd waits for a container ID via channel input, and gathers information
func SysCmd(cmdChan, runcChan <-chan string) {
	_, err := chrootPath("/host")
	if err != nil {
		fmt.Println("Error getting chroot on host in ProcessContainer due to: ", err)
	}

	for {
		select {
		case containerID := <-cmdChan:
			fmt.Println("running this: ", containerinfo.Path+" inspect "+containerID)
			cmd := exec.Command(containerinfo.Path, "inspect", containerID)

			var out bytes.Buffer
			cmd.Stdout = &out

			if cErr := cmd.Run(); err != nil {
				fmt.Println("Error running inspect command: ", cErr)
			}

			//sStr := out.String()
			//fmt.Println("Command output was", sStr)
			models.ChrootOut <- out.Bytes()

		case scanContainer := <-runcChan:
			//fmt.Println("running runc inspect command")
			runCmd := exec.Command("/usr/bin/runc", "state", scanContainer)

			var runOut bytes.Buffer
			runCmd.Stdout = &runOut

			if runcErr := runCmd.Run(); err != nil {
				fmt.Println("Error running state command: ", runcErr)
			}

			//runcStr := runOut.String()
			//fmt.Println("Command output was", runcStr)
			models.RuncOut <- runOut.Bytes()
		}
	}
}

// Chroot provides chroot access to the mounted host filesystem
func chrootPath(chrPath string) (func() error, error) {
	root, err := os.Open("/")
	if err != nil {
		fmt.Println("Error getting root FD", err)
		return nil, err
	}

	if err := syscall.Chroot(chrPath); err != nil {
		root.Close()
		return nil, err
	}

	return func() error {
		defer root.Close()
		if err := root.Chdir(); err != nil {
			return err
		}
		return syscall.Chroot(".")
	}, nil
}
