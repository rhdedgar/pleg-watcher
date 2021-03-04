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

	jbyte := dial.GetActiveContainers()

	if len(jbyte) > 0 {
		fmt.Println("Container list returned empty.")
		return
	}

	if err := json.Unmarshal(jbyte, &crictlOutput); err != nil {
		fmt.Println("Error unmarshalling crictl list output json: ", err)
		return
	}

	for _, container := range crictlOutput.Containers {
		containerID := container.ID
		if containerID == "" {
			fmt.Println("crictilOutput container.ID is empty. Has crictl output changed?")
			continue
		}
		containerinfo.ProcessContainer(containerID)
	}
}

// ScheduledHostScan performs a malware scan on the node/host OS files
func ScheduledHostScan() {
	scanDirs := strings.Split(config.ScanDirs, ",")

	fmt.Printf("%v top-level directories to scan\n", len(scanDirs))

	scannerOptions := clscmd.NewDefaultManagedScannerOptions()
	scannerOptions.PostResultURL = config.PostResultURL
	scannerOptions.OutFile = config.OutFile

	for _, scanDir := range scanDirs {
		fmt.Println("Scanning directory:", scanDir)

		scannerOptions.ScanDir = scanDir

		if err := scannerOptions.Validate(); err != nil {
			fmt.Println("Error validating scanner options: ", err)
			continue
		}

		scanner := mainscan.NewDefaultManagedScanner(*scannerOptions)

		if err := scanner.AcquireAndScan(); err != nil {
			fmt.Println("Error returned from scanner: ", err)
		}
	}
}
