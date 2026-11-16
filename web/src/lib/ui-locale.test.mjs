import assert from "node:assert/strict";
import test from "node:test";
import { parseUiLocale, uiText, navigationTitle } from "./ui-locale.ts";
test("defaults to Simplified Chinese; accepts the supported English and Traditional Chinese locales", () => {
  assert.equal(parseUiLocale(null), "zh-CN");
  assert.equal(parseUiLocale("fr"), "zh-CN");
  assert.equal(parseUiLocale("en"), "en");
  assert.equal(parseUiLocale("zh-TW"), "zh-TW");
});
test("labels switch languages without changing stable navigation identifiers", () => {
  assert.equal(uiText("zh-CN", "搜索", "Search"), "搜索");
  assert.equal(uiText("en", "搜索", "Search"), "Search");
  assert.equal(navigationTitle("en", "tasks", "任务"), "Tasks");
  assert.equal(navigationTitle("zh-CN", "tasks", "任务"), "任务");
  assert.equal(navigationTitle("en", "custom", "Custom"), "Custom");
  assert.equal(navigationTitle("en", "constructor", "Custom"), "Custom");
});

test("Traditional Chinese navigation and report labels have a safe fallback", () => {
  assert.equal(uiText("zh-TW", "报告中心", "Report center"), "報告中心");
  assert.equal(uiText("zh-TW", "保存", "Save"), "儲存");
  assert.equal(uiText("zh-TW", "未翻译的界面", "Untranslated UI"), "未翻译的界面");
  assert.equal(uiText("zh-TW", "constructor", "Constructor"), "constructor");
  assert.equal(navigationTitle("zh-TW", "reports", "报告中心"), "報告中心");
  assert.equal(navigationTitle("zh-TW", "constructor", "Custom"), "Custom");
});
