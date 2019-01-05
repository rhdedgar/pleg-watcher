package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rhdedgar/pleg-watcher/models"
	"io/ioutil"
	"net/http"
)

// SendData Marshals and POSTs json data to the URL designated in the config file.
func SendData(mLabel models.Label) {
	// TODO replace with variables from config/config.go package read from a config file
	url := "http://127.0.0.1:8080"
	fmt.Println("URL:>", url)

	jsonStr, err := json.Marshal(mLabel)
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
