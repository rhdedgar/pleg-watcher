package watcher

import (
	"encoding/json"
	"fmt"
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

// QuoteVar will quote the first occurrence of substring r found in original string s.
func QuoteVar(s, r string) string {
	return strings.Replace(s, r, "\""+r+"\"", 1)
}

// Format converts systemd output into a usable, JSON-compatible go struct
func Format(inputStr string) (PLEGEvent, error) {
	var plegEvent PLEGEvent

	if strings.Contains(inputStr, "ContainerStarted") {
		fmt.Println("Found container started event", inputStr)

		// Gather only the unquoted json of the PLEG Event.
		out := strings.SplitAfter(inputStr, "&pleg.PodLifecycleEvent")[1]

		// Quote the json so it can be Unmarshaled into a struct
		for _, item := range []string{"ID", "Type", "Data"} {
			out = QuoteVar(out, item)
		}

		if err := json.Unmarshal([]byte(out), &plegEvent); err != nil {
			return plegEvent, fmt.Errorf("Error unmarshalling plegEvent json: %v\n", err)
		}

		if plegEvent == (PLEGEvent{}) {
			return plegEvent, fmt.Errorf("The PLEGEvent structure is empty. Journalctl hyperkube may have changed.\n")
		}
	}
	return plegEvent, nil
}

// CheckOutput filters through new systemd lines as they're received from a string channel.
func CheckOutput(line <-chan string) {
	for {
		select {
		case inputStr := <-line:
			plegEvent, err := Format(inputStr)
			if err != nil {
				fmt.Println("Error returned from Format:", err)
				continue
			}

			if err := containerinfo.ProcessContainer(plegEvent.Data); err != nil {
				fmt.Println("Error returned from ProcessContainer: ", err)
			}
		}
	}
}

// PLEGWatch initalizes a new JournalReader and starts following systemd output
func PLEGWatch(out *models.LineInfo) {
	path := os.Getenv("JOURNAL_PATH")

	fmt.Println("Journal path:", path)

	jcfg := sdjournal.JournalReaderConfig{
		NumFromTail: uint64(1),
		Path:        path,
		Matches: []sdjournal.Match{
			{
				Field: sdjournal.SD_JOURNAL_FIELD_SYSLOG_IDENTIFIER,
				Value: "hyperkube",
			},
		},
	}

	r, err := sdjournal.NewJournalReader(jcfg)
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
