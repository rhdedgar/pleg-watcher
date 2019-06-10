package chroot

import (
	"os"
	"syscall"
)

// Chroot provides chroot access to the mounted host filesystem
func Chroot(chrPath string) (func() error, error) {
	root, err := os.Open("/")
	if err != nil {
		return nil, err
	}

	if err := syscall.Chroot(chrPath); err != nil {
		root.Close()
		return nil, err
	}

	return func() error {
		defer root.Close()
		if err := root.Chdir(); err != nil {
			return err
		}
		return syscall.Chroot(".")
	}, nil
}
