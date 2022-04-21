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
	fmt.Println("pleg-watcher v0.0.71")

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
