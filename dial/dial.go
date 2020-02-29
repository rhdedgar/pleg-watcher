package dial

import (
	"fmt"
	"net/rpc"
	"os"
)

var (
	// sockPath is the path to the Unix Domain Socket that listens for containerIDs and returns container info.
	sockPath = os.Getenv("INFO_SOCK_PATH")
)

// InfoSrv passes information about a container (arg1)
// to an RPC function (arg2), and returns information about that container.
func InfoSrv(containerID, functionName string) []byte {
	var reply []byte

	// functionName is the name of the RPC function to call from the container info server.
	functionName = "InfoSrv." + functionName

	client, err := rpc.DialHTTP("unix", sockPath)
	if err != nil {
		fmt.Println("Error dialing container info socket:", err)
	}

	err = client.Call(functionName, &containerID, &reply)
	if err != nil {
		fmt.Println("Error calling server function", err)
	}

	fmt.Printf("InfoSrv got '%v'\n", reply)

	return reply
}
