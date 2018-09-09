package stub

import (
	"fmt"
	"reflect"
	"context"

	v1alpha1 "github.com/madorn/plex-operator/pkg/apis/plex/v1alpha1"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"github.com/sirupsen/logrus"



)







//NewHandler returns the Handler type
func NewHandler() sdk.Handler {
	return &Handler{}
}

//Handler is an empty struct
type Handler struct {
}

//Handle function is where our event logic lives
func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {
	switch o := event.Object.(type) {
	case *v1alpha1.Plex:
		plex := o

		// Ignore the delete event since the garbage collector will clean up all secondary resources for the CR
		// All secondary resources must have the CR set as their OwnerReference for this to be the case
		if event.Deleted {
			return nil
		}

		// Create the deployment if it doesn't exist
		dep := deploymentForPlex(plex)
		err := sdk.Create(dep)
		if err != nil && !apierrors.IsAlreadyExists(err) {
			return fmt.Errorf("failed to create deployment: %v", err)
		}

		// Ensure the deployment size is the same as the spec
		err = sdk.Get(dep)
		if err != nil {
			return fmt.Errorf("failed to get deployment: %v", err)
		}
		size := plex.Spec.Size
		if *dep.Spec.Replicas != size {
			dep.Spec.Replicas = &size
			err = sdk.Update(dep)
			if err != nil {
				return fmt.Errorf("failed to update deployment: %v", err)
			}
		}

		// Update the Plex status with the pod names
		podList := podList()
		labelSelector := labels.SelectorFromSet(labelsForPlex(plex.Name)).String()
		listOps := &metav1.ListOptions{LabelSelector: labelSelector}
		err = sdk.List(plex.Namespace, podList, sdk.WithListOptions(listOps))
		if err != nil {
			return fmt.Errorf("failed to list pods: %v", err)
		}
		podNames := getPodNames(podList.Items)
		if !reflect.DeepEqual(podNames, plex.Status.Pods) {
			plex.Status.Pods = podNames
			err := sdk.Update(plex)
			if err != nil {
				return fmt.Errorf("failed to update memcached status: %v", err)
			}
		}
	}
	return nil
}

