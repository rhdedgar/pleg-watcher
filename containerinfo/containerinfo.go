package containerinfo

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/rhdedgar/pleg-watcher/config"
	"github.com/rhdedgar/pleg-watcher/containerscan"
	"github.com/rhdedgar/pleg-watcher/dial"
	"github.com/rhdedgar/pleg-watcher/docker"
	"github.com/rhdedgar/pleg-watcher/models"
	"github.com/rhdedgar/pleg-watcher/sender"
)

var (
	// Path is the path to the container runtime interface utility
	Path = "/usr/bin/crictl"
	// UseDocker if crictl not found
	UseDocker    = false
	skipPrefix   []string
	skipNS       = make(map[string]struct{})
	nsPrefixList = config.SkipNamespacePrefixes
	nsList       = config.SkipNamespaces
)

// skipNamespace checks our env var provided skip lists for any namespaces
// that we don't want to scan.
func skipNamespace(ns string) bool {
	// If the provided NS is in the skip map, return true
	if _, ok := skipNS[ns]; ok {
		return true
	}

	// If the provided NS starts with a prefix in the restricted prefix list, return true
	for _, i := range skipPrefix {
		if strings.HasPrefix(ns, i) {
			return true
		}
	}

	return false
}

// ProcessContainer takes a containerID string and retrieves
// info about it from crictl. Then sends the information to
// pod-logger if found.
func ProcessContainer(containerID string) error {
	var dCon docker.DockerContainer
	var cCon models.Status

	jbyte := dial.CallInfoSrv(containerID, "GetContainerInfo")

	if len(jbyte) > 0 {
		fmt.Println("Bytes returned empty")
		return nil
	}

	if UseDocker {
		fmt.Println("docker enabled, dCon is:", dCon)
		if err := json.Unmarshal(jbyte, &dCon); err != nil {
			return fmt.Errorf("Error unmarshalling docker output json: %v\n", err)
		}

		podNS := dCon[0].Config.Labels.IoKubernetesPodNamespace

		if skipNamespace(podNS) {
			return fmt.Errorf("Container is in the %v namespace, skipping \n", podNS)

		} else if dCon[0].State.Status == "running" {
			fmt.Println("container state is running")
			go containerscan.PrepDockerScan(dCon)
			sender.SendDockerData(dCon)
		}
	} else {
		if err := json.Unmarshal(jbyte, &cCon); err != nil {
			fmt.Println("bytes look like: ", string(jbyte))
			return fmt.Errorf("Error unmarshalling crictl output json: %v\n", err)
		}
		podNS := cCon.Status.Labels.IoKubernetesPodNamespace

		if skipNamespace(podNS) {
			return fmt.Errorf("Container is in the %v namespace, skipping\n", podNS)

		} else if cCon.Status.State == "CONTAINER_RUNNING" {
			go containerscan.PrepCrioScan(cCon)
			sender.SendCrioData(cCon)
		}
	}
	return nil
}

func init() {
	if _, err := os.Stat("/host/usr/bin/crictl"); os.IsNotExist(err) {
		fmt.Println("Cannot find /host/usr/bin/crictl, using /host/usr/bin/docker")
		Path = "/usr/bin/docker"
		UseDocker = true
	}
	skipPrefix = strings.Split(config.SkipNamespacePrefixes, ",")

	for _, i := range strings.Split(config.SkipNamespaces, ",") {
		skipNS[i] = struct{}{}
	}
}
