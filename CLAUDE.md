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

Release = push a `vX.Y.Z` tag (only after the change is merged to `main`). The `release-on-tag` workflow builds the executables + pushes the ghcr image (`ghcr.io/thetillhoff/webscan:X.Y.Z`, no `v`); Flux image-automation on hydra then scans, bumps the deployment, and rolls it out — usually within minutes, no manual step. Force it early with the `ci-build-then-flux-reconcile` skill. Renovate dep PRs auto-tag a patch bump via `tag-on-main`; human releases are tagged manually.
