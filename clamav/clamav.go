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

package clamav

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/openshift/clam-scanner/pkg/clamav"
	"github.com/rhdedgar/pleg-watcher/api"
)

// ScannerName is a string of the name of the scanner
const ScannerName = "clamav"

// ClamScanner is a structure of two vars
// Socket is the location of the clamav socket.
// clamd is a new clamav ClamdSession
type ClamScanner struct {
	Socket string

	clamd clamav.ClamdSession
}

var _ api.Scanner = &ClamScanner{}

// NewScanner initializes a new clamd session
func NewScanner(socket string) (api.Scanner, error) {
	clamSession, err := clamav.NewClamdSession(socket, true)
	if err != nil {
		fmt.Println("NewScanner error")
		return nil, err
	}
	return &ClamScanner{
		Socket: socket,
		clamd:  clamSession,
	}, nil
}

// Scan will scan the image
func (s *ClamScanner) Scan(ctx context.Context, path string, filter api.FilesFilter) ([]api.Result, interface{}, error) {
	scanResults := []api.Result{}
	// Useful for debugging
	scanStarted := time.Now()

	defer func() {
		fmt.Printf("clamav scan took %ds (%d problems found)\n", int64(time.Since(scanStarted).Seconds()), len(scanResults))
	}()

	if err := s.clamd.ScanPath(ctx, path, clamav.FilterFiles(filter)); err != nil {
		return nil, nil, err
	}

	s.clamd.WaitTillDone()
	defer s.clamd.Close()

	clamResults := s.clamd.GetResults()

	for _, r := range clamResults.Files {
		r := api.Result{
			Name:           ScannerName,
			ScannerVersion: "0.101.2",
			Timestamp:      scanStarted,
			Reference:      fmt.Sprintf("file://%s", strings.TrimPrefix(r.Filename, path)),
			Description:    r.Result,
		}
		scanResults = append(scanResults, r)
	}
	fmt.Println("clamav results: ", scanResults)
	return scanResults, nil, nil
}

// Name returns the const ScannerName
func (s *ClamScanner) Name() string {
	return ScannerName
}
