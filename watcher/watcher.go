package watcher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/coreos/go-systemd/sdjournal"
	"github.com/rhdedgar/pleg-watcher/container"
	"github.com/rhdedgar/pleg-watcher/models"
)

type PLEGEvent struct {
	ID   string `json:"ID"`
	Type string `json:"Type"`
	Data string `json:"Data"`
}

type PLEGBuffer struct {
	bLine bytes.Buffer
}

func quoteVar(s string, r string) string {
	return strings.Replace(s, r, "\""+r+"\"", 1)
}

func CheckOutput(line <-chan string) {
	var plegEvent PLEGEvent

	for {
		select {
		case inputStr := <-line:
			fmt.Println(inputStr)
			//fmt.Println(line)

			if strings.Contains(inputStr, "ContainerStarted") {
				// Gather only the unquoted json of the PLEG Event
				out := strings.SplitAfter(inputStr, "&pleg.PodLifecycleEvent")[1]

				// Quote the json so it can be Unmarshaled into a struct
				for _, item := range []string{"ID", "Type", "Data"} {
					out = quoteVar(out, item)
				}

				if err := json.Unmarshal([]byte(out), &plegEvent); err != nil {
					fmt.Println("error unmarshaling json: ", err)
				}

				fmt.Println("Data key:\n", plegEvent.Data)

				container.ProcessContainer(plegEvent.Data)
			}
			fmt.Println("Not a creation event, skipping")
			//return ""
		}
	}
}

func PLEGWatch(out *models.LineInfo) {
	path := os.Getenv("JOURNAL_PATH")

	fmt.Println(path)

	jrcfg := sdjournal.JournalReaderConfig{
		NumFromTail: 10,
		Path:        path,
		//		Matches: []sdjournal.Match{
		//			{
		//				Field: sdjournal.SD_JOURNAL_FIELD_SYSLOG_IDENTIFIER,
		//				Value: "atomic-openshift-node",
		//			},
		//		},
	}

	jr, err := sdjournal.NewJournalReader(jrcfg)
	if err != nil {
		log.Printf("[ERROR] journal: %v", err)
		return
	}
	defer jr.Close()

	fmt.Println("=== Watching journal ===")

	until := make(chan time.Time)

	if err := jr.Follow(until, out); err != nil {
		log.Fatalf("Could not read from journal: %s", err)
	}
	fmt.Println("done following")
}
