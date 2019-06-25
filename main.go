package main

import (
	"github.com/rhdedgar/pleg-watcher/chroot"
	"github.com/rhdedgar/pleg-watcher/models"
	"github.com/rhdedgar/pleg-watcher/watcher"
)

func main() {
	var line models.LineInfo
	line = make(chan string)

	go chroot.SysCmd(models.ChrootChan)
	go watcher.PLEGWatch(&line)
	watcher.CheckOutput(line)
}
