import assert from "node:assert/strict";
import test from "node:test";
import { skillLabel, skillSummary, toolLabel, toolSummary, knowledgeCategory, knowledgeKind, knowledgeCategorySummary } from "./resource-labels.ts";

test("shipped resources have Chinese display labels without rewriting their identifiers", () => {
  assert.equal(skillLabel("playwright-cli"), "Playwright 浏览器自动化");
  assert.equal(toolLabel("kb_search"), "检索知识库");
  assert.equal(knowledgeCategory("info-leak"), "信息泄露");
  assert.equal(knowledgeKind("script"), "脚本");
});

test("unknown keys and custom tools retain their original names", () => {
  assert.equal(skillLabel("my-custom-skill"), "my-custom-skill");
  assert.equal(toolLabel("kb_search", false), "kb_search");
  assert.equal(toolLabel("custom_tool"), "custom_tool");
  assert.equal(knowledgeCategory("new-category"), "new-category");
  assert.equal(knowledgeKind(""), "");
  assert.equal(skillLabel("constructor"), "constructor");
});

test("Chinese introductions are separate from source descriptions", () => {
  assert.match(skillSummary("playwright-cli"), /浏览器操作/);
  assert.match(skillSummary("finding-triage"), /原始证据/);
  assert.equal(skillSummary("my-custom-skill"), undefined);
  assert.equal(skillSummary("constructor"), undefined);
});


test("every shipped skill has a Chinese catalog introduction", () => {
  const keys = ["active-scan","apk-recon","asset-correlation","c2-supershell","chains-ad","chains-app","chains-meta","chains-pivot","chains-priv","chains-pwn","chains-recon","chains-rev","chains-web","chisel-tunnel","cloud-security","cn-proxy-pool","credential-attack","db-utilization","dir-bruteforce","fofa-recon","frp-tunnel","fscan-intranet","gogo-scan","host-enum","hunter-recon","infra-pivot","lateral-movement","memshell-fallback","nuclei-scan","passive-recon","payout-flow-analysis","quake-recon","recon-pipeline","redteam-kit-host-notes","redteam-setup","shell-handler","shodan-recon","suo5-tunnel","tunnel-selection","unauth-exploit","vps-reverse-shell","web-fingerprint","webshell-toolkit","api-recon","api-security","bbot-recon","codeaudit","edr-bypass-re","finding-triage","patch-diff","php2shell","playwright-cli","rea","scopesentry"];
  assert.equal(keys.length, 54);
  for (const key of keys) {
    assert.notEqual(skillLabel(key), key);
    assert.match(skillSummary(key), /[\u3400-\u9fff]/u, key);
  }
});

 test("tool catalog summaries are display-only and do not apply to custom tools", () => {
  assert.match(toolSummary("kb_search"), /知识库/);
  assert.match(toolSummary("get_finding_traffic"), /流量证据/);
  assert.equal(toolSummary("kb_search", false), undefined);
  assert.equal(toolSummary("custom_tool"), undefined);
  assert.equal(toolSummary("constructor"), undefined);
});

 test("every catalog tool has a Chinese name and introduction", () => {
  const keys = ["graph_overview","list_findings","list_facts","node_detail","add_intent","list_goals","prove_goal","goal_met","report_finding","record_fact","set_goals","set_constraints","add_hint","kill_work","steer_work","get_worker_output","get_worker_trace","search_all_worker_traces","list_worker_traces","expand_digest","insert_assets","add_company_scope","add_task_scope","list_untested_assets","list_assets","list_companies","save_task_credential","list_task_credentials","delete_assets_by_host","create_skill","update_skill_file","create_custom_tool","update_custom_tool","mcp_status","create_mcp","update_mcp","kb_search","kb_get","list_tasks","spawn_task","pause_task","add_task_hint","get_task_graph","list_task_findings","get_task_worker_trace","list_task_worker_traces","search_task_worker_traces","get_task_node_detail","update_finding_report","get_finding_traffic","traffic_search","traffic_get","read_task_purge_intel","asset_neighbors"];
  assert.equal(keys.length, 54);
  for (const key of keys) {
    assert.notEqual(toolLabel(key), key);
    assert.match(toolSummary(key), /[\u3400-\u9fff]/u, key);
    assert.equal(toolSummary(key, false), undefined);
  }
});

test("knowledge category guides cover the shipped taxonomy and preserve unknown categories", () => {
  const keys = ["info-leak","auth-bypass","unauthorized","rce","recon","methodology","file-upload","weak-password","frontend-reverse","app-audit","env","waf","path-traversal","ssrf","tunnel","sqli","deserialization","protocol","privesc"];
  for (const key of keys) {
    assert.notEqual(knowledgeCategory(key), key);
    assert.match(knowledgeCategorySummary(key), /[\u3400-\u9fff]/u, key);
  }
  assert.equal(knowledgeCategorySummary("unknown"), undefined);
  assert.equal(knowledgeCategorySummary("constructor"), undefined);
});
