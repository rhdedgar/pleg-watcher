package main

import (
	"github.com/rhdedgar/pleg-watcher/models"
	"github.com/rhdedgar/pleg-watcher/watcher"
)

func main() {
	var line models.LineInfo
	line = make(chan string)

	go watcher.PLEGWatch(&line)
	watcher.CheckOutput(line)
}
