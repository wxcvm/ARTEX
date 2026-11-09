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

export function buildFindingReportMarkdown({ finding, lineage, traffic, packets, retests }: FindingReportData): string {
  const lines = [
    "# " + inline(finding.name || finding.vulnclass || "未命名漏洞"), "",
    "- 漏洞 ID：" + inline(finding.id),
    "- 类型：" + inline(finding.vulnclass || "未分类"),
    "- 严重程度：" + inline(finding.severity),
    "- 状态：" + inline(finding.status),
    "- 发现时间：" + inline(finding.ts),
    "- 任务：" + inline(finding.task_description || finding.task_id || "未知"),
  ];
  if (finding.source_task_id) lines.push("- 来源任务：" + inline(finding.source_task_id));
  if (finding.inherited) lines.push("- 继承漏洞：是");
  if (finding.assets?.length) lines.push("- 资产：" + finding.assets.map((a) => inline(a.type + " " + a.label)).join("；"));
  lines.push("", "## 漏洞摘要", "", finding.summary || "（暂无）", "",
    "## 证据 / PoC", "", finding.evidence ? fence(finding.evidence, "text") : "（暂无）", "",
    "## 详细报告", "", finding.report || "（暂无）", "");
  if (finding.report_stale) lines.push("**注意：** 证据已发生变更，详细报告可能尚未更新。", "");
  lines.push("## 关联 HTTP 流量证据", "");
  if (!traffic.bindings.length) lines.push("（暂无）", "");
  for (const binding of traffic.bindings) {
    const item = packets.find((packet) => packet.binding.id === binding.id);
    const snap = binding.snapshot;
    lines.push("### " + inline(binding.role + " · " + snap.method + " " + snap.url), "",
      "- 流量 ID：" + inline(snap.source_traffic_id),
      "- 状态码：" + snap.status,
      "- 说明：" + inline(binding.note || "无"), "");
    if (snap.req_head) lines.push("**请求头**", "", fence(snap.req_head, "http"), "");
    if (item) lines.push("**请求正文**", "",
      item.request.binary ? "（二进制正文未展开）" : fence(item.request.content || "（空）", "http"),
      item.request.truncated ? "（预览已截断，请在 ARTEX 中下载原始流量）" : "", "");
    if (snap.resp_head) lines.push("**响应头**", "", fence(snap.resp_head, "http"), "");
    if (item) lines.push("**响应正文**", "",
      item.response.binary ? "（二进制正文未展开）" : fence(item.response.content || "（空）", "http"),
      item.response.truncated ? "（预览已截断，请在 ARTEX 中下载原始流量）" : "", "");
  }
  lines.push("## 复测记录", "");
  if (!retests.length) lines.push("（暂无）", "");
  for (const r of retests) {
    lines.push("### #" + r.id, "",
      "- 时间：" + inline(r.created_at), "- 状态：" + inline(r.status),
      "- 结论：" + inline(r.verdict || "未给出"), "");
    if (r.notes) lines.push("**备注**", "", r.notes, "");
    if (r.summary) lines.push("**摘要**", "", r.summary, "");
    if (r.evidence) lines.push("**证据**", "", fence(r.evidence, "text"), "");
    if (r.error) lines.push("**错误**", "", inline(r.error), "");
  }
  lines.push("## 漏洞溯源", "");
  if (!lineage.nodes.length) lines.push("（暂无）", "");
  for (const node of lineage.nodes) lines.push("- " + inline(node.id) + " · " + inline(node.type) +
    " · " + inline(node.state) + (node.payload ? " · " + inline(node.payload) : ""));
  if (lineage.edges.length) lines.push("", "### 溯源关系", "",
    ...lineage.edges.map((edge) => "- " + inline(edge.src) + " → " + inline(edge.rel) + " → " + inline(edge.dst)));
  return lines.join("\n").trimEnd() + "\n";
}
