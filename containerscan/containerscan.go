package containerscan

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rhdedgar/pleg-watcher/channels"
	clscmd "github.com/rhdedgar/pleg-watcher/cmd"
	"github.com/rhdedgar/pleg-watcher/models"
	"github.com/rhdedgar/pleg-watcher/runcspec"
	mainscan "github.com/rhdedgar/pleg-watcher/scanner"
)

var (
	scanResultsDir = os.Getenv("SCANRESULTSDIR")
	postResultURL  = os.Getenv("POSTRESULTURL")
	outFile        = os.Getenv("OUTFILE")
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

func getRootFS(containerID string) (string, error) {
	var runcState runcspec.RuncState

	fmt.Println("Getting root container layer for: ", containerID)

	go channels.SetStringChan(models.RuncChan, containerID)

	jbyte := <-models.RuncOut

	if err := json.Unmarshal(jbyte, &runcState); err != nil {
		fmt.Println("Output returned from runc state: ", string(jbyte))
		return "", fmt.Errorf("Error unmarshalling runc output json: %v", err)
	}

	if runcState.RootFS != "" {
		return runcState.RootFS, nil
	}

	return "", fmt.Errorf("Output of runc state RootFS was empty")
}

func getCrioLayers(containerID string) []string {
	var layers []string
	var crioLayers []string
	var runcState runcspec.RuncState

	fmt.Println("Getting cri-o layers: ", containerID)

	go channels.SetStringChan(models.RuncChan, containerID)

	jbyte := <-models.RuncOut

	//fmt.Println("Channel returned: ", string(jbyte))
	if err := json.Unmarshal(jbyte, &runcState); err != nil {
		fmt.Println("Error unmarshalling runc output json:", err)
		fmt.Println(string(jbyte))
		return crioLayers
	}

	pid := runcState.Pid
	//rootPath := runcState.RootFS
	//dirPath := filepath.Dir(rootPath)
	//IDPath := filepath.Base(rootPath)

	mountPath := "/proc/" + strconv.Itoa(pid) + "/mountinfo"
	//mountOutput := ""

	f, err := os.Open(mountPath)
	if err != nil {
		fmt.Println("Error opening file, waiting 5 seconds in case it just hasn't been created yet: ", mountPath, err)
		time.Sleep(5 * time.Second)
		f, err = os.Open(mountPath)
	}

	defer f.Close()

	bufScan := bufio.NewScanner(f)
	bufScan.Scan()
	scanOut := bufScan.Text()

	if err := bufScan.Err(); err != nil {
		fmt.Println("Error reading layer", err)
		return crioLayers
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
	//fmt.Println("returning layers")
	return crioLayers
}

// PrepCrioScan gets a slice of container filesystem layers from getCrioLayers
// and then initiates a scan for each of the returned layers.
func PrepCrioScan(cCon models.Status) {
	fmt.Println("In scan block")
	scannerOptions := clscmd.NewDefaultContainerLayerScannerOptions()
	cID := cCon.Status.ID

	//cLayers := getCrioLayers(cID)

	rootFS, err := getRootFS(cID)

	if err != nil {
		fmt.Println(err)
		return
	}

	//if len(cLayers) == 0 {
	//	fmt.Println("layers returned empty")
	//	return
	//}
	//fmt.Println(cLayers)

	scannerOptions.ScanResultsDir = scanResultsDir
	scannerOptions.PostResultURL = postResultURL
	scannerOptions.OutFile = outFile

	//for _, l := range cLayers {
	scannerOptions.ScanDir = rootFS

	if err := scannerOptions.Validate(); err != nil {
		fmt.Println("Error validating scanner options: ", err)
	}

	scanner := mainscan.NewDefaultContainerLayerScanner(*scannerOptions)
	scanner.ScanOutputs.ScanResults.NameSpace = cCon.Status.Labels.IoKubernetesPodNamespace
	scanner.ScanOutputs.ScanResults.PodName = cCon.Status.Labels.IoKubernetesPodName

	if err := scanner.AcquireAndScan(); err != nil {
		fmt.Println("Error returned from scanner: ", err)
	}
	//}
}
