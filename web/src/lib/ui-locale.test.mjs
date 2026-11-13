import assert from "node:assert/strict";
import test from "node:test";
import { parseUiLocale, uiText, navigationTitle } from "./ui-locale.ts";
test("defaults to Simplified Chinese; accepts only the supported English locale", () => {
  assert.equal(parseUiLocale(null), "zh-CN");
  assert.equal(parseUiLocale("fr"), "zh-CN");
  assert.equal(parseUiLocale("en"), "en");
});
test("labels switch languages without changing stable navigation identifiers", () => {
  assert.equal(uiText("zh-CN", "搜索", "Search"), "搜索");
  assert.equal(uiText("en", "搜索", "Search"), "Search");
  assert.equal(navigationTitle("en", "tasks", "任务"), "Tasks");
  assert.equal(navigationTitle("zh-CN", "tasks", "任务"), "任务");
  assert.equal(navigationTitle("en", "custom", "Custom"), "Custom");
  assert.equal(navigationTitle("en", "constructor", "Custom"), "Custom");
});
