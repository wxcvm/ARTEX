---
name: mobile-binary-review
description: "Review authorized Android/iOS hybrid-app packages and desktop binaries beyond initial endpoint extraction: manifests, runtime boundaries, native components, client storage, and evidence-based validation. Use with apk-recon and rea."
---

# Mobile / desktop binary security review

This is the **follow-on investigation** after `apk-recon` finds package metadata or endpoints, or `rea` opens a compiled file. Use only provided or otherwise in-scope artifacts. The goal is to identify testable application-security issues, not to assume every embedded string is a credential or a vulnerability.

Research input: the **topic coverage** of `binary-mobile-reversing` in the user-supplied `mojin-2.0.0-Windows-x64.zip`. This skill is independently written; no source scripts, hardcoded targets, extraction offsets or exploit chains were copied.

## 1. Identify the artifact, not just the filename

Capture SHA-256, package ID/product identity, version/build number, platform/architecture, signing information where available, provenance, and which part of the artifact was actually inspected.

| Platform | High-signal components | Important distinctions |
| --- | --- | --- |
| Android APK/AAB/XAPK | Manifest, DEX, native `.so`, assets, packaged resources | A string from a library can belong to an SDK, a test endpoint or a real runtime configuration |
| UniApp / DCloud | `assets/apps/`, `app-service.js`, native plugins, app manifest | Hybrid JavaScript may carry application logic that JADX will not reveal |
| Flutter | `libapp.so`, app assets, native bridges and platform channels | Obfuscated Dart AOT changes what static tooling can prove; URLs in strings are leads only |
| iOS IPA | plist entitlements, frameworks, embedded web content, app transport settings | Missing keys or unusable binaries must be reported, not guessed |
| Windows PE / Electron | PE imports/resources, executable manifest, ASAR/JS bundle, native modules | Distinguish bundled dependencies from code actually loaded in production |
| Firmware or other binaries | ELF/PE sections, filesystem layout, executable entrypoints, embedded configs | Static recovery alone is not runtime exploitability |

Prefer ARTEX's existing `rea` skill for read-only static inspection. Use `apk-recon` for initial infrastructure seeds, and `codeaudit` for source made available by the project. Confirm tool availability before proposing dynamic instrumentation.

## 2. Map trust boundaries

- **Client-side authority:** exported activities/services/providers/receivers, URL/deep link handlers, intent routing, JS↔native bridges, WebView file access, and runtime feature flags.
- **Network:** API base URL resolution, environment switches, TLS validation and certificate pinning behavior, proxy use, and whether credentials are sent outside the expected origin.
- **Local data:** preferences, databases, cache, backups, logging, key stores, clipboard and inter-app sharing.
- **Native code:** JNI/FFI/native plugin boundaries, untrusted parsers, dynamic libraries and their update paths.
- **Release:** signing, debug flags, symbols, dependency versions, packaging of development-only configuration.

For each suspected issue, note **attacker-controlled input → boundary → sensitive operation → observed result**. A manifest declaration or an embedded token-like pattern alone is not a confirmed exploit.

## 3. Focused verification

1. Record a concrete hypothesis and which app version or binary implements it.
2. Inspect the smallest relevant file, function or runtime event. For dynamically selected endpoints, a static URL list is not proof of the endpoint actually used.
3. On an isolated, authorized test device or emulator, reproduce with test accounts and non-sensitive fixtures where applicable. Keep raw traffic, logs and minimal reproduction inputs.
4. Compare expected behavior with the observation. For access-control defects, verify authorization on the **server side**, not merely the visibility of client UI elements.
5. Distinguish `observed`, `inferred` and `not tested` outcomes. Secrets and tokens in evidence should be redacted before reporting.

Do not run untrusted bundled code or unexplained “deobfuscation helpers” on the user's host without review. If tooling or a device is unavailable, complete the static portion and state the boundary explicitly rather than fabricating a dynamic test.

## 4. ARTEX handoff

Record the artifact fingerprint and inspected components, source file/offset or `rea` evidence ID, affected version, reproducibility, relevant asset/task references, potential impact, and minimal remediation. Use existing `finding-triage` and evidence/report exports for follow-up. Avoid a parallel package database or new execution service.

This skill provides guidance only. It does not install reverse-engineering tools, enable debugging on a real device, start scans, or change tool permissions.
