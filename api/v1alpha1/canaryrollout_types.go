/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0
*/

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=canaryrollouts,scope=Namespaced

package v1alpha1

import (
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CanaryRollout is the Schema for the canaryrollouts API
type CanaryRollout struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`

    Spec   CanaryRolloutSpec   `json:"spec,omitempty"`
    Status CanaryRolloutStatus `json:"status,omitempty"`
}

// CanaryRolloutSpec defines the desired state of CanaryRollout
type CanaryRolloutSpec struct {
    // StableIngress is the name of the Ingress handling 100% traffic pre-rollout
    StableIngress string `json:"stableIngress"`

    // CanaryIngress is the name of the Ingress whose "canary-weight" annotation will be patched
    CanaryIngress string `json:"canaryIngress"`

    // Steps defines the weight+pause sequence for the rollout
    Steps []RolloutStep `json:"steps"`
}

// RolloutStep describes one weight change and optional pause
type RolloutStep struct {
    // Weight is the percentage of traffic to send to the canary Ingress
    Weight int32 `json:"weight"`

    // PauseSeconds is how long (in seconds) to wait after setting this weight
    // +optional
    PauseSeconds *int32 `json:"pauseSeconds,omitempty"`
}

// CanaryRolloutStatus defines the observed state of CanaryRollout
type CanaryRolloutStatus struct {
    // CurrentStep is the index into Spec.Steps that has most recently been applied
    // +optional
    CurrentStep *int32 `json:"currentStep,omitempty"`

    // Completed will be true after the final weight has been applied
    // +optional
    Completed bool `json:"completed,omitempty"`
}

// +kubebuilder:object:root=true

// CanaryRolloutList contains a list of CanaryRollout
type CanaryRolloutList struct {
    metav1.TypeMeta `json:",inline"`
    metav1.ListMeta `json:"metadata,omitempty"`
    Items           []CanaryRollout `json:"items"`
}

