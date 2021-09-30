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
	// SkipNamespaces is a comma separated list of specific namespaces
	// to exclude from scan operations.
	SkipNamespaces = os.Getenv("SKIP_NAMESPACES")
	// SkipNamespacePrefixes is comma separated list of namespace prefixes
	// to exclude from scanning. e.g. ("openshift-")
	SkipNamespacePrefixes = os.Getenv("SKIP_NAMESPACE_PREFIXES")
	// ScanResultsDir is the optional directory in which to write out positive scan results.
	// This could be directory that Splunk searches for container log files if not using PostResultURL.
	ScanResultsDir = os.Getenv("SCAN_RESULTS_DIR")
	// PostResultURL is the OpenShift service URL or route where we send positive scan results.
	PostResultURL = os.Getenv("POST_RESULT_URL")
	// OutFile is an optional parameter for writing positive scan results locally.
	OutFile = os.Getenv("OUT_FILE")
	// ScanDirs is a comma-separated list of directories to include in scheduled host filesystem scans.
	ScanDirs = os.Getenv("HOST_SCAN_DIRS")
	// MinConDay is the minimum number of days a container has been running in order to be included in scheduled container scans.
	MinConDay = os.Getenv("MIN_CON_DAY")
	// ActiveScan is set to true if this instance of the program is to actively watch for new containers to scan.
	ActiveScan = os.Getenv("ACTIVE_SCAN")
	// ScheduledScan is set to true if this instance of the program is to scan existing containers on a specified basis.
	ScheduledScan = os.Getenv("SCHEDULED_SCAN")
)
