package dial

import (
	"fmt"
	"net/rpc"

	"github.com/rhdedgar/pleg-watcher/config"
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
