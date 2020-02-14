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

	fmt.Println("pleg-watcher v0.0.22, v0.0.17 was successful. testing SD_JOURNAL_FIELD_SYSTEMD_UNIT kubelet.go filtering now")
	line = make(chan string)

	// This gets set up first so that chroot doesn't interfere with libraries loading.
	// It runs in a separate goroutine as it provides the actual watcher functionality.
	go watcher.PLEGWatch(&line)
	time.Sleep(5 * time.Second)

	// Another goroutine to wait for container IDs, gather info about the container, and return it.
	go chroot.SysCmd(models.ChrootChan, models.RuncChan)

	// Continuously filter out irrelevant kubelet output.
	watcher.CheckOutput(line)
}
