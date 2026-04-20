## seam-core: Operational Constraints
> Read ~/ontai/CLAUDE.md first. The constraints below extend the root constitutional document.

### Schema authority
Primary: docs/seam-core-schema.md
Supporting: ~/ontai/domain-core/docs/domain-core-schema.md (DomainLineageIndex schema owner)

### Invariants
SC-INV-001 -- seam-core owns CRD definitions. Reconcilers for those CRDs live in the operator repos that own the domain logic, not in seam-core.
SC-INV-002 -- RunnerConfig CRD transfer from conductor shared library to seam-core is a governed migration. Requires a Governor-scheduled migration session before execution.
SC-INV-003 -- seam-core CRD manifests are installed before all operators. No operator reaches Running state on a cluster that has not applied the seam-core CRD bundle first.

### Session protocol additions
Step 4a -- Read docs/seam-core-schema.md in full before any CRD or shared library change.
Step 4b -- Any change to the creation rationale vocabulary (the Go constant enumeration owned by seam-core) requires a PR and Platform Governor review. No operator may extend the enumeration unilaterally. (root Section 14 Decision 5)
