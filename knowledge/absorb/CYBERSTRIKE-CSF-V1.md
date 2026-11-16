# CyberStrike CSF reference pack (v1)

This first curated pack contains nine original ARTEX defensive reference documents. Source topics were reviewed in [CyberStrike](https://github.com/CyberStrikeus/CyberStrike/tree/9303b210ab82bd6362d7ad7fbd2b2c8f2e335450) and mapped to the [NIST CSF 2.0](https://doi.org/10.6028/NIST.CSWP.29) framework. The upstream repository LICENSE is AGPL-3.0. The references are independent ARTEX summaries and operational review guidance, not verbatim imported skills or executable material. In particular, a CSF outcome does not establish a vulnerability, severity, CWE, or independent verification.

## Source index

| Control | Theme | Upstream source |
| --- | --- | --- |
| ID.AM-01 | 资产台账与硬件清单 / Asset inventory | [SKILL.md](https://github.com/CyberStrikeus/CyberStrike/blob/9303b210ab82bd6362d7ad7fbd2b2c8f2e335450/.cyberstrike/skill/NIST/CSF_v2.0/ID_identify/ID.AM-01_idam-01/SKILL.md) |
| ID.RA-01 | 漏洞验证与风险台账 / Vulnerability risk assessment | [SKILL.md](https://github.com/CyberStrikeus/CyberStrike/blob/9303b210ab82bd6362d7ad7fbd2b2c8f2e335450/.cyberstrike/skill/NIST/CSF_v2.0/ID_identify/ID.RA-01_idra-01/SKILL.md) |
| PR.AA-05 | 权限治理与最小权限 / Least privilege access review | [SKILL.md](https://github.com/CyberStrikeus/CyberStrike/blob/9303b210ab82bd6362d7ad7fbd2b2c8f2e335450/.cyberstrike/skill/NIST/CSF_v2.0/PR_protect/PR.AA-05_praa-05/SKILL.md) |
| PR.DS-11 | 备份保护与恢复测试 / Backup integrity and restore testing | [SKILL.md](https://github.com/CyberStrikeus/CyberStrike/blob/9303b210ab82bd6362d7ad7fbd2b2c8f2e335450/.cyberstrike/skill/NIST/CSF_v2.0/PR_protect/PR.DS-11_prds-11/SKILL.md) |
| PR.PS-02 | 软件生命周期与补丁管理 / Software patch lifecycle | [SKILL.md](https://github.com/CyberStrikeus/CyberStrike/blob/9303b210ab82bd6362d7ad7fbd2b2c8f2e335450/.cyberstrike/skill/NIST/CSF_v2.0/PR_protect/PR.PS-02_prps-02/SKILL.md) |
| DE.CM-01 | 网络监测与日志覆盖 / Network continuous monitoring | [SKILL.md](https://github.com/CyberStrikeus/CyberStrike/blob/9303b210ab82bd6362d7ad7fbd2b2c8f2e335450/.cyberstrike/skill/NIST/CSF_v2.0/DE_detect/DE.CM-01_decm-01/SKILL.md) |
| RS.AN | 事件分析与证据保全 / Incident analysis | [SKILL.md](https://github.com/CyberStrikeus/CyberStrike/blob/9303b210ab82bd6362d7ad7fbd2b2c8f2e335450/.cyberstrike/skill/NIST/CSF_v2.0/RS_respond/incident-analysis-rs-an-incident/SKILL.md) |
| RC.RP | 事件恢复执行与验收 / Incident recovery plan execution | [SKILL.md](https://github.com/CyberStrikeus/CyberStrike/blob/9303b210ab82bd6362d7ad7fbd2b2c8f2e335450/.cyberstrike/skill/NIST/CSF_v2.0/RC_recover/incident-recovery-plan-execution-rc-rp/SKILL.md) |
| GV.SC | 供应链风险与依赖来源 / Supply chain risk management | [SKILL.md](https://github.com/CyberStrikeus/CyberStrike/blob/9303b210ab82bd6362d7ad7fbd2b2c8f2e335450/.cyberstrike/skill/NIST/CSF_v2.0/GV_govern/cybersecurity-supply-chain-risk/SKILL.md) |

The machine-readable source of truth is `db/knowledgepacks/cyberstrike-csf-v1.json`; each entry records its upstream path and blob SHA. Entries are grouped under `defensive-csf` and contain Chinese review guidance plus English search keywords, expected evidence and ARTEX integration notes.

## Runtime

The pack is embedded in the backend, and startup seeds it once into the existing global knowledge document library using stable UUIDs and a transaction with a settings marker. No new database table, schema migration, dependency, model call, or change to agent tool permissions is required. The existing global document UI and scoped `kb_search` / `kb_get` can read it. Keyword search is immediately available; semantic indexing follows the existing configuration and indexing lifecycle. Existing user-edited or deleted documents are not reinstalled on later startup. This does not modify the bundled SQLite knowledge.db or its numeric IDs.

## Verification

The manifest test checks all nine unique IDs, source metadata and searchable text. The PostgreSQL integration test seeds a unique fixture pack, verifies global document search and provenance, edits/deletes records and verifies a second seed does not overwrite or restore them. Existing user documents are not used as fixtures.

## Scope

Nine reviewed reference documents are the first delivery, not a claim of complete CyberStrike knowledge absorption. No automatic CWE/severity mapping, active testing instructions, or wholesale upstream copy is included. Expand the catalog after source and duplication review.
