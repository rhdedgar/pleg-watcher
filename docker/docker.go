package docker

import "time"

type State struct {
	Status     string    `json:"Status"`
	Running    bool      `json:"Running"`
	Paused     bool      `json:"Paused"`
	Restarting bool      `json:"Restarting"`
	OOMKilled  bool      `json:"OOMKilled"`
	Dead       bool      `json:"Dead"`
	Pid        int       `json:"Pid"`
	ExitCode   int       `json:"ExitCode"`
	Error      string    `json:"Error"`
	StartedAt  time.Time `json:"StartedAt"`
	FinishedAt time.Time `json:"FinishedAt"`
}

type LogConfig struct {
	Type    string  `json:"Type"`
	ConfigA ConfigA `json:"Config"`
}

type ConfigA struct {
	MaxSize string `json:"max-size"`
}

type RestartPolicy struct {
	Name              string `json:"Name"`
	MaximumRetryCount int    `json:"MaximumRetryCount"`
}

type HostConfig struct {
	Binds                []string      `json:"Binds"`
	ContainerIDFile      string        `json:"ContainerIDFile"`
	LogConfig            LogConfig     `json:"LogConfig"`
	NetworkMode          string        `json:"NetworkMode"`
	PortBindings         interface{}   `json:"PortBindings"`
	RestartPolicy        RestartPolicy `json:"RestartPolicy"`
	AutoRemove           bool          `json:"AutoRemove"`
	VolumeDriver         string        `json:"VolumeDriver"`
	VolumesFrom          interface{}   `json:"VolumesFrom"`
	CapAdd               interface{}   `json:"CapAdd"`
	CapDrop              interface{}   `json:"CapDrop"`
	DNS                  interface{}   `json:"Dns"`
	DNSOptions           interface{}   `json:"DnsOptions"`
	DNSSearch            interface{}   `json:"DnsSearch"`
	ExtraHosts           interface{}   `json:"ExtraHosts"`
	GroupAdd             interface{}   `json:"GroupAdd"`
	IpcMode              string        `json:"IpcMode"`
	Cgroup               string        `json:"Cgroup"`
	Links                interface{}   `json:"Links"`
	OomScoreAdj          int           `json:"OomScoreAdj"`
	PidMode              string        `json:"PidMode"`
	Privileged           bool          `json:"Privileged"`
	PublishAllPorts      bool          `json:"PublishAllPorts"`
	ReadonlyRootfs       bool          `json:"ReadonlyRootfs"`
	SecurityOpt          []string      `json:"SecurityOpt"`
	UTSMode              string        `json:"UTSMode"`
	UsernsMode           string        `json:"UsernsMode"`
	ShmSize              int           `json:"ShmSize"`
	Runtime              string        `json:"Runtime"`
	ConsoleSize          []int         `json:"ConsoleSize"`
	Isolation            string        `json:"Isolation"`
	CPUShares            int           `json:"CpuShares"`
	Memory               int           `json:"Memory"`
	NanoCpus             int           `json:"NanoCpus"`
	CgroupParent         string        `json:"CgroupParent"`
	BlkioWeight          int           `json:"BlkioWeight"`
	BlkioWeightDevice    interface{}   `json:"BlkioWeightDevice"`
	BlkioDeviceReadBps   interface{}   `json:"BlkioDeviceReadBps"`
	BlkioDeviceWriteBps  interface{}   `json:"BlkioDeviceWriteBps"`
	BlkioDeviceReadIOps  interface{}   `json:"BlkioDeviceReadIOps"`
	BlkioDeviceWriteIOps interface{}   `json:"BlkioDeviceWriteIOps"`
	CPUPeriod            int           `json:"CpuPeriod"`
	CPUQuota             int           `json:"CpuQuota"`
	CPURealtimePeriod    int           `json:"CpuRealtimePeriod"`
	CPURealtimeRuntime   int           `json:"CpuRealtimeRuntime"`
	CpusetCpus           string        `json:"CpusetCpus"`
	CpusetMems           string        `json:"CpusetMems"`
	Devices              []interface{} `json:"Devices"`
	DiskQuota            int           `json:"DiskQuota"`
	KernelMemory         int           `json:"KernelMemory"`
	MemoryReservation    int           `json:"MemoryReservation"`
	MemorySwap           int           `json:"MemorySwap"`
	MemorySwappiness     int           `json:"MemorySwappiness"`
	OomKillDisable       bool          `json:"OomKillDisable"`
	PidsLimit            int           `json:"PidsLimit"`
	Ulimits              interface{}   `json:"Ulimits"`
	CPUCount             int           `json:"CpuCount"`
	CPUPercent           int           `json:"CpuPercent"`
	IOMaximumIOps        int           `json:"IOMaximumIOps"`
	IOMaximumBandwidth   int           `json:"IOMaximumBandwidth"`
}

