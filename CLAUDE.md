# CLAUDE.md — seam-core
> Inherits from: ~/ontai/CLAUDE.md (read that first, always)

---

## Project Identity

seam-core is the schema controller. It owns cross-operator CRD definitions that
no single operator owns. It is the authoritative home for CRDs that multiple
operators depend on.

CRDs owned here:
- RunnerConfig — produced by Platform and Wrapper, reconciled by Conductor
- InfrastructurePolicy — produced by humans/guardian, reconciled by Guardian
- InfrastructureProfile — reconciled by Guardian

seam-core has no execution logic and no capability engine. It installs CRD
definitions and runs a schema controller that validates CRD schema versions.

Component:    seam-core
Org:          github.com/ontai-dev
Domain:       infrastructure.ontai.dev
Image:        registry.ontai.dev/ontai-dev/seam-core:<semver>
Dev image:    10.20.0.1:5000/ontai-dev/seam-core:dev
Working dir:  ~/ontai/seam-core

---

## Schema Reference

Primary: All schema documents in ~/ontai/ — seam-core is the CRD registry.
The seam-core engineer reads ALL schema documents before any implementation work.

---

## Project Invariants

SC-INV-001 — seam-core owns CRD definitions. It does not own reconcilers.
Each CRD's reconciler lives in the appropriate component operator.

SC-INV-002 — RunnerConfig CRD ownership transfer from conductor shared library
to seam-core is a governed migration. It requires a dedicated Controller Engineer
session scheduled by the Platform Governor. It is not a documentation gap — it has
implementation implications across conductor, platform, and wrapper.

SC-INV-003 — seam-core installs before all operators. No operator starts without
its CRDs present.

---

## Session Protocol Addition

Step 4a — Read all schema documents before any CRD definition work.
Step 4b — Before adding a CRD, confirm it has no existing owner in another operator.
Step 4c — CRD schema changes require a major version bump in the CRD group version.

---

*seam-core component constitution*
*Amendments appended below with date and rationale.*
