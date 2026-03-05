* [Developer documentation](#developer-documentation)
* [Contribution guidelines](#contribution-guidelines)
* [Release management](#release-management)

# Developer Documentation

## Development Tools

The build system assumes the following binaries are in the `PATH`:

```
make
git
go
npm
kind
podman (or docker)
```

Make sure you have installed on your local machine the Go version mentioned in
the root `go.mod` file

Once these tools are installed, run `make tools` to install all required
project dependencies to ``tmp/bin``

```sh
make tools
```

## Environment Setup

To setup the environment, first install project-specific tools, then run the unified environment setup script
`hack/setup-e2e-env.sh`. This script provides a consistent setup process used by both
local development and CI environments to prevent config drift.

```sh
# First, install project tools
make tools

# Then run the environment setup
./hack/setup-e2e-env.sh
```

The script does the following:
* Validates that project tools (operator-sdk, oc, etc.) are available from `make tools`
* Installs kind and kubectl if not already available
* Sets up a local Kind cluster
* Installs the Operator Lifecycle Manager (OLM) in the cluster
* Sets up a local registry to push the local operator and bundle images
* Installs monitoring CRDs

For advanced usage or CI integration, the script supports many options:

```sh
# First, install project tools
make tools

# Full setup with defaults (typical local development)
./hack/setup-e2e-env.sh

# Only validate prerequisites without setting up
./hack/setup-e2e-env.sh --validate-only

# Install additional packages (any system packages)
./hack/setup-e2e-env.sh curl jq tree htop git-lfs

# Use custom versions
./hack/setup-e2e-env.sh --kind-version v0.23.0 --kind-image kindest/node:v1.25.0
```

See `./hack/setup-e2e-env.sh --help` for all available options.

Once done, the cluster can be deleted by running:

```
kind delete cluster --name obs-operator
```

**Note:** The old `hack/kind/setup.sh` script is deprecated but still works for backward compatibility - it forwards to the new unified script.

## Running End to End tests

To run the E2E tests locally against the kind cluster that was setup following
the instructions above:

```sh
./test/run-e2e.sh
```

**NOTE:** `./test/run-e2e.sh --help` shows options that are useful when
rerunning tests.

## Running the Operator locally

Observability Operator relies heavily on the (forked) Prometheus Operator to do
most of the heavy lifting of creation of Prometheus and Alertmanager.  The
easiest way to use deploy prometheus operator is to run the
`observability-operator` bundle which installs both `observability-operator`
and `prometheus-operator`,  and then scale the `observability-operator`
deployment to 0, so that the operator can be  run out of cluster using `go run`

### Create the development Operator Bundle

The command below builds the operator + OLM bundle and pushes them to the
local-registry running in Kind cluster:

```sh
make operator-image bundle-image operator-push bundle-push  \
    IMAGE_BASE="local-registry:30000/observability-operator" \
    VERSION=0.0.0-dev  \
    PUSH_OPTIONS=--tls-verify=false
```

### Deploy the development Operator Bundle

Use `operator-sdk` to deploy the operator bundle:

```sh
./tmp/bin/operator-sdk run bundle \
    local-registry:30000/observability-operator-bundle:0.0.0-dev \
    --install-mode AllNamespaces \
    --namespace operators --skip-tls

```
Running the above should deploy operator and show

```
INFO[0044] OLM has successfully installed "observability-operator.v0.0.0-dev"

```

### Run the Operator from your local machine

Scale down the operator currently deployed in cluster:

```sh
kubectl scale --replicas=0 -n operators deployment/observability-operator
```

Start the operator locally:

```sh
# replace ~/.kube/config with your own KUBECONFIG path if different.
go run ./cmd/operator/... --zap-devel  --zap-log-level=100 --kubeconfig ~/.kube/config 2>&1 |
  tee tmp/operator.log
```

# Contribution guidelines

## Manifests and code generation

The Kubernetes CRDs and the ClusterRole needed for their management are
generated from the Go types in `pkg/apis`. Run `make generate` to regenerate the
Kubernetes manifests when changing these files.

The project uses [controller-gen](https://github.com/kubernetes-sigs/controller-tools/tree/master/cmd/controller-gen)
for code generation. For detailed information on the available code generation
markers, please refer to the controller-gen CLI page in
the [kubebuilder documentation](https://book.kubebuilder.io/reference/markers.html)

## Commit message convention

Commit messages need to comply to the [Conventional Commits specification](https://www.conventionalcommits.org/en/v1.0.0/)
and should be structured as follows:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

The type and description are used to generate a changelog and determine the
next release version.
Most commonly used types are:

* `fix:` a commit of the type fix patches a bug in your codebase. This
  correlates with PATCH in Semantic Versioning.

* `feat:` a commit of the type feat introduces a new feature to the codebase.
  This correlates with MINOR in Semantic Versioning.

* `BREAKING CHANGE:` a commit that has a footer BREAKING CHANGE:, or appends a
 `!` after the type/scope, introduces a breaking API change (correlating with
 MAJOR in Semantic Versioning).
 A BREAKING CHANGE can be part of commits of any type.

Other than `fix:` and `feat:`, the following type can also be used: `build:`,
`chore:`, `ci:`, `docs:`, `style:`, `refactor:`, `perf:` and `test:`.

# Release management

The project follows [SemVer 2.0.0](https://semver.org/)

```
Given a version number MAJOR.MINOR.PATCH, increment the:

MAJOR version when you make incompatible API changes,
MINOR version when you add functionality in a backwards compatible manner, and
PATCH version when you make backwards compatible bug fixes.
Additional labels for pre-release and build metadata are available as extensions to the MAJOR.MINOR.PATCH format.
```

Creating new releases is fully automated and requires minimal human
interaction. The changelog, release notes and release version are generated by
the CI based on the commits added since the latest release.

## Release branches

Releases are cut from dedicated `release-MAJOR.MINOR` branches (e.g. `release-1.4`).
Each release branch holds its own `olm/index-template.yaml` which accumulates
the catalog history for that release line (RC and stable entries).
The `main` branch holds the development catalog history.

CI publishes three channels:

| Channel | Trigger | Source branch |
|---|---|---|
| `development` | any push to `main` | `main` |
| `candidate` | `chore(release):` commit on `release-*` | `release-X.Y` |
| `stable` | GitHub release promoted from pre-release | `release-X.Y` |

For `candidate` and `stable`, CI opens a PR against the release branch with the
updated `olm/index-template.yaml`. Merging this PR records the catalog change.

## How to create a new release

### 1. Create a release branch (if it doesn't exist)

For a new minor version, branch off `main`:

```sh
git checkout main && git pull
git checkout -b release-1.5
git push origin release-1.5
```

For a patch release the `release-MAJOR.MINOR` branch already exists.

### 2. Commit the release on the release branch

```sh
git checkout release-1.5
git pull
git checkout -b cut-1.5.0
make initiate-release
```

This creates a `chore(release): X.Y.Z` commit that updates `CHANGELOG.md` and
`VERSION`. Review it, then push and open a PR **against the release branch**:

```sh
git push origin cut-1.5.0
gh pr create --base release-1.5
```

### 3. Merge the PR and watch CI

Once the PR is merged into the release branch, CI automatically:

1. Creates a git tag and a **GitHub pre-release** for the new version.
2. Builds and publishes OLM images for the `candidate` channel.
3. Opens a PR against the release branch with the updated `olm/index-template.yaml`.
   Merge this PR to record the catalog change in the repository.

### 4. Promote to stable

Once testing is complete, uncheck `Set as a pre-release` on the GitHub release
page to mark it as production-ready. This triggers CI to:

1. Build and publish OLM images for the `stable` channel.
2. Open a PR against the release branch with the updated catalog. Merge it.

### How to force a release version

```sh
RELEASE_VERSION=1.5.0
make initiate-release-as RELEASE_VERSION=$RELEASE_VERSION
```

## How to publish a new release to the Community Catalog

After a new stable release has been published, update the operator version in
the [OpenShift community catalog](https://github.com/redhat-openshift-ecosystem/community-operators-prod).

Assumptions:

* You have already forked and cloned `https://github.com/redhat-openshift-ecosystem/community-operators-prod`.
* The `origin` remote refers to the upstream repository and the `fork` remote to the forked repository.

1. Check out the release branch locally (the stable catalog PR should already
   be merged, so the `bundle/` directory reflects the release):

```sh
VERSION=1.5.0
git checkout release-1.5
git pull
```

2. Copy the `bundle/` directory to your community-catalog fork:

```sh
cd ../../redhat-openshift-ecosystem/community-operators-prod
git checkout main
git fetch && git reset --hard origin/main

git checkout -b observability-operator-$VERSION
mkdir -p operators/observability-operator/$VERSION
cp -r ../../rhobs/observability-operator/bundle operators/observability-operator/$VERSION
```

3. Validate the bundle (this should already have been done in CI):

```sh
operator-sdk bundle validate operators/observability-operator/$VERSION \
	--select-optional name=operatorhub \
	--optional-values=k8s-version=1.21 \
	--select-optional suite=operatorframework
```

4. Commit (signed) and push for review:

NOTE: The commit message follows a convention (see `git log`) and must be signed.

```sh
git add operators/observability-operator/$VERSION
git commit -sS -m "operator observability-operator ($VERSION)"
git push -u fork HEAD
```

5. Submit the pull request, e.g: https://github.com/redhat-openshift-ecosystem/community-operators-prod/pull/3084

6. There may be some changes required to fix the bundle. Make those changes and
   ensure the fixes are ported back to the Observability Operator repo.
   E.g.: https://github.com/rhobs/observability-operator/pull/333

## How to update the forked prometheus-operator

The observability operator uses a forked (downstream) version of the upstream
Prometheus operator to ensure that it can be installed alongside the upstream
operator without conflict. The forked operator is maintained at
(https://github.com/rhobs/obo-prometheus-operator/) which contains the
instructions to synchronize from upstream.

When a new downstream version is available (e.g. `v0.69.0-rhobs1`), you need to
update these 2 files and replace the old version by the new one:

* `go.mod`
* `deploy/dependencies/kustomization.yaml`

Then regenerate all the manifests:

```sh
make generate
```

Finally submit a pull request with all the changes.

Example: (https://github.com/rhobs/observability-operator/pull/380)
