// STUB — InfrastructureLineageIndex is the infrastructure domain instantiation
// of DomainLineageIndex from core.ontai.dev. See seam-core-schema.md §3 and
// domain-core-schema.md for the full design rationale.
//
// The LineageController that manages InfrastructureLineageIndex CR lifecycle is
// a deferred implementation milestone. This stub defines the types only.
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/ontai-dev/seam-core/pkg/lineage"
)

// InfrastructureLineageIndexRootBinding records the root declaration that anchors
// this lineage index. All fields are immutable after admission.
type InfrastructureLineageIndexRootBinding struct {
	// RootKind is the kind of the root declaration (e.g., TalosCluster, PackExecution).
	RootKind string `json:"rootKind"`

	// RootName is the name of the root declaration.
	RootName string `json:"rootName"`

	// RootNamespace is the namespace of the root declaration.
	RootNamespace string `json:"rootNamespace"`

	// RootUID is the UID of the root declaration at time of index creation.
	RootUID types.UID `json:"rootUID"`

	// RootObservedGeneration is the metadata.generation of the root declaration
	// when this index was created.
	RootObservedGeneration int64 `json:"rootObservedGeneration"`
}

// DescendantEntry records a single derived object in the lineage index.
// Entries are appended monotonically. An entry is never modified or removed.
type DescendantEntry struct {
	// Kind is the kind of the derived object.
	Kind string `json:"kind"`

	// Name is the name of the derived object.
	Name string `json:"name"`

	// Namespace is the namespace of the derived object.
	Namespace string `json:"namespace"`

	// UID is the UID of the derived object.
	UID types.UID `json:"uid"`

	// SeamOperator is the name of the Seam Operator that created this derived
	// object (e.g., platform, guardian, wrapper, conductor).
	SeamOperator string `json:"seamOperator"`

	// CreationRationale is the reason this derived object was created, drawn from
	// the Seam Core controlled vocabulary (pkg/lineage.CreationRationale).
	//
	// +kubebuilder:validation:Enum=ClusterProvision;ClusterDecommission;SecurityEnforcement;PackExecution;VirtualizationFulfillment;ConductorAssignment;VortexBinding
	CreationRationale lineage.CreationRationale `json:"creationRationale"`

	// RootGenerationAtCreation is the metadata.generation of the root declaration
	// at the time this derived object was created.
	RootGenerationAtCreation int64 `json:"rootGenerationAtCreation"`
}

// InfrastructurePolicyBindingStatus records the InfrastructurePolicy and
// InfrastructureProfile bound to the root declaration at last evaluation.
type InfrastructurePolicyBindingStatus struct {
	// DomainPolicyRef is the name of the InfrastructurePolicy bound to the root
	// declaration.
	// +optional
	DomainPolicyRef string `json:"domainPolicyRef,omitempty"`

	// DomainProfileRef is the name of the InfrastructureProfile bound to the
	// root declaration.
	// +optional
	DomainProfileRef string `json:"domainProfileRef,omitempty"`

	// PolicyGenerationAtLastEvaluation is the metadata.generation of the bound
	// InfrastructurePolicy at the time of the last policy evaluation cycle.
	// +optional
	PolicyGenerationAtLastEvaluation int64 `json:"policyGenerationAtLastEvaluation,omitempty"`

	// DriftDetected is true if the controller detected drift between the expected
	// state derived from the InfrastructurePolicy and the observed state of
	// derived objects at the last evaluation.
	// +optional
	DriftDetected bool `json:"driftDetected,omitempty"`
}

// InfrastructureLineageIndexSpec is the spec of an InfrastructureLineageIndex.
type InfrastructureLineageIndexSpec struct {
	// RootBinding records the root declaration that anchors this lineage index.
	// Immutable after admission. The admission webhook rejects any update that
	// modifies a field in this section.
	RootBinding InfrastructureLineageIndexRootBinding `json:"rootBinding"`

	// DescendantRegistry is the list of all objects derived from the root
	// declaration. Appended monotonically as new derived objects are created.
	// Entries are never modified or removed.
	// +optional
	DescendantRegistry []DescendantEntry `json:"descendantRegistry,omitempty"`

	// PolicyBindingStatus records the InfrastructurePolicy and InfrastructureProfile
	// bound to the root declaration at last evaluation.
	// +optional
	PolicyBindingStatus *InfrastructurePolicyBindingStatus `json:"policyBindingStatus,omitempty"`
}

// InfrastructureLineageIndexStatus is the observed state of an
// InfrastructureLineageIndex.
type InfrastructureLineageIndexStatus struct {
	// ObservedGeneration is the last generation of the InfrastructureLineageIndex
	// that the controller has processed.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Conditions holds the standard Kubernetes condition array for this resource.
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,shortName=ili

// InfrastructureLineageIndex is the sealed causal chain index for a root
// declaration in the Seam infrastructure domain. It instantiates the abstract
// DomainLineageIndex schema from core.ontai.dev.
//
// One InfrastructureLineageIndex is created per root declaration. All derived
// objects carry a reference to their root's index; they do not carry their own
// index instances. This is the Lineage Index Pattern — seam-core-schema.md §3.
//
// STUB: The LineageController that manages this CR's lifecycle is a deferred
// implementation milestone. See seam-core-schema.md §6.
type InfrastructureLineageIndex struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   InfrastructureLineageIndexSpec   `json:"spec,omitempty"`
	Status InfrastructureLineageIndexStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// InfrastructureLineageIndexList contains a list of InfrastructureLineageIndex.
type InfrastructureLineageIndexList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []InfrastructureLineageIndex `json:"items"`
}

func init() {
	SchemeBuilder.Register(&InfrastructureLineageIndex{}, &InfrastructureLineageIndexList{})
}
