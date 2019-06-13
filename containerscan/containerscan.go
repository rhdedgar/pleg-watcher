package containerscan

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/rhdedgar/pleg-watcher/chroot"
	clscmd "github.com/rhdedgar/pleg-watcher/cmd"
	"github.com/rhdedgar/pleg-watcher/models"
	"github.com/rhdedgar/pleg-watcher/runcspec"
	mainscan "github.com/rhdedgar/pleg-watcher/scanner"
)

// custSplit takes 3 parameters and returns a string.
// s is the string to split.
// d is the delimiter by which to split s.
// i is the slice index of the string to return, if applicable. Usually 1 or 0.
// If the string was not split, the original string is returned idempotently.
func custSplit(s, d string, i int) string {
	tempS := s
	splits := strings.Split(s, d)

	if len(splits) >= i+1 {
		tempS = splits[i]
	}

	return tempS
}

// custReg takes 2 arguments and returns a string slice.
//
// scanOut is the string output from the crio /proc/$PID/mountinfo file.
//
// regString is the `raw string` containing the regex match pattern to use.
func custReg(scanOut, regString string) []string {
	var newLayers []string

	reg := regexp.MustCompile(regString)
	matched := reg.FindAllString(scanOut, -1)

	if matched != nil {
		for _, layer := range matched {
			newLayers = append(newLayers, layer)
		}
	}

	return newLayers
}

func getCrioLayers(containerID string) []string {
	var layers []string
	var crioLayers []string

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
	//rootPath := runcState.RootFS
	//dirPath := filepath.Dir(rootPath)
	//IDPath := filepath.Base(rootPath)

	mountPath := "/proc/" + string(pid) + "/mountinfo"
	//mountOutput := ""

	f, err := os.Open(mountPath)
	if err != nil {
		fmt.Println("error opening file:", mountPath, err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Scan()
	scanOut := scanner.Text()

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading layer", err)
	}

	layers = append(layers, custReg(scanOut, `lowerdir=(.*),upperdir`)...)
	layers = append(layers, custReg(scanOut, `upperdir=(.*),workdir`)...)

	for _, l := range layers {
		items := strings.Split(l, ":")
		for _, i := range items {
			j := custSplit(i, ",", 0)
			j = custSplit(j, "=", 1)

			crioLayers = append(crioLayers, j)
		}
	}

	//if runcState.Status == "CONTAINER_RUNNING" {
	//	sender.SendCrioData(cCon)
	//}

	return crioLayers
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
