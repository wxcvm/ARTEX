export const UI_LOCALE_KEY = "artex_ui_locale";
export const DEFAULT_UI_LOCALE = "zh-CN" as const;
export type UiLocale = "zh-CN" | "en";
export function parseUiLocale(value: string | null | undefined): UiLocale {
  return value === "en" ? "en" : DEFAULT_UI_LOCALE;
}
export function uiText(locale: UiLocale, zh: string, en: string): string {
  return locale === "en" ? en : zh;
}
const englishNavigation: Record<string, string> = {
  dashboard: "Overview", chat: "Chat", tasks: "Tasks", findings: "Findings",
  traffic: "Traffic records", commands: "Tool execution", "llm-records": "Model records",
  assets: "Assets", sync: "Asset sync", workspace: "Workspace",
  welcome: "Getting started / Setup", knowledge: "Knowledge base", llm: "Model configuration",
  agents: "Agents", mcp: "MCP servers", skills: "Skills", tools: "Tool management",
  notify: "Notifications", intercept: "Interception rules", "asset-intercept": "Asset interception",
  approvals: "Approval records", logs: "Logs", privacy: "Privacy & exposure",
  settings: "System settings", help: "User guide",
};
export function navigationTitle(locale: UiLocale, id: string, original: string): string {
  return locale === "en" && Object.hasOwn(englishNavigation, id) ? englishNavigation[id] : original;
}
