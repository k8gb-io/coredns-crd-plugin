# Dependency management

Prefer actively maintained dependencies with compatible licenses and a clear
security reporting process. Review new dependencies for necessity, provenance,
maintenance status, license obligations, and added transitive dependencies.
Avoid adding a dependency where existing project code or the standard library
already provides the required behavior.

## Acquisition and tracking

Go modules are declared in [go.mod](go.mod) and [terratest/test/go.mod](terratest/test/go.mod).
Their `go.sum` files record module content checksums. `go mod download` obtains
modules through the configured Go proxy; the default public configuration uses
`proxy.golang.org` and the Go checksum database at `sum.golang.org`. Do not
disable checksum verification to bypass a mismatch. `go mod verify` checks
downloaded module contents against the module cache's recorded hashes.
Use `go list -m all` in each module to inspect its resolved dependency graph.

GitHub Actions and tools are declared in [.github/workflows](.github/workflows),
container bases in [Dockerfile](Dockerfile), and chart dependencies in
[charts/coredns/Chart.yaml](charts/coredns/Chart.yaml). `helm dependency update
charts/coredns` retrieves the configured chart dependencies. Docker retrieves the
base image while building. These are separate from Go module checksums and must
be reviewed independently.

[Renovate](renovate.json5) proposes dependency updates as pull requests, including
Go modules and supported workflow, container, and chart dependencies. Actions
are pinned to commit SHAs. The baseline publisher is the documented exception:
grc.store requires its `v1` workflow identity, so Renovate ignores that dependency.
Maintainers must review upstream changes to that mutable tag.

## Reviewing updates

Review release notes, security advisories, compatibility, transitive dependency
changes, and checksum changes. Update module manifests and checksums together;
run unit tests, lint, and the integration suite when dependency behavior affects
DNS or Kubernetes compatibility. Renovate updates follow the same review and
required-check policy as other changes. Security updates take priority through
the process in [SECURITY.md](SECURITY.md).
