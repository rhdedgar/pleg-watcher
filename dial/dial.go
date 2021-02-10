package dial

import (
	"fmt"
	"net/rpc"
	"strconv"
	"time"

	"github.com/rhdedgar/pleg-watcher/config"
)

var (
	minConDay = 0
)

// CallInfoSrv passes information about a container (arg1)
// to an RPC function (arg2), and returns information about that container.
func CallInfoSrv(containerID, functionName string) []byte {
	var reply []byte

	// functionName is the name of the RPC function to call from the container info server.
	functionName = "InfoSrv." + functionName

	client, err := rpc.DialHTTP("unix", config.SockPath)
	if err != nil {
		fmt.Println("Error dialing container info socket:", config.SockPath, err)
	}

	err = client.Call(functionName, &containerID, &reply)
	if err != nil {
		fmt.Println("Error calling server function", err)
	}

	if len(reply) > 0 {
		fmt.Println("A reply was returned")
	} else {
		fmt.Println("The reply was empty")
	}
	return reply
}

// GetContainerInfo queries for containers running on the host for more than 24 hours
// and populates a []btye slice of container IDs
func GetContainerInfo() []byte {
	var reply []byte

	// functionName is the name of the RPC function to call from the container info server.
	functionName := "InfoSrv.GetContainers"

	client, err := rpc.DialHTTP("unix", config.SockPath)
	if err != nil {
		fmt.Println("Error dialing container info socket:", config.SockPath, err)
	}

	curTime := time.Now()
	dayAgo := curTime.AddDate(0, 0, minConDay).String()

	err = client.Call(functionName, dayAgo, &reply)
	if err != nil {
		fmt.Println("Error calling server function", err)
	}

	if len(reply) > 0 {
		fmt.Println("A reply was returned")
	} else {
		fmt.Println("The reply was empty")
	}

	return reply
}

func init() {
	var err error

	minConDay, err = strconv.Atoi(config.MinConDay)
	if err != nil {
		fmt.Println(err)
	}
}
