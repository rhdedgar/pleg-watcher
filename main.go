package main

import (
	"fmt"

	"github.com/rhdedgar/pleg-watcher/chroot"
	"github.com/rhdedgar/pleg-watcher/models"
	"github.com/rhdedgar/pleg-watcher/watcher"
)

func main() {
	var line models.LineInfo

	fmt.Println("pleg-watcher v0.0.7")
	line = make(chan string)

	// This gets set up first so that chroot doesn't interfere.
	watcher.PLEGWatch(&line)

	go chroot.SysCmd(models.ChrootChan, models.RuncChan)
	watcher.CheckOutput(line)
}
