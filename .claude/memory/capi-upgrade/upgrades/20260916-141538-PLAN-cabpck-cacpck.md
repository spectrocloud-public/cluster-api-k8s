# CAPI reconcile — cabpck-cacpck — 2026-09-16 — mode=plan

| field | value |
|--|--|
| source S | `origin/archive-spectro-master-4.10.0` @ 0f9fa4f7 |
| target T | `v1.13.2` @ a264cd90 |
| fork-point F | dad1d76b |
| upstream U | upstream/main |
| CAPI | v1.10.4 -> v1.13.2 · go 1.25.0 · contract v1beta1 · axis 1(version-bump) · n-3 PASS(gap=3) |
| new branch | `spectro-v1.13.2-master` |
| engine hash | 2a24e78d29bf65808c95669f60bf123e |
| ⚠ merges in F..S (evil-merge check) | 5933d9c39028bf9e8deb429ba963800feba00b50 17d5993cd02970bbb359d634736c63fe514253d2 3b95a4c2847bc464003d356ad495423e43dd1022 21cedf56a659678f664cfaf6a38d44296e802d2d 686f98b7a742d310ca6ce804a36894ec22e8ec01 fe4d3866f4f903d97f9be6252ed68ab3c4e20e34 472f6fd26a5a4349a6eb23229db0e4d5e98ef4c3 4a5a9c5392d803a375bec59a0fa8fd2e5fc09e2e d3f4e96167a76d1b51d7cfc05c14ff397bced778 62ebac9495e696a8721c1238577d77a9bd1a91de  |

## ⚠️ Regression Watch
12 parity-review candidate(s) — reviewer/human must confirm the target covers the fork's behavior; set risk L/M/H. **High risk blocks auto-progression and requires explicit human sign-off.**

| sha | decision | files (parity-check) | risk | confirm |
|--|--|--|--|--|
| `29002e71` | PICK-RESOLVE | go.mod   | High | yes |
| `79576dd8` | MIXED | spectro/base/bootstrap/delete-namespace.yaml spectro/base/bo | High | yes |
| `991aa7c0` | PICK-VERIFY | bootstrap/main.go controlplane/main.go  | Low | yes |
| `e71debeb` | MIXED | spectro/base/bootstrap/kustomization.yaml spectro/base/contr | Med | yes |
| `0d3289d3` | MIXED | spectro/base/bootstrap/kustomization.yaml spectro/base/contr | Med | yes |
| `526f280f` | PICK-RESOLVE | go.mod go.sum   | High | yes |
| `10e0483a` | PICK-RESOLVE | Dockerfile   | Low | yes |
| `90bd017a` | MIXED | spectro/base/bootstrap/kustomization.yaml spectro/base/contr | Med | yes |
| `b3398d6e` | PICK-RESOLVE | Makefile   | Med | yes |
| `fd330ce1` | PICK-RESOLVE | Dockerfile   | Med | yes |
| `f207d93f` | PICK-RESOLVE | pkg/ck8s/config_init.go   | Med | yes |
| `90e1543c` | PICK-RESOLVE | .github/workflows/spectro-release.yaml Makefile go.mod go.su | High | yes |

