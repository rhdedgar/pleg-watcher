/*
Copyright 2019 Doug Edgar.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
	if nsList != "" {
		if _, ok := skipNS[ns]; ok {
			return true
		}
	}

	// If the provided NS starts with a prefix in the restricted prefix list, return true
	if nsPrefixList != "" {
		for _, i := range skipPrefix {
			if strings.HasPrefix(ns, i) {
				return true
			}
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

	if containerID == "" {
		return fmt.Errorf("Cannot process empty containerID.\n")
	}

	jbyte := dial.CallInfoSrv(containerID, "GetContainerInfo")

	if len(jbyte) <= 0 {
		return fmt.Errorf("Bytes returned empty.\n")
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
	// Make the environment variable into something usable (a string array) that we can check for skippable prefixes.
	skipPrefix = strings.Split(config.SkipNamespacePrefixes, ",")

	// Pre-populate a map which will be used like a set to determine if an NS should be skipped.
	for _, i := range strings.Split(config.SkipNamespaces, ",") {
		skipNS[i] = struct{}{}
	}
}
