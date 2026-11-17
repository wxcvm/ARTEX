---
name: ai-llm-app-security
description: "Use for authorized security review of LLM applications, RAG, tool-using agents, MCP integrations, and model artifact loading. Trace trust boundaries and validate findings with observable effects."
---

# AI / Agent Application Security

Use this skill when the **application under test** uses a language model, retrieval system, agent tools, or MCP servers. This is an application-security workflow, **not** a rule for changing ARTEX's own model policies. Reuse the current task, asset, traffic, evidence and findings features.

Research input: the AI/LLM threat-surface categories in the user-provided `mojin-2.0.0-Windows-x64` package (`skills/ai-llm-app-attack/SKILL.md`). The checklist below is independently written for ARTEX; no executable code or default configuration from that package is imported.

## Map the application before testing

Record the visible data flow, with a source and trust level for each hop:

```text
user / external page / document / email / tool result
             ↓
prompt assembly / retrieval / memory
             ↓
model response / tool selection
             ↓
MCP or local tools / remote APIs / output renderer
```

For each boundary, note **who controls the content**, which identity and permissions are in use, what operations it can reach, and where the data goes. Unknown connections remain *unknown*; never invent private system prompts, hidden tools, or backend behavior.

## Inspection matrix

| Surface | Question to verify | Evidence to retain |
| --- | --- | --- |
| Direct and indirect prompt injection | Can lower-trust input change the intended workflow or trigger an unintended tool call? | Source text, model/tool event, outcome |
| Tool privileges | Can a low-privilege request cause reads, writes, commands or outbound requests outside its intended authority? | Tool name, arguments, authorization context, observed result |
| MCP integration | Are server identities, resource provenance, tool descriptions, permissions, and updates trusted appropriately? | Connection config without secrets, installed version, tool invocation record |
| RAG and memory | Can an untrusted retrieved item be treated as an instruction, persist across turns, or contaminate another user's answer? | Document ID, retrieval trace, session boundaries and response |
| Tenant and session separation | Does information or authorization cross accounts, organizations, or conversations? | Two test identities, controlled test record, request/response |
| Output handling | Does model-produced text flow into HTML, templates, SQL, executable code or command interpreters without suitable handling? | Sink, input, observed application behavior |
| Resource usage | Are retries, loops and retrieval breadth bounded to the task's stated cost/impact budget? | Token/tool counts, stop reason, recorded cost |
| Model artifacts | Are model files and loaders obtained from trusted sources and are unsafe deserializers used? | Manifest, loader configuration, dependency version |

**Do not claim exploitation from a plausible prompt, static pattern, or model's description of what it supposedly did.** A confirmed finding needs the relevant request/response, tool event, state change or repeatable application behavior. A model merely printing a credential-like string is not proof it accessed a real secret.

## Minimal verification sequence

1. Select one boundary and form a falsifiable hypothesis, such as “retrieved document instructions may override the user's task.”
2. Use a controlled document, test account, mock secret, and isolated test tool where possible. Preserve the exact benign baseline.
3. Compare the same workflow with and without the untrusted input. Keep request identifiers, document identifiers, timestamps, model/provider where observed, and actual tool events.
4. Confirm any side effect at the destination rather than relying on assistant narration. Mark absent tool evidence as **not verified**.
5. Record the result through ARTEX's existing evidence/finding workflow, including a negative test if it rules out the hypothesis. Keep test data and real customer data separate.
6. Recommend the smallest appropriate control at the demonstrated boundary: authorization at the tool, retrieval provenance, output encoding, constrained deserialization, or resource budgeting.

## Result contract

For each candidate, return: **boundary → hypothesis → test context → observed evidence → outcome → affected asset/version → mitigation → confidence**. Use `confirmed` only when the claimed behavior was actually observed. Otherwise use `pending`, `not reproduced`, or `not tested` as appropriate. Do not copy a generic LLM risk list into confirmed findings.

## ARTEX integration notes

- Reuse existing traffic capture, task evidence, asset references, findings and reporting. Do **not** create a second database or scan engine.
- Treat MCP and model-provider metadata as observations, not proof of reachability or approval.
- This skill only supplies reasoning guidance. It does not install tools, change permissions, start scans, or enable an external service.
