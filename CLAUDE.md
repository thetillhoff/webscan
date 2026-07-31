# webscan

CLI + web scanner. `cmd/webscan` (CLI), `cmd/webscan-web` (server). Scan engine: `pkg/webscan/`, per-scan packages `pkg/*Scan/`.

## Hosted deployment (hydra)

Runs on the hydra cluster. Connect: see the `kubernetes-kind` skill (`KUBECONFIG=~/code/thetillhoff/infra/pulumi/kubeconfig`, `--context=admin@hydra`).

- Namespace: `webscan`. Manifests: `~/code/thetillhoff/infra/kubernetes/apps/hydra/webscan/` (Flux/GitOps, image bumped by image-automation).
- Web mode queues scan jobs; each job's status/result/output lives in **redis** (`deploy/redis`), keyed `webscan:job:<id>`.

Check a stuck/silent scan:

```bash
kubectl logs deploy/webscan -n webscan --context=admin@hydra --tail=100
kubectl exec deploy/redis -n webscan --context=admin@hydra -- redis-cli KEYS 'webscan:job:*'
kubectl exec deploy/redis -n webscan --context=admin@hydra -- redis-cli HGETALL webscan:job:<id>
```

On timeout the worker writes `status=timeout` to redis but logs nothing — an apparent "silent break" in the logs. Check the job's redis fields (`status`, `error`, `duration`) before assuming a crash.

## Version & releasing

Released tags can run **ahead** of local `main` (deployed prod may be a newer `vX.Y.Z` than `git describe`). Before editing to fix prod behavior, `git fetch --tags` and check out the deployed tag — the fix may already exist.

`main` is a **protected branch**: direct `git push origin main` is rejected (`protected branch hook declined`; `enforce_admins` on, required checks `go-test` and `verify-version`). All changes land via PR — `gh pr merge <n> --merge --delete-branch`. There is no local-merge-and-push path.

Release = push a `vX.Y.Z` tag (only after the change is merged to `main`). The `release-on-tag` workflow builds the executables + pushes the ghcr image (`ghcr.io/thetillhoff/webscan:X.Y.Z`, no `v`); Flux image-automation on hydra then scans, bumps the deployment, and rolls it out — usually within minutes, no manual step. Force it early with the `ci-build-then-flux-reconcile` skill.

**Renovate dep updates release themselves.** `automergeType: branch` fast-forwards `main` to a SHA that already passed `verify-version` on the `renovate/**` branch, which is why protection does not block them; `tag-on-main` then rebuilds, tags the patch bump and calls the release workflow. Those releases get **no CHANGELOG entry** — the commit adding one could not land on protected `main` (see below) — so `tag-on-main` passes `release_body: "Updated dependencies."` to the release workflow instead. Consequence: `CHANGELOG.md` skips the version numbers of dependency-only patches. Feature and fix releases are still tagged manually and do carry a CHANGELOG section, which `release_body` being empty makes the release workflow read.

**A deleted tag's version number is spent.** `tag-on-main`'s `verify-build` builds and version-checks a commit before tagging it, and the release workflow's `cleanup-on-failure` deletes the tag if the release fails before publishing anything. Neither licenses re-pushing that tag: `proxy.golang.org` permanently caches the first SHA it saw for `vX.Y.Z`, so pointing the version at a different commit gives every `go install` a checksum mismatch nothing on our side can clear. After a cleanup, fix the cause and release the next patch number.

That is the general rule worth remembering here: **`GITHUB_TOKEN` never triggers workflows.** A commit, branch or PR it creates gets no CI, so on this repo it can never satisfy the required check — no bot-authored commit or PR is mergeable — and a tag it pushes does not fire `on: push: tags:`, which is why `tag-on-main` calls the release workflow directly via `trigger-release`. Automating anything that has to land on `main` needs a PAT or GitHub App token instead of `GITHUB_TOKEN`.
