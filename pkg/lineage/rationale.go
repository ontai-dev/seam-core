// Package lineage defines the sealed causal chain types used across all Seam
// Operators. All operators import this package when populating the lineage fields
// on objects they create. No operator extends this vocabulary unilaterally — new
// values require a Pull Request to seam-core and Platform Governor review.
package lineage

// CreationRationale is the compile-time controlled vocabulary for the reason
// a Seam-managed object was created. Every SealedCausalChain carries exactly
// one CreationRationale value, drawn from this enumeration.
//
// This enumeration is the single source of truth for creation rationale strings
// across the entire Seam platform. It is not a free-text field and it is not
// a per-operator registry. Machine reasoning and audit tooling depend on this
// vocabulary remaining stable and typed.
//
// +kubebuilder:validation:Enum=ClusterProvision;ClusterDecommission;SecurityEnforcement;PackExecution;VirtualizationFulfillment;ConductorAssignment;VortexBinding
type CreationRationale string

const (
	// ClusterProvision is used by the Platform operator when a TalosCluster or
	// related cluster lifecycle root declaration is created.
	ClusterProvision CreationRationale = "ClusterProvision"

	// ClusterDecommission is used by the Platform operator when a cluster
	// decommission root declaration (e.g., ClusterReset) is created.
	ClusterDecommission CreationRationale = "ClusterDecommission"

	// SecurityEnforcement is used by the Guardian operator when a security plane
	// root declaration (e.g., RBACPolicy, PermissionSet) causes a derived object
	// to be created.
	SecurityEnforcement CreationRationale = "SecurityEnforcement"

	// PackExecution is used by the Wrapper operator when a pack delivery or
	// execution root declaration (e.g., PackExecution, ClusterPack) is created.
	PackExecution CreationRationale = "PackExecution"

	// VirtualizationFulfillment is used by the Screen operator (future) when a
	// virtualization workload root declaration (e.g., VirtCluster) is created.
	VirtualizationFulfillment CreationRationale = "VirtualizationFulfillment"

	// ConductorAssignment is used by the Conductor binary in agent mode when an
	// operational assignment object is created by the management cluster Conductor.
	ConductorAssignment CreationRationale = "ConductorAssignment"

	// VortexBinding is used by the Vortex operator (future) when a portal policy
	// binding root declaration is created.
	VortexBinding CreationRationale = "VortexBinding"
)
