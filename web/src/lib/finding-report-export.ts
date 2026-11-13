import type { Edge, Finding, FindingRetest, FindingTraffic, FindingTrafficDetail, TaskNode } from "./types";

type Lineage = { nodes: TaskNode[]; edges: Edge[] };

export interface FindingReportSource {
  getFinding(id: string, contextTask?: string): Promise<Finding>;
  findingLineage(id: string): Promise<Lineage>;
  findingTraffic(id: string, contextTask?: string): Promise<FindingTraffic>;
  findingTrafficDetail(id: string, bindingId: string, contextTask?: string): Promise<FindingTrafficDetail>;
  findingRetests(id: string): Promise<FindingRetest[]>;
}

export interface FindingReportData {
  finding: Finding;
  lineage: Lineage;
  traffic: FindingTraffic;
  packets: FindingTrafficDetail[];
  retests: FindingRetest[];
}

// All records must be present; partial API failures are not silently exported.
export async function collectFindingReport(source: FindingReportSource, id: string, contextTask?: string): Promise<FindingReportData> {
  const [finding, lineage, traffic, retests] = await Promise.all([
    source.getFinding(id, contextTask), source.findingLineage(id),
    source.findingTraffic(id, contextTask), source.findingRetests(id),
  ]);
  const packets = await Promise.all(traffic.bindings.map((binding) =>
    source.findingTrafficDetail(id, binding.id, contextTask),
  ));
  return { finding, lineage, traffic, packets, retests };
}

const tick = String.fromCharCode(96);

