package main

import (
	"fmt"

	"github.com/rhdedgar/pleg-watcher/chroot"
	"github.com/rhdedgar/pleg-watcher/models"
	"github.com/rhdedgar/pleg-watcher/watcher"
)

func main() {
	var line models.LineInfo

	fmt.Println("pleg-watcher v0.9.1")
	line = make(chan string)

	go chroot.SysCmd(models.ChrootChan, models.RuncChan)
	go watcher.PLEGWatch(&line)
	watcher.CheckOutput(line)
}
