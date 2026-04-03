package lineage

import (
	"k8s.io/apimachinery/pkg/types"
)

// SealedCausalChain is the immutable causal chain field embedded in every
// Seam-managed CRD spec. It records the sealed derivation from the root
// declaration that caused this object to exist.
//
// IMMUTABILITY CONTRACT: This field is authored once at object creation time
// and is sealed permanently at that point. The admission webhook will reject
// any update request that modifies any field within SealedCausalChain after
// the object has been created. No controller, no human operator, and no
// automation pipeline may alter this field post-admission.
//
// All Seam Operators embed this type by reference in their CRD specs rather
// than redefining its fields. The single definition here is the authoritative
// source of the field shape for the entire platform.
type SealedCausalChain struct {
	// RootKind is the kind of the root declaration that caused this object to
	// exist (e.g., TalosCluster, PackExecution, RBACPolicy).
	RootKind string `json:"rootKind"`

	// RootName is the name of the root declaration.
	RootName string `json:"rootName"`

	// RootNamespace is the namespace of the root declaration.
	RootNamespace string `json:"rootNamespace"`

	// RootUID is the UID of the root declaration at the time this object was
	// created. Used to verify that no root declaration replacement has occurred.
	RootUID types.UID `json:"rootUID"`

	// CreatingOperator identifies the Seam Operator that created this object.
	// This is a structured identity carrying the operator name and its deployed
	// version at creation time.
	CreatingOperator OperatorIdentity `json:"creatingOperator"`

	// CreationRationale is the reason this object was created, drawn from the
	// Seam Core controlled vocabulary defined in rationale.go. It is not a
	// free-text field.
	CreationRationale CreationRationale `json:"creationRationale"`

	// RootGenerationAtCreation is the metadata.generation of the root declaration
	// at the time this object was created. Together with RootUID, it provides a
	// complete temporal anchor for the derivation record.
	RootGenerationAtCreation int64 `json:"rootGenerationAtCreation"`
}

// OperatorIdentity identifies the Seam Operator that authored a derived object.
// It is embedded in SealedCausalChain and is subject to the same immutability
// contract — it is sealed at object creation and never modified.
type OperatorIdentity struct {
	// Name is the canonical name of the Seam Operator (e.g., platform, guardian,
	// wrapper, conductor).
	Name string `json:"name"`

	// Version is the deployed version of the operator at the time the object was
	// created (e.g., v1.26.5-r3). This allows audit tooling to correlate objects
	// with the operator version that produced them.
	Version string `json:"version"`
}
