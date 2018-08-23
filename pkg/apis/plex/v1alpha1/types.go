package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/api/core/v1"
)

const (
	// DefaultBaseImage is the default docker image to use for Plex Pods.
	DefaultBaseImage = "plexinc/pms-docker"
	// DefaultBaseImageVersion is the default version to use for Plex Pods.
	DefaultBaseImageVersion = "1.13.0.5023-31d3c0c65"
	// DefaultTimeZone is the default time zone to use for Plex Pods.
	DefaultTimeZone = "America/New_York"
	// DefaultClaimToken is the default claim token to use for Plex Pods.
	DefaultClaimToken = ""
	// DefaultConfigMapName is the default ConfigMap name to use for Plex preferences.xml
	DefaultConfigMapName = "plex-preferences-cm"
	// DefaultConfigMountPath is the default configuration path for Plex Pods.
	DefaultConfigMountPath = "/config"
	// DefaultTranscodeMountPath is the default transcode path for Plex Pods.
	DefaultTranscodeMountPath = "/transcode"
	// DefaultDataMountPath is the default data path for Plex Pods.
	DefaultDataMountPath = "/data"
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
	// BaseImage image to use for a Plex deployment.
	BaseImage string `json:"baseImage"`

	// BaseImageVersion is the version of base image to use.
	BaseImageVersion string `json:"baseImageVersion"`

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

	// Pod defines the policy for pods owned by Plex operator.
	// This field cannot be updated once the CR is created.
	Pod *PlexPodPolicy `json:"pod,omitempty"`
}

// PlexPodPolicy defines the policy for pods owned by Plex operator.
type PlexPodPolicy struct {
	// Resources is the resource requirements for the plex container.
	Resources v1.ResourceRequirements `json:"resources,omitempty"`
}

// SetDefaults sets the default vaules for the Jira spec and returns true if the spec was changed
func (p *Plex) SetDefaults() bool {
	changed := false
	if len(p.Spec.BaseImage) == 0 {
		p.Spec.BaseImage = DefaultBaseImage
		changed = true
	}
	if len(p.Spec.BaseImageVersion) == 0 {
		p.Spec.BaseImageVersion = DefaultBaseImageVersion
		changed = true
	}
	if len(p.Spec.TimeZone) == 0 {
		p.Spec.TimeZone = DefaultTimeZone
		changed = true
	}
	if len(p.Spec.ClaimToken) == 0 {
		p.Spec.ClaimToken = DefaultTimeZone
		changed = true
	}
	if len(p.Spec.ConfigMapName) == 0 {
		p.Spec.ConfigMapName = DefaultConfigMapName
	}
	if len(p.Spec.ConfigMountPath) == 0 {
		p.Spec.ConfigMountPath = DefaultConfigMountPath
		changed = true
	}
	if len(p.Spec.DataMountPath) == 0 {
		p.Spec.DataMountPath = DefaultDataMountPath
		changed = true
	}
	if len(p.Spec.TranscodeMountPath) == 0 {
		p.Spec.DataMountPath = DefaultTranscodeMountPath
		changed = true
	}
	return changed
}

// IsPodPolicySet shortcut function to determine Pod policy status.
func (p *Plex) IsPodPolicySet() bool {
	return p.Spec.Pod != nil
}
//PlexStatus defines the status for the Plex custom resource
type PlexStatus struct {
	// Fill me
}
