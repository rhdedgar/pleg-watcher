package containerinfo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/rhdedgar/pleg-watcher/chroot"
	"github.com/rhdedgar/pleg-watcher/containerscan"
	"github.com/rhdedgar/pleg-watcher/docker"
	"github.com/rhdedgar/pleg-watcher/models"
	"github.com/rhdedgar/pleg-watcher/sender"
)

var (
	path = "/usr/bin/crictl"
	// UseDocker if crictl not found
	UseDocker = false
)

// ProcessContainer takes a containerID string and retrieves
// info about it from crictl. Then sends the information to
// pod-logger if found.
func ProcessContainer(containerID string) {
	var dCon docker.DockerContainer
	var cCon models.Status

	fmt.Println("inspecting: ", containerID)

	exit, err := chroot.Chroot("/host")
	if err != nil {
		panic(err)
	}

	// Using command output in lieu of a wrapper
	cmd := exec.Command(path, "inspect", containerID)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		fmt.Println("Error running command: ", err)
	}

	if err := exit(); err != nil {
		panic(err)
	}

	jbyte := out.Bytes()

	if UseDocker {
		if err := json.Unmarshal(jbyte, &dCon); err != nil {
			fmt.Println("Error unmarshalling docker output json:", err)
		}
		if strings.HasPrefix(dCon[0].State.Status.Label.IoKubernetesPodNamespace, "openshift-") {
			return
		} else if dCon[0].State.Status == "running" {
			sender.SendDockerData(dCon)
		}
	} else {
		if err := json.Unmarshal(jbyte, &cCon); err != nil {
			fmt.Println("Error unmarshalling crictl output json:", err)
		}
		if strings.HasPrefix(dCon[0].State.Status.Label.IoKubernetesPodNamespace, "openshift-") {
			return
		} else if cCon.Status.State == "CONTAINER_RUNNING" {
			go sender.SendCrioData(cCon)
			containerscan.PrepCrioScan(cCon)
		}
	}
}

func init() {
	if _, err := os.Stat("/host/usr/bin/crictl"); os.IsNotExist(err) {
		fmt.Println("Cannot find /host/usr/bin/crictl, using /host/usr/bin/docker")
		path = "/usr/bin/docker"
		UseDocker = true
	}
}
