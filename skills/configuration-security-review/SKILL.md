---
name: configuration-security-review
description: "Review application API specifications, infrastructure-as-code, Kubernetes manifests, container images, and dependency configuration for security misconfigurations with optional installed scanners and independently verified findings."
---

# Configuration and software supply-chain review

Choose this skill for **source- or artifact-based** configuration review: OpenAPI/AsyncAPI schema files, Terraform, CloudFormation, Kubernetes manifests, container definitions and lockfiles. It complements `api-security`, `codeaudit`, `cloud-security` and `component-advisory-triage`; do not treat scanner output as automatically confirmed vulnerabilities.

Research input: the `api-schema-analyzer`, `trivy` and `checkov` YAML **tool descriptions** in the user-provided `mojin-2.0.0-Windows-x64.zip`. This is newly written guidance, not a direct import of their free-form command templates, and introduces no new executable tool.

## Intake and routing

1. Record source revision, artifact digest, build environment, target platform, and why the configuration is expected to be deployed.
2. Match the artifact to the review method below; avoid running unrelated checks simply because a scanner exists.
3. Use an already installed scanner if its version and invocation are actually verified. Do not infer that a tool exists merely because a YAML recipe in another product marks it `enabled`.

| Artifact | Relevant questions | Optional tool examples |
| --- | --- | --- |
| OpenAPI/AsyncAPI, GraphQL schemas | Undocumented auth requirements, inconsistent request validation, unexpected public endpoints, unsafe examples | Spectral lint, API-spec diff, manual route review |
| Terraform / CloudFormation | Public resources, identity policy breadth, transport encryption, secrets in plans, logging and network ingress | Checkov, Trivy configuration scan |
| Kubernetes YAML / Helm rendered manifests | Workload identity, RBAC, pod security, host mounts, secret handling, exposed services | Checkov, Trivy configuration scan |
| Dockerfile and container filesystem/image | Base-image version, package vulnerabilities, root permissions, copied secrets, writable paths, exposed interfaces | Trivy image/filesystem scan |
| Lockfiles and SBOMs | Exact dependency versions, affected ranges, ownership and patch status | Trivy filesystem scan, ecosystem advisory review |

The tool examples are **candidates**, not a promise of support. ARTEX's currently managed external-tool catalog covers a different set of 7 reconnaissance binaries; these new tools must not be advertised as installed, auto-installable or natively managed without compatible implementation and tests.

## Evidence-first review

- Compare configuration with the actual deployment or documented build path. A permissive example file does not prove production exposure.
- Link specific manifest field/file:line, image digest, dependency PURL or API operation ID to the hypothesis.
- Verify scope-specific controls (platform policy, deployment overlays, ingress, cloud IAM, exceptions) before confirming.
- Separate a **policy warning**, a **version advisory match** and an **observed security defect**.
- In multi-stage configurations, check effective values after rendering and overrides rather than just source defaults.
- Record negative findings and scan gaps: a clean scanner result is not a comprehensive security guarantee.

If offline, no compatible scanner is installed or dependencies cannot be resolved, provide the manual static review and explicitly state unavailable checks. Do not silently download tool binaries, modify PATH, run remote scanning APIs or add a new daemon.

## Deliverables

For each candidate: **artifact/version → affected setting → intended vs observed behavior → source evidence → deployment applicability → severity rationale → remediation/retest**.

Use existing ARTEX task evidence, reports and `finding-triage` for findings; use `component-advisory-triage` for CVE/version applicability and `api-security` for runtime API authorization checks. This is a Skill-only addition and does not change execution, authorization or installation behavior.
