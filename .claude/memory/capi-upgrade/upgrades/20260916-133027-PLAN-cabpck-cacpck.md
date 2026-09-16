# CAPI reconcile — cabpck-cacpck — 2026-09-16 — mode=plan

> **Fresh-baseline curated replay — not a naive full replay.** Special case: the "upstream target"
> `v1.13.2` is a local tag pinned to the tip of the unmerged upstream PR
> [canonical/cluster-api-k8s#236](https://github.com/canonical/cluster-api-k8s/pull/236) (CAPI v1beta1→v1beta2
> bump), captured as `origin/capi-v1.13.2-base` so this cycle survives PR #236's eventual fate.
> `spectro-master` was reset to the same commit; the previous spectro fork history lives on
> `origin/archive-spectro-master-4.10.0` (SPECTRO_SRC). The maintainers explicitly want a **fresh
> baseline at v1.13.2 with commits landed in priority order over time — not a blind full replay**.
>
> ### Priority tiers (curation intent, first-time triage)
>
> **Tier 1 — Webhook splits (PCP-4817 series, land first):**
> - `29002e71` PCP-4817 Controller & Webhook Seperation — NEEDS-DECISION (go.mod conflict)
> - `79576dd8` PCP-4817 Manifest Generation — NEEDS-DECISION (MIXED — split code vs regen)
> - `991aa7c0` PCP-4817 Watch Namespace change — PICK-VERIFY
> - `e71debeb` PCP-4817 Additional Fixes — NEEDS-DECISION (MIXED — split code vs regen)
>
> **Tier 2 — Essentials (CI, Dockerfile, Makefile, build/CVE bumps):**
> - CI/build: `526f280f` PCP-4889 GH action · `10e0483a` PCP-4889 remove govulncheck · `f9fdfee8`/`3c6215a3`/`41d153d1` Go builder bumps · `fd330ce1` golang alpine image · `facc2ecf` PCP-7232 Go 1.26.5 · `90e1543c` PCP-7446 Go 1.26.6 · `b3398d6e` Makefile image tags
> - Deps/CVEs: `adec40b9` PCP-5293 controller-runtime v0.20.4 · `dd9205f2` PCP-6722 vulns · `a42ff341` PCP-6787 CVEs · `5ff9002d` PCP-6891 vulns · `cac840ae` PCP-7446 grpc bump
> - Fixes/policy: `122375ba` PCP-6722 non-default CNI CP-ready · `ca8a59b6` PCP-7027 mask containerd/kubelet [PICK] · `f207d93f` PCP-7250 map-assignment fix · `347a119b` and `0f9fa4f7` CODEOWNERS [PICK]
>
> **Tier 3 — Feature-specific / special (defer or evaluate separately):**
> - `0d3289d3` PCP-4817 sync-period base manifests — NEEDS-DECISION (MIXED)
> - `90bd017a` tolerations to CK8s deployments — NEEDS-DECISION (MIXED)
>
> **What Gate-1 approves:** the 21 NEEDS-DECISION rows below + the 5 Regression-Watch entries + the
> frozen `source_sha`. Reviewer/human should focus on Tier 1 first; Tier 2 largely resolves once the
> Go/deps baseline is settled at v1.13.2; Tier 3 can defer.


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
5 parity-review candidate(s) — reviewer/human must confirm the target covers the fork's behavior; set risk L/M/H. **High risk blocks auto-progression and requires explicit human sign-off.**

| sha | decision | files (parity-check) | risk | confirm |
|--|--|--|--|--|
| `79576dd8` | MIXED | spectro/base/bootstrap/delete-namespace.yaml spectro/base/bo | REVIEW | yes |
| `991aa7c0` | PICK-VERIFY | bootstrap/main.go controlplane/main.go  | REVIEW | yes |
| `e71debeb` | MIXED | spectro/base/bootstrap/kustomization.yaml spectro/base/contr | REVIEW | yes |
| `0d3289d3` | MIXED | spectro/base/bootstrap/kustomization.yaml spectro/base/contr | REVIEW | yes |
| `90bd017a` | MIXED | spectro/base/bootstrap/kustomization.yaml spectro/base/contr | REVIEW | yes |

