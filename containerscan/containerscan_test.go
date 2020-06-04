package containerscan_test

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rhdedgar/pleg-watcher/config"
	. "github.com/rhdedgar/pleg-watcher/containerscan"
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

var _ = Describe("Containerscan", func() {
	var (
		testContainerID  = "testcontainerid"
		rootFS           = "/var/lib/containers/storage/overlay/6e98981bb10fb2404f039c04d999de9c3a23e9362dfbca96b6d9385c8b052f30/merged"
		procPIDMountInfo = `overlay /var/lib/containers/storage/overlay/5db062515b858b32faf8b28fa8f478cb8934da8a98ff88f6229ec390ff2575c2/merged overlay rw,context="system_u:object_r:container_file_t:s0:c22,c23",relatime,lowerdir=/var/lib/containers/storage/overlay/l/KGVVTLSJFMVP3COI6MMEZ3J6SZ:/var/lib/containers/storage/overlay/l/BF37XRKHAR45TXDVLROK4WUP62:/var/lib/containers/storage/overlay/l/FLF3UD5KETBQWMJKX3S4WUMICU:/var/lib/containers/storage/overlay/l/RNICQWYK23FI3NPIBEXJBVAAP3,upperdir=/var/lib/containers/storage/overlay/5db062515b858b32faf8b28fa8f478cb8934da8a98ff88f6229ec390ff2575c2/diff,workdir=/var/lib/containers/storage/overlay/5db062515b858b32faf8b28fa8f478cb8934da8a98ff88f6229ec390ff2575c2/work 0 0`
		subString        = "lowerdir=/var/lib/containers/storage/overlay/l/KGVVTLSJFMVP3COI6MMEZ3J6SZ:/var/lib/containers/storage/overlay/l/BF37XRKHAR45TXDVLROK4WUP62:/var/lib/containers/storage/overlay/l/FLF3UD5KETBQWMJKX3S4WUMICU:/var/lib/containers/storage/overlay/l/RNICQWYK23FI3NPIBEXJBVAAP3,upperdir"
		subSlice         = []string{
			"/var/lib/containers/storage/overlay/l/KGVVTLSJFMVP3COI6MMEZ3J6SZ",
			"/var/lib/containers/storage/overlay/l/BF37XRKHAR45TXDVLROK4WUP62",
			"/var/lib/containers/storage/overlay/l/FLF3UD5KETBQWMJKX3S4WUMICU",
			"/var/lib/containers/storage/overlay/l/RNICQWYK23FI3NPIBEXJBVAAP3",
		}
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

	Describe("CustSplit", func() {
		Context("Validate string splitting works as expected", func() {
			It("Should correctly split strings", func() {
				items := strings.Split(subString, ":")
				Expect(len(items)).To(Equal(4))

				layer := CustSplit(items[0], ",", 0)
				Expect(layer).To(Equal("lowerdir=/var/lib/containers/storage/overlay/l/KGVVTLSJFMVP3COI6MMEZ3J6SZ"))

				result := CustSplit(layer, "=", 1)
				Expect(result).To(Equal("/var/lib/containers/storage/overlay/l/KGVVTLSJFMVP3COI6MMEZ3J6SZ"))
			})
		})
	})

	Describe("CustReg", func() {
		Context("Validate regex works as expected", func() {
			It("Should extract the substring between lowerdir and upperdir", func() {
				result := CustReg(procPIDMountInfo, `lowerdir=(.*),upperdir`)
				Expect(result[0]).To(Equal(subString))
			})
		})
	})

	Describe("MountOverlayFS", func() {
		Context("Validate mounting functions work as expected", func() {
			It("Should create dirs and attempt to mount them as OverlayFS", func() {
				fmt.Println(subSlice)
				err := os.MkdirAll("/mnt/", os.ModePerm)
				if err != nil {
					fmt.Println(err)
				}

				for _, l := range subSlice {
					err := os.MkdirAll("/host/"+l, os.ModePerm)
					if err != nil {
						fmt.Println(err)
					}
				}

				result, err := MountOverlayFS(subSlice, testContainerID)
				Expect(result).To(Equal("/mnt/" + testContainerID))
				Expect(err).To(Not(BeNil()))
			})
		})
	})

	Describe("GetLayerInfo", func() {
		Context("Validate we can read from /host/proc/<PID>/mountinfo", func() {
			It("Should read mount point data from the mountinfo file", func() {
				result, err := GetLayerInfo("./example_mountinfo")
				Expect(result).To(Equal(procPIDMountInfo))
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("GetRootFS", func() {
		Context("Validate we can read runc state output", func() {
			It("Should find and return the rootfs field of the example runc state output", func() {
				result, err := GetRootFS("./runc_state_example.json")
				Expect(result).To(Equal(rootFS))
				Expect(err).To(BeNil())
			})
		})
	})
})
