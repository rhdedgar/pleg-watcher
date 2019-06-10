package containerscan

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"

	"github.com/rhdedgar/pleg-watcher/chroot"
	clscmd "github.com/rhdedgar/pleg-watcher/cmd"
	"github.com/rhdedgar/pleg-watcher/models"
	"github.com/rhdedgar/pleg-watcher/runcspec"
	mainscan "github.com/rhdedgar/pleg-watcher/scanner"
)

func getCrioLayers(containerID string) []string {
	var layers []string

	//var dCon docker.DockerContainer
	//var cCon models.Status
	var runcState runcspec.RuncState

	fmt.Println("inspecting: ", containerID)

	exit, err := chroot.Chroot("/host")
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("/usr/bin/runc", "inspect", containerID)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		fmt.Println("Error running command: ", err)
	}

	if err := exit(); err != nil {
		panic(err)
	}

	jbyte := out.Bytes()

	if err := json.Unmarshal(jbyte, &runcState); err != nil {
		fmt.Println("Error unmarshalling crictl output json:", err)
	}

	pid := runcState.Pid
	mountPath := "/proc/" + string(pid) + "/mountinfo"

	//if runcState.Status == "CONTAINER_RUNNING" {
	//	sender.SendCrioData(cCon)
	//}

	return layers
}

// PrepCrioScan gets a slice of container filesystem layers from getCrioLayers
// and then initiates a scan for each of the returned layers.
func PrepCrioScan(cCon models.Status) {
	scannerOptions := clscmd.NewDefaultContainerLayerScannerOptions()
	cID := cCon.Status.ID

	cLayers := getCrioLayers(cID)

	scannerOptions.ScanResultsDir = ""
	scannerOptions.PostResultURL = ""
	scannerOptions.OutFile = ""

	for _, l := range cLayers {
		scannerOptions.ScanDir = l

		if err := scannerOptions.Validate(); err != nil {
			log.Fatalf("Error: %v", err)
		}

		scanner := mainscan.NewDefaultContainerLayerScanner(*scannerOptions)
		if err := scanner.ClamScanner(); err != nil {
			log.Fatalf("Error: %v", err)
		}
	}
}
