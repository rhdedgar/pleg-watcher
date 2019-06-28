package models

var (
	// ChrootChan is the input channel for containerIDs
	ChrootChan = make(chan string)
	// ChrootOut is the output channel for crictl commands run by chroot.SysCmd
	ChrootOut = make(chan []byte)
	// RuncChan is the input channel for containerIDs used in runc inspection commands
	RuncChan = make(chan string)
	// RuncOut is the output channel for runc commands run by chroot.SysCmd
	RuncOut = make(chan []byte)
)
