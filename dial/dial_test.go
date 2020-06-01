package dial_test

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/rpc"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/rhdedgar/pleg-watcher/config"
	. "github.com/rhdedgar/pleg-watcher/dial"
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

var _ = Describe("Dial", func() {
	var (
		testContainerID = "testcontainerid"
	)

	config.SockPath = "@testSock"

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

	Describe("PostDockerPodLog", func() {
		Context("Validate the ability to dial and receive data from an RPC server", func() {
			It("Should receive a reply from the RPC server", func() {
				reply := CallInfoSrv(testContainerID, "GetContainerInfo")

				Expect(reply).ToNot(BeNil())
			})
		})
	})
})