function inline(value: string): string {
  return value.replace(/[\\_*{}[\]()#+!|>~]/g, "\\$&")
    .split(tick).join("\\" + tick).replace(/\r?\n/g, " ");
}

function fence(value: string, lang = ""): string {
  const runs = value.match(new RegExp(tick + "+", "g")) ?? [];
  const marker = tick.repeat(Math.max(3, ...runs.map((run) => run.length + 1)));
  return marker + lang + "\n" + value.trimEnd() + "\n" + marker;
}

export function findingReportFilename(finding: Finding): string {
  const title = (finding.name || finding.vulnclass || "finding").normalize("NFKC")
    .replace(/[\\/:*?"<>|\x00-\x1f]/g, "-").replace(/\s+/g, " ").trim().slice(0, 80) || "finding";
  return "finding-" + finding.id.replace(/[^a-z0-9_-]/gi, "-") + "-" + title + ".md";
}

export function buildFindingReportMarkdown({ finding, lineage, traffic, packets, retests }: FindingReportData, locale: "zh-CN" | "en" = "zh-CN"): string {
  const tr = (zh: string, en: string) => locale === "en" ? en : zh;
  const lines = [
    "# " + inline(finding.name || finding.vulnclass || tr("未命名漏洞", "Untitled finding")), "",
    tr("- 漏洞 ID：", "- Finding ID: ") + inline(finding.id),
    tr("- 类型：", "- Class: ") + inline(finding.vulnclass || tr("未分类", "Uncategorized")),
    tr("- 严重程度：", "- Severity: ") + inline(finding.severity),
    tr("- 状态：", "- Status: ") + inline(finding.status),
    tr("- 发现时间：", "- Discovered at: ") + inline(finding.ts),
    tr("- 任务：", "- Task: ") + inline(finding.task_description || finding.task_id || tr("未知", "Unknown")),
  ];
  if (finding.source_task_id) lines.push(tr("- 来源任务：", "- Source task: ") + inline(finding.source_task_id));
  if (finding.inherited) lines.push(tr("- 继承漏洞：是", "- Inherited finding: yes"));
  if (finding.assets?.length) lines.push(tr("- 资产：", "- Assets: ") + finding.assets.map((a) => inline(a.type + " " + a.label)).join(tr("；", "; ")));
  lines.push("", tr("## 漏洞摘要", "## Finding summary"), "", finding.summary || tr("（暂无）", "(None)"), "",
    tr("## 证据 / PoC", "## Evidence / PoC"), "", finding.evidence ? fence(finding.evidence, "text") : tr("（暂无）", "(None)"), "",
    tr("## 详细报告", "## Detailed report"), "", finding.report || tr("（暂无）", "(None)"), "");
  if (finding.report_stale) lines.push(tr("**注意：** 证据已发生变更，详细报告可能尚未更新。", "**Note:** Evidence changed; the detailed report may not yet be updated."), "");
  lines.push(tr("## 关联 HTTP 流量证据", "## Related HTTP traffic evidence"), "");
  if (!traffic.bindings.length) lines.push(tr("（暂无）", "(None)"), "");
  for (const binding of traffic.bindings) {
    const item = packets.find((packet) => packet.binding.id === binding.id);
    const snap = binding.snapshot;
    lines.push("### " + inline(binding.role + " · " + snap.method + " " + snap.url), "",
      tr("- 流量 ID：", "- Traffic ID: ") + inline(snap.source_traffic_id),
      tr("- 状态码：", "- Status code: ") + snap.status,
      tr("- 说明：", "- Note: ") + inline(binding.note || tr("无", "None")), "");
    if (snap.req_head) lines.push(tr("**请求头**", "**Request headers**"), "", fence(snap.req_head, "http"), "");
    if (item) lines.push(tr("**请求正文**", "**Request body**"), "",
      item.request.binary ? tr("（二进制正文未展开）", "(Binary body not expanded)") : fence(item.request.content || tr("（空）", "(Empty)"), "http"),
      item.request.truncated ? tr("（预览已截断，请在 ARTEX 中下载原始流量）", "(Preview truncated; download original traffic in ARTEX)") : "", "");
    if (snap.resp_head) lines.push(tr("**响应头**", "**Response headers**"), "", fence(snap.resp_head, "http"), "");
    if (item) lines.push(tr("**响应正文**", "**Response body**"), "",
      item.response.binary ? tr("（二进制正文未展开）", "(Binary body not expanded)") : fence(item.response.content || tr("（空）", "(Empty)"), "http"),
      item.response.truncated ? tr("（预览已截断，请在 ARTEX 中下载原始流量）", "(Preview truncated; download original traffic in ARTEX)") : "", "");
  }
  lines.push(tr("## 复测记录", "## Retest records"), "");
  if (!retests.length) lines.push(tr("（暂无）", "(None)"), "");
  for (const r of retests) {
    lines.push("### #" + r.id, "",
      tr("- 时间：", "- Time: ") + inline(r.created_at), tr("- 状态：", "- Status: ") + inline(r.status),
      tr("- 结论：", "- Verdict: ") + inline(r.verdict || tr("未给出", "Not provided")), "");
    if (r.notes) lines.push(tr("**备注**", "**Notes**"), "", r.notes, "");
    if (r.summary) lines.push(tr("**摘要**", "**Summary**"), "", r.summary, "");
    if (r.evidence) lines.push(tr("**证据**", "**Evidence**"), "", fence(r.evidence, "text"), "");
    if (r.error) lines.push(tr("**错误**", "**Error**"), "", inline(r.error), "");
  }
  lines.push(tr("## 漏洞溯源", "## Finding lineage"), "");
  if (!lineage.nodes.length) lines.push(tr("（暂无）", "(None)"), "");
  for (const node of lineage.nodes) lines.push("- " + inline(node.id) + " · " + inline(node.type) +
    " · " + inline(node.state) + (node.payload ? " · " + inline(node.payload) : ""));
  if (lineage.edges.length) lines.push("", tr("### 溯源关系", "### Lineage relations"), "",
    ...lineage.edges.map((edge) => "- " + inline(edge.src) + " → " + inline(edge.rel) + " → " + inline(edge.dst)));
  return lines.join("\n").trimEnd() + "\n";
}
