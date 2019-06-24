package models

var (
	// ChrootChan is the input channel for containerIDs
	ChrootChan = make(chan string)
	// ChrootOut is the output channel for commands run by chroot.SysCmd
	ChrootOut = make(chan []byte)
)
