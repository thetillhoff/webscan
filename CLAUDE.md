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

`main` is a **protected branch**: direct `git push origin main` is rejected (`protected branch hook declined`; `enforce_admins` on, required check `verify-version`). All changes land via PR — `gh pr merge <n> --merge --delete-branch`. There is no local-merge-and-push path.

Release = push a `vX.Y.Z` tag (only after the change is merged to `main`). The `release-on-tag` workflow builds the executables + pushes the ghcr image (`ghcr.io/thetillhoff/webscan:X.Y.Z`, no `v`); Flux image-automation on hydra then scans, bumps the deployment, and rolls it out — usually within minutes, no manual step. Force it early with the `ci-build-then-flux-reconcile` skill.

**All releases are tagged manually.** Renovate dep updates land on `main` by themselves — `automergeType: branch` fast-forwards `main` to a SHA that already passed `verify-version` on the `renovate/**` branch, which is why protection does not block them — but the auto **patch release** does not happen. `tag-on-main`'s `update-changelog` job tries to push a new CHANGELOG commit to `main`, and protection rejects it, because that commit has no `verify-version` run and never can.

That is the general rule worth remembering here: **`GITHUB_TOKEN` never triggers workflows.** A commit, branch or PR it creates gets no CI, so on this repo it can never satisfy the required check — no bot-authored commit or PR is mergeable — and a tag it pushes does not fire `on: push: tags:`, which is why `tag-on-main` calls the release workflow directly via `trigger-release`. Automating anything that has to land on `main` needs a PAT or GitHub App token instead of `GITHUB_TOKEN`.
