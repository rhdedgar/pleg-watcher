package containerinfo

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/rhdedgar/pleg-watcher/channels"
	"github.com/rhdedgar/pleg-watcher/containerscan"
	"github.com/rhdedgar/pleg-watcher/docker"
	"github.com/rhdedgar/pleg-watcher/models"
	"github.com/rhdedgar/pleg-watcher/sender"
)

var (
	// Path is the path to the container runtime interface utility
	Path = "/usr/bin/crictl"
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

	//models.ChrootChan <- containerID
	channels.SetStringChan(models.ChrootChan, containerID)
	jbyte := <-models.ChrootOut

	if UseDocker {
		if err := json.Unmarshal(jbyte, &dCon); err != nil {
			fmt.Println("Error unmarshalling docker output json:", err)
		}
		if strings.HasPrefix(dCon[0].Config.Labels.IoKubernetesPodNamespace, "openshift-") {
			return
		} else if dCon[0].State.Status == "running" {
			sender.SendDockerData(dCon)
		}
	} else {
		if err := json.Unmarshal(jbyte, &cCon); err != nil {
			fmt.Println("Error unmarshalling crictl output json:", err)
		}
		if strings.HasPrefix(cCon.Status.Labels.IoKubernetesPodNamespace, "openshift-") {
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
		Path = "/usr/bin/docker"
		UseDocker = true
	}
}
