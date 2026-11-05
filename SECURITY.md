# Security policy

coredns-crd-plugin is part of the k8gb project. The
[k8gb security policy](https://github.com/k8gb-io/k8gb/blob/master/SECURITY.md)
and coordinated disclosure process cover this repository and its `k8s_crd` image.

## Reporting a vulnerability

Email [cncf-k8gb-maintainers@lists.cncf.io](mailto:cncf-k8gb-maintainers@lists.cncf.io)
or use this repository's [private vulnerability reporting form](https://github.com/k8gb-io/coredns-crd-plugin/security/advisories/new).
Do not disclose suspected vulnerabilities in public issues or pull requests.
Include the plugin version or image digest, k8gb and Kubernetes versions where
applicable, the affected configuration, impact, and steps to reproduce. Remove
credentials and personal data. State prominently if exploitation is ongoing.

The project policy specifies a response in less than 14 days. Maintainers
coordinate confirmation, remediation, disclosure, and a patch release with the
reporter. Include your GitHub username if you want to collaborate on a private
advisory. Public disclosures are published in
[k8gb security advisories](https://github.com/k8gb-io/k8gb/security/advisories);
advisories should identify affected plugin versions and fixed image tags.
Reports filed in this repository can also be coordinated and disclosed through
its [security advisories](https://github.com/k8gb-io/coredns-crd-plugin/security/advisories).

## Security updates and support

Use the latest compatible plugin release and monitor the
[releases](https://github.com/k8gb-io/coredns-crd-plugin/releases) and project
advisories for security updates. Older versions have no guaranteed backport or
support period. Report suspected vulnerabilities even if your version is old;
maintainers assess the affected versions during triage. Community support and
defect reporting are described in [CONTRIBUTING.md](CONTRIBUTING.md).

## CRA stewardship

This project is supported under the Linux Foundation CRA stewardship framework,
as described at https://www.linuxfoundation.org/security. Security vulnerabilities
should be reported through the mechanisms described above, which we will
coordinate with our CRA steward. For actively exploited vulnerabilities and
severe incidents that may require CRA escalation, please use the project's
emergency security reporting mechanisms as appropriate. See
[steward escalation](CONTRIBUTING.md#escalating-to-the-cra-steward).

## Baseline review evidence

The [OSPS scan](.github/workflows/publish-to-grc-store.yml) evaluates Maturity
Level 2. A successful workflow means the evaluation was published; it does not
mean every control passed. Inspect the result and any controls needing review at
https://grc.store/targets/k8gb-io/coredns-crd-plugin.

### Build and release trust (OSPS-AC-04, OSPS-BR-01.03)

Pull request workflows use GitHub hosted runners, read repository contents,
and run lint, unit tests, integration tests, and CodeQL. They do not publish
releases. CodeQL's analysis job additionally writes code scanning results.
[release.yaml](.github/workflows/release.yaml) runs on `v*` tag pushes and grants
only its release job permission to create releases and publish GHCR packages.
A separate read-only job verifies the tagged commit belongs to `main` before
the release job starts; both checkouts use the event SHA. The proposed
[release tag ruleset](.github/release-tags-ruleset.json) restricts `v*` tag
creation, updates and deletion to repository administrators.
Tag creation and changes to this workflow are privileged maintainer actions.
Review workflow triggers, token permissions, checkout references, and third-party
action changes whenever this boundary changes.

The baseline publisher is invoked by its `v1` tag because grc.store verifies that
workflow identity. This exception to SHA pinning is documented in the workflow
and Renovate configuration. The caller grants the reusable workflow OIDC access
to publish signed evaluation results; this does not sign plugin release images.

### Interfaces and assessment (OSPS-SA-01, OSPS-SA-02, OSPS-SA-03)

[ARCHITECTURE.md](ARCHITECTURE.md) maps DNS requests, Kubernetes watches,
configuration, and GeoIP data to implementation files and trust boundaries.
The [project self-assessment](https://github.com/k8gb-io/k8gb/blob/master/self-assessment.md)
covers this plugin as k8gb's embedded CoreDNS. That assessment is project-wide;
it is not a separate audit of this repository or of every plugin release.

### Vulnerability publication (OSPS-VM-04.01)

The disclosure destination and process are recorded above even when there is no
published advisory history from which the scanner can infer them.

### Access and secrets (OSPS-GV-01.01)

Project roles and infrastructure ownership follow
[k8gb governance](https://github.com/k8gb-io/k8gb/blob/master/GOVERNANCE.md).
Repository administrators can manage repository settings and Actions secrets;
write access can affect release workflows and tags. Project-wide administrator
metadata in Security Insights should not be treated as a complete inventory of
this repository's effective GitHub permissions. Maintainers must review direct,
team, organization-owner, and automation access when granting or removing roles.
On 2026-11-04, the GitHub collaborator API reported repository administrator
access for `caniszczyk`, `ytsarev`, `jkremser`, `k0da`, `donovanmuller`,
`thelinuxfoundation`, and `kuritka`; `elohmrow` and `abaguas` had maintain access.
This observation does not inventory application permissions or every external
infrastructure credential. Recheck effective permissions when roles change.

Keep credentials in GitHub Actions secrets or short-lived job tokens, never in
source, images, issues, or logs. Revoke and rotate exposed credentials, remove
them from active configuration, and report security impact through the private
channel above. Deleting a secret from Git history does not revoke it.
