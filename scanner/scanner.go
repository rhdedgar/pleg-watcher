package scanner

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rhdedgar/pleg-watcher/api"
	"github.com/rhdedgar/pleg-watcher/clamav"
	"github.com/rhdedgar/pleg-watcher/cmd"
)

// ContainerLayerScanner is the interface for all image containerLayerScanners.
type ContainerLayerScanner interface {
	// Inspect inspects and serves the image based on the ContainerLayerScannerOptions.
	ClamScanner() error
}

// scanOutputs is a struct to hold all the scan outputs that needs to be served
type scanOutputs struct {
	ScanReport     []byte
	HtmlScanReport []byte
	ScanResults    api.ScanResult
}

// defaultContainerLayerScanner is the default implementation of ContainerLayerScanner.
type defaultContainerLayerScanner struct {
	opts        cmd.ContainerLayerScannerOptions
	scanOutputs scanOutputs
}

// NewDefaultContainerLayerScanner provides a new default scanner.
func NewDefaultContainerLayerScanner(opts cmd.ContainerLayerScannerOptions) ContainerLayerScanner {
	containerLayerScanner := &defaultContainerLayerScanner{
		opts: opts,
	}

	containerLayerScanner.scanOutputs.ScanResults = api.ScanResult{
		APIVersion: api.DefaultResultsAPIVersion,
		Results:    []api.Result{},
	}

	return containerLayerScanner
}

// Inspect inspects and serves the image based on the ImageInspectorOptions.
func (i *defaultContainerLayerScanner) ClamScanner() error {
	err := i.acquireAndScan()
	if err != nil {
		return fmt.Errorf("failed to acquire and scan: %v", err.Error())
	}

	return err
}

// AcquireAndScan acquires and scans the image based on the ContainerLayerScannerOptions.
func (i *defaultContainerLayerScanner) acquireAndScan() error {
	var (
		scanner  api.Scanner
		err      error
		filterFn api.FilesFilter
	)

	ctx := context.Background()

	scanner, err = clamav.NewScanner(i.opts.ClamSocket)
	if err != nil {
		return fmt.Errorf("failed to initialize clamav scanner: %v", err)
	}
	results, _, err := scanner.Scan(ctx, i.opts.ScanDir, filterFn)
	if err != nil {
		log.Printf("DEBUG: Unable to scan directory %q with ClamAV: %v", i.opts.ScanDir, err)
		return err
	}
	i.scanOutputs.ScanResults.Results = append(i.scanOutputs.ScanResults.Results, results...)

	if len(i.opts.PostResultURL) > 0 {
		if err := i.postResults(i.scanOutputs.ScanResults); err != nil {
			log.Printf("Error posting results: %v", err)
			return err
		}
	}

	if len(i.opts.OutFile) > 0 {
		if err := i.writeFile(i.scanOutputs.ScanResults); err != nil {
			log.Printf("Error writing file: %v", err)
			return err
		}
	}

	return nil
}

func (i *defaultContainerLayerScanner) postResults(scanResults api.ScanResult) error {
	url := i.opts.PostResultURL
	log.Printf("Posting results to %q ...", url)
	resultJSON, err := json.Marshal(scanResults)
	if err != nil {
		return err
	}
	client := http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewReader(resultJSON))
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	log.Printf("DEBUG: Success: %v", resp)
	return nil
}

func (i *defaultContainerLayerScanner) writeFile(scanResults api.ScanResult) error {
	outFile := i.opts.OutFile
	log.Printf("Writing results to %q ...", outFile)

	openFile, err := os.Create(outFile)
	if err != nil {
		return err
	}

	defer openFile.Close()

	jOut, err := json.Marshal(scanResults)
	if err != nil {
		return err
	}

	fileWrite, err := openFile.WriteString(string(jOut))
	if err != nil {
		return err
	}
	fmt.Printf("wrote %d bytes\n", fileWrite)

	return nil
}
