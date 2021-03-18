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

// CallInfoSrv passes a container ID (arg1)
// to an RPC function (arg2), and returns information about that container.
func CallInfoSrv(containerID, functionName string) []byte {
	var reply []byte

	if containerID == "" {
		fmt.Println("containerID cannot be empty.")
		return reply
	}

	// functionName is the name of the RPC function to call from the container info server.
	functionName = "InfoSrv." + functionName

	fmt.Printf("Calling %v for %v\n", functionName, containerID)

	client, err := rpc.DialHTTP("unix", config.SockPath)
	if err != nil {
		fmt.Println("Error dialing container info socket: ", config.SockPath, err)
		fmt.Println("Skipping this container.")
		return reply
	}

	err = client.Call(functionName, &containerID, &reply)
	if err != nil {
		fmt.Println("Error calling server function: ", functionName, err)
	}

	if len(reply) > 0 {
		fmt.Printf("A reply was returned from %v: %v\n", functionName, string(reply))
	} else {
		fmt.Printf("The reply from %v was empty.\n", functionName)
	}
	return reply
}

// GetActiveContainers queries for containers running on the host for more than 24 hours
// and populates a []btye slice of container IDs
func GetActiveContainers() []byte {
	var reply []byte

	// functionName is the name of the RPC function to call from the container info server.
	functionName := "InfoSrv.GetContainers"

	client, err := rpc.DialHTTP("unix", config.SockPath)
	if err != nil {
		fmt.Println("Error dialing container info socket: ", config.SockPath, err)
		return reply
	}

	curTime := time.Now()
	dayAgo := curTime.AddDate(0, 0, minConDay).String()

	err = client.Call(functionName, dayAgo, &reply)
	if err != nil {
		fmt.Println("Error calling server function: ", functionName, err)
	}

	if len(reply) > 0 {
		fmt.Printf("A reply was returned from %v: %v\n", functionName, string(reply))
	} else {
		fmt.Printf("The reply from %v was empty.", functionName)
	}

	return reply
}

func init() {
	if config.ScheduledScan == "true" {
		if config.MinConDay == "" {
			return
		}

		var err error

		minConDay, err = strconv.Atoi(config.MinConDay)
		if err != nil {
			fmt.Println(err)
		}
	}
}
