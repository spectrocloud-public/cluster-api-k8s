# cabpck-cacpck — CAPI Upgrade memory (Protocol U)

Per-repo state for the CAPI-upgrade reconcile engine
(`spectro-capx-upgrade:capi-repo-upgrade`).

## Layout

- `conf/capi-upgrade.conf` — identity, sources, classifier inputs, integration tail
- `decisions/known-commits.tsv` — authoritative SHA→category (+ optional decision col)
- `decisions/skip-commits.txt` — explicit human skip list
- `decisions/resolve-commits.txt` — explicit PICK-RESOLVE overrides
- `decisions/regwatch-risk.tsv` — durable Regression-Watch risk (survives re-plans)
- `notes/`, `runtime/` — free-form + generated artifacts

## First cycle — fresh integration line

The previous `spectro-master` (v0.4.2 / palette 4.10.0 line) was archived
server-side as `origin/archive-spectro-master-4.10.0`. The new
`origin/spectro-master` starts fresh at CAPI v1.13.2 (v1beta1→v1beta2 base,
same commit as `origin/capi-bump` = canonical/cluster-api-k8s#236 HEAD; a
durable copy `origin/capi-v1.13.2-base` is pinned so the target survives
that PR's lifecycle). Essential spectro commits (webhook splits first, then
other essentials, then feature-specific tweaks) are cherry-picked onto this
fresh baseline in a curated order over time.

- **Source** (`SPECTRO_SRC`): `origin/archive-spectro-master-4.10.0` — the
  archive is enumerated by the reconcile engine as the pool of candidate
  commits; humans curate the priority.
- **Target**: local tag `v1.13.2` at `origin/capi-v1.13.2-base` HEAD. If
  canonical later publishes a real `v1.13.2` tag at the same SHA (they will
  if PR #236 merges as-is) the local tag lines up.

## Review reference

Every commit on `origin/archive-spectro-master-4.10.0` is a candidate for
this cycle. The planner surfaces each with a PICK/PICK-VERIFY/SKIP/NEEDS-
DECISION verdict — treat those as suggestions, not the final priority.
Essential first (webhook splits), then others as required.
