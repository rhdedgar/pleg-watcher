package container

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rhdedgar/pleg-watcher/models"
	"github.com/rhdedgar/pleg-watcher/sender"
	"log"
	"os/exec"
)

func ProcessContainer(containerID string) {
	var container models.Status

	// Using raw command output in lieu of a proper wrapper
	cmd := exec.Command("/usr/bin/crictl", "inspect", containerID)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	jbyte := out.Bytes()

	err := json.Unmarshal(jbyte, &container)
	if err != nil {
		fmt.Println("There was an error:", err)
	}

	if container.Status.State == "CONTAINER_RUNNING" {
		sender.SendData(container)
	}
}
