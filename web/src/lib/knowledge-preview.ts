export type AttachmentStatus = "none" | "match" | "conflict" | "missing" | "ambiguous";
export const attachmentLabels: Record<AttachmentStatus, string> = {
  none: "无附件", match: "正文一致", conflict: "正文有差异", missing: "附件缺失", ambiguous: "附件不唯一",
};
export interface KnowledgeSource { name: string; sha256: string }
export interface KnowledgeStatus {
  enabled: boolean;
  mode: "read_only_preview";
  source?: KnowledgeSource;
  summary?: { total: number; conflicts: number; missing: number; with_content: number };
}
export interface KnowledgeItem {
  id: string; code: string; title: string; kind: string; category: string;
  component: string; cve: string; tags: string; source_verified: boolean;
  attachment_status: AttachmentStatus;
}
export interface KnowledgeList {
  items: KnowledgeItem[]; total: number; offset: number; limit: number; source: KnowledgeSource;
}
export interface KnowledgeEntry extends Omit<KnowledgeItem, "attachment_status"> {
  versions: string; source: string; source_url: string; description: string; usage: string;
  content: string; content_sha256: string;
  attachment: { path: string; content: string; sha256: string; status: AttachmentStatus };
}
export function sourceClaimLabel(claim: boolean): string {
  return claim ? "来源声称已验证 · 本项目未复核" : "来源未验证 · 本项目未复核";
}
