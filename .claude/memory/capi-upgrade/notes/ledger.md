# cluster-api-k8s upgrade ledger

## PCP-7766 — v1.13.4 (canonical/main @ 25b67ac, untagged upstream; local tag v1.13.4)

### Target anchor (Gate-1 approves this) — latest upstream commit, NOT a release tag
- T = `canonical/main@25b67ac55cbc75d149f69bff5aa33ca1963ba53b` (#257, on top of #240 CAPI 1.13.4 + #256).
  Untagged by choice: upstream's latest release v0.6.3 (8c4f3de) predates #240; we do not wait for a release.
- Engine needs a name for T (it derives `capi-plan-<t>`, `spectro-<t>-master` and Makefile TAG from it), so
  T is aliased by a LOCAL-ONLY lightweight tag `v1.13.4` (named after the CAPI version, matching the
  `v1.13.2-fork-spectro-*` scheme). Never pushed — the fork only receives capi-plan-v1.13.4,
  spectro-v1.13.4-master and the Spectro release tags (v1.13.4-fork-spectro-<N>).
- Reproduce in any clone:  `git fetch upstream && git tag v1.13.4 25b67ac55cbc75d149f69bff5aa33ca1963ba53b`
  (summary.json / plan-latest.json pin target_sha, so a mis-pointed alias is caught).
- ALWAYS pass `--target v1.13.4` explicitly: with no --target the engine auto-picks the highest local vX.Y.Z
  tag, which in a clone carrying this alias is v1.13.4, not upstream's v0.6.3.
- Remove after Gate 2: `git tag -d v1.13.4`.
- 458fa8b/a264cd9 (our closed canonical#236) SKIP: superseded by canonical#239+#240+#256.
  Upstream added CK8s API v1beta3 (storage, CAPI v1beta2 contract); v1beta2 kept served at the
  old shape. Our v1beta2 schema == upstream v1beta2 schema (all 4 CRDs); v1beta3 spec == our v1beta2
  spec, status-only diffs. Full audit: posted on the plan PR.
- Execute-time fork fix: metadata.yaml release series 0.7 -> contract v1beta2 (upstream lists "v1beta3").
- 36a7693 REGENERATE: upstream already at controller-gen v0.21.0; re-derive spectro/generated.
- 5462415/a962eb3 PICK-RESOLVE: conflict on go.mod/go.sum vs T; resolve = per-module max(fork, T) (all fork bumps >= T, no downgrades) + keep Go 1.26.7 release pins; go mod tidy.

### 658dd33 (PCP-7634 omitzero) — human decision 2026-11-03: FORK-ONLY, not upstreamed
- Ledger: `658dd334 F4 PICK-RESOLVE` (lands the original v1beta2 tag unchanged; harmless on T).
- Why a fork delta is needed: on T, `minProperties: 1` on machineTemplate.metadata exists ONLY in the
  v1beta3 schemas (ck8scontrolplanes .spec.machineTemplate.metadata; ck8scontrolplanetemplates
  .spec.template.spec.machineTemplate.metadata). Upstream v1beta2 now embeds CAPI core/v1beta1
  ObjectMeta (no minProperties), so the original commit alone is a no-op on T. Without omitzero on
  v1beta3, an empty ObjectMeta marshals as `metadata: {}` (conversion webhook v1beta2->v1beta3,
  defaulting responses, and any client on v1beta3 types) and fails validation.

## Execute-time fork deltas (approved at Gate 1; executor applies, reviewed at Gate 2)
1. `PCP-7634: omitzero on v1beta3 CK8sControlPlaneMachineTemplate metadata` — in
   controlplane/api/v1beta3/ck8scontrolplane_types.go (CK8sControlPlaneMachineTemplate.ObjectMeta):
   `json:"metadata,omitempty"` -> `json:"metadata,omitempty,omitzero"`. Covers both CRDs (shared Go
   field). Then `make generate manifests` + spectro/generated regen (expect no CRD schema change).
   Verify (T2): create CK8sControlPlane AND CK8sControlPlaneTemplate with no machineTemplate.metadata
   via v1beta2 and via v1beta3 -> accepted; stored object has no `metadata: {}`.
2. `metadata.yaml`: release series 0.7 -> `contract: v1beta2` (upstream lists the non-contract "v1beta3").
