/*
SPDX-FileCopyrightText: 2026 SAP SE or an SAP affiliate company and image-mapper-cop contributors
SPDX-License-Identifier: Apache-2.0
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"

	"github.com/sap/component-operator-runtime/pkg/component"
	componentoperatorruntimetypes "github.com/sap/component-operator-runtime/pkg/types"
)

// ImageMapperSpec defines the desired state of ImageMapper.
type ImageMapperSpec struct {
	component.Spec `json:",inline"`
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:default=1
	ReplicaCount int `json:"replicaCount,omitempty"`
	// +optional
	Image                          component.ImageSpec `json:"image"`
	component.KubernetesProperties `json:",inline"`
	ObjectSelector                 *metav1.LabelSelector `json:"objectSelector,omitempty"`
	NamespaceSelector              *metav1.LabelSelector `json:"namespaceSelector,omitempty"`
	Mapping                        []MappingRule         `json:"mapping,omitempty"`
	LabelsAddedIfModified          map[string]string     `json:"labelsAddedIfModified,omitempty"`
	AnnotationsAddedIfModified     map[string]string     `json:"annotationsAddedIfModified,omitempty"`
	LogLevel                       int                   `json:"logLevel,omitempty"`
}

// MappingRule describes how images are transformed.
type MappingRule struct {
	Pattern     string `json:"pattern"`
	Replacement string `json:"replacement"`
}

// ImageMapperStatus defines the observed state of ImageMapper.
type ImageMapperStatus struct {
	component.Status `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="State",type=string,JSONPath=`.status.state`
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +genclient

// ImageMapper is the Schema for the imagemappers API.
type ImageMapper struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ImageMapperSpec `json:"spec,omitempty"`
	// +kubebuilder:default={"observedGeneration":-1}
	Status ImageMapperStatus `json:"status,omitempty"`
}

var _ component.Component = &ImageMapper{}

// +kubebuilder:object:root=true

// ImageMapperList contains a list of ImageMapper.
type ImageMapperList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ImageMapper `json:"items"`
}

func (s *ImageMapperSpec) ToUnstructured() map[string]any {
	result, err := runtime.DefaultUnstructuredConverter.ToUnstructured(s)
	if err != nil {
		panic(err)
	}
	return result
}

func (c *ImageMapper) GetDeploymentNamespace() string {
	if c.Spec.Namespace != "" {
		return c.Spec.Namespace
	}
	return c.Namespace
}

func (c *ImageMapper) GetDeploymentName() string {
	if c.Spec.Name != "" {
		return c.Spec.Name
	}
	return c.Name
}

func (c *ImageMapper) GetSpec() componentoperatorruntimetypes.Unstructurable {
	return &c.Spec
}

func (c *ImageMapper) GetStatus() *component.Status {
	return &c.Status.Status
}

func init() {
	SchemeBuilder.Register(&ImageMapper{}, &ImageMapperList{})
}
