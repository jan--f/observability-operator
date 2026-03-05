# Release Workflow

## Important pointers

* Releases are cut from `release-MAJOR.MINOR` branches (e.g. `release-1.4`).
* The `olm/index-template.yaml` file in each branch is the source of truth for
  that branch's OLM catalog history. It lives directly in the repository (both
  `main` and `release-*` branches) and is updated via automated PRs opened by CI.
* `olm/update-channels.sh` modifies the template in place during CI; the
  calling workflow commits the result and opens a PR.

## Channel mapping

| Channel | Trigger | Source branch | PR target |
|---|---|---|---|
| `development` | push to `main` (non-release commit) | `main` | `main` |
| `candidate` | `chore(release):` push to `release-*` | `release-X.Y` | `release-X.Y` |
| `stable` | GitHub release promoted from pre-release | `release-X.Y` | `release-X.Y` |

## Release Workflow

![Release Workflow](./assets/release.png)

NOTE: the source for the UML can be found under [assets directory](./assets/release.uml)
