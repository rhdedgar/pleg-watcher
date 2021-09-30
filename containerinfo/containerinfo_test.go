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

package containerinfo_test

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/rpc"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/rhdedgar/pleg-watcher/config"
	. "github.com/rhdedgar/pleg-watcher/containerinfo"
)

// InfoSrv is the base type that needs to be exported for RPC to work.
type InfoSrv struct {
}

// GetContainerInfo is the RPC-exported method that returns docker or crictl info about a container.
func (g InfoSrv) GetContainerInfo(containerID *string, reply *[]byte) error {
	crictlFilePath := "./crictl_inspect_example.json"

	*reply = loadExample(crictlFilePath)

	return nil
}

// GetRuncInfo is the RPC-exported method that returns runc info about a container.
func (g InfoSrv) GetRuncInfo(containerID *string, reply *[]byte) error {
	runcFilePath := "./runc_state_example.json"

	*reply = loadExample(runcFilePath)

	return nil
}

// loadExmple reads an example file path string, and returns its contents as a byte string.
func loadExample(filePath string) []byte {
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error loading secrets json: ", err)
		return []byte{}
	}

	return fileBytes
}

var _ = Describe("Containerinfo", func() {
	var (
		testContainerID = "testcontainerid"
	)

	config.SockPath = "@testSock"
	UseDocker = false

	go func() {
		InfoSrv := new(InfoSrv)

		rpc.Register(InfoSrv)
		rpc.HandleHTTP()

		listener, err := net.Listen("unix", config.SockPath)
		if err != nil {
			fmt.Println("Error starting listener:", err)
		}

		http.Serve(listener, nil)
	}()

	Describe("ProcessContainer", func() {
		Context("Validate container data gets processed", func() {
			It("Should detect the openshift namespace from the container output", func() {
				err := ProcessContainer(testContainerID)
				exErr := fmt.Errorf("Container is in openshift-* namespace, skipping")

				Expect(err).To(Equal(exErr))
			})
		})
	})
})
