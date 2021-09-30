/*
Copyright 2019 Doug Edgar.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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

// ScanOutputs is a struct to hold all the scan outputs that needs to be served
type ScanOutputs struct {
	ScanReport     []byte
	HtmlScanReport []byte
	ScanResults    api.ScanResult
}

// defaultManagedScanner is the default implementation of ManagedScanner.
type defaultManagedScanner struct {
	opts        cmd.ManagedScannerOptions
	ScanOutputs ScanOutputs
}

// NewDefaultManagedScanner provides a new default scanner.
func NewDefaultManagedScanner(opts cmd.ManagedScannerOptions) *defaultManagedScanner {
	ManagedScanner := &defaultManagedScanner{
		opts: opts,
	}

	ManagedScanner.ScanOutputs.ScanResults = api.ScanResult{
		APIVersion: api.DefaultResultsAPIVersion,
		Results:    []api.Result{},
	}

	return ManagedScanner
}

// AcquireAndScan acquires and scans the image based on the ManagedScannerOptions.
func (i *defaultManagedScanner) AcquireAndScan() error {
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
		if err := i.WriteFile(i.ScanOutputs.ScanResults); err != nil {
			fmt.Printf("Error writing file: %v", err)
			return err
		}
	}

	return nil
}

func (i *defaultManagedScanner) WriteFile(scanResults api.ScanResult) error {
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
