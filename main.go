package main

import (
	"fmt"
	"time"

	"github.com/rhdedgar/pleg-watcher/config"
	"github.com/rhdedgar/pleg-watcher/models"
	"github.com/rhdedgar/pleg-watcher/scheduler"
	"github.com/rhdedgar/pleg-watcher/watcher"
)

var (
	activeScan    = config.ActiveScan
	scheduledScan = config.ScheduledScan
)

func main() {
	fmt.Println("pleg-watcher v0.0.61")

	if activeScan != "" {
		var line models.LineInfo

		line = make(chan string)

		fmt.Println("Starting active scanner")

		// This runs in a separate goroutine as it provides the actual watcher functionality.
		go watcher.PLEGWatch(&line)
		time.Sleep(5 * time.Second)

		// Continuously filter out irrelevant kubelet output.
		watcher.CheckOutput(line)
	} else if scheduledScan != "" {
		fmt.Println("Starting scheduled container scans.")
		scheduler.ScheduledContainerScan()

		fmt.Println("Starting scheduled host scan.")
		scheduler.ScheduledHostScan()

		fmt.Println("Scheduled scanning has completed")
	}
}
