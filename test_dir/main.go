package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	//"time"
)

type PLEGEvent struct {
	ID   string `json:"ID"`
	Type string `json:"Type"`
	Data string `json:"Data"`
}

func quoteVar(s string, r string) string {
	return strings.Replace(s, r, "\""+r+"\"", 1)
}

func main() {
	//p := fmt.Println

	//now := time.Now()
	//p(now)

	//out := strings.TrimLeft(strings.TrimRight(file_str, "{"), "}")
	//out := strings.SplitAfter(strings.SplitAfter(file_str, "{")[1], "}")

	var plegEvent PLEGEvent

	filePath := "./containerstarted.txt"

	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println(err)
	}

	fileStr := string(fileBytes)

	if strings.Contains(fileStr, "ContainerStarted") {

		// Gather only the unquoted json of the PLEG Event
		out := strings.SplitAfter(fileStr, "&pleg.PodLifecycleEvent")[1]

		//		locJson, err := json.Marshal(out)
		//		if err != nil {
		//			fmt.Println("error marshalling json:\n", err)
		//		}

		//		fmt.Println(string(locJson))
		//		fmt.Println(json.Unmarshal(locJson, &plegEvent))

		// Quote the json so it can be Unmarshaled into a struct
		for _, item := range []string{"ID", "Type", "Data"} {
			out = quoteVar(out, item)
		}

		if err := json.Unmarshal([]byte(out), &plegEvent); err != nil {
			fmt.Println("error unmarshaling json: ", err)
		}

		fmt.Println("Data key:\n", plegEvent.Data)

	} else {
		fmt.Println("moving along, no containers here")
	}
}
