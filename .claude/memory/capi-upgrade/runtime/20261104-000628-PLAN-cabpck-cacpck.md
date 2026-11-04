# CAPI reconcile — cabpck-cacpck — 2026-11-04 — mode=plan

| field | value |
|--|--|
| source S | `origin/spectro-master` @ a5c91c47 |
| target T | `v1.13.4` @ 25b67ac5 |
| fork-point F | 11419404 |
| upstream U | upstream/main |
| CAPI | v1.13.2 -> v1.13.4 · go 1.26.7 · contract v1beta3 · axis 2(CONTRACT-MIGRATION: flag, do NOT auto-perform) · n-3 PASS(gap=0) |
| new branch | `spectro-v1.13.4-master` |
| engine hash | a0585e3f5714ec09e54d560657fff58a |

## ⚠️ Regression Watch
12 parity-review candidate(s) — reviewer/human must confirm the target covers the fork's behavior; set risk L/M/H. **High risk blocks auto-progression and requires explicit human sign-off.**

| sha | decision | files (parity-check) | risk | confirm |
|--|--|--|--|--|
| `458fa8b9` | SKIP-upstreamed | .idea/.gitignore .idea/cluster-api-k8s.iml .idea/go.imports. | High | yes |
| `a264cd90` | SKIP-upstreamed | bootstrap/api/v1beta2/zz_generated.deepcopy.go bootstrap/con | Low | yes |
| `7f14a99e` | PICK-VERIFY | bootstrap/main.go controlplane/main.go  | Low | yes |
| `a34d16ab` | PICK-VERIFY | bootstrap/main.go controlplane/main.go  | Med | yes |
| `1f998988` | PICK-VERIFY | .github/workflows/spectro-release.yaml Dockerfile Makefile  | Low | yes |
| `54624153` | PICK-RESOLVE | go.mod go.sum   | High | yes |
| `ddde6e01` | PICK-VERIFY | Makefile  | Low | yes |
| `ff5d4956` | PICK-VERIFY | pkg/ck8s/config_init.go  | Low | yes |
| `658dd334` | PICK-RESOLVE | controlplane/api/v1beta2/ck8scontrolplane_types.go  | Med | yes |
| `658dd334` | CRD/API | controlplane/api/v1beta2/ck8scontrolplane_types.go  | Med | yes |
| `f974bfde` | PICK-VERIFY | Makefile  | Low | yes |
| `a962eb3c` | PICK-RESOLVE | .github/workflows/spectro-release.yaml Makefile go.mod go.su | Low | yes |

## Decisions (baseline=T · cumulative rehearsal)
| # | sha | cat | subject | decision | reason |
|--|--|--|--|--|--|
| 1 | `458fa8b9` | F1 | Bump CAPI to v1.13.2 | SKIP | upstreamed (verify) — authoritative map |
| 2 | `a264cd90` | F1 | Fix e2e tests | SKIP | upstreamed (verify) — authoritative map |
| 3 | `7f14a99e` | F4? | PCP-4817 Watch Namespace change | PICK-VERIFY | clean apply but touches upstream-changed/F4-hotspot path — verify |
| 4 | `d5f7260e` | F2? | PCP-7027: mask existing containerd and kubelet to make sure it doesn't | PICK | applies cleanly (fork-only) [confidence=heuristic] |
| 5 | `c5b6fec4` | M* | chore: set CODEOWNERS to @snehala27 @vishu2498 | PICK | applies cleanly (fork-only) [confidence=heuristic] |
| 6 | `d1a2de7c` | M* | chore: set CODEOWNERS to @snehala27 @vishu2498 | PICK | applies cleanly (fork-only) [confidence=heuristic] |
| 7 | `a34d16ab` | F4? | resolve 29002e71: PCP-4817 Controller & Webhook Seperation | PICK-VERIFY | clean apply but touches upstream-changed/F4-hotspot path — verify |
| 8 | `6e1d824d` | F2 | resolve 79576dd8: PCP-4817 Manifest Generation (overlays only; spectro | PICK | applies cleanly (fork-only) |
| 9 | `9eac332a` | F2 | resolve e71debeb: PCP-4817 Additional Fixes | PICK | applies cleanly (fork-only) |
| 10 | `38bfeac8` | F2 | resolve 0d3289d3: PCP-4817 Adding sync period to base manifests | PICK | applies cleanly (fork-only) |
| 11 | `ca50d75c` | F2 | resolve 90bd017a: Adding tolerations to CK8s deployments | PICK | applies cleanly (fork-only) |
| 12 | `1f998988` | F4? | resolve 526f280f: PCP-4889 Create github action for generating CK8s bo | PICK-VERIFY | clean apply but touches upstream-changed/F4-hotspot path — verify |
| 13 | `ce16a90e` | M* | resolve 10e0483a: PCP-4889 Removing govulncheck scan | PICK | applies cleanly (fork-only) |
| 14 | `3b657c0a` | M* | resolve fd330ce1: Update golang alpine image path | PICK | applies cleanly (fork-only) |
| 15 | `54624153` | M* | resolve 90e1543c: PCP-7446 bump Go toolchain to 1.26.6 and update deps | PICK-RESOLVE | conflict — human-approved resolve; resolver 3-way merge PENDING: go.mod go.sum  |
| 16 | `ddde6e01` | F4? | resolve b3398d6e: Makefile changes for image tags | PICK-VERIFY | clean apply but touches upstream-changed/F4-hotspot path — verify |
| 17 | `ff5d4956` | F4? | resolve f207d93f: PCP-7250 fixed the direct map assignment of ck8sconf | PICK-VERIFY | clean apply but touches upstream-changed/F4-hotspot path — verify |
| 18 | `36a7693f` | REGEN | regen: bump controller-gen v0.16.1->v0.21.0 (Go 1.26 compat) + re-deri | REGENERATE | generated artifact — regenerate via make (not cherry-picked) |
| 19 | `658dd334` | F4 | Add omitzero to CK8sControlPlaneMachineTemplate metadata (#36) | PICK-RESOLVE | force-pick (human-approved; applied cleanly) |
| 20 | `f974bfde` | F4? | Add fork in tag (#37) | PICK-VERIFY | clean apply but touches upstream-changed/F4-hotspot path — verify |
| 21 | `3d31fde2` | F2? | capx: seed .claude/memory/capi-upgrade (Protocol U engine memory) (#33 | PICK | applies cleanly (fork-only) [confidence=heuristic] |
| 22 | `a962eb3c` | M* | Fix vulns (#40) | PICK-RESOLVE | conflict — human-approved resolve; resolver 3-way merge PENDING: .github/workflows/spectro-release.yaml Makefile go.mod go.sum  |
| 23 | `a5c91c47` | M* | chore(capx): refresh the memory tree scaffold + correct the Spectro ba | PICK | applies cleanly (fork-only) [confidence=heuristic] |

## Summary
- PICK=18 (PICK-VERIFY=6 · PICK-RESOLVE=3) · SKIP=2 · REGENERATE=1 · NEEDS-DECISION=0
- INCOMPLETE: unresolved=2 · regen/verify-pending=1 · regression-watch=12 · commit-failures=0
- Branch NOT mergeable while INCOMPLETE>0 or any High regression unconfirmed.
- Legend: PICK auto · PICK-VERIFY auto+verify · PICK-RESOLVE force+resolver · SKIP · NEEDS-DECISION human · REGEN make · MIXED split(human)
