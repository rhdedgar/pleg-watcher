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

package watcher_test

import (
	"encoding/json"
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/rhdedgar/pleg-watcher/watcher"
)

var _ = Describe("Watcher", func() {
	var (
		systemdOut       = `Dec 24 20:54:27 ip-123-12-12-123.us-east-1.compute.internal atomic-openshift-node[3322]: I0918 20:54:27.606604    3322 kubelet.go:1865] SyncLoop (PLEG): "django-psql-persistent-2-vzpl7_test-application(2e5c1fda-b089-11e8-aa1a-02ac3a1f9d61)", event: &pleg.PodLifecycleEvent{ID:"2e5c1fda-b089-11e8-aa1a-02ac3a1f9d61", Type:"ContainerStarted", Data:"c769a55aced1373b74df3bf1c2c283dcdc4bddbf914793b6b4eebd7f3ad6f718"}`
		systemdFormatted = `{"ID":"2e5c1fda-b089-11e8-aa1a-02ac3a1f9d61", "Type":"ContainerStarted", "Data":"c769a55aced1373b74df3bf1c2c283dcdc4bddbf914793b6b4eebd7f3ad6f718"}`
	)

	Describe("QuoteVar", func() {
		Context("Validate systemd output can be formatted as JSON", func() {
			It("Should correctly quote and format a string", func() {
				result := strings.SplitAfter(systemdOut, "&pleg.PodLifecycleEvent")[1]

				for _, item := range []string{"ID", "Type", "Data"} {
					result = QuoteVar(result, item)
				}

				Expect(result).To(Equal(systemdFormatted))
			})
		})
	})

	Describe("Format", func() {
		Context("Validate systemd output converted to golang struct", func() {
			It("Should correctly quote, format and structure a string into a struct", func() {
				var resultMock PLEGEvent

				err := json.Unmarshal([]byte(systemdFormatted), &resultMock)
				if err != nil {
					fmt.Println("Error unmarshalling plegEvent json: ", err)
				}

				result, err := Format(systemdOut)
				Expect(err).To(BeNil())

				Expect(result).To(Equal(resultMock))
			})
		})
	})
})
