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

package scanner_test

import (
	"bufio"
	"encoding/json"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gotest.tools/assert/cmp"

	"github.com/rhdedgar/pleg-watcher/api"
	"github.com/rhdedgar/pleg-watcher/cmd"
	. "github.com/rhdedgar/pleg-watcher/scanner"
)

var _ = Describe("Scanner", func() {
	var (
		curTime     = time.Now()
		outFile     = "output.txt"
		scanResults = &api.ScanResult{
			APIVersion:  "v1alpha",
			ContainerID: "5968c9dde7f9a837021d4b855a3ae75528a4322f9478ab166dd58b9f2a4b8a66",
			ImageID:     "sha256:eea8df720e7ef1f0b3333532a5407addf7bc2c7fe211ecdc7685ec9fd367f57a",
			ImageName:   "",
			NameSpace:   "testNS",
			PodName:     "testPod",
			Results: []api.Result{{
				Name:           "clamav",
				ScannerVersion: "0.99.99",
				Timestamp:      curTime,
				Reference:      "file://home/jboss/testdir/phishy.php",
				Description:    "Phish.Phishy.A(OpenShift).UNOFFICIAL FOUND",
			}},
		}
	)

	Describe("NewDefaultManagedScanner", func() {
		Context("Validate new container layer scanners get created", func() {
			It("Should return a new *defaultManagedScanner", func() {
				scannerOptions := cmd.NewDefaultManagedScannerOptions()

				scanner := NewDefaultManagedScanner(*scannerOptions)

				Expect(scanner.ScanOutputs.ScanResults.APIVersion).To(Equal("v1alpha"))
			})
		})
	})

	Describe("WriteFile", func() {
		Context("Validate scan results can be written to disk", func() {
			It("Should write an api.ScanResult to a file on disk", func() {
				var result api.ScanResult

				scannerOptions := cmd.NewDefaultManagedScannerOptions()
				scannerOptions.OutFile = outFile

				scanner := NewDefaultManagedScanner(*scannerOptions)
				scanner.ScanOutputs.ScanResults = *scanResults

				// use WriteFile, check that it exists, open it
				err := scanner.WriteFile(scanner.ScanOutputs.ScanResults)
				Expect(err).To(BeNil())

				_, err = os.Stat(outFile)
				Expect(err).To(BeNil())

				file, err := os.Open(outFile)
				Expect(err).To(BeNil())

				// defer stack is LIFO
				defer file.Close()
				defer os.Remove(outFile)

				// reading one line from outFile
				bufScanner := bufio.NewScanner(file)
				bufScanner.Scan()

				err = bufScanner.Err()
				Expect(err).To(BeNil())

				// finally, compare with our example struct
				err = json.Unmarshal(bufScanner.Bytes(), &result)
				Expect(err).To(BeNil())

				Expect(cmp.Equal(result, *scanResults))
			})
		})
	})
})
