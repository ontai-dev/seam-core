# seam-core: Codebase Reference

## 1. Purpose

Seam-core is the exclusive schema authority for all cross-operator CRD definitions in the ONT platform (Decision G, locked April 2026). All types consumed by more than one operator are defined here under `infrastructure.ontai.dev/v1alpha1`. Operators import seam-core types and implement reconciliation behavior; they do not define CRDs. Seam-core also runs 4 controllers: `LineageReconciler`, `DescendantReconciler`, `DSNSReconciler`, `OutcomeReconciler`. Seam-core does NOT implement business logic for any operator.

---

## 2. CRD Types (`api/v1alpha1/`)

### Complete type inventory

| File | Size | Kind | Reconciler owner |
|------|------|------|-----------------|
| `taloscluster_types.go` | 260L | `InfrastructureTalosCluster` | platform |
| `runnerconfig_types.go` | 205L | `InfrastructureRunnerConfig` | platform (creates), conductor (populates status) |
| `clusterpack_types.go` | 194L | `InfrastructureClusterPack` | wrapper |
| `packexecution_types.go` | 103L | `InfrastructurePackExecution` | wrapper |
| `packinstance_types.go` | 148L | `InfrastructurePackInstance` | wrapper |
| `packbuild_types.go` | 112L | `InfrastructurePackBuild` | compiler (read-only input) |
| `packreceipt_types.go` | 132L | `InfrastructurePackReceipt` | conductor agent (tenant) |
| `packoperationresult_types.go` | 235L | `PackOperationResult` (POR) | conductor execute-mode (creates), wrapper (reads) |
| `driftsignal_types.go` | 120L | `DriftSignal` | conductor (writes), wrapper/platform (handle) |
| `infrastructurelineageindex_types.go` | 285L | `InfrastructureLineageIndex` | seam-core LineageReconciler |
| `talosclusteroperationresult_types.go` | 147L | `TalosClusterOperationResult` | platform |
| `seammembership_types.go` | 116L | `SeamMembership` | guardian |

Total: 12 CRD types. All under `infrastructure.ontai.dev/v1alpha1`. Resource names use `infrastructure` prefix (e.g., `infrastructuretalosclusters`, `infrastructureclusterpacks`).

### Key struct fields by type

**`InfrastructurePackReceiptSpec`** (`packreceipt_types.go`):
- `ClusterPackRef string` (L40)
- `TargetClusterRef string` (L43)
- `RBACDigest string` (L47)
- `WorkloadDigest string` (L51)
- `ChartVersion string` (L55)
- `ChartURL string` (L59)
- `ChartName string` (L63)
- `HelmVersion string` (L67)
- `DeployedResources []PackReceiptDeployedResource` (L74)

**`PackReceiptDeployedResource`** (`packreceipt_types.go:11`): `APIVersion string`, `Kind string`, `Namespace string` (omitempty for cluster-scoped), `Name string`. Each entry is one Kubernetes resource applied to the tenant cluster. Used by conductor `checkDrift()` and `teardownOrphanedReceipt()`.

**`PackOperationResultSpec`** (`packoperationresult_types.go:105`):
- `Revision int64` (L110) -- monotonically increasing; single-active-revision pattern (Decision E)
- `PreviousRevisionRef string` (L116) -- name of superseded POR CR; absent for revision 1
- `TalosClusterOperationResultRef string` (L122) -- stub field, not populated by any current controller
- `PackExecutionRef string` (L126)
- `ClusterPackRef string` (L130)

**`DriftSignalSpec`** (`driftsignal_types.go:48`):
- `State DriftSignalState` (L51) -- enum: `"pending"` (L15), `"delivered"` (L19), `"queued"` (L23), `"confirmed"` (L27)
- `CorrelationID string` (L55)
- `AffectedCRRef DriftAffectedCRRef` (L61) -- struct at L31: Group, Version, Kind, Namespace, Name
- `EscalationCounter int32` (L76)

**`InfrastructurePackBuildSpec`** (`packbuild_types.go:48`):
- `Category InfrastructurePackBuildCategory` (L54) -- enum at L9: `"helm"`, `"kustomize"`, `"raw"`
- `HelmSource *InfrastructurePackHelmSource` (L58)
- `KustomizeSource *InfrastructurePackKustomizeSource` (L62) -- struct at L34 (present in schema but no compiler implementation yet, T-12 open)
- `RawSource *InfrastructurePackRawSource` (L66)

### TalosCluster CRD CEL validation state

No `x-kubernetes-validations` rules exist in `config/crd/infrastructure.ontai.dev_infrastructuretalosclusters.yaml`. The YAML includes comments "Mandatory on mode=import" (L267) for the `role` field, but this is NOT enforced at admission. **T-04a open**: CEL rule `self.mode != 'import' || (has(self.role) && self.role != '')` on `InfrastructureTalosClusterSpec` is required. Run `make generate-crd` after adding `+kubebuilder:validation:XValidation` marker.

---

## 3. Controllers (`internal/controller/`)