type Data struct {
	LowerDir  string `json:"LowerDir"`
	MergedDir string `json:"MergedDir"`
	UpperDir  string `json:"UpperDir"`
	WorkDir   string `json:"WorkDir"`
}

type GraphDriver struct {
	Name string `json:"Name"`
	Data Data   `json:"Data"`
}

type Mounts struct {
	Type        string `json:"Type"`
	Source      string `json:"Source"`
	Destination string `json:"Destination"`
	Mode        string `json:"Mode"`
	RW          bool   `json:"RW"`
	Propagation string `json:"Propagation"`
}

type Healthcheck struct {
	Test []string `json:"Test"`
}

type Labels struct {
	AnnotationIoKubernetesContainerHash                     string `json:"annotation.io.kubernetes.container.hash"`
	AnnotationIoKubernetesContainerRestartCount             string `json:"annotation.io.kubernetes.container.restartCount"`
	AnnotationIoKubernetesContainerTerminationMessagePath   string `json:"annotation.io.kubernetes.container.terminationMessagePath"`
	AnnotationIoKubernetesContainerTerminationMessagePolicy string `json:"annotation.io.kubernetes.container.terminationMessagePolicy"`
	AnnotationIoKubernetesPodTerminationGracePeriod         string `json:"annotation.io.kubernetes.pod.terminationGracePeriod"`
	Architecture                                            string `json:"architecture"`
	AuthoritativeSourceURL                                  string `json:"authoritative-source-url"`
	BuildDate                                               string `json:"build-date"`
	ComRedhatBuildHost                                      string `json:"com.redhat.build-host"`
	ComRedhatComponent                                      string `json:"com.redhat.component"`
	Description                                             string `json:"description"`
	DistributionScope                                       string `json:"distribution-scope"`
	IoK8SDescription                                        string `json:"io.k8s.description"`
	IoK8SDisplayName                                        string `json:"io.k8s.display-name"`
	IoKubernetesContainerLogpath                            string `json:"io.kubernetes.container.logpath"`
	IoKubernetesContainerName                               string `json:"io.kubernetes.container.name"`
	IoKubernetesDockerType                                  string `json:"io.kubernetes.docker.type"`
	IoKubernetesPodName                                     string `json:"io.kubernetes.pod.name"`
	IoKubernetesPodNamespace                                string `json:"io.kubernetes.pod.namespace"`
	IoKubernetesPodUID                                      string `json:"io.kubernetes.pod.uid"`
	IoKubernetesSandboxID                                   string `json:"io.kubernetes.sandbox.id"`
	IoOpenshiftBuildCommitAuthor                            string `json:"io.openshift.build.commit.author"`
	IoOpenshiftBuildCommitDate                              string `json:"io.openshift.build.commit.date"`
	IoOpenshiftBuildCommitID                                string `json:"io.openshift.build.commit.id"`
	IoOpenshiftBuildCommitMessage                           string `json:"io.openshift.build.commit.message"`
	IoOpenshiftBuildCommitRef                               string `json:"io.openshift.build.commit.ref"`
	IoOpenshiftBuildName                                    string `json:"io.openshift.build.name"`
	IoOpenshiftBuildNamespace                               string `json:"io.openshift.build.namespace"`
	IoOpenshiftBuildSourceContextDir                        string `json:"io.openshift.build.source-context-dir"`
	IoOpenshiftBuildSourceLocation                          string `json:"io.openshift.build.source-location"`
	IoOpenshiftTags                                         string `json:"io.openshift.tags"`
	Name                                                    string `json:"name"`
	Release                                                 string `json:"release"`
	Summary                                                 string `json:"summary"`
	URL                                                     string `json:"url"`
	VcsRef                                                  string `json:"vcs-ref"`
	VcsType                                                 string `json:"vcs-type"`
	Vendor                                                  string `json:"vendor"`
	Version                                                 string `json:"version"`
}

