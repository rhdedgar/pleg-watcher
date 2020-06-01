package config

import "os"

var (
	// SockPath represents the local pod's abstract Unix Domain socket.
	// The container-info container uses this to listen for RPC calls.
	// It listens for containerIDs, and returns information about that container.
	SockPath = os.Getenv("INFO_SOCKET")
)
