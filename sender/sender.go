package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rhdedgar/pleg-watcher/api"
	"github.com/rhdedgar/pleg-watcher/config"
	"github.com/rhdedgar/pleg-watcher/docker"
	"github.com/rhdedgar/pleg-watcher/models"
)

// SendDockerData Marshals and POSTs json data to the pod-logger service.
func SendDockerData(dCon docker.DockerContainer) (int, error) {
	jsonStr, err := json.Marshal(dCon)
	if err != nil {
		return 0, fmt.Errorf("Error marshalling docker json to send to pod-logger: %v\n", err)
	}
	resp, err := sendLog(jsonStr, config.DockerURL)
	if err != nil {
		return 0, fmt.Errorf("Error sending log: %v\n", err)
	}
	return resp, nil
}

// SendCrioData Marshals and POSTs json data to the pod-logger service.
func SendCrioData(mStat models.Status) (int, error) {
	jsonStr, err := json.Marshal(mStat)
	if err != nil {
		return 0, fmt.Errorf("Error marshalling crio json to send to pod-logger: %v\n", err)
	}
	resp, err := sendLog(jsonStr, config.CrioURL)
	if err != nil {
		return 0, fmt.Errorf("Error sending log: %v\n", err)
	}
	return resp, nil
}

// SendClamData Marshals and POSTs json data to the pod-logger service.
func SendClamData(sRes api.ScanResult) (int, error) {
	jsonStr, err := json.Marshal(sRes)
	if err != nil {
		return 0, fmt.Errorf("Error marshalling clam json to send to pod-logger: %v\n", err)
	}
	resp, err := sendLog(jsonStr, config.ClamURL)
	if err != nil {
		return 0, fmt.Errorf("Error sending log: %v\n", err)
	}
	return resp, nil
}

func sendLog(jsonStr []byte, url string) (int, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Println("Error creating new HTTP request:", err)
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending to pod-logger at %v: %v \n", url, err)
		fmt.Printf("Could not send %v \n", string(jsonStr[:]))
		return 0, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}
