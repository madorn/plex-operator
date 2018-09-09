package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PlexList resource
type PlexList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Plex `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

//Plex resource
type Plex struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              PlexSpec   `json:"spec"`
	Status            PlexStatus `json:"status,omitempty"`
}
// PlexSpec is the spec for the Plex Custom Resource
type PlexSpec struct {
	
	// Size is the number of replicas for the Plex deployment.
	Size int32	`json:"size"`

	// TimeZone is the time zone to use.
	TimeZone string `json:"timeZone"`

	// ClaimToken is the claim token for the server to obtain a real server token.
	ClaimToken string `json:"claimToken"`

	// ConfigMapName is the name of ConfigMap to use or create.
	ConfigMapName string `json:"configMapName"`

	// ConfigMountPath 
	ConfigMountPath string `json:"configMountPath"`

	// DataMountPath path for configuration path for Plex
	DataMountPath string `json:"dataMountPath"`

	// TranscodeMountPath path for transcode path for Plex
	TranscodeMountPath string `json:"transcodeMountPath"`

}

//PlexStatus defines the status for the Plex custom resource
type PlexStatus struct {
	Pods             []string          `json:"pods"`
	ExternalAddresses map[string]string `json:"externalAddresses"`
}
