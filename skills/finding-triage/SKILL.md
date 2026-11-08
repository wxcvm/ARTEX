---
name: finding-triage
description: "Use this skill when consolidating, reviewing, de-duplicating, or retesting security findings and preparing evidence-backed deliverables. Distinguish observed facts from hypotheses, preserve original evidence, and prepare machine-readable handoff."
---

# Finding triage and evidence handoff

Use for **authorized** assessment results that already exist. This is a review and reporting workflow, not an instruction to expand testing scope or run further attacks. It complements `codeaudit` and `api-security`: those discover and verify candidates; this skill turns their evidence into defensible findings and retest decisions.

## 1. Intake and provenance

For each candidate record the existing finding ID (if any), source task/session, asset reference, collection time, collection method, tool name and version (if known), and immutable artifact reference. Do not substitute a similar-looking asset's traffic as proof that the original test generated it.

Separate:

- **Observed:** actual response, code path, state change, or executed test outcome.
- **Inferred:** the explanation connecting those observations to a root cause.
- **Unknown:** missing environment, authentication context, reproduction, or impact evidence.

Never invent request/response bodies, screenshots, CVEs, CWE IDs, severity, reproduction counts, or a successful retest. If source materials are absent, report the gap.

## 2. Review and consolidate

1. Read the original finding and linked evidence before changing triage status.
2. Confirm that the alleged behavior is reachable in the assessed version and within the task's scope.
3. Check existing mitigations and legitimate alternative explanations.
4. Compare candidates by root cause, affected asset, entry point, authorization context, and independent proof. Similar titles alone do **not** justify merging.
5. Preserve every independent finding and every original artifact, including evidence that contradicts the hypothesis.
6. Choose a status supported by evidence: pending/in-progress, confirmed, false positive, duplicate, ignored, risk accepted, fixed, or resolved. Record the rationale; lack of evidence is **pending**, not automatically false positive.

Only attach a CVE or CWE identifier when the evidence and authoritative taxonomy support the exact mapping.

## 3. Verification and retest

For a new confirmation, use the smallest non-destructive validation that answers the outstanding question. A retest must identify its target version, environment, tested conditions, timestamp and observed result. If testing was not run, write **not tested**, not **fixed**.

Distinguish three outcomes:

- **Still reproducible:** the same root cause and impact remain observable.
- **Fixed:** the relevant test no longer reproduces *and* the patch/mitigation path was checked.
- **Cannot confirm:** environment or evidence did not support a meaningful decision.

Keep the original evidence and the new retest evidence linked; do not overwrite the first report or erase history. For uncertain outcomes, retain the prior triage decision until it can be reviewed.

## 4. Handoff and exports

A useful summary for each finding contains: title, affected asset and task, severity with rationale, status, observed behavior, root cause if established, minimal reproduction conditions, expected-versus-actual, evidence references, remediation guidance, and retest status. Redact credentials, user data and raw tokens before sharing.

Reuse ARTEX's existing findings, evidence, retest and export features rather than introducing a parallel store:

- **Markdown:** human review and remediation context.
- **JSON:** structured details and downstream integrations where available.
- **CSV:** summary and triage spreadsheets; keep untrusted text as text.
- **SARIF:** machine-readable scanner handoff. Never fabricate source file or line locations for network-only findings, and never mark pending issues as confirmed.

Check the output against the persisted finding IDs and status. An export succeeding is not evidence the underlying vulnerability is confirmed.

## Completion rule

Return a concise triage table or structured summary with explicit counts for confirmed, pending, rejected/duplicate and retest outcomes. If none were verified, say so. Never collapse **no findings**, **scan incomplete**, and **no evidence available** into the same result.
