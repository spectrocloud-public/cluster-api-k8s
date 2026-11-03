# cluster-api-k8s upgrade ledger

## PCP-7766 — v1.13.4 (canonical/main @ 25b67ac, untagged upstream; local tag v1.13.4)
- 458fa8b/a264cd9 (our closed canonical#236) SKIP: superseded by canonical#239+#240+#256.
  Upstream added CK8s API v1beta3 (storage, CAPI v1beta2 contract); v1beta2 kept served at the
  old shape. Our v1beta2 schema == upstream v1beta2 schema (all 4 CRDs); v1beta3 spec == our v1beta2
  spec, status-only diffs. Full audit: posted on the plan PR.
- 658dd33 omitzero: PICK + execute-time PORT to controlplane/api/v1beta3 (upstream v1beta3 keeps
  minProperties:1 on machineTemplate.metadata without omitzero).
- Execute-time fork fix: metadata.yaml release series 0.7 -> contract v1beta2 (upstream lists "v1beta3").
- 36a7693 REGENERATE: upstream already at controller-gen v0.21.0; re-derive spectro/generated.
- 5462415/a962eb3 PICK-RESOLVE: conflict on go.mod/go.sum vs T; resolve = per-module max(fork, T) (all fork bumps >= T, no downgrades) + keep Go 1.26.7 release pins; go mod tidy.
