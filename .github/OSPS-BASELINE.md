# OSPS Baseline Level 2 rollout

This repository is a codebase of k8gb and is in scope for
[k8gb#2563](https://github.com/k8gb-io/k8gb/issues/2563).
[PR #175](https://github.com/k8gb-io/coredns-crd-plugin/pull/175) supplies the
scheduled scanner, inherited project Security Insights metadata, and the
publisher's Renovate exception. Repository documentation and enforcement
complement that workflow; publishing a scan alone is not baseline compliance.

## Deployment order

The 2026-11-04 audit enabled and verified secret scanning, push protection,
private vulnerability reporting, read-only default Actions tokens, and disabled
Actions approval of pull requests. Organization MFA and an HTTPS organization
homepage were already configured. The official v0.33.2 scanner against GitHub
`main`, using the local authenticated account and catalog 2026-08-28 at Level 2,
changed from 17 passed / 13 warnings / 11 failed to 21 passed / 13 warnings /
7 failed after those settings changes. This is an evaluation of the published
repository, not of unmerged changes. The overall result still fails; no badge
completion is claimed. Its token permissions differ from the workflow job token.

1. Merge the scan workflow and repository guidance, including CI that reports
   on every pull request. Confirm `lint`, `terratest`, `Analyze (go)`, `CodeQL`,
   and `DCO` report successfully. In particular, verify the CodeQL result check
   is emitted before requiring it.
2. **Done (2026-11-04):** private vulnerability reporting, secret scanning and
   push protection are enabled; Actions defaults to read-only and cannot approve
   pull requests. Keep release write permissions scoped to the release job.
3. The staged `main protection` ruleset (ID `26573234`) is active with one
   approval, resolved conversations, DCO, deletion protection, force-push
   protection and no bypass actors. Classic protection was removed after
   equivalent ruleset enforcement was verified. After the workflow changes land,
   update it with the full CI requirements in [main-ruleset.json](main-ruleset.json):

   ```sh
   gh api repos/k8gb-io/coredns-crd-plugin/rulesets/26573234 \
     --method PUT --input .github/main-ruleset.json
   ```

   Confirm the ruleset ID and current settings before updating it. Require all
   five check contexts only once their workflows are on `main`.
4. A grc.store namespace administrator must add this repository as a trusted
   publisher in `k8gb-io`, restricted to `refs/heads/main`. This cannot be
   configured by a repository file or GitHub's repository settings.
5. Dispatch `gh workflow run publish-to-grc-store.yml -R k8gb-io/coredns-crd-plugin`.
   Use `-f dry-run=true` first to test without publication if needed. Download
   `pvtr-results` and inspect the evaluation, catalog version, evaluator version,
   failed controls, and controls needing human review. A green Actions run means
   publication succeeded, even if its evaluation says `Failed`.

The workflow's `osps-baseline` alias follows the scanner's latest bundled catalog.
The original Slam work plan specifies catalog 2026-08-28; record the actual
catalog in completion evidence and verify the event's accepted version if the
alias advances. The workflow needs no privileged personal access token.

## Release tags

[release-tags-ruleset.json](release-tags-ruleset.json) restricts creation,
updates and deletion of `v*` tags to repository administrators. Its bypass actor
is repository role `5` (admin). Import it as a separate tag ruleset:

```sh
gh api repos/k8gb-io/coredns-crd-plugin/rulesets \
  --method POST --input .github/release-tags-ruleset.json
```

Inspect existing rulesets first and use PUT with the existing ID when updating
this rule. Admin-controlled tags remain a trusted release input. The
[release workflow](workflows/release.yaml) additionally checks that the tagged
commit is an ancestor of `main` in a read-only job before the job with release
and package write permissions starts. Both jobs check out the event's immutable
SHA. GoReleaser is pinned to v2.20.0 and its configuration is checked by PR CI.

## Review evidence and remaining work

| Area | Evidence or action |
| --- | --- |
| Cleaner | [SECURITY-INSIGHTS.yml](../SECURITY-INSIGHTS.yml), schema v2.2.0, with project inheritance and the repository's full core team |
| CRA Readiness | [SECURITY.md](../SECURITY.md) and [steward escalation](../CONTRIBUTING.md#escalating-to-the-cra-steward) |
| Chronicler | [CONTRIBUTING.md](../CONTRIBUTING.md), [DEPENDENCY.md](../DEPENDENCY.md), build and test commands |
| Mechanizer | Scheduled [publisher](workflows/publish-to-grc-store.yml), a successful published evaluation, and configured trusted publisher |
| Inspector | [Project self-assessment](https://github.com/k8gb-io/k8gb/blob/master/self-assessment.md), including embedded CoreDNS; human review must confirm continued coverage |
| Defender | Effective repository protections, published evaluation, and [manual-review evidence](../SECURITY.md#baseline-review-evidence) |
| Release authenticity | The existing release workflow publishes checksums but does not sign plugin artifacts or generate provenance. Add and verify those before claiming signed releases; inherited k8gb verification instructions apply to k8gb artifacts |
| Permission inventory | Verify actual repository, organization-owner, team, and automation access. The inherited project administrator list is not the effective repository ACL |
| Reviewable test inputs (QA-05.02) | Unit tests generate GeoIP data in a temporary directory from Go source; `make geodata` generates the ignored integration fixture before deployment |
| Organization MFA | Organization setting; verify with organization administrators when the workflow token cannot observe it |
| Badge submission | Maintainers submit completion evidence after implementation and review; neither a green scan workflow nor these documents award a badge |

Scanner versions affected by
[ossf/pvtr-github-repo-scanner#528](https://github.com/ossf/pvtr-github-repo-scanner/issues/528)
can falsely fail QA-03.01 when rulesets cover checks absent from classic
protection. The upstream fix is in
[#529](https://github.com/ossf/pvtr-github-repo-scanner/pull/529).
Record the evaluator version and inspect effective rules when interpreting that
finding; do not weaken protections to accommodate a scanner result.
