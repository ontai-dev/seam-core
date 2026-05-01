package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PackReceiptDeployedResource records a single Kubernetes resource that was applied
// to the tenant cluster as part of a pack-deploy Job. Used by conductor role=tenant
// to detect drift between declared and actual cluster state.
// CLUSTERPACK-BL-VERSION-CLEANUP. conductor-schema.md.
type PackReceiptDeployedResource struct {
	// APIVersion is the full API version string (e.g., "apps/v1").
	APIVersion string `json:"apiVersion"`

	// Kind is the resource kind (e.g., "Deployment").
	Kind string `json:"kind"`

	// Namespace is the namespace the resource was applied to. Empty for cluster-scoped resources.
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// Name is the resource name.
	Name string `json:"name"`
}

// InfrastructurePackReceiptSpec defines the desired state of an InfrastructurePackReceipt.
// Written by the packinstance pull loop on the tenant cluster conductor after
// Ed25519 signature verification. INV-026. conductor-schema.md.
type InfrastructurePackReceiptSpec struct {
	// PackInstanceRef is the name of the PackInstance CR this receipt acknowledges.
	// +optional
	PackInstanceRef string `json:"packInstanceRef,omitempty"`

	// SignatureRef is the name of the signed artifact Secret on the management cluster
	// (seam-pack-signed-{cluster}-{packInstance}) from which this receipt was derived.
	// +optional
	SignatureRef string `json:"signatureRef,omitempty"`

	// ClusterPackRef is the name of the ClusterPack CR this receipt acknowledges.
	ClusterPackRef string `json:"clusterPackRef"`

	// TargetClusterRef is the name of the cluster this receipt was generated on.
	TargetClusterRef string `json:"targetClusterRef"`

	// RBACDigest is the OCI digest of the RBAC layer. Carried from ClusterPack for audit.
	// +optional
	RBACDigest string `json:"rbacDigest,omitempty"`

	// WorkloadDigest is the OCI digest of the workload layer. Carried from ClusterPack.
	// +optional
	WorkloadDigest string `json:"workloadDigest,omitempty"`

	// ChartVersion is the Helm chart version. Carried from ClusterPack.
	// +optional
	ChartVersion string `json:"chartVersion,omitempty"`

	// ChartURL is the Helm chart repository URL. Carried from ClusterPack.
	// +optional
	ChartURL string `json:"chartURL,omitempty"`

	// ChartName is the Helm chart name. Carried from ClusterPack.
	// +optional
	ChartName string `json:"chartName,omitempty"`

	// HelmVersion is the Helm SDK version. Carried from ClusterPack.
	// +optional
	HelmVersion string `json:"helmVersion,omitempty"`

	// DeployedResources is the inventory of Kubernetes resources applied to the tenant cluster
	// during the pack-deploy Job. Conductor role=tenant uses this list to detect drift by
	// verifying each resource still exists with the expected state.
	// CLUSTERPACK-BL-VERSION-CLEANUP, conductor-schema.md.
	// +optional
	DeployedResources []PackReceiptDeployedResource `json:"deployedResources,omitempty"`
}

// InfrastructurePackReceiptStatus is the observed state of an InfrastructurePackReceipt.
// Written by the packinstance pull loop after signature verification. INV-026.
type InfrastructurePackReceiptStatus struct {
	// Verified indicates whether the Ed25519 signature on the PackInstance artifact
	// was successfully verified against the management cluster's public key. INV-026.
	// +optional
	Verified bool `json:"verified,omitempty"`

	// Signature is the base64-encoded Ed25519 signature from the signed artifact Secret.
	// Stored for auditability and idempotency checking. INV-026.
	// +optional
	Signature string `json:"signature,omitempty"`

	// VerificationFailedReason is set when Verified=false and describes the
	// specific verification failure (e.g., "Ed25519 signature verification failed (INV-026)").
	// +optional
	VerificationFailedReason string `json:"verificationFailedReason,omitempty"`

	// ObservedGeneration is the generation most recently reconciled.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Conditions is the list of status conditions for this PackReceipt.
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,shortName=ipr
// +kubebuilder:printcolumn:name="Pack",type=string,JSONPath=".spec.clusterPackRef"
// +kubebuilder:printcolumn:name="Verified",type=boolean,JSONPath=".status.verified"
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=".metadata.creationTimestamp"

// InfrastructurePackReceipt is the seam-core CRD for pack delivery acknowledgement on a tenant cluster.
// Written by conductor agent after signature verification. INV-026.
type InfrastructurePackReceipt struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   InfrastructurePackReceiptSpec   `json:"spec,omitempty"`
	Status InfrastructurePackReceiptStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// InfrastructurePackReceiptList contains a list of InfrastructurePackReceipt.
type InfrastructurePackReceiptList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []InfrastructurePackReceipt `json:"items"`
}

func init() {
	SchemeBuilder.Register(&InfrastructurePackReceipt{}, &InfrastructurePackReceiptList{})
}
