import assert from "node:assert/strict";
import test from "node:test";

const { buildFindingReportMarkdown, collectFindingReport, findingReportFilename } =
  await import("./finding-report-export.ts");

function sample() {
  const finding = {
    id: "17", name: "SQL 注入", vulnclass: "SQL Injection", severity: "high",
    status: "confirmed", ts: "2026-11-09", summary: "摘要",
    evidence: "line\n" + String.fromCharCode(96).repeat(4) + "\nbody",
    report: "## 复现方法\n内容", report_stale: true,
    assets: [{ id: "a", type: "ip", label: "192.0.2.1" }],
  };
  const binding = { id: "55", role: "proof", note: "证据", snapshot: {
    method: "POST", url: "https://example.test", source_traffic_id: "tr-1", status: 500,
    req_head: "Header: one", resp_head: "HTTP/1.1 500",
  }};
  const traffic = { bindings: [binding] };
  const packets = [{ binding,
    request: { content: "payload", binary: false, truncated: false },
    response: { content: "raw", binary: true, truncated: true },
  }];
  const lineage = { nodes: [{ id: "1", type: "intent", state: "done", payload: "probe" }],
    edges: [{ src: "1", rel: "yields", dst: "17" }] };
  const retests = [{ id: 2, status: "completed", verdict: "reproduced", created_at: "today",
    notes: "check", summary: "still present", evidence: "POC", error: "" }];
  return { finding, traffic, packets, lineage, retests };
}

test("exports full report with traffic, retests and lineage", () => {
  const data = sample();
  const md = buildFindingReportMarkdown(data);
  for (const needle of ["SQL 注入", "复现方法", "tr-1", "still present",
    "probe", "二进制正文未展开", "预览已截断", "报告可能尚未更新",
    String.fromCharCode(96).repeat(5) + "text"]) {
    assert.ok(md.includes(needle), "missing " + needle);
  }
  assert.ok(!md.includes("raw\n"));
  assert.equal(findingReportFilename({ ...data.finding, name: "../report|test" }),
    "finding-17-..-report-test.md");
});

test("preserves inherited task context and rejects partial failure", async () => {
  const data = sample();
  const calls = [];
  const source = {
    getFinding: async (...args) => (calls.push(["finding", ...args]), data.finding),
    findingLineage: async (...args) => (calls.push(["lineage", ...args]), data.lineage),
    findingTraffic: async (...args) => (calls.push(["traffic", ...args]), data.traffic),
    findingTrafficDetail: async (...args) => (calls.push(["detail", ...args]), data.packets[0]),
    findingRetests: async (...args) => (calls.push(["retests", ...args]), data.retests),
  };
  const out = await collectFindingReport(source, "17", "t-2");
  assert.deepEqual(out.packets, data.packets);
  assert.deepEqual(calls.at(-1), ["detail", "17", "55", "t-2"]);
  source.findingTrafficDetail = async () => { throw Error("offline"); };
  await assert.rejects(collectFindingReport(source, "17"), /offline/);
});

test("empty evidence sections remain explicit", () => {
  const d = sample();
  d.traffic.bindings = [];
  d.packets = [];
  d.retests = [];
  d.lineage.nodes = [];
  assert.ok(buildFindingReportMarkdown(d).includes("## 复测记录\n\n（暂无）"));
});


test("localized standalone export changes headings but preserves report and evidence", () => {
  const data = sample();
  const zh = buildFindingReportMarkdown(data, "zh-CN");
  const en = buildFindingReportMarkdown(data, "en");
  assert.match(zh, /## 详细报告/);
  assert.match(en, /## Detailed report/);
  for (const md of [zh, en]) {
    assert.ok(md.includes(data.finding.report));
    assert.ok(md.includes(data.finding.evidence));
    assert.ok(md.includes(data.packets[0].request.content));
    assert.ok(md.includes("Header: one"));
  }
});
