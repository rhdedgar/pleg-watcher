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

	clscmd "github.com/rhdedgar/pleg-watcher/cmd"
	"github.com/rhdedgar/pleg-watcher/config"
	"github.com/rhdedgar/pleg-watcher/dial"
	"github.com/rhdedgar/pleg-watcher/docker"
	"github.com/rhdedgar/pleg-watcher/models"
	"github.com/rhdedgar/pleg-watcher/runcspec"
	mainscan "github.com/rhdedgar/pleg-watcher/scanner"
	"golang.org/x/sys/unix"
)

var (
	scanResultsDir = config.ScanResultsDir
	postResultURL  = config.PostResultURL
	outFile        = config.OutFile
)

// CustSplit takes 3 parameters and returns a string.
// s is the string to split.
// d is the delimiter by which to split s.
// i is the slice index of the string to return, if applicable. Usually 1 or 0.
// If the string was not split, the original string is returned idempotently.
func CustSplit(s, d string, i int) string {
	tempS := s
	splits := strings.Split(s, d)

	if len(splits) >= i+1 {
		tempS = splits[i]
	}

	return tempS
}

// CustReg takes 2 arguments and returns a string slice.
//
// scanOut is the string output from the container /proc/$PID/mountinfo file.
//
// regString is the `raw string` containing the regex match pattern to use.
func CustReg(scanOut, regString string) []string {
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

// MountOverlayFS takes a slice of strings containing OverlayFS layer paths
// and mounts them to a read-only /mnt dir named after their container ID.
func MountOverlayFS(layers []string, cID string) (string, error) {
	scanDir := "/mnt/" + cID

	err := os.MkdirAll(scanDir, 0700)
	if err != nil {
		return scanDir, fmt.Errorf("Error creating scanDir: %v", err)
	}

	overlayPath := "lowerdir=/host" + strings.Join(layers, ":/host")
	fmt.Println("Trying to mount: ", overlayPath)

	err = unix.Mount("overlay", scanDir, "overlay", unix.MS_NODEV|unix.MS_NOEXEC|unix.MS_RDONLY, overlayPath)
	if err != nil {
		return scanDir, fmt.Errorf("Error mounting scanDir: %v", err)
	}

	return scanDir, nil
}

// GetRootFS returns the root file system path for a container returned by "runc state <containerID>"
func GetRootFS(containerID string) (string, error) {
	var runcState runcspec.RuncState

	fmt.Println("Getting root container layer for: ", containerID)

	// Avoid race condition with container info not being available yet
	for i := 1; i <= 6; i++ {
		jbyte := dial.CallInfoSrv(containerID, "GetRuncInfo")

		if len(jbyte) == 0 {
			if i > 5 {
				return "", fmt.Errorf("GetRootFS: Error getting root FS")
			}
			time.Sleep(time.Duration(i) * time.Second)
			continue
		} else {
			if err := json.Unmarshal(jbyte, &runcState); err != nil {
				fmt.Println("Output returned from runc state: ", string(jbyte))
				return "", fmt.Errorf("Error unmarshalling runc output json: %v", err)
			}

			if runcState.RootFS != "" {
				return runcState.RootFS, nil
			}
		}
	}
	return "", fmt.Errorf("Output of runc state RootFS was empty")
}

// GetLayerInfo reads /host/proc/<PID>/mountinfo and returns the line containing OverlayFS mount directories.
func GetLayerInfo(mountPath string) (string, error) {
	var scanOut string

	for i := 0; i <= 5; i++ {
		f, err := os.Open(mountPath)
		if err != nil {
			if i >= 5 {
				return "", fmt.Errorf("GetLayerInfo: Error returning layers")
			}
			fmt.Println("GetLayerInfo: Error opening file, waiting 5 seconds in case it just hasn't been created yet: ", mountPath, err)
			time.Sleep(5 * time.Second)
			continue
		} else {
			defer f.Close()

			bufScan := bufio.NewScanner(f)
			bufScan.Scan()
			scanOut = bufScan.Text()

			if err := bufScan.Err(); err != nil {
				return "", fmt.Errorf("GetLayerInfo: Error reading layer %v", err)
			}
			break
		}
	}

	return scanOut, nil
}

// getCrioLayers takes a containerID string and queries crictl inspect and runc state
// for overlayFS mount data found in /proc/PID/mountinfo.
func getCrioLayers(containerID string) ([]string, error) {
	var layers []string
	var crioLayers []string
	var runcState runcspec.RuncState

	fmt.Println("Getting cri-o layers: ", containerID)

	jbyte := dial.CallInfoSrv(containerID, "GetRuncInfo")

	if err := json.Unmarshal(jbyte, &runcState); err != nil {
		fmt.Println(string(jbyte))
		return crioLayers, fmt.Errorf("Error unmarshalling runc output json: %v", err)
	}

	pid := runcState.Pid

	mountPath := "/host/proc/" + strconv.Itoa(pid) + "/mountinfo"

	scanOut, err := GetLayerInfo(mountPath)
	if err != nil {
		fmt.Println(err)
	}

	layers = append(layers, CustReg(scanOut, `lowerdir=(.*),upperdir`)...)

	for _, l := range layers {
		items := strings.Split(l, ":")
		for _, i := range items {
			j := CustSplit(i, ",", 0)
			j = CustSplit(j, "=", 1)

			crioLayers = append(crioLayers, j)
		}
	}
	return crioLayers, nil
}

// PrepCrioScan gets a slice of container filesystem layers from getCrioLayers
// and then initiates a scan for each of the returned layers.
func PrepCrioScan(cCon models.Status) {
	fmt.Printf("In scan block for %v in %v\n",
		cCon.Status.Labels.IoKubernetesPodName,
		cCon.Status.Labels.IoKubernetesPodNamespace)

	scannerOptions := clscmd.NewDefaultManagedScannerOptions()
	cID := cCon.Status.ID

	cLayers, err := getCrioLayers(cID)

	if err != nil {
		fmt.Println(err)
		return
	}

	if len(cLayers) == 0 {
		fmt.Println("layers returned empty")
		return
	}
	fmt.Println(cLayers)

	scanDir, err := MountOverlayFS(cLayers, cID)
	if err != nil {
		fmt.Println(err)
		return
	}

	scannerOptions.ScanResultsDir = scanResultsDir
	scannerOptions.PostResultURL = postResultURL
	scannerOptions.OutFile = outFile
	scannerOptions.ScanDir = scanDir

	if err := scannerOptions.Validate(); err != nil {
		fmt.Println("Error validating scanner options: ", err)
	}

	scanner := mainscan.NewDefaultManagedScanner(*scannerOptions)
	scanner.ScanOutputs.ScanResults.NameSpace = cCon.Status.Labels.IoKubernetesPodNamespace
	scanner.ScanOutputs.ScanResults.PodName = cCon.Status.Labels.IoKubernetesPodName

	if err := scanner.AcquireAndScan(); err != nil {
		fmt.Println("Error returned from scanner: ", err)
	}

	err = unix.Unmount(scanDir, 0)
	if err != nil {
		fmt.Println("Error unmounting scanDir after scanning: ", err)
	}

	os.Remove(scanDir)
	if err != nil {
		fmt.Println("Error removing scanDir after unmounting: ", err)
	}
}

// getDockerLayers takes a containerID string and queries docker inspect
// for overlayFS mount data found in /proc/PID/mountinfo.
func getDockerLayers(containerID string, procID int) ([]string, error) {
	var layers []string
	var dockerLayers []string

	fmt.Println("Getting docker layers: ", containerID)

	mountPath := "/host/proc/" + strconv.Itoa(procID) + "/mountinfo"

	fmt.Println("Going to open:", mountPath)

	f, err := os.Open(mountPath)
	if err != nil {
		fmt.Println("getDockerLayers: Error opening file, waiting 10 seconds in case it just hasn't been created yet: ", mountPath, err)
		time.Sleep(10 * time.Second)
		f, err = os.Open(mountPath)
	}

	defer f.Close()

	bufScan := bufio.NewScanner(f)
	bufScan.Scan()
	scanOut := bufScan.Text()

	if err := bufScan.Err(); err != nil {
		return dockerLayers, fmt.Errorf("Error reading layer %v", err)
	}

	layers = append(layers, CustReg(scanOut, `lowerdir=(.*),upperdir`)...)

	for _, l := range layers {
		items := strings.Split(l, ":")
		for _, i := range items {
			j := CustSplit(i, ",", 0)
			j = CustSplit(j, "=", 1)

			dockerLayers = append(dockerLayers, j)
		}
	}
	return dockerLayers, nil
}

// PrepDockerScan gets a slice of container filesystem layers from getDockerLayers
// and then initiates a scan for each of the returned layers.
func PrepDockerScan(dCon docker.DockerContainer) {
	fmt.Printf("In scan block for %v in %v\n",
		dCon[0].Config.Labels.IoKubernetesPodName,
		dCon[0].Config.Labels.IoKubernetesPodNamespace)

	scannerOptions := clscmd.NewDefaultManagedScannerOptions()
	cID := dCon[0].ID
	pID := dCon[0].State.Pid

	cLayers, err := getDockerLayers(cID, pID)

	if err != nil {
		fmt.Println(err)
		return
	}

	if len(cLayers) == 0 {
		fmt.Println("layers returned empty")
		return
	}
	fmt.Println(cLayers)

	scanDir, err := MountOverlayFS(cLayers, cID)
	if err != nil {
		fmt.Println(err)
		return
	}

	scannerOptions.ScanResultsDir = scanResultsDir
	scannerOptions.PostResultURL = postResultURL
	scannerOptions.OutFile = outFile
	scannerOptions.ScanDir = scanDir

	if err := scannerOptions.Validate(); err != nil {
		fmt.Println("Error validating scanner options: ", err)
	}

	scanner := mainscan.NewDefaultManagedScanner(*scannerOptions)
	scanner.ScanOutputs.ScanResults.NameSpace = dCon[0].Config.Labels.IoKubernetesPodNamespace
	scanner.ScanOutputs.ScanResults.PodName = dCon[0].Config.Labels.IoKubernetesPodName

	if err := scanner.AcquireAndScan(); err != nil {
		fmt.Println("Error returned from scanner: ", err)
	}

	err = unix.Unmount(scanDir, 0)
	if err != nil {
		fmt.Println("Error unmounting scanDir after scanning: ", err)
	}

	os.Remove(scanDir)
	if err != nil {
		fmt.Println("Error removing scanDir after unmounting: ", err)
	}
}
