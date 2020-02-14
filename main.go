package main

import (
	"fmt"
	"time"

	"github.com/rhdedgar/pleg-watcher/chroot"
	"github.com/rhdedgar/pleg-watcher/models"
	"github.com/rhdedgar/pleg-watcher/watcher"
)

func main() {
	var line models.LineInfo

	fmt.Println("pleg-watcher v0.0.18, v0.0.17 was successful. testing SYSTEMD_UNIT kubelet filtering now")
	line = make(chan string)

	// This gets set up first so that chroot doesn't interfere with libraries loading.
	go watcher.PLEGWatch(&line)
	time.Sleep(5 * time.Second)

	go chroot.SysCmd(models.ChrootChan, models.RuncChan)
	watcher.CheckOutput(line)
}
