package main

import (
	"fmt"

	"github.com/rhdedgar/pleg-watcher/chroot"
	"github.com/rhdedgar/pleg-watcher/models"
	"github.com/rhdedgar/pleg-watcher/watcher"
)

func main() {
	var line models.LineInfo

	fmt.Println("pleg-watcher v0.0.6")
	line = make(chan string)

	watcher.PLEGWatch(&line)
	go chroot.SysCmd(models.ChrootChan, models.RuncChan)
	watcher.CheckOutput(line)
}
