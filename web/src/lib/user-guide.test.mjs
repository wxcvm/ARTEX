import assert from "node:assert/strict";
import { readFileSync } from "node:fs";
import test from "node:test";
import { guideHref, guideLinks, guideSlug, parseGuide } from "./user-guide.ts";

test("Chinese headings retain GitHub-style anchors", () => {
  assert.equal(guideSlug("2. Windows 11 从源码安装（推荐使用 fork 最新代码）"),
    "2-windows-11-从源码安装推荐使用-fork-最新代码");
  assert.equal(guideSlug("3. Linux / macOS 安装"), "3-linux--macos-安装");
  assert.equal(guideSlug("10. 历史经验（Lesson）与效果评估"),
    "10-历史经验lesson与效果评估");
});

test("markdown guide is split once into sections without copying its body", () => {
  const source = [
    "# ARTEX 中文使用手册",
    "> 第一段介绍。",
    "",
    "## 目录",
    "- [第一节](#1-第一节)",
    "## 1. 第一节",
    "模型配置和工作步骤。",
    "### 注意",
    "一段详细说明。",
    "## 2. 第二节",
    "MCP 说明。",
  ].join("\n");
  const guide = parseGuide(source);
  assert.match(guide.introduction, /第一段介绍/);
  assert.doesNotMatch(guide.introduction, /^# /m);
  assert.deepEqual(guide.sections.map((section) => section.id), ["目录", "1-第一节", "2-第二节"]);
  assert.match(guide.sections[1].content, /### 注意/);
  assert.doesNotMatch(guide.sections[1].content, /MCP 说明/);
});

test("relative links resolve to repository files, not /help; anchors stay local", () => {
  const original = [
    "[README](../README.md#backups)",
    "[经验](rsi-lessons.md)",
    "![任务截图](../screenshots/tasks.png)",
    "[本地](#1-第一节)",
    "[网站](https://example.org/guide)",
  ].join("\n");
  const linked = guideLinks(original);
  assert.match(linked, /\]\(https:\/\/github\.com\/wxcvm\/ARTEX\/blob\/main\/README\.md#backups\)/);
  assert.match(linked, /\]\(https:\/\/github\.com\/wxcvm\/ARTEX\/blob\/main\/docs\/rsi-lessons\.md\)/);
  assert.match(linked, /!\[任务截图\]\(https:\/\/raw\.githubusercontent\.com\/wxcvm\/ARTEX\/main\/screenshots\/tasks\.png\)/);
  assert.match(linked, /\[本地\]\(#1-第一节\)/);
  assert.match(linked, /\[网站\]\(https:\/\/example\.org\/guide\)/);
});


test("contextual help links resolve to real section or subsection headings in the manual", () => {
  const manual = readFileSync(new URL("../../../docs/USER_GUIDE_ZH.md", import.meta.url), "utf8");
  const headingIds = new Set(
    [...manual.matchAll(/^#{2,3} (.+)\r?$/gm)].map((match) => guideSlug(match[1])),
  );
  for (const path of [
    "/function/tasks", "/function/tasks/detail", "/system/llm", "/system/skills",
    "/system/mcp", "/system/agents", "/system/tools", "/system/settings",
    "/system/logs", "/function/traffic", "/function/assets", "/function/findings",
    "/function/sync",
  ]) {
    const href = guideHref(path);
    assert.ok(href.startsWith("/help/#"), path + " has no specific help link");
    assert.ok(headingIds.has(href.split("#")[1]), path + " links to a non-existent heading: " + href);
    assert.equal(guideHref(path + "/"), href, "trailing slash changes help target");
  }
  assert.equal(guideHref("/chat"), "/help/");
  for (const heading of [
    "1. 快速选择安装方式", "5.1 创建模型 Profile", "6.1 创建任务",
    "6.3 任务详情页逐项看", "8.2 从 Findings 批量导出",
  ]) {
    assert.ok(headingIds.has(guideSlug(heading)), "quick start has no target: " + heading);
  }
});
