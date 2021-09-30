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

package api

import (
	"context"
	"os"
	"time"
)

// The default version for the result API object
const DefaultResultsAPIVersion = "v1alpha"

// ScanResult represents the compacted result of all scans performed on the target
type ScanResult struct {
	// APIVersion represents an API version for this result
	APIVersion string `json:"apiVersion"`
	// ContainerID contains the containerspec container to inspect.
	ContainerID string `json:"containerID,omitempty"`
	// ImageIUD is a SHA256 identifier of the scanned image
	// Note that we don't set the imageID when a container is the target of the scan.
	ImageID string `json:"imageID,omitempty"`
	// ImageName is a full pull spec of the input image
	ImageName string `json:"imageName,omitempty"`
	// NameSpace is the namespace in which the container was created
	NameSpace string `json:"nameSpace,omitempty"`
	// PodName is the name of the pod in which the container was created
	PodName string `json:"podName,omitempty"`
	// Results contains compacted results of various scans performed on the image.
	// Empty results means no problems were found with the given image.
	Results []Result `json:"results,omitempty"`
}

// Result represents the compacted result of a single scan
type Result struct {
	// Name is the name of the scanner that produced this result
	Name string `json:"name"`
	// ScannerVersion is the scanner version
	ScannerVersion string `json:"scannerVersion"`
	// Timestamp is the exact time this scan was performed
	Timestamp time.Time `json:"timestamp"`
	// Reference contains URL to more details about the given result
	Reference string `json:"reference"`
	// Description describes the result in human readable form
	Description string `json:"description,omitempty"`
	// Summary contains a list of severities for the given result
	Summary []Summary `json:"summary,omitempty"`
}

type Severity string

var (
	SeverityLow       Severity = "low"
	SeverityModerate  Severity = "moderate"
	SeverityImportant Severity = "important"
	SeverityCritical  Severity = "critical"
)

// Summary represents a severy of a given result. The result can have multiple severieties
// defined.
type Summary struct {
	// Label is the human readable severity (high, important, etc.)
	Label Severity
}

var (
	ScanOptions = []string{"clamav"}
)

// APIVersions holds a slice of supported API versions.
type APIVersions struct {
	// Versions is the supported API versions
	Versions []string `json:"versions"`
}

// FilesFilter desribes callback to filter files.
type FilesFilter func(string, os.FileInfo) bool

// Scanner interface that all scanners should define.
type Scanner interface {
	// Scan will perform a scan on the given path for the given Image.
	// It should return compacted results for JSON serialization and additionally scanner
	// specific results with more details. The context object can be used to cancel the scanning process.
	Scan(ctx context.Context, path string, filter FilesFilter) ([]Result, interface{}, error)

	// Name is the scanner's name
	Name() string
}
