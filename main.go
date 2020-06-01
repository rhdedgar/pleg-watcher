package main

import (
	"fmt"
	"time"

	"github.com/rhdedgar/pleg-watcher/models"
	"github.com/rhdedgar/pleg-watcher/watcher"
)

func main() {
	var line models.LineInfo

	fmt.Println("pleg-watcher v0.0.47")
	line = make(chan string)

	// This runs in a separate goroutine as it provides the actual watcher functionality.
	go watcher.PLEGWatch(&line)
	time.Sleep(5 * time.Second)

	// Continuously filter out irrelevant kubelet output.
	watcher.CheckOutput(line)
}