## Decisions (baseline=T · cumulative rehearsal)
| # | sha | cat | subject | decision | reason |
|--|--|--|--|--|--|
| 1 | `29002e71` | F4? | PCP-4817 Controller & Webhook Seperation | NEEDS-DECISION | conflict on: go.mod  |
| 2 | `79576dd8` | MIXED | PCP-4817 Manifest Generation | NEEDS-DECISION | mixed generated+code — split: regen the generated part, pick the code part (human) |
| 3 | `991aa7c0` | F4? | PCP-4817 Watch Namespace change | PICK-VERIFY | clean apply but touches upstream-changed/F4-hotspot path — verify |
| 4 | `e71debeb` | MIXED | PCP-4817 Additional Fixes | NEEDS-DECISION | mixed generated+code — split: regen the generated part, pick the code part (human) |
| 5 | `0d3289d3` | MIXED | PCP-4817 Adding sync period to base manifests | NEEDS-DECISION | mixed generated+code — split: regen the generated part, pick the code part (human) |
| 6 | `526f280f` | F4? | PCP-4889 Create github action for generating CK8s bootstrap and contro | NEEDS-DECISION | conflict on: go.mod go.sum  |
| 7 | `10e0483a` | M* | PCP-4889 Removing govulncheck scan | NEEDS-DECISION | conflict on: Dockerfile  |
| 8 | `90bd017a` | MIXED | Adding tolerations to CK8s deployments | NEEDS-DECISION | mixed generated+code — split: regen the generated part, pick the code part (human) |
| 9 | `b3398d6e` | F4? | Makefile changes for image tags | NEEDS-DECISION | conflict on: Makefile  |
| 10 | `adec40b9` | F4? | PCP-5293 Upgrade controller-runtime version to v0.20.4 (#5) | NEEDS-DECISION | conflict on: bootstrap/api/v1beta2/ck8sconfig_webhook.go bootstrap/api/v1beta2/ck8sconfigtemplate_webhook.go controlplane/api/v1beta2/ck8scontrolplane_webhook.go go.mod go.sum  |
| 11 | `f9fdfee8` | M* | Update Go version in spectro-release workflow to 1.24.12 | NEEDS-DECISION | conflict on: .github/workflows/spectro-release.yaml  |
| 12 | `3c6215a3` | M* | Update builder golang version to 1.24.13 (#7) | NEEDS-DECISION | conflict on: .github/workflows/spectro-release.yaml  |
| 13 | `122375ba` | F4? | PCP-6722 Support for non default CNI by letting CK8s CP to be ready wh | NEEDS-DECISION | conflict on: controlplane/controllers/ck8scontrolplane_controller.go pkg/ck8s/config_init.go  |
| 14 | `dd9205f2` | M* | PCP-6722 Fixing vulns | NEEDS-DECISION | conflict on: go.mod go.sum  |
| 15 | `41d153d1` | M* | Update go version for build | NEEDS-DECISION | conflict on: .github/workflows/spectro-release.yaml  |
| 16 | `fd330ce1` | M* | Update golang alpine image path | NEEDS-DECISION | conflict on: Dockerfile  |
| 17 | `a42ff341` | M* | PCP-6787 Fix CVEs | NEEDS-DECISION | conflict on: go.mod go.sum  |
| 18 | `5ff9002d` | M* | PCP-6891 Updating go build version & fixing vulns | NEEDS-DECISION | conflict on: .github/workflows/spectro-release.yaml go.mod go.sum  |
| 19 | `ca8a59b6` | F2? | PCP-7027: mask existing containerd and kubelet to make sure it doesn't | PICK | applies cleanly (fork-only) [confidence=heuristic] |
| 20 | `facc2ecf` | F4? | PCP-7232: Bump Go builder to 1.26.5 | NEEDS-DECISION | conflict on: .github/workflows/spectro-release.yaml Makefile  |
| 21 | `f207d93f` | F4? | PCP-7250: fixed the direct map assignment of ck8sconfig object to boot | NEEDS-DECISION | conflict on: pkg/ck8s/config_init.go  |
| 22 | `90e1543c` | M* | PCP-7446 : bump Go toolchain to 1.26.6 and update deps to fix CVEs | NEEDS-DECISION | conflict on: .github/workflows/spectro-release.yaml Makefile go.mod go.sum  |
| 23 | `cac840ae` | M* | PCP-7446 : bump google.golang.org/grpc to 1.82.1 | NEEDS-DECISION | conflict on: go.mod go.sum  |
| 24 | `347a119b` | M* | chore: set CODEOWNERS to @snehala27 @vishu2498 | PICK | applies cleanly (fork-only) [confidence=heuristic] |
| 25 | `0f9fa4f7` | M* | chore: set CODEOWNERS to @snehala27 @vishu2498 | PICK | applies cleanly (fork-only) [confidence=heuristic] |

## Summary
- PICK=4 (PICK-VERIFY=1 · PICK-RESOLVE=0) · SKIP=0 · REGENERATE=0 · NEEDS-DECISION=21
- INCOMPLETE: unresolved=21 · regen/verify-pending=0 · regression-watch=5 · commit-failures=0
- Branch NOT mergeable while INCOMPLETE>0 or any High regression unconfirmed.
- Legend: PICK auto · PICK-VERIFY auto+verify · PICK-RESOLVE force+resolver · SKIP · NEEDS-DECISION human · REGEN make · MIXED split(human)
