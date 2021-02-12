package scheduler

import (
	"encoding/json"
	"fmt"
	"strings"

	clscmd "github.com/rhdedgar/pleg-watcher/cmd"
	"github.com/rhdedgar/pleg-watcher/config"
	"github.com/rhdedgar/pleg-watcher/containerinfo"
	"github.com/rhdedgar/pleg-watcher/crictilspec"
	"github.com/rhdedgar/pleg-watcher/dial"
	mainscan "github.com/rhdedgar/pleg-watcher/scanner"
)

// ScheduledContainerScan gets a list of currently running containers from
// the container-info sidecar, and queus them for scanning.
func ScheduledContainerScan() {
	var crictlOutput crictilspec.Containers

	jbyte := dial.GetContainerInfo()

	if len(jbyte) > 0 {
		fmt.Println("Container list returned empty.")
	}

	if err := json.Unmarshal(jbyte, &crictlOutput); err != nil {
		fmt.Println("Error unmarshalling crictl list output json: ", err)
	}

	for _, container := range crictlOutput.Containers {
		containerID := container.ID
		containerinfo.ProcessContainer(containerID)
	}
}

// ScheduledHostScan performs a malware scan on the node/host OS files
func ScheduledHostScan() {
	scanDirs := strings.Split(config.ScanDirs, ",")

	fmt.Printf("%v top-level directories to scan\n", len(scanDirs))

	for _, scanDir := range scanDirs {
		fmt.Println("Scanning directory:", scanDir)

		scannerOptions := clscmd.NewDefaultManagedScannerOptions()
		scannerOptions.PostResultURL = config.PostResultURL
		scannerOptions.OutFile = config.OutFile
		scannerOptions.ScanDir = scanDir

		if err := scannerOptions.Validate(); err != nil {
			fmt.Println("Error validating scanner options: ", err)
		}

		scanner := mainscan.NewDefaultManagedScanner(*scannerOptions)

		if err := scanner.AcquireAndScan(); err != nil {
			fmt.Println("Error returned from scanner: ", err)
		}
		/*
			err = unix.Unmount(scanDir, 0)
			if err != nil {
				fmt.Println("Error unmounting scanDir after scanning: ", err)
			}

			os.Remove(scanDir)
			if err != nil {
				fmt.Println("Error removing scanDir after unmounting: ", err)
			}*/
	}
}
