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
	"golang.org/x/sys/unix"
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

// mountOverlayFS takes a slice of strings containing OverlayFS layer paths
// and mounts them to a read-only /mnt dir named after their container ID.
func mountOverlayFS(layers []string, cID string) (string, error) {
	scanDir := "/mnt/" + cID

	err := os.Mkdir(scanDir, 0700)
	if err != nil {
		return "", fmt.Errorf("Error creating scanDir: %v", err)
	}

	overlayPath := strings.Join(layers, ":")

	err = unix.Mount("overlay", scanDir, "overlay", unix.MS_NODEV|unix.MS_NOEXEC|unix.MS_RDONLY, overlayPath)
	if err != nil {
		return "", fmt.Errorf("Error mounting scanDir: %v", err)
	}

	return scanDir, nil
}

func getRootFS(containerID string) (string, error) {
	var runcState runcspec.RuncState

	fmt.Println("Getting root container layer for: ", containerID)

	// Avoid race condition with container layers not being written yet
	time.Sleep(15 * time.Second)

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

// getCrioLayers takes a containerID string and queries crictl inspect and runc state
// for overlayFS mount data found in /proc/PID/mountinfo.
func getCrioLayers(containerID string) ([]string, error) {
	var layers []string
	var crioLayers []string
	var runcState runcspec.RuncState

	fmt.Println("Getting cri-o layers: ", containerID)

	go channels.SetStringChan(models.RuncChan, containerID)

	jbyte := <-models.RuncOut

	//fmt.Println("Channel returned: ", string(jbyte))
	if err := json.Unmarshal(jbyte, &runcState); err != nil {
		fmt.Println(string(jbyte))
		return crioLayers, fmt.Errorf("Error unmarshalling runc output json: %v", err)
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
		return crioLayers, fmt.Errorf("Error reading layer %v", err)
	}

	layers = append(layers, custReg(scanOut, `lowerdir=(.*),upperdir`)...)
	//layers = append(layers, custReg(scanOut, `upperdir=(.*),workdir`)...)

	for _, l := range layers {
		items := strings.Split(l, ":")
		for _, i := range items {
			j := custSplit(i, ",", 0)
			j = custSplit(j, "=", 1)

			crioLayers = append(crioLayers, j)
		}
	}
	//fmt.Println("returning layers")
	return crioLayers, nil
}

// PrepCrioScan gets a slice of container filesystem layers from getCrioLayers
// and then initiates a scan for each of the returned layers.
func PrepCrioScan(cCon models.Status) {
	fmt.Printf("In scan block for %v in %v",
		cCon.Status.Labels.IoKubernetesPodName,
		cCon.Status.Labels.IoKubernetesPodNamespace)

	scannerOptions := clscmd.NewDefaultContainerLayerScannerOptions()
	cID := cCon.Status.ID

	cLayers, err := getCrioLayers(cID)
	//rootFS, err := getRootFS(cID)

	if err != nil {
		fmt.Println(err)
		return
	}

	if len(cLayers) == 0 {
		fmt.Println("layers returned empty")
		return
	}
	fmt.Println(cLayers)

	scanDir, err := mountOverlayFS(cLayers, cID)
	if err != nil {
		fmt.Println(err)
		return
	}

	scannerOptions.ScanResultsDir = scanResultsDir
	scannerOptions.PostResultURL = postResultURL
	scannerOptions.OutFile = outFile

	//fmt.Printf("Scanning both %v and %v \n", rootFS, "/host"+rootFS)

	//for _, scandir := range []string{rootFS, "/host" + rootFS} {

	scannerOptions.ScanDir = scanDir //filepath.Dir(rootFS)

	if err := scannerOptions.Validate(); err != nil {
		fmt.Println("Error validating scanner options: ", err)
	}

	scanner := mainscan.NewDefaultContainerLayerScanner(*scannerOptions)
	scanner.ScanOutputs.ScanResults.NameSpace = cCon.Status.Labels.IoKubernetesPodNamespace
	scanner.ScanOutputs.ScanResults.PodName = cCon.Status.Labels.IoKubernetesPodName

	if err := scanner.AcquireAndScan(); err != nil {
		fmt.Println("Error returned from scanner: ", err)
	}

	err = unix.Unmount(scanDir, 0)
	if err != nil {
		fmt.Println("Error unmounting scanDir after scanning: ", err)
	}
}
