package cmd

import (
	"fmt"

	"os"
)

// DefaultClamSocketLocation is the default location of the clamd socket on the host
var DefaultClamSocketLocation = os.Getenv("CLAM_SOCKET")

// MultiStringVar is implementing flag.Value
type MultiStringVar struct {
	Values []string
}

// Set appends values to a set
func (sv *MultiStringVar) Set(s string) error {
	sv.Values = append(sv.Values, s)
	return nil
}

func (sv *MultiStringVar) String() string {
	return fmt.Sprintf("%v", sv.Values)
}

// ManagedScannerOptions is the main scanner implementation and holds the configuration
// for a clam scanner.
type ManagedScannerOptions struct {
	// ScanDir is the name of the directory to be scanned.
	ScanDir string
	// ScanResultsDir is the directory that will contain the results of the scan
	ScanResultsDir string
	// ClamSocket is the location of clamav socket file
	ClamSocket string
	// PostResultURL represents an URL where the image-inspector should post the results of
	// the scan.
	PostResultURL string
	// OutFile is the name of the file on disk to write
	OutFile string
}

// NewDefaultManagedScannerOptions provides a new ManagedScannerOptions with default values.
func NewDefaultManagedScannerOptions() *ManagedScannerOptions {
	return &ManagedScannerOptions{
		ScanDir:        "",
		ScanResultsDir: "",
		ClamSocket:     DefaultClamSocketLocation,
		PostResultURL:  "",
		OutFile:        "",
	}
}

// Validate performs validation on the field settings.
func (i *ManagedScannerOptions) Validate() error {
	if len(i.ScanDir) == 0 {
		return fmt.Errorf("a directory to scan must be specified")
	}
	if len(i.ScanResultsDir) > 0 {
		fi, err := os.Stat(i.ScanResultsDir)
		if err == nil && !fi.IsDir() {
			return fmt.Errorf("scan-results-dir %q is not a directory", i.ScanResultsDir)
		}
	}
	return nil
}
