package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rhdedgar/pleg-watcher/models"
	"io/ioutil"
	"net/http"
	"os"
)

// SendData Marshals and POSTs json data to the pod-logger service.
func SendData(mStat models.Status) {
	url := os.Getenv("POD_LOGGER_URL")

	jsonStr, err := json.Marshal(mStat)
	if err != nil {
		fmt.Println(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// TODO Prometheus to check header response
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
