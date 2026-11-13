import { navigationTitle, uiText, type UiLocale } from "@/lib/ui-locale";
import {
  Activity,
  Ban,
  BellRing,
  Bot,
  BookOpen,
  Brain,
  Bug,
  ClipboardList,
  FolderOpen,
  FolderSync,
  LayoutDashboard,
  type LucideIcon,
  MessageSquare,
  Network,
  Plug,
  Radio,
  ScrollText,
  Settings2,
  ShieldAlert,
  ShieldCheck,
  Sparkles,
  Target,
  Terminal,
  Wrench,
} from "lucide-react";

export type NavBadge = "new" | "soon";

export interface NavSubItem {
  id: string;
  title: string;
  url: string;
  icon?: LucideIcon;
  badge?: NavBadge;
  disabled?: boolean;
  newTab?: boolean;
}

interface NavItemBase {
  id: string;
  title: string;
  icon?: LucideIcon;
  badge?: NavBadge;
  disabled?: boolean;
  newTab?: boolean;
}

export interface NavMainLinkItem extends NavItemBase {
  url: string;
  subItems?: never;
}

export interface NavMainParentItem extends NavItemBase {
  subItems: NavSubItem[];
}

export type NavMainItem = NavMainLinkItem | NavMainParentItem;

export interface NavGroup {
  id: number;
  label?: string;
  items: NavMainItem[];
}

export const sidebarItems: NavGroup[] = [
  {
    id: 1,
    label: "工作台",
    items: [
      { id: "dashboard", title: "总览", url: "/dashboard", icon: LayoutDashboard },
      { id: "chat", title: "对话", url: "/chat", icon: MessageSquare },
      { id: "tasks", title: "任务", url: "/function/tasks", icon: Target },
      { id: "findings", title: "发现结果", url: "/function/findings", icon: Bug },
      { id: "traffic", title: "流量记录", url: "/function/traffic", icon: Activity },
      { id: "commands", title: "工具执行", url: "/function/commands", icon: Terminal },
      { id: "llm-records", title: "模型记录", url: "/function/llm-records", icon: Radio },
      { id: "assets", title: "资产", url: "/function/assets", icon: Network },
      { id: "sync", title: "资产同步", url: "/function/sync", icon: FolderSync },
      { id: "workspace", title: "工作区", url: "/function/workspace", icon: FolderOpen },
    ],
  },
  {
    id: 2,
    label: "系统与配置",
    items: [
      { id: "welcome", title: "开始使用 / 配置引导", url: "/welcome", icon: BookOpen },
      { id: "knowledge", title: "浏览知识库", url: "/system/knowledge", icon: BookOpen },
      { id: "llm", title: "模型配置", url: "/system/llm", icon: Brain },
      { id: "agents", title: "智能体", url: "/system/agents", icon: Bot },
      { id: "mcp", title: "MCP 服务", url: "/system/mcp", icon: Plug },
      { id: "skills", title: "技能库", url: "/system/skills", icon: Sparkles },
      { id: "tools", title: "工具管理", url: "/system/tools", icon: Wrench },
      { id: "notify", title: "通知", url: "/system/notify", icon: BellRing },
      { id: "intercept", title: "拦截规则", url: "/system/intercept", icon: ShieldAlert },
      { id: "asset-intercept", title: "资产拦截", url: "/system/intercept/assets", icon: Ban },
      { id: "approvals", title: "审批记录", url: "/system/intercept/approvals", icon: ClipboardList },
      { id: "logs", title: "日志", url: "/system/logs", icon: ScrollText },
      { id: "privacy", title: "隐私与暴露面", url: "/system/privacy", icon: ShieldCheck },
      { id: "settings", title: "系统设置", url: "/system/settings", icon: Settings2 },
      { id: "help", title: "使用指南", url: "/help", icon: BookOpen },
    ],
  },
];

export function localizedSidebarItems(locale: UiLocale): NavGroup[] {
  return sidebarItems.map(group => ({
    ...group,
    label: group.id === 1 ? uiText(locale, "工作台", "Workspace") : group.id === 2 ? uiText(locale, "系统与配置", "System & configuration") : group.label,
    items: group.items.map(item => item.subItems ? { ...item, title: navigationTitle(locale, item.id, item.title), subItems: item.subItems.map(sub => ({ ...sub, title: navigationTitle(locale, sub.id, sub.title) })) } : { ...item, title: navigationTitle(locale, item.id, item.title) }),
  }));
}