type Config struct {
	Hostname     string      `json:"Hostname"`
	Domainname   string      `json:"Domainname"`
	User         string      `json:"User"`
	AttachStdin  bool        `json:"AttachStdin"`
	AttachStdout bool        `json:"AttachStdout"`
	AttachStderr bool        `json:"AttachStderr"`
	Tty          bool        `json:"Tty"`
	OpenStdin    bool        `json:"OpenStdin"`
	StdinOnce    bool        `json:"StdinOnce"`
	Env          []string    `json:"Env"`
	Cmd          []string    `json:"Cmd"`
	Healthcheck  Healthcheck `json:"Healthcheck"`
	ArgsEscaped  bool        `json:"ArgsEscaped"`
	Image        string      `json:"Image"`
	Volumes      interface{} `json:"Volumes"`
	WorkingDir   string      `json:"WorkingDir"`
	Entrypoint   interface{} `json:"Entrypoint"`
	OnBuild      interface{} `json:"OnBuild"`
	Labels       Labels      `json:"Labels"`
}

type NetworkSettings struct {
	Bridge                 string      `json:"Bridge"`
	SandboxID              string      `json:"SandboxID"`
	HairpinMode            bool        `json:"HairpinMode"`
	LinkLocalIPv6Address   string      `json:"LinkLocalIPv6Address"`
	LinkLocalIPv6PrefixLen int         `json:"LinkLocalIPv6PrefixLen"`
	Ports                  interface{} `json:"Ports"`
	SandboxKey             string      `json:"SandboxKey"`
	SecondaryIPAddresses   interface{} `json:"SecondaryIPAddresses"`
	SecondaryIPv6Addresses interface{} `json:"SecondaryIPv6Addresses"`
	EndpointID             string      `json:"EndpointID"`
	Gateway                string      `json:"Gateway"`
	GlobalIPv6Address      string      `json:"GlobalIPv6Address"`
	GlobalIPv6PrefixLen    int         `json:"GlobalIPv6PrefixLen"`
	IPAddress              string      `json:"IPAddress"`
	IPPrefixLen            int         `json:"IPPrefixLen"`
	IPv6Gateway            string      `json:"IPv6Gateway"`
	MacAddress             string      `json:"MacAddress"`
	Networks               struct{}    `json:"Networks"`
}

type DockerContainer []struct {
	ID              string          `json:"Id"`
	Created         time.Time       `json:"Created"`
	Path            string          `json:"Path"`
	Args            []string        `json:"Args"`
	State           State           `json:"State"`
	Image           string          `json:"Image"`
	ResolvConfPath  string          `json:"ResolvConfPath"`
	HostnamePath    string          `json:"HostnamePath"`
	HostsPath       string          `json:"HostsPath"`
	LogPath         string          `json:"LogPath"`
	Name            string          `json:"Name"`
	RestartCount    int             `json:"RestartCount"`
	Driver          string          `json:"Driver"`
	MountLabel      string          `json:"MountLabel"`
	ProcessLabel    string          `json:"ProcessLabel"`
	AppArmorProfile string          `json:"AppArmorProfile"`
	ExecIDs         interface{}     `json:"ExecIDs"`
	HostConfig      HostConfig      `json:"HostConfig"`
	GraphDriver     GraphDriver     `json:"GraphDriver"`
	Mounts          []Mounts        `json:"Mounts"`
	Config          Config          `json:"Config"`
	NetworkSettings NetworkSettings `json:"NetworkSettings"`
}
