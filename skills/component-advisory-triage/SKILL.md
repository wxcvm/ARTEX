---
name: component-advisory-triage
description: "Use when an identified software component and version require evidence-based CVE/advisory triage, version-range matching, patch validation, and actionable handoff."
---

# Component advisory triage

When an authorized assessment identifies a product, framework, library, container image, or package version, use this workflow to connect **inventory evidence → authoritative advisories → affected version → observed exposure → verification**. Search results alone never establish a vulnerability.

Research input: the topic coverage of `skills/component-vuln-intel/SKILL.md` in the user-provided `mojin-2.0.0-Windows-x64` package. This is a new ARTEX-oriented workflow, not a direct copy of that file's exhaustive search or proxy instructions.

## 1. Establish an exact component identity

Record name, ecosystem, vendor, package identifier (PURL/CPE if grounded), observed version, version source, deployed context, and confidence in the fingerprint. Distinguish a banner, a declared dependency, a loaded runtime version, and a version inferred from unrelated infrastructure. If the version is unknown, leave it unknown rather than inventing one.

## 2. Search authoritative sources first

Choose sources that fit the component; **do not mechanically query every site**.

- Vendor security advisory and fixed-version notes.
- GitHub Security Advisories, OSV, NVD/CVE records, ecosystem security trackers (for example npm, PyPI, Go, Maven, container or Linux distribution trackers).
- Upstream issue, commit or patch when it resolves a disputed affected-version claim.
- Community reports and PoCs only as secondary leads; require independent validation of their claims and never execute unreviewed code automatically.

Search by exact package/vendor and version first, then CVE/GHSA identifier and affected version range. Record advisory URL, retrieval time and its stated publication/update date. “No result from available sources” means **no matching advisory found in those sources**, not “the component is safe.”

## 3. Decide applicability

For every candidate, compare:

1. Affected-version range versus the **deployed** version, not just a package filename.
2. Relevant module, optional feature, runtime/configuration, platform and authentication preconditions.
3. Required exposure and whether the affected code path is reachable in this environment.
4. Whether a vendor patch, mitigation, or backport changes applicability.

An advisory that merely mentions the product is a **lead**, not a confirmed finding. Different vendors and distributions may use different patch versions.

## 4. Verify with existing ARTEX features

- Link the advisory to the corresponding task and asset; use existing knowledge lookup, traffic/evidence records, finding triage and retest features.
- Prefer controlled, low-impact reproduction, a non-destructive version/configuration check, or vendor-supplied verification steps.
- Separate **version-affected**, **exposure observed**, and **impact reproduced**. State which claim each artifact actually supports.
- Save negative evidence and unresolved cases without converting them into confirmed vulnerabilities.
- If no network or advisory source is reachable, state the gap and preserve the existing task; do not force proxy installation or an automatic scan.

## Handoff record

| Field | Required meaning |
| --- | --- |
| Component | Product/package, version and fingerprint source |
| Advisory | CVE/GHSA/vendor ID, source and checked date |
| Applicability | Affected range and relevant preconditions |
| Evidence | Exact artifact IDs or requests/responses supporting the claim |
| Decision | Confirmed / pending validation / not applicable / cannot assess |
| Remediation | Fixed version, vendor mitigation, or verification still required |

Use ARTEX's existing finding, task, evidence and knowledge storage. Do not introduce a new database, mandatory remote query service, or a new authorization gate. This skill contains guidance only; it installs and executes nothing.
