package watcher

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/coreos/go-systemd/v22/sdjournal"
	"github.com/rhdedgar/pleg-watcher/containerinfo"
	"github.com/rhdedgar/pleg-watcher/models"
)

// PLEGEvent represents relevant data from Kubernetes Pod Lifecycle Event Generator messages.
type PLEGEvent struct {
	ID   string `json:"ID"`
	Type string `json:"Type"`
	Data string `json:"Data"`
}

func quoteVar(s string, r string) string {
	return strings.Replace(s, r, "\""+r+"\"", 1)
}

// CheckOutput filters through new systemd lines as they're received from a string channel.
func CheckOutput(line <-chan string) {
	var plegEvent PLEGEvent

	for {
		select {
		case inputStr := <-line:
			if strings.Contains(inputStr, "ContainerStarted") {
				// Gather only the unquoted json of the PLEG Event.
				out := strings.SplitAfter(inputStr, "&pleg.PodLifecycleEvent")[1]

				// Quote the json so it can be Unmarshaled into a struct
				for _, item := range []string{"ID", "Type", "Data"} {
					out = quoteVar(out, item)
				}

				if err := json.Unmarshal([]byte(out), &plegEvent); err != nil {
					fmt.Println("error unmarshalling json: ", err)
				}
				containerinfo.ProcessContainer(plegEvent.Data)
			}
		}
	}
}

// PLEGWatch initalizes a new JournalReader and starts following systemd output
func PLEGWatch(out *models.LineInfo) {
	path := os.Getenv("JOURNAL_PATH")

	fmt.Println(path)

	r, err := sdjournal.NewJournalReader(sdjournal.JournalReaderConfig{
		Since: time.Duration(time.Millisecond),
		Path:  path,
		Matches: []sdjournal.Match{
			{
				Field: sdjournal.SD_JOURNAL_FIELD_SYSLOG_IDENTIFIER,
				Value: "kubelet",
			},
		},
	})

	if err != nil {
		fmt.Printf("Error opening journal: %v\n", err)
	}

	if r == nil {
		fmt.Println("Error: got a nil reader.")
	}

	defer r.Close()

	fmt.Println("=== Watching journal ===")

	until := make(chan time.Time)

	if err := r.Follow(until, out); err != nil {
		fmt.Printf("Could not read from journal: %s\n", err)
	}

}

func isEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}