## Decisions (baseline=T · cumulative rehearsal)
| # | sha | cat | subject | decision | reason |
|--|--|--|--|--|--|
| 1 | `29002e71` | F4? | PCP-4817 Controller & Webhook Seperation | PICK-RESOLVE | conflict — human-approved resolve; resolver 3-way merge PENDING: go.mod  |
| 2 | `79576dd8` | MIXED | PCP-4817 Manifest Generation | NEEDS-DECISION | mixed generated+code — split: regen the generated part, pick the code part (human) |
| 3 | `991aa7c0` | F4? | PCP-4817 Watch Namespace change | PICK-VERIFY | clean apply but touches upstream-changed/F4-hotspot path — verify |
| 4 | `e71debeb` | MIXED | PCP-4817 Additional Fixes | NEEDS-DECISION | mixed generated+code — split: regen the generated part, pick the code part (human) |
| 5 | `0d3289d3` | MIXED | PCP-4817 Adding sync period to base manifests | NEEDS-DECISION | mixed generated+code — split: regen the generated part, pick the code part (human) |
| 6 | `526f280f` | F4? | PCP-4889 Create github action for generating CK8s bootstrap and contro | PICK-RESOLVE | conflict — human-approved resolve; resolver 3-way merge PENDING: go.mod go.sum  |
| 7 | `10e0483a` | M* | PCP-4889 Removing govulncheck scan | PICK-RESOLVE | conflict — human-approved resolve; resolver 3-way merge PENDING: Dockerfile  |
| 8 | `90bd017a` | MIXED | Adding tolerations to CK8s deployments | NEEDS-DECISION | mixed generated+code — split: regen the generated part, pick the code part (human) |
| 9 | `b3398d6e` | F4? | Makefile changes for image tags | PICK-RESOLVE | conflict — human-approved resolve; resolver 3-way merge PENDING: Makefile  |
| 10 | `adec40b9` | F4? | PCP-5293 Upgrade controller-runtime version to v0.20.4 (#5) | SKIP | skip-list (explicit human decision) |
| 11 | `f9fdfee8` | M* | Update Go version in spectro-release workflow to 1.24.12 | SKIP | skip-list (explicit human decision) |
| 12 | `3c6215a3` | M* | Update builder golang version to 1.24.13 (#7) | SKIP | skip-list (explicit human decision) |
| 13 | `122375ba` | F4? | PCP-6722 Support for non default CNI by letting CK8s CP to be ready wh | SKIP | skip-list (explicit human decision) |
| 14 | `dd9205f2` | M* | PCP-6722 Fixing vulns | SKIP | skip-list (explicit human decision) |
| 15 | `41d153d1` | M* | Update go version for build | SKIP | skip-list (explicit human decision) |
| 16 | `fd330ce1` | M* | Update golang alpine image path | PICK-RESOLVE | conflict — human-approved resolve; resolver 3-way merge PENDING: Dockerfile  |
| 17 | `a42ff341` | M* | PCP-6787 Fix CVEs | SKIP | skip-list (explicit human decision) |
| 18 | `5ff9002d` | M* | PCP-6891 Updating go build version & fixing vulns | SKIP | skip-list (explicit human decision) |
| 19 | `ca8a59b6` | F2? | PCP-7027: mask existing containerd and kubelet to make sure it doesn't | PICK | applies cleanly (fork-only) [confidence=heuristic] |
| 20 | `facc2ecf` | F4? | PCP-7232: Bump Go builder to 1.26.5 | SKIP | skip-list (explicit human decision) |
| 21 | `f207d93f` | F4? | PCP-7250: fixed the direct map assignment of ck8sconfig object to boot | PICK-RESOLVE | conflict — human-approved resolve; resolver 3-way merge PENDING: pkg/ck8s/config_init.go  |
| 22 | `90e1543c` | M* | PCP-7446 : bump Go toolchain to 1.26.6 and update deps to fix CVEs | PICK-RESOLVE | conflict — human-approved resolve; resolver 3-way merge PENDING: .github/workflows/spectro-release.yaml Makefile go.mod go.sum  |
| 23 | `cac840ae` | M* | PCP-7446 : bump google.golang.org/grpc to 1.82.1 | SKIP | skip-list (explicit human decision) |
| 24 | `347a119b` | M* | chore: set CODEOWNERS to @snehala27 @vishu2498 | PICK | applies cleanly (fork-only) [confidence=heuristic] |
| 25 | `0f9fa4f7` | M* | chore: set CODEOWNERS to @snehala27 @vishu2498 | PICK | applies cleanly (fork-only) [confidence=heuristic] |

## Summary
- PICK=4 (PICK-VERIFY=1 · PICK-RESOLVE=7) · SKIP=10 · REGENERATE=0 · NEEDS-DECISION=4
- INCOMPLETE: unresolved=11 · regen/verify-pending=0 · regression-watch=12 · commit-failures=0
- Branch NOT mergeable while INCOMPLETE>0 or any High regression unconfirmed.
- Legend: PICK auto · PICK-VERIFY auto+verify · PICK-RESOLVE force+resolver · SKIP · NEEDS-DECISION human · REGEN make · MIXED split(human)
