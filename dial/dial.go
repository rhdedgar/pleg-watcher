package dial

import (
	"fmt"
	"net/rpc"
	"strconv"
	"time"

	"github.com/rhdedgar/pleg-watcher/config"
)

var (
	client    *rpc.Client
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

	err := client.Call(functionName, &containerID, &reply)
	if err != nil {
		fmt.Println("Error calling server function: ", functionName, err)
	}

	if len(reply) > 0 {
		// fmt.Printf("A reply was returned from %v: %v\n", functionName, string(reply))
		fmt.Printf("A reply was returned from %v.\n", functionName)
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

	curTime := time.Now()
	dayAgo := curTime.AddDate(0, 0, minConDay).String()

	err := client.Call(functionName, dayAgo, &reply)
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
	var err error

	for i := 0; i <= 5; i++ {
		if i >= 5 {
			fmt.Println("Error: dialing containerinfo socket has reached maximum number of allowed retry attempts.")
			break
		}

		client, err = rpc.Dial("unix", config.SockPath)
		if err != nil {
			fmt.Println("Error dialing container info socket: ", config.SockPath, err)
			fmt.Printf("Waiting another %v seconds before trying again. \n", i)
			time.Sleep(time.Duration(i) * time.Second)
			continue
		}
		break
	}

	if config.ScheduledScan == "true" {
		if config.MinConDay == "" {
			return
		}

		minConDay, err = strconv.Atoi(config.MinConDay)
		if err != nil {
			fmt.Println("Error parsing minimum container age environment variable: ", err)
		}
	}
}
