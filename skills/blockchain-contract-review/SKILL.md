---
name: blockchain-contract-review
description: "Use for authorized source-level security review of Solidity/EVM smart contracts, protocol integrations, and deployment controls. Validate privilege, asset-flow, accounting and upgrade findings with reproducible evidence."
---

# Smart contract and DeFi security review

This skill focuses on auditing deployed or supplied contracts in an explicitly authorized project. It reviews code, deployment metadata and controlled test transactions; it is **not** an instruction to trade, extract funds, or attack unrelated on-chain systems.

Research input: the `blockchain-contract-attack` topic entry in the user-provided `mojin-2.0.0-Windows-x64.zip`. This workflow is newly written for ARTEX and imports no external scripts, binary executables, or runtime permissions.

## 1. Identify what is actually deployed

Record chain/network, chain ID, contract address, implementation and proxy address (where applicable), verified source and compiler settings, deployment block, dependencies and privileged roles. Contract names and explorer labels alone are insufficient for establishing identity. Compare the checked-out source to the actual deployed bytecode or verified build metadata whenever possible.

Define the trust boundaries: users, owners, guardians, timelocks, upgrade administrators, external protocols, oracle feeds, cross-chain messaging, token contracts, and keepers. Treat any off-chain operator or oracle response as external data.

## 2. Review matrix

| Surface | Review questions |
| --- | --- |
| Authorization | Who may mint, burn, move funds, pause, change roles, update dependencies or upgrade implementation? Can role changes bypass timelock/governance assumptions? |
| Asset flow | Do deposits, withdrawals, fees and refunds preserve balances and token units? Are returned values, fee-on-transfer tokens and nonstandard ERC-20 behaviors handled? |
| Reentrancy and callbacks | Are state transitions and invariants safe across external calls and callbacks, including token hooks? Check actual guards and call order rather than pattern-matching alone. |
| Arithmetic and accounting | Are decimal conversions, rounding directions, overflow conditions, underflow assumptions and share-price calculations consistent across functions? |
| Oracle and price inputs | Are freshness, fallback, decimal normalization, manipulation resistance and circuit breakers appropriate to the application? |
| Signatures and approvals | Are domain separation, chain ID, nonces, expiry, signer identity and replay handling correct? |
| Upgrades and delegatecall | Are storage layouts stable and implementation initializer/upgrade authorization correct across versions? |
| Cross-chain flows | Are message origin, chain ID, replay protection, ordering, finality and bridge signer validation verifiable? |
| Availability and griefing | Can unbounded loops, paused dependencies, revert propagation or fee exhaustion lock legitimate users? |
| Deployment and operations | Do admin keys, ownership, keeper assumptions, emergency procedures and monitoring match the intended production settings? |

Prioritize questions by the actual protocol architecture, not by a generic vulnerability-count target. Do not equate any external call or `delegatecall` with a confirmed vulnerability.

## 3. Verify locally before claims

1. State a precise invariant or authorization expectation (e.g., only authorized roles can alter critical parameters, or total shares remain consistent with underlying accounting).
2. Record the relevant contract version, source files/functions and configuration. Distinguish static analysis hypotheses from verified exploitability.
3. Where the development environment supports it, run existing compiler/tests, static checks (such as Slither), or a locally forked/fake-chain test. Never run an unfamiliar script or issue a real-value on-chain transaction automatically.
4. Compare baseline behavior with the suspected failing condition. Keep transaction traces, logs, failed assertion output and the exact test inputs.
5. Link the result to ARTEX's existing finding/asset/evidence records. If runtime evidence is unavailable, retain a pending code-review finding with explicit uncertainty.

## 4. Report

For each issue, provide chain and contract identifier, deployed version, relevant authorization assumptions, affected function/source location, expected invariant, observed deviation, actual evidence, conditions for impact, affected assets, severity rationale, and a proposed fix/retest. **Do not label loss of funds or cross-chain compromise as observed unless it was independently demonstrated.**

Use the existing `codeaudit`, `finding-triage`, `patch-diff` and reporting paths for subsequent investigation. This Skill is guidance only: it does not register a scanner, install tools, expand scanning scope, or change approval behavior.