| File | Key struct / function | What it does |
|------|-----------------------|--------------|
| `lineage_controller.go:94` | `LineageReconciler`, `Reconcile()` L101 | Creates one `InfrastructureLineageIndex` per root declaration CR. `buildILI()` L203 constructs ILI from root object. `writeGovernanceAnnotation()` L248 annotates root. `ensureLineageSyncedTrue()` L272 sets condition. `pruneStaleDescendants()` L343 removes stale registry entries. `lineageIndexName()` L451: returns `"lineage-{kind}-{name}"` (same pattern as `pkg/lineage/descendant.go:33`). |
| `descendant_reconciler.go:55` | `DescendantReconciler`, `Reconcile()` L62 | Appends derived object entries to the parent ILI `spec.descendantRegistry`. SetupWithManager at L153. |
| `dsns_reconciler.go:142` | `DSNSReconciler`, `Reconcile()` L202 | Generates DSNS (Domain Semantic Name Service) zone file entries from CRD events. `deriveRecords()` L274 dispatches to per-GVK derivation functions: `deriveTalosClusterRecords()` L319, `derivePackInstanceRecords()` L436, `deriveRunnerConfigRecords()` L481, `deriveIdentityBindingRecords()` L379, `deriveIdentityProviderRecords()` L409. `refreshNSGlue()` L170 reloads NS glue records. |
| `outcome_reconciler.go:36` | `OutcomeReconciler`, `Reconcile()` L43 | Watches terminal outcome conditions on CRs; `classifyTerminalOutcome()` L124 determines outcome type (success/failure). |

---

## 4. Shared Libraries (`pkg/`)

### `pkg/lineage/`

| File | Key symbols |
|------|-------------|
| `descendant.go:33` | `IndexName(kind, name string) string` -- returns `"lineage-{kind}-{name}"` (lowercase). Mirrors `lineageIndexName()` in controller. |
| `descendant.go:52` | `SetDescendantLabels(obj, iliName, iliNamespace, operator string, rationale CreationRationale, actorRef string)` -- labels derived objects for lineage tracking |
| `chain.go:20` | `SealedCausalChain` struct -- immutable causal chain field embedded in every root declaration CRD |
| `chain.go:54` | `OperatorIdentity` struct -- identifies the Seam Operator that authored a derived object; embedded in `SealedCausalChain` |
| `rationale.go:17` | `CreationRationale` type + const enumeration: `ClusterProvision`, `ClusterDecommission`, `SecurityEnforcement`, `PackExecution`, `VirtualizationFulfillment`, `ConductorAssignment`, `VortexBinding`. New values require seam-core PR + Platform Governor review (Decision 5, SC-INV step 4b). |

### `pkg/conditions/`

| File | Key symbols |
|------|-------------|
| `conditions.go:22` | `SetCondition(conditions *[]metav1.Condition, conditionType, status, reason, message string, observedGeneration int64)` |
| `conditions.go:54` | `FindCondition(conditions []metav1.Condition, conditionType string) *metav1.Condition` |
| `validate.go:272` | `ValidateCondition(conditionType, reason string) error` -- used in unit tests to assert every condition emission uses a registered type+reason pair |
| `validate.go:293` | `KnownConditionTypes() []string` |
| `validate.go:304` | `ValidReasonsFor(conditionType string) []string` |

---

## 5. CRD Generation

`make generate-crd` runs controller-gen, outputs to `config/crd/`. CRD YAML files are the installed artifacts. seam-core CRD bundle is applied before all operators (SC-INV-003).

Enable bundle phase 00: `lab/configs/ccs-mgmt/compiled/enable/00-infrastructure-dependencies/seam-core-crds.yaml` installs all seam-core CRDs on management cluster.

---

## 6. Invariants

| ID | Rule | Location |
|----|------|----------|
| SC-INV-001 | seam-core owns all cross-operator CRD definitions under `infrastructure.ontai.dev` | `api/v1alpha1/` |
| SC-INV-002 | Phase 2B complete (2026-04-25); no further governed migrations required | seam-core CLAUDE.md |
| SC-INV-003 | seam-core CRD manifests installed before all operators | enable bundle phase 00 |
| Decision G | All cross-operator CRD schemas exclusively owned by seam-core | `api/v1alpha1/*.go` |
| Decision 5 | `CreationRationale` is a compile-time enum; new values require PR + Governor review | `pkg/lineage/rationale.go` |

---

## 7. Open Items

**T-04a**: No `x-kubernetes-validations` in `InfrastructureTalosCluster` CRD. CEL rule for `mode=import` requiring non-empty `role` is absent. Required: add `+kubebuilder:validation:XValidation` marker to `InfrastructureTalosClusterSpec`, run `make generate-crd`.

**T-05 (conductor-side)**: `InfrastructurePackBuildCategory` enum (L9-14 in `packbuild_types.go`) and `KustomizeSource` struct (L34) exist in seam-core schema. `PackBuildInput` struct in `conductor/cmd/compiler/compile.go:498` has NO `Category` field -- dispatch is by nil-check on HelmSource/RawSource pointers. Schema is ahead of implementation.
