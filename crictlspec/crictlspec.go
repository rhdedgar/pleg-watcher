package crictlspec

type Containers struct {
	Containers []Container `json:"containers"`
}

type Container struct {
	ID           string     `json:"id"`
	PodSandboxID string     `json:"podSandboxId"`
	Metadata     Metadata   `json:"metadata"`
	Image        Image      `json:"image"`
	ImageRef     string     `json:"imageRef"`
	State        string     `json:"state"`
	CreatedAt    string     `json:"createdAt"`
	Labels       Label      `json:"labels"`
	Annotations  Annotation `json:"annotations,omitempty"`
}

type Metadata struct {
	Name    string `json:"name"`
	Attempt int    `json:"attempt"`
}
type Image struct {
	Image string `json:"image"`
}
type Label struct {
	IoKubernetesContainerName string `json:"io.kubernetes.container.name"`
	IoKubernetesPodName       string `json:"io.kubernetes.pod.name"`
	IoKubernetesPodNamespace  string `json:"io.kubernetes.pod.namespace"`
	IoKubernetesPodUID        string `json:"io.kubernetes.pod.uid"`
}
type Annotation struct {
	IoKubernetesContainerHash                     string `json:"io.kubernetes.container.hash"`
	IoKubernetesContainerRestartCount             string `json:"io.kubernetes.container.restartCount"`
	IoKubernetesContainerTerminationMessagePath   string `json:"io.kubernetes.container.terminationMessagePath"`
	IoKubernetesContainerTerminationMessagePolicy string `json:"io.kubernetes.container.terminationMessagePolicy"`
	IoKubernetesPodTerminationGracePeriod         string `json:"io.kubernetes.pod.terminationGracePeriod"`
}
