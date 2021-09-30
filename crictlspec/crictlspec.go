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