// deploymentForPlex returns a Plex Deployment object
func deploymentForPlex(p *v1alpha1.Plex) *appsv1.Deployment {
	ls := labelsForPlex(p.Name)
	replicas := p.Spec.Size

	dep := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      p.Name,
			Namespace: p.Namespace,
		},

		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},

			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},

				Spec: v1.PodSpec{
					ServiceAccountName:	"useroot",
					Containers: []v1.Container{{
						Name: "plex",
						Resources: v1.ResourceRequirements{
							Requests: v1.ResourceList{
								"cpu":	*resource.NewMilliQuantity(250, resource.BinarySI),
								"memory": *resource.NewMilliQuantity(64, resource.BinarySI),
							},
							Limits: v1.ResourceList{
								"cpu":	*resource.NewMilliQuantity(500, resource.BinarySI),
								"memory": *resource.NewMilliQuantity(128, resource.BinarySI),
							},
						},
						Image:   "plexinc/pms-docker:1.13.0.5023-31d3c0c65",
						ImagePullPolicy: v1.PullPolicy("Always"),
						VolumeMounts:	[]v1.VolumeMount{{
							Name: "plex-config",
							MountPath: p.Spec.ConfigMountPath,
						},
						{	
							Name: "plex-transcode",
							MountPath: p.Spec.TranscodeMountPath,
						},
						{
							Name: "plex-data",
							MountPath: p.Spec.DataMountPath,
						},
					},
						Ports:	[]v1.ContainerPort{{
								Name:	"plex-ui",
								Protocol:	v1.ProtocolTCP,
								ContainerPort:	32400,
							},
							{
								Name:	"plex-home",
								Protocol:	v1.ProtocolTCP,
								ContainerPort: 3005,
							},
							{
								Name:	"plex-roku",
								Protocol:	v1.ProtocolTCP,
								ContainerPort: 8324,
							},
							{
								Name:	"plex-dlna-tcp",
								Protocol:	v1.ProtocolTCP,
								ContainerPort: 32469,
							},
							{
								Name: "plex-dlna-udp",
								Protocol:	v1.ProtocolUDP,
								ContainerPort: 1900,
							},
							{
								Name: "plex-discovery1",
								Protocol:	v1.ProtocolUDP,
								ContainerPort: 32410,
							},
							{
								Name: "plex-discovery2",
								Protocol: v1.ProtocolUDP,
								ContainerPort: 32412,
							},
							{
								Name: "plex-discovery3",
								Protocol:	v1.ProtocolUDP,
								ContainerPort: 32413,
							},
							{
								Name: "plex-discovery4",
								Protocol: v1.ProtocolUDP,
								ContainerPort: 32414,
							},
						},
						Env:	[]v1.EnvVar{{
								Name: "TZ",
								Value: p.Spec.TimeZone,
							},
							{
								Name: "CLAIM_TOKEN",
								Value: p.Spec.ClaimToken,
							},
						},
						ReadinessProbe: &v1.Probe{
							Handler: v1.Handler{
								Exec: &v1.ExecAction{
									Command: []string{
										"/bin/sh",
										"-c",
										"LD_LIBRARY_PATH=/usr/lib/plexmediaserver '/usr/lib/plexmediaserver/Plex Media Server Tests' --gtest_filter=SanityChecks",
									},
								},
							},
								InitialDelaySeconds: 5,
								PeriodSeconds: 10,
							},
						LivenessProbe: &v1.Probe{
							Handler: v1.Handler{
									HTTPGet: &v1.HTTPGetAction{
										Path: "/",
										Port: intstr.FromInt(32400),	
										Scheme:	v1.URISchemeHTTP,
									},
								},
								InitialDelaySeconds: 15,
								PeriodSeconds:       20,
								},
							}},
					InitContainers: []v1.Container{{
						Name: "init",
						Resources: v1.ResourceRequirements{
							Requests: v1.ResourceList{
								"cpu":	*resource.NewMilliQuantity(250, resource.BinarySI),
								"memory": *resource.NewMilliQuantity(64, resource.BinarySI),
							},
							Limits: v1.ResourceList{
								"cpu":	*resource.NewMilliQuantity(500, resource.BinarySI),
								"memory": *resource.NewMilliQuantity(128, resource.BinarySI),
							},
						},
						Image:   "busybox:1.29",
						ImagePullPolicy: v1.PullPolicy("Always"),
						Command:	[]string{
							"/bin/sh",
							"-c",
							"mkdir -p /config/Library/Application Support/Plex Media Server && cp /etc/plex/Preferences.xml /config/Library/Application Support/Plex Media Server",
						},
						VolumeMounts:	[]v1.VolumeMount{{
							Name: "plex-config",
							MountPath: p.Spec.ConfigMountPath,
						},
						{	
							Name: "plex-preferences",
							MountPath: "/etc/plex",
						},
					},
				}},
					Volumes:	[]v1.Volume{{
						Name:	"plex-config",
						VolumeSource:	v1.VolumeSource{
								PersistentVolumeClaim:	&v1.PersistentVolumeClaimVolumeSource{
									ClaimName:	"plex-config",
								},
							},
						},
						{
						Name: "plex-transcode",
						VolumeSource: v1.VolumeSource{
								PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
									ClaimName: "plex-transcode",
								},
							},
						},
						{
						Name: "plex-data",
						VolumeSource:	v1.VolumeSource{
								PersistentVolumeClaim:	&v1.PersistentVolumeClaimVolumeSource{
									ClaimName:	"plex-data",
								},
							},
						},
						{
						Name: "plex-preferences",
						VolumeSource:	v1.VolumeSource{
								ConfigMap:	&v1.ConfigMapVolumeSource{
									LocalObjectReference:	v1.LocalObjectReference{
										Name: p.Spec.ConfigMapName,
									},
						},
					},
				},
			},
		},
	},
},
	}
addOwnerRefToObject(dep, asOwner(p))
	return dep
}

// labelsForPlex returns the labels for selecting the resources
// belonging to the given memcached CR name.
func labelsForPlex(name string) map[string]string {
	return map[string]string{"app": "plex", "tier": "frontend", "environment": "prod"}
	}

// addOwnerRefToObject appends the desired OwnerReference to the object
func addOwnerRefToObject(obj metav1.Object, ownerRef metav1.OwnerReference) {
	obj.SetOwnerReferences(append(obj.GetOwnerReferences(), ownerRef))
}

// asOwner returns an OwnerReference set as the memcached CR
func asOwner(p *v1alpha1.Plex) metav1.OwnerReference {
	trueVar := true
	return metav1.OwnerReference{
		APIVersion: p.APIVersion,
		Kind:       p.Kind,
		Name:       p.Name,
		UID:        p.UID,
		Controller: &trueVar,
	}
}

// podList returns a v1.PodList object
func podList() *v1.PodList {
	return &v1.PodList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
	}
}

// getPodNames returns the pod names of the array of pods passed in
func getPodNames(pods []v1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}