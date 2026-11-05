This program generates the synthetic `geoip.mmdb` fixture used by integration
tests. Its two networks map to `site1` and `site2`; it contains no production
GeoIP data. The database is generated from this reviewable source and is not
committed to version control.

Run `make geodata` from the repository root, or `go run .` in this directory.
`make deploy-app` and `make terratest` generate it automatically before creating
the Kubernetes ConfigMap. The fixed build epoch makes the output reproducible.
Unit tests separately generate their own database in a temporary directory from
[`common/k8sctrl/geodata_test.go`](../../common/k8sctrl/geodata_test.go).
