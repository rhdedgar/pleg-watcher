package scanner

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/rhdedgar/pleg-watcher/api"
	"github.com/rhdedgar/pleg-watcher/clamav"
	"github.com/rhdedgar/pleg-watcher/cmd"
	"github.com/rhdedgar/pleg-watcher/sender"
)

// ContainerLayerScanner is the interface for all image containerLayerScanners.
//type ContainerLayerScanner interface {
// Inspect inspects and serves the image based on the ContainerLayerScannerOptions.
//	ClamScanner() error
//}

// ScanOutputs is a struct to hold all the scan outputs that needs to be served
type ScanOutputs struct {
	ScanReport     []byte
	HtmlScanReport []byte
	ScanResults    api.ScanResult
}

// defaultContainerLayerScanner is the default implementation of ContainerLayerScanner.
type defaultContainerLayerScanner struct {
	opts        cmd.ContainerLayerScannerOptions
	ScanOutputs ScanOutputs
}

// NewDefaultContainerLayerScanner provides a new default scanner.
func NewDefaultContainerLayerScanner(opts cmd.ContainerLayerScannerOptions) *defaultContainerLayerScanner {
	containerLayerScanner := &defaultContainerLayerScanner{
		opts: opts,
	}

	containerLayerScanner.ScanOutputs.ScanResults = api.ScanResult{
		APIVersion: api.DefaultResultsAPIVersion,
		Results:    []api.Result{},
	}

	return containerLayerScanner
}

// AcquireAndScan acquires and scans the image based on the ContainerLayerScannerOptions.
func (i *defaultContainerLayerScanner) AcquireAndScan() error {
	var (
		scanner  api.Scanner
		err      error
		filterFn api.FilesFilter
	)

	ctx := context.Background()

	scanner, err = clamav.NewScanner(i.opts.ClamSocket)
	if err != nil {
		fmt.Println("Error initializing clam:")
		return fmt.Errorf("failed to initialize clamav scanner: %v", err)
	}

	results, _, err := scanner.Scan(ctx, i.opts.ScanDir, filterFn)
	if err != nil {
		fmt.Printf("DEBUG: Unable to scan directory %q with ClamAV: %v", i.opts.ScanDir, err)
		return err
	}

	i.ScanOutputs.ScanResults.Results = append(i.ScanOutputs.ScanResults.Results, results...)
	if len(i.opts.PostResultURL) > 0 && len(i.ScanOutputs.ScanResults.Results) > 0 {
		fmt.Println("Infected files found, sending: ", i.ScanOutputs.ScanResults.Results)
		sender.SendClamData(i.ScanOutputs.ScanResults)
	}

	fmt.Println("The results slice: ", i.ScanOutputs.ScanResults.Results)
	if len(i.opts.OutFile) > 0 {
		if err := i.writeFile(i.ScanOutputs.ScanResults); err != nil {
			fmt.Printf("Error writing file: %v", err)
			return err
		}
	}

	return nil
}

func (i *defaultContainerLayerScanner) writeFile(scanResults api.ScanResult) error {
	outFile := i.opts.OutFile
	fmt.Printf("Writing results to %q ...", outFile)

	openFile, err := os.OpenFile(outFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	defer openFile.Close()

	jOut, err := json.Marshal(scanResults)
	if err != nil {
		return err
	}

	fileWrite, err := openFile.WriteString(string(jOut) + "\n")
	if err != nil {
		return err
	}
	fmt.Printf("wrote %d bytes\n", fileWrite)

	return nil
}
