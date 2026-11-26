# webscan

CLI + web scanner. `cmd/webscan` (CLI), `cmd/webscan-web` (server). Scan engine: `pkg/webscan/`, per-scan packages `pkg/*Scan/`.

## Hosted deployment (hydra)

Runs on the hydra cluster. Use `KUBECONFIG=~/code/infra/pulumi/kubeconfig` and `--context=admin@hydra`.

- Namespace: `webscan`. Manifests: `~/code/infra/kubernetes/apps/hydra/webscan/` (Flux/GitOps, image bumped by image automation).
- Web mode queues scan jobs. Each job's status, result and output live in **redis** (`deploy/redis`), keyed `webscan:job:<id>`.

Check a stuck or silent scan:

```bash
kubectl logs deploy/webscan -n webscan --context=admin@hydra --tail=100
kubectl exec deploy/redis -n webscan --context=admin@hydra -- redis-cli KEYS 'webscan:job:*'
kubectl exec deploy/redis -n webscan --context=admin@hydra -- redis-cli HGETALL webscan:job:<id>
```

On timeout the worker writes `status=timeout` to redis but logs nothing. Check the job's redis fields (`status`, `error`, `duration`) before you assume a crash.

## Version and releasing

- Released tags can run **ahead** of local `main`. Before you fix prod behavior, run `git fetch --tags` and check out the deployed tag: the fix may exist.
- `main` is protected (`enforce_admins`, required checks `go-test` and `verify-version`). Every change lands through a PR: `gh pr merge <n> --merge --delete-branch`. A direct push is rejected.
- Release: push a `vX.Y.Z` tag after the change is merged to `main`. The `release-on-tag` workflow builds the executables and pushes `ghcr.io/thetillhoff/webscan:X.Y.Z` (no `v`). Flux image automation then rolls it out. To force it early, use the `ci-build-then-flux-reconcile` skill.
- Renovate dependency updates release themselves: `tag-on-main` tags the patch bump and passes `release_body: "Updated dependencies."`. They get no `CHANGELOG.md` entry, so the changelog skips those version numbers. Feature and fix releases are tagged by hand and do carry a CHANGELOG section.
- A deleted tag's version number is spent. `proxy.golang.org` caches the first SHA of each `vX.Y.Z`. After a failed release, fix the cause and release the next patch number. Never re-push the same tag.
- `GITHUB_TOKEN` never triggers workflows. A bot-authored commit or PR gets no CI, so it can never pass the required check. That is why `tag-on-main` calls the release workflow directly (`trigger-release`). For automation that must land on `main`, use a PAT or GitHub App token.
