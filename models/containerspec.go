package models

import (
	"time"
)

type Mount struct {
	ContainerPath  string `json:"ContainerPath,omitempty" yaml:"ContainerPath,omitempty"`
	HostPath       string `json:"HostPath,omitempty" yaml:"HostPath,omitempty"`
	Propagation    string `json:"Propagation,omitempty" yaml:"Propagation,omitempty"`
	Readonly       bool   `json:"Readonly,omitempty" yaml:"Readonly,omitempty"`
	SelinuxRelabel bool   `json:"SelinuxRelabel,omitempty" yaml:"SelinuxRelabel,omitempty"`
}

type MetaData struct {
	Attempt int    `json:"Attempt,omitempty" yaml:"Attempt,omitempty"`
	Name    string `json:"Name,omitempty" yaml:"Name,omitempty"`
}

type Label struct {
	IoKubernetesContainerName string `json:"io.kubernetes.container.name" yaml:"io.kubernetes.container.name"`
	IoKubernetesPodName       string `json:"io.kubernetes.pod.name" yaml:"io.kubernetes.pod.name"`
	IoKubernetesPodNamespace  string `json:"io.kubernetes.pod.namespace" yaml:"io.kubernetes.pod.namespace"`
	IoKubernetesPodUID        string `json:"io.kubernetes.pod.uid" yaml:"io.kubernetes.pod.uid"`
}

type Port struct {
	ContainerPort int    `json:"ContainerPort,omitempty" yaml:"ContainerPort,omitempty"`
	Protocol      string `json:"Protocol,omitempty" yaml:"Protocol,omitempty"`
}

type Annotation struct {
	IoKubernetesContainerHash                     string `json:"io.kubernetes.container.hash" yaml:"io.kubernetes.container.hash"`
	IoKubernetesContainerPorts                    string `json:"io.kubernetes.container.ports,omitempty" yaml:"io.kubernetes.container.ports,omitempty"`
	IoKubernetesContainerRestartCount             string `json:"io.kubernetes.container.restartCount,omitempty" yaml:"io.kubernetes.container.restartCount,omitempty"`
	IoKubernetesContainerTerminationMessagePath   string `json:"io.kubernetes.container.terminationMessagePath,omitempty" yaml:"io.kubernetes.container.terminationMessagePath,omitempty"`
	IoKubernetesContainerTerminationMessagePolicy string `json:"io.kubernetes.container.terminationMessagePolicy,omitempty" yaml:"io.kubernetes.container.terminationMessagePolicy,omitempty"`
	IoKubernetesPodTerminationGracePeriod         string `json:"io.kubernetes.pod.terminationGracePeriod,omitempty" yaml:"io.kubernetes.pod.terminationGracePeriod,omitempty"`
}

type Status struct {
	Status Container `json:"Status" yaml:"Status"`
}

type Container struct {
	ID         string    `json:"Id" yaml:"Id"`
	Metadata   MetaData  `json:"Metadata,omitempty" yaml:"Metadata,omitempty"`
	State      string    `json:"State,omitempty" yaml:"State,omitempty"`
	CreatedAt  time.Time `json:"CreatedAt,omitempty" yaml:"CreatedAt,omitempty"`
	StartedAt  time.Time `json:"StartedAt,omitempty" yaml:"StartedAt,omitempty"`
	FinishedAt time.Time `json:"FinishedAt,omitempty" yaml:"FinishedAt,omitempty"`
	ExitCode   int       `json:"ExitCode,omitempty" yaml:"ExitCode,omitempty"`
	Image      struct {
		Image string `json:"Image,omitempty" yaml:"Image,omitempty"`
	}
	ImageRef    string     `json:"ImageRef,omitempty" yaml:"ImageRef,omitempty"`
	Reason      string     `json:"Reason,omitempty" yaml:"Reason,omitempty"`
	Message     string     `json:"Message,omitempty" yaml:"Message,omitempty"`
	Labels      Label      `json:"Labels,omitempty" yaml:"Labels,omitempty"`
	Annotations Annotation `json:"Annotations,omitempty" yaml:"Annotations,omitempty"`
	Mounts      []Mount    `json:"Mounts,omitempty" yaml:"Mounts,omitempty"`
	LogPath     string     `json:"LogPath,omitempty" yaml:"LogPath,omitempty"`
}

type Image struct {
	ID          string   `json:"Id" yaml:"Id"`
	RepoTags    []string `json:"RepoTags,omitempty" yaml:"RepoTags,omitempty"`
	RepoDigests []string `json:"RepoDigests,omitempty" yaml:"RepoDigests,omitempty"`
	Size        int64    `json:"Size,omitempty" yaml:"Size,omitempty"`
	UID         UID      `json:"Uid,omitempty" yaml:"Uid,omitempty"`
	Username    string   `json:"Username,omitempty" yaml:"Username,omitempty"`
}

type UID struct {
	Value string `json:"Value,omitempty" yaml:"Value,omitempty"`
}

