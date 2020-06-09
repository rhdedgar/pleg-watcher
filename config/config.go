package config

import "os"

var (
	// SockPath represents the local pod's abstract Unix Domain socket.
	// The container-info container uses this to listen for RPC calls.
	// It listens for containerIDs, and returns information about that container.
	SockPath = os.Getenv("INFO_SOCKET")
	// DockerURL is the URL path of the server that listens for POSTed JSON
	// data consisting of `docker inspect` output.
	DockerURL = os.Getenv("DOCKER_LOG_URL")
	// CrioURL is the URL path of the server that listens for POSTed JSON
	// data consisting of `crictl inspect` output.
	CrioURL = os.Getenv("CRIO_LOG_URL")
	// ClamURL is the URL path of the server that listens for POSTed JSON
	// data consisting of positive clam scan output.
	ClamURL = os.Getenv("CLAM_LOG_URL")
)
