# Contributing

This repository is part of [k8gb](https://github.com/k8gb-io/k8gb).
The project's [governance](https://github.com/k8gb-io/k8gb/blob/master/GOVERNANCE.md),
[code of conduct](https://github.com/k8gb-io/k8gb/blob/master/CODE_OF_CONDUCT.md),
and [contribution process](https://github.com/k8gb-io/k8gb/blob/master/CONTRIBUTING.md)
apply here.

## Reporting defects and getting help

Search [existing issues](https://github.com/k8gb-io/coredns-crd-plugin/issues)
before opening a new report. Include the plugin version or container digest,
Go version for build failures, Kubernetes and k8gb versions, expected and actual
DNS answers, a minimal Corefile and DNSEndpoint example, reproduction commands,
and relevant logs. Describe frequency and impact. Redact credentials, private
addresses where necessary, and personal information. For usage questions, use
[k8gb discussions](https://github.com/k8gb-io/k8gb/discussions).
Report suspected vulnerabilities privately as described in [SECURITY.md](SECURITY.md).

## Building from source

Install Git and the Go toolchain specified by the `go` directive in [go.mod](go.mod).
From a fresh checkout:

```sh
git clone https://github.com/k8gb-io/coredns-crd-plugin.git
cd coredns-crd-plugin
go mod download
go mod verify
go build -o coredns ./cmd/coredns.go
./coredns -plugins
```

The plugin list must include `k8s_crd`. `make build` builds a Linux binary for the
container image. With Docker installed, `make image TAG=dev` builds
`ghcr.io/k8gb-io/k8s_crd:dev` locally. See [README.md](README.md) for configuration.

## Tests and review

Run `go test ./...` for unit tests and `make lint` with the golangci-lint version
pinned in [.github/workflows/check.yml](.github/workflows/check.yml).
Add or update tests when changing behavior, covering failures and compatibility
as well as successful requests. Keep the tests passing when dependencies change.
The CI lint job also checks license headers with `golic`.

Integration tests require Docker, k3d, kubectl, and Helm:

```sh
make create-local-cluster
make terratest TAG=dev
make destroy-local-cluster
```

These commands create a disposable `coredns-crd` cluster, build and import the
image, install a chart and CRD, and run the tests in `terratest/test`. The deploy
target changes the current kubectl context to this cluster. Run the cleanup
command after failures too. The cluster config expects the Docker network
`k3d-action-bridge-network`; create it with `docker network create
k3d-action-bridge-network` if it does not already exist. CI provisions the network
and cluster with k3d-action. The GeoIP test database is generated from source by
`make geodata`, which the deployment target runs automatically.

Open pull requests against `main`. Explain the behavior change, tests performed,
and any operational impact. Sign off commits with `git commit -s` for the DCO
check. Changes require a non-author maintainer's review and successful required
checks. CI runs for documentation changes too so required checks always report.
Review security policy, dependency policy, architecture, and interface
documentation when a change affects them.

## Dependencies and releases

[DEPENDENCY.md](DEPENDENCY.md) describes dependency acquisition and updates.
The [release workflow](.github/workflows/release.yaml) invokes GoReleaser on
maintainer-created `v*` tags. [.goreleaser.yml](.goreleaser.yml) defines Linux
amd64/arm64 builds, checksums, release notes, and GHCR image manifests.
Release archives and checksums are published under
[GitHub Releases](https://github.com/k8gb-io/coredns-crd-plugin/releases), and
images under `ghcr.io/k8gb-io/k8s_crd`.

Download the checksum file and archive from the same release, then verify the
archive with `sha256sum --check checksums.txt` (or `shasum -a 256 --check
checksums.txt` on macOS). A checksum detects corruption; it does not establish a
publisher's signing identity. This repository's current release workflow does
not generate release signatures or provenance attestations. k8gb's verification
instructions for its own signed artifacts do not verify this plugin's artifacts.

## Escalating to the CRA steward

For severe incidents or active exploitation, use the private reporting channels
in [SECURITY.md](SECURITY.md) and identify the urgency and impact. Maintainers
coordinate escalation to the Linux Foundation CRA steward under the
[project escalation process](https://github.com/k8gb-io/k8gb/blob/master/CONTRIBUTING.md#escalating-to-the-cra-steward).
Keep technical vulnerability details private until coordinated disclosure.