type ImageOld struct {
	ID              string    `json:"Id" yaml:"Id"`
	RepoTags        []string  `json:"RepoTags,omitempty" yaml:"RepoTags,omitempty"`
	Parent          string    `json:"Parent,omitempty" yaml:"Parent,omitempty"`
	Comment         string    `json:"Comment,omitempty" yaml:"Comment,omitempty"`
	Created         time.Time `json:"Created,omitempty" yaml:"Created,omitempty"`
	Container       string    `json:"Container,omitempty" yaml:"Container,omitempty"`
	ContainerConfig Config    `json:"ContainerConfig,omitempty" yaml:"ContainerConfig,omitempty"`
	DockerVersion   string    `json:"DockerVersion,omitempty" yaml:"DockerVersion,omitempty"`
	Author          string    `json:"Author,omitempty" yaml:"Author,omitempty"`
	Config          *Config   `json:"Config,omitempty" yaml:"Config,omitempty"`
	Architecture    string    `json:"Architecture,omitempty" yaml:"Architecture,omitempty"`
	Size            int64     `json:"Size,omitempty" yaml:"Size,omitempty"`
	VirtualSize     int64     `json:"VirtualSize,omitempty" yaml:"VirtualSize,omitempty"`
	RepoDigests     []string  `json:"RepoDigests,omitempty" yaml:"RepoDigests,omitempty"`
}

type Config struct {
	Hostname          string              `json:"Hostname,omitempty" yaml:"Hostname,omitempty"`
	Domainname        string              `json:"Domainname,omitempty" yaml:"Domainname,omitempty"`
	User              string              `json:"User,omitempty" yaml:"User,omitempty"`
	Memory            int64               `json:"Memory,omitempty" yaml:"Memory,omitempty"`
	MemorySwap        int64               `json:"MemorySwap,omitempty" yaml:"MemorySwap,omitempty"`
	MemoryReservation int64               `json:"MemoryReservation,omitempty" yaml:"MemoryReservation,omitempty"`
	KernelMemory      int64               `json:"KernelMemory,omitempty" yaml:"KernelMemory,omitempty"`
	CPUShares         int64               `json:"CpuShares,omitempty" yaml:"CpuShares,omitempty"`
	CPUSet            string              `json:"Cpuset,omitempty" yaml:"Cpuset,omitempty"`
	AttachStdin       bool                `json:"AttachStdin,omitempty" yaml:"AttachStdin,omitempty"`
	AttachStdout      bool                `json:"AttachStdout,omitempty" yaml:"AttachStdout,omitempty"`
	AttachStderr      bool                `json:"AttachStderr,omitempty" yaml:"AttachStderr,omitempty"`
	PortSpecs         []string            `json:"PortSpecs,omitempty" yaml:"PortSpecs,omitempty"`
	ExposedPorts      map[Port]struct{}   `json:"ExposedPorts,omitempty" yaml:"ExposedPorts,omitempty"`
	StopSignal        string              `json:"StopSignal,omitempty" yaml:"StopSignal,omitempty"`
	Tty               bool                `json:"Tty,omitempty" yaml:"Tty,omitempty"`
	OpenStdin         bool                `json:"OpenStdin,omitempty" yaml:"OpenStdin,omitempty"`
	StdinOnce         bool                `json:"StdinOnce,omitempty" yaml:"StdinOnce,omitempty"`
	Env               []string            `json:"Env,omitempty" yaml:"Env,omitempty"`
	Cmd               []string            `json:"Cmd" yaml:"Cmd"`
	DNS               []string            `json:"Dns,omitempty" yaml:"Dns,omitempty"`
	Image             string              `json:"Image,omitempty" yaml:"Image,omitempty"`
	Volumes           map[string]struct{} `json:"Volumes,omitempty" yaml:"Volumes,omitempty"`
	VolumeDriver      string              `json:"VolumeDriver,omitempty" yaml:"VolumeDriver,omitempty"`
	VolumesFrom       string              `json:"VolumesFrom,omitempty" yaml:"VolumesFrom,omitempty"`
	WorkingDir        string              `json:"WorkingDir,omitempty" yaml:"WorkingDir,omitempty"`
	MacAddress        string              `json:"MacAddress,omitempty" yaml:"MacAddress,omitempty"`
	Entrypoint        []string            `json:"Entrypoint" yaml:"Entrypoint"`
	NetworkDisabled   bool                `json:"NetworkDisabled,omitempty" yaml:"NetworkDisabled,omitempty"`
	SecurityOpts      []string            `json:"SecurityOpts,omitempty" yaml:"SecurityOpts,omitempty"`
	OnBuild           []string            `json:"OnBuild,omitempty" yaml:"OnBuild,omitempty"`
	Mounts            []Mount             `json:"Mounts,omitempty" yaml:"Mounts,omitempty"`
	Labels            map[string]string   `json:"Labels,omitempty" yaml:"Labels,omitempty"`
}
